package users

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/OJOMB/fightpicker/internal/grpc/gen/go/users/v1"
	usersservice "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

var (
	defaultPageSize = 20
	maxPageSize     = 100
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

type Paginator interface {
	PageSizeGetter
	LastSeenIDGetter
}

type PageSizeGetter interface {
	GetPageSize() uint32
}

type LastSeenIDGetter interface {
	GetLastSeenId() string
}

// ParsePaginationParams is a helper that parses pagination parameters from the HTTP request.
func (s *Server) ParsePaginationParams(p Paginator) (pageSize int, lastSeenID *id.UUID7, err error) {
	pageSize, err = s.parsePageSize(p)
	if err != nil {
		return 0, nil, err
	}

	lastSeenID, err = s.parseLastSeenID(p)
	if err != nil {
		return 0, nil, err
	}

	return pageSize, lastSeenID, nil
}

// parsePageSize is a pagination helper that parses the page size from the query parameter string.
func (s *Server) parsePageSize(p PageSizeGetter) (int, error) {
	reqPageSize := p.GetPageSize()

	if reqPageSize == 0 {
		return defaultPageSize, nil
	}

	if reqPageSize > uint32(maxPageSize) {
		return maxPageSize, nil
	}

	return int(reqPageSize), nil
}

// parseLastSeenID is a pagination helper that parses the last seen ID from the query parameter string.
func (s *Server) parseLastSeenID(p LastSeenIDGetter) (*id.UUID7, error) {
	lastSeenIDStr := p.GetLastSeenId()
	lastSeenID, err := s.id.ParseString(lastSeenIDStr)
	if err != nil {
		return nil, errors.Wrap(usersservice.ErrInvalidParameter, "last_seen_user_id")
	}

	return &lastSeenID, nil
}

// CreateUser registers a new user in the system
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	user, err := s.service.CreateUser(ctx, createUserRequestDTOtoIDO(req))
	if err != nil {
		return nil, s.toStatus(err)
	}

	return userIDOtoDTO(user), nil
}

// GetUserByID retrieves a user by their ID.
func (s *Server) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.User, error) {
	userID, err := s.id.ParseString(req.GetUserId())
	if err != nil {
		return nil, s.toStatus(errors.Wrap(usersservice.ErrInvalidParameter, "user_id"))
	}

	user, err := s.service.GetUserByID(ctx, userID)
	if err != nil {
		return nil, s.toStatus(err)
	}

	return userIDOtoDTO(user), nil
}

// GetUserByEmail retrieves a user by their email.
func (s *Server) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.User, error) {
	user, err := s.service.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, s.toStatus(err)
	}

	return userIDOtoDTO(user), nil
}

// GetUserByUsername retrieves a user by their username.
func (s *Server) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.User, error) {
	user, err := s.service.GetUserByUsername(ctx, req.GetUsername())
	if err != nil {
		return nil, s.toStatus(err)
	}

	return userIDOtoDTO(user), nil
}

// UpdateUser updates an existing user's information.
func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	userID, err := s.id.ParseString(req.GetUserId())
	if err != nil {
		return nil, s.toStatus(errors.Wrap(usersservice.ErrInvalidParameter, "user_id"))
	}

	updatedUser, err := s.service.UpdateUser(ctx, userID, updateUserRequestDTOtoIDO(req))
	if err != nil {
		return nil, s.toStatus(err)
	}

	return userIDOtoDTO(updatedUser), nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	userID, err := s.id.ParseString(req.GetUserId())
	if err != nil {
		return nil, s.toStatus(errors.Wrap(usersservice.ErrInvalidParameter, "user_id"))
	}

	if err := s.service.DeleteUserByID(ctx, userID); err != nil {
		return nil, s.toStatus(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	pageSize, lastSeenId, err := s.ParsePaginationParams(req)
	if err != nil {
		return nil, s.toStatus(err)
	}

	users, totalCount, err := s.service.ListUsers(ctx, pageSize, lastSeenId)
	if err != nil {
		return nil, s.toStatus(err)
	}

	return &pb.ListUsersResponse{
		PageSize:   uint32(pageSize),
		Users:      usersIDOtoDTOs(users),
		TotalCount: uint64(totalCount),
	}, nil
}
