package users

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/twmb/franz-go/pkg/kgo"

	usersservice "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
)

const defaultUserRole = "user"

// CreateUser converts a User IDO to a DBO and calls the repo to create the user in the database
// It also assigns the default "user" role to the newly created user.
func (r *Repo) CreateUser(ctx context.Context, user usersservice.User) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	// rollback is a no-op if the transaction is already committed so this is safe
	defer tx.Rollback(ctx)

	qs := r.dbClient.WithTx(tx)

	dbParams := UserIDOtoCreateUserParamsDBO(user)
	if err := qs.CreateUser(ctx, dbParams); err != nil {
		return dbErrorToServiceError(err)
	}

	pgTypeNow := pgtype.Timestamptz{
		Time:  r.dateTimeTool.Now(),
		Valid: true,
	}

	pgTypeUser := pgtype.UUID{
		Bytes: user.ID,
		Valid: true,
	}

	// Assign default role to user
	assignRoleParams := postgres.AssignRoleToUserByRoleNameParams{
		UserID:    user.ID,
		Name:      defaultUserRole,
		CreatedAt: pgTypeNow,
		CreatedBy: pgTypeUser,
		UpdatedAt: pgTypeNow,
		UpdatedBy: pgTypeUser,
	}
	if err := qs.AssignRoleToUserByRoleName(ctx, assignRoleParams); err != nil {
		r.logger.ErrorContext(ctx, "default role not in roles table")
		return dbErrorToServiceError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return dbErrorToServiceError(err)
	}

	// publish message to kafka topic for further processing (email welcome and validation etc)
	// TODO: add outbox pattern for reliable delivery
	r.publishUserCreated(user)

	return nil
}

func (r *Repo) publishUserCreated(user usersservice.User) {
	if r.events.client == nil {
		return
	}

	key := user.ID.Bytes()
	value, _ := json.Marshal(struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
	}{
		Email:     user.Email,
		FirstName: user.FirstName,
	})

	// use a fresh context - we don't want the req ctx to time out because the user creation request has finished before the callback is invoked
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	r.events.client.Produce(ctx, &kgo.Record{
		Topic: r.events.topicPostUserCreate,
		Key:   key,
		Value: value,
	}, func(record *kgo.Record, err error) {
		defer cancel()
		if err != nil {
			r.logger.ErrorContext(ctx, "failed to produce user.created", "error", err, "record_key", record.Key)
		}
	})
}
