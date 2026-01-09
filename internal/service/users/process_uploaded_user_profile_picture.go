package users

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type Objects3Getter interface {
	GetS3Object(ctx context.Context, bucketName, objectKey string) ([]byte, error)
}

type ProfilePictureAndThumbnailerS3Putter interface {
	PutProfilePictureAndThumbnailToS3(ctx context.Context, userID uuid.UUID, profilePicBytes, thumbnailBytes []byte) error
}

func (s *Service) ProcessUploadedUserProfilePicture(ctx context.Context, userID uuid.UUID, bucketName, objectKey string) error {
	// fetch the object from S3
	objBytes, err := s.repo.GetS3Object(ctx, bucketName, objectKey)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to fetch S3 object for profile picture processing", "error", err, "bucket", bucketName, "object_key", objectKey)
		return err
	}

	profilePic, thumbnail, err := s.imageProcessor.ProcessUserProfilePicture(objBytes)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to process profile picture image", "error", err, "user_id", userID)
		return err
	}

	// store the processed images back to S3 or another storage as needed
	err = s.repo.PutProfilePictureAndThumbnailToS3(ctx, userID, profilePic, thumbnail)
	if err = s.repo.PutProfilePictureAndThumbnailToS3(ctx, userID, profilePic, thumbnail); err != nil {
		s.logger.ErrorContext(ctx, "failed to store processed profile picture and thumbnail", "error", err, "user_id", userID)
		return err
	}

	s.logger.InfoContext(ctx, "successfully processed and stored user profile picture", "user_id", userID)
	return nil
}
