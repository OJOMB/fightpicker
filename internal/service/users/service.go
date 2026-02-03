package users

import (
	"context"
	"fmt"
	"regexp"

	usermetrics "github.com/OJOMB/fightpicker/internal/metrics/users"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/datetimes"
	"github.com/OJOMB/fightpicker/pkg/id"
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

type ImageProcessor interface {
	ProcessUserProfilePicture(imageBytes []byte) ([]byte, []byte, error)
}

// Service responsible for user-related business logic.
type Service struct {
	repo                Repo
	domain              string
	authTool            AuthTool
	regexEmail          *regexp.Regexp
	idTool              id.UUID7GeneratorParser
	ctxTool             contextual.ContextProvider
	dateTimeTool        datetimes.Now
	imageProcessor      ImageProcessor
	emailAddressNoReply string

	metrics *usermetrics.Metrics
	logger  logs.Logger
}

// New creates a new instance of the user service.
func New(
	repo Repo,
	idTool id.UUID7GeneratorParser,
	now datetimes.Now,
	ctxTool contextual.ContextProvider,
	authTool AuthTool,
	imageProcessor ImageProcessor,
	domain,
	emailAddressNoReply string,
	logger logs.Logger,
) (*Service, error) {
	if repo == nil {
		return nil, errNilRepo
	}

	if ctxTool == nil {
		return nil, errNilContextTool
	}

	if logger == nil {
		return nil, errNilLogger
	}

	if idTool == nil {
		return nil, errNilIDTool
	}

	if now == nil {
		return nil, errNilNowFunc
	}

	if authTool == nil {
		return nil, errNilPasswordHasher
	}

	if imageProcessor == nil {
		return nil, errNilImageProcessor
	}

	emailRegex, err := regexp.Compile(emailRegexPattern)
	if err != nil {
		return nil, errEmailRegexCompile
	}

	metrics, err := usermetrics.New()
	if err != nil {
		return nil, err
	}

	return &Service{
		repo:                repo,
		regexEmail:          emailRegex,
		idTool:              idTool,
		ctxTool:             ctxTool,
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
	return fmt.Sprintf("%s/static/users/default-profile-picture-%s.webp", s.domain, user.Gender.String())
}
