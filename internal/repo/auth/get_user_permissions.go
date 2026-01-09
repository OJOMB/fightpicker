package auth

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"

	service "github.com/OJOMB/fightpicker/internal/service/auth"
)

func (r *Repo) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, service.Permissions, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, nil, err
	}

	// rollback is a no-op if the transaction is already committed so this is safe
	defer tx.Rollback(ctx)

	roles, err := r.client.GetUserRolesByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	if len(roles) == 0 {
		r.logger.ErrorContext(ctx, "user without any roles", "user_id", userID)
		return nil, service.Permissions{}, ErrInternalError
	}

	permissionRows, err := r.client.GetUserPermissionsByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	// map[version]map[resource]map[operation]name
	var perms = make(service.Permissions)
	for _, row := range permissionRows {
		version := row.Version
		resource := row.Resource
		operation := row.Operation
		name := row.Name

		if _, ok := perms[version]; !ok {
			perms[version] = make(map[string]map[string]map[string]struct{})
		}

		if _, ok := perms[version][string(resource)]; !ok {
			perms[version][string(resource)] = make(map[string]map[string]struct{})
		}

		if _, ok := perms[version][string(resource)][operation]; !ok {
			perms[version][string(resource)][operation] = make(map[string]struct{})
		}

		perms[version][string(resource)][operation][name] = struct{}{}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	return roles, perms, nil
}
