package users

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/OJOMB/fightpicker/internal/grpc/gen/go/auth/v1"
	authservice "github.com/OJOMB/fightpicker/internal/service/auth"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

// Server implements pb.AuthServiceServer
type Server struct {
	pb.AuthServiceServer

	service *authservice.Service
	logger  logs.Logger
}

func extractRefreshToken(cookies []string) string {
	for _, cookie := range cookies {
		for part := range strings.SplitSeq(cookie, ";") {
			if after, found := strings.CutPrefix(strings.TrimSpace(part), "refresh_token="); found {
				return after
			}
		}
	}

	return ""
}

// NewServer is a constructor for the Auth gRPC server.
func NewServer(service *authservice.Service, logger logs.Logger) *Server {
	return &Server{
		service: service,
		logger:  logger,
	}
}

func (s *Server) setRefreshTokenInMetadata(ctx context.Context, refreshToken string, refreshTokenTTL time.Time) error {
	maxAge := time.Until(refreshTokenTTL).Seconds()

	return grpc.SetHeader(
		ctx,
		metadata.Pairs(
			"set-cookie",
			fmt.Sprintf("refresh_token=%s; HttpOnly; Secure; SameSite=Strict; Path=/; Max-Age=%d", refreshToken, maxAge),
		),
	)
}

// Login authenticates a user and returns an access token and a refresh token in the metadata.
func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	email := req.GetEmail()
	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	password := req.GetPassword()
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// ... your auth logic ...
	accessToken, refreshToken, refreshTokenTTL, err := s.service.Login(ctx, email, password)
	if err != nil {
		return nil, s.toStatus(ctx, err)
	}

	// Set refresh token in metadata
	if err := s.setRefreshTokenInMetadata(ctx, refreshToken, refreshTokenTTL); err != nil {
		return nil, s.toStatus(ctx, err)
	}

	// Return access token in response body
	return &pb.LoginResponse{AccessToken: accessToken}, nil
}

func (s *Server) Refresh(ctx context.Context, _ *emptypb.Empty) (*pb.LoginResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata in context")
	}

	refreshTokens := md.Get("cookie")
	var refreshToken string
	if len(refreshTokens) > 0 {
		refreshToken = refreshTokens[0]
	}

	accessToken, newRefreshToken, newRefreshTokenTTL, err := s.service.Refresh(ctx, refreshToken)
	if err != nil {
		return nil, s.toStatus(ctx, err)
	}

	// Set new refresh token in metadata
	if err := s.setRefreshTokenInMetadata(ctx, newRefreshToken, newRefreshTokenTTL); err != nil {
		return nil, s.toStatus(ctx, err)
	}

	// Return new access token in response body
	return &pb.LoginResponse{AccessToken: accessToken}, nil
}

func (s *Server) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata in context")
	}

	cookies := md.Get("cookie")
	refreshToken := extractRefreshToken(cookies)
	if refreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "missing refresh token in cookies")
	}

	err := s.service.Logout(ctx, refreshToken)
	if err != nil {
		return nil, s.toStatus(ctx, err)
	}

	// Clear refresh token cookie
	if err := s.setRefreshTokenInMetadata(ctx, "", time.Time{}); err != nil {
		return nil, s.toStatus(ctx, err)
	}

	return new(emptypb.Empty), nil
}
