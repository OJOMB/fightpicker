package users

import (
	"context"

	"github.com/pkg/errors"

	pb "github.com/OJOMB/fightpicker/internal/grpc/gen/go/users/v1"
	usersservice "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

// Server implements pb.UsersServiceServer
type Server struct {
	pb.UsersServiceServer

	id id.UUID7GeneratorParser

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

// CreateUser registers a new user in the system
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	user, err := s.service.CreateUser(ctx, createUserRequestDTOtoIDO(req))
	if err != nil {
		return nil, s.toStatus(err)
	}

	return userIDOtoDTO(user), nil
}

// GetUser retrieves a user by their ID.
func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	userID, err := s.id.ParseString(req.GetUserId())
	if err != nil {
		return nil, s.toStatus(errors.Wrap(usersservice.ErrInvalidParameter, "invalid user ID format"))
	}

	user, err := s.service.GetUserByID(ctx, userID)
	if err != nil {
		return nil, s.toStatus(err)
	}

	return userIDOtoDTO(user), nil
}

// UpdateUser updates an existing user's information.
func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	userID, err := s.id.ParseString(req.GetUserId())
	if err != nil {
		return nil, s.toStatus(errors.Wrap(usersservice.ErrInvalidParameter, "invalid user ID format"))
	}

	updatedUser, err := s.service.UpdateUser(ctx, userID, updateUserRequestDTOtoIDO(req))
	if err != nil {
		return nil, s.toStatus(err)
	}

	return userIDOtoDTO(updatedUser), nil
}
