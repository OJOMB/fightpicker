package users

import (
	"context"
	"fmt"
	"regexp"

	"github.com/gofrs/uuid/v5"

	usermetrics "github.com/OJOMB/fightpicker/internal/metrics/users"
	"github.com/OJOMB/fightpicker/pkg/datetimes"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

// emailRegexPattern is the regex pattern used to validate email addresses.
// should match the regex in the database schema.
const emailRegexPattern = `^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`

type AuthTool interface {
	GenerateVerificationToken() (string, error)
	HashVerificationToken(token string) ([]byte, error)
	HashPassword(password string) (string, error)
}

// Repo defines the interface for the user repository.
type Repo interface {
	UserCreator
	UserByIDGetter
	UserByIDDeleter
	UserByEmailGetter
	UserByUsernameGetter
	UserUpdater
	UserLister
	UserFollowerLister
	UserFolloweeLister
	UserFollower
	UserUnfollower
	PresignedGetURLGenerator
	PresignedPutURLGenerator
	Objects3Getter
	ProfilePictureAndThumbnailerS3Putter
	VerificationEmailSender
	EmailByTokenVerifier
}

// IDGenerator defines the interface for the ID generator used in the service.
type IDGenerator interface {
	Generate() uuid.UUID
}

type ImageProcessor interface {
	ProcessUserProfilePicture(imageBytes []byte) ([]byte, []byte, error)
}

// Service responsible for user-related business logic.
type Service struct {
	repo                Repo
	domain              string
	authTool            AuthTool
	regexEmail          *regexp.Regexp
	idGen               IDGenerator
	dateTimeTool        datetimes.Now
	imageProcessor      ImageProcessor
	emailAddressNoReply string

	metrics *usermetrics.Metrics
	logger  logs.Logger
}

// NewService creates a new instance of the user service.
func NewService(repo Repo, idGen IDGenerator, now datetimes.Now, authTool AuthTool, imageProcessor ImageProcessor, domain, emailAddressNoReply string, logger logs.Logger) (*Service, error) {
	if repo == nil {
		return nil, ErrNilRepo
	}

	if logger == nil {
		return nil, ErrNilLogger
	}

	if idGen == nil {
		return nil, ErrNilIDGenerator
	}

	if now == nil {
		return nil, ErrNilNowFunc
	}

	if authTool == nil {
		return nil, ErrNilPasswordHasher
	}

	if imageProcessor == nil {
		return nil, ErrNilImageProcessor
	}

	emailRegex, err := regexp.Compile(emailRegexPattern)
	if err != nil {
		return nil, ErrEmailRegexCompile
	}

	metrics, err := usermetrics.New()
	if err != nil {
		return nil, err
	}

	return &Service{
		repo:                repo,
		regexEmail:          emailRegex,
		idGen:               idGen,
		dateTimeTool:        now,
		imageProcessor:      imageProcessor,
		authTool:            authTool,
		domain:              domain,
		emailAddressNoReply: emailAddressNoReply,

		metrics: metrics,
		logger:  logger,
	}, nil
}

func (s *Service) injectPresignedGetURLIfNeeded(ctx context.Context, user *User) {
	if user.ProfilePicture != "" {
		presignedURL, _, err := s.repo.GeneratePresignedGetURL(ctx, user.ID)
		if err != nil {
			s.logger.ErrorContext(ctx, "failed to generate presigned URL for user profile picture", "error", err, "user_id", user.ID)
			// TODO: inject static URL for default profile picture
		}

		user.ProfilePicture = presignedURL

		return
	}

	user.ProfilePicture = s.getDefaultProfilePictureURL(user)
}

func (s *Service) getDefaultProfilePictureURL(user *User) string {
	return fmt.Sprintf("%s/static/images/default-profile-picture-%s.webp", s.domain, user.Gender.String())
}
