package users

import (
	"context"
)

// UserCreator defines the interface for creating a new user.
type UserCreator interface {
	CreateUser(ctx context.Context, user User) error
}

// CreateUser creates a new user in the system (without profile picture).
func (svc *Service) CreateUser(ctx context.Context, user User) (User, error) {
	if err := svc.validateCreationReq(&user); err != nil {
		return User{}, err
	}

	user.ID = svc.idTool.Generate()
	now := svc.dateTimeTool.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.UpdatedBy = user.ID

	// NB profile picture is handled in a separate process after creation

	hashedPassword, err := svc.authTool.HashPassword(user.PasswordHash)
	if err != nil {
		svc.logger.DebugContext(ctx, "failed to hash password", "error", err, "user_id", user.ID)
		return User{}, err
	}

	user.PasswordHash = hashedPassword

	if err := svc.repo.CreateUser(ctx, user); err != nil {
		svc.logger.DebugContext(ctx, "failed to create user", "error", err, "user_id", user.ID)
		return User{}, err
	}

	// we now want to publish a message to the user created topic for further processing (e.g. verification email)
	svc.logger.DebugContext(ctx, "successfully created user", "user_id", user.ID)

	return user, nil
}
