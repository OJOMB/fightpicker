package users

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/pkg/errors"

	pb "github.com/OJOMB/fightpicker/internal/grpc/gen/go/users/v1"
	usersservice "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

// Server implements pb.UsersServiceServer
type Server struct {
	pb.UsersServiceServer

	service *usersservice.Service
	logger  *logs.Logger
}

// NewServer is a constructor for the Users gRPC server.
func NewServer(service *usersservice.Service, logger *logs.Logger) *Server {
	return &Server{
		service: service,
		logger:  logger,
	}
}

// GetUser retrieves a user by their ID.
func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	userID, err := uuid.FromString(req.UserId.Value)
	if err != nil {
		return nil, s.toStatus(errors.Wrap(usersservice.ErrInvalidParameter, "invalid user ID format"))
	}

	user, err := s.service.GetUserByID(ctx, userID)
	if err != nil {
		return nil, s.toStatus(err)
	}

	return toProtoUser(user), nil
}
