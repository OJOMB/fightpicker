package users

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twmb/franz-go/pkg/kgo"
	"gopkg.in/mail.v2"

	usermetrics "github.com/OJOMB/fightpicker/internal/metrics/users"
	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/datetimes"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type KafkaClient interface {
	ProduceMessage(ctx context.Context, key, value []byte) error
}

type DBClient interface {
	AssignRoleToUserByRoleName(ctx context.Context, arg postgres.AssignRoleToUserByRoleNameParams) error
	CountFollowers(ctx context.Context, followeeID uuid.UUID) (int64, error)
	CountFollowing(ctx context.Context, followerID uuid.UUID) (int64, error)
	CreateFighter(ctx context.Context, arg postgres.CreateFighterParams) error
	CreateUser(ctx context.Context, arg postgres.CreateUserParams) error
	DeleteUserByID(ctx context.Context, id uuid.UUID) error
	FollowUser(ctx context.Context, arg postgres.FollowUserParams) error
	GetFighterByID(ctx context.Context, id uuid.UUID) (postgres.Fighter, error)
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (postgres.RefreshToken, error)
	GetUserByEmail(ctx context.Context, email string) (postgres.GetUserByEmailRow, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (postgres.GetUserByIDRow, error)
	GetUserByUsername(ctx context.Context, username string) (postgres.GetUserByUsernameRow, error)
	GetUserPermissionsAndRolesByID(ctx context.Context, id uuid.UUID) ([]postgres.GetUserPermissionsAndRolesByIDRow, error)
	GetUserPermissionsByID(ctx context.Context, id uuid.UUID) ([]postgres.GetUserPermissionsByIDRow, error)
	GetUserRBAC(ctx context.Context, id uuid.UUID) (postgres.GetUserRBACRow, error)
	GetUserRolesByID(ctx context.Context, userID uuid.UUID) ([]string, error)
	IsFollowing(ctx context.Context, arg postgres.IsFollowingParams) (bool, error)
	IsUserEmailVerifiedByID(ctx context.Context, id uuid.UUID) (bool, error)
	ListFollowers(ctx context.Context, arg postgres.ListFollowersParams) ([]postgres.ListFollowersRow, error)
	ListFollowees(ctx context.Context, arg postgres.ListFolloweesParams) ([]postgres.ListFolloweesRow, error)
	ListUsers(ctx context.Context, arg postgres.ListUsersParams) ([]postgres.ListUsersRow, error)
	RevokeAndRotateRefreshTokenByHash(ctx context.Context, arg postgres.RevokeAndRotateRefreshTokenByHashParams) error
	RevokeRefreshTokenByHash(ctx context.Context, arg postgres.RevokeRefreshTokenByHashParams) error
	StoreRefreshToken(ctx context.Context, arg postgres.StoreRefreshTokenParams) error
	UnfollowUser(ctx context.Context, arg postgres.UnfollowUserParams) error
	UpdateEmailVerificationTokenHashByUserID(ctx context.Context, arg postgres.UpdateEmailVerificationTokenHashByUserIDParams) error
	UpdateUserByID(ctx context.Context, arg postgres.UpdateUserByIDParams) error
	UpdateUserProfilePictureByID(ctx context.Context, arg postgres.UpdateUserProfilePictureByIDParams) error
	VerifyUserEmailByTokenHash(ctx context.Context, arg postgres.VerifyUserEmailByTokenHashParams) (uuid.UUID, error)
	WithTx(tx pgx.Tx) *postgres.Queries
}

type Repo struct {
	pool         *pgxpool.Pool
	host         string
	dbClient     DBClient
	dateTimeTool datetimes.Now
	awsS3        repoS3
	email        repoEmail
	events       repoEvents

	metrics *usermetrics.Metrics
	logger  logs.Logger
}

type repoS3 struct {
	client             *s3.Client
	presigner          *s3.PresignClient
	BucketNameMedia    string
	BucketNameMediaRaw string
	presignedPutURLTTL time.Duration
	presignedGetURLTTL time.Duration
}

type repoEvents struct {
	client                *kgo.Client
	topicPostUserCreate   string
	topicPostUserDelete   string
	topicPostUserVerified string
}

type repoEmail struct {
	dialer      *mail.Dialer
	AddrNoReply string
}

func New(
	pool *pgxpool.Pool, host string, client DBClient, s3Client *s3.Client, dateTimeTool datetimes.Now,
	bucketNameMedia, bucketNameMediaRaw string,
	presignedGetURLTTL, presignedPutURLTTL time.Duration,
	kafkaClient *kgo.Client, topicPostUserCreate, topicPostUserVerified string,
	emailDialer *mail.Dialer, emailAddressNoReply string,
	logger logs.Logger,
) (*Repo, error) {
	presigner := s3.NewPresignClient(s3Client)

	metrics, err := usermetrics.New()
	if err != nil {
		return nil, err
	}

	return &Repo{
		pool:         pool,
		host:         host,
		dbClient:     client,
		dateTimeTool: dateTimeTool,
		awsS3: repoS3{
			client:             s3Client,
			presigner:          presigner,
			BucketNameMedia:    bucketNameMedia,
			BucketNameMediaRaw: bucketNameMediaRaw,
			presignedGetURLTTL: presignedGetURLTTL,
			presignedPutURLTTL: presignedPutURLTTL,
		},
		events: repoEvents{
			client:                kafkaClient,
			topicPostUserCreate:   topicPostUserCreate,
			topicPostUserVerified: topicPostUserVerified,
		},
		email: repoEmail{
			dialer:      emailDialer,
			AddrNoReply: emailAddressNoReply,
		},
		metrics: metrics,
		logger:  logger.With("component", "users_repo"),
	}, nil
}
