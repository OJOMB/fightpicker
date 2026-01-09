package users

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
)

// PutProfilePictureAndThumbnailToS3 updates the users profile picture in the DB and uploads the processed profile picture and thumbnail to S3.
func (r *Repo) PutProfilePictureAndThumbnailToS3(ctx context.Context, userID uuid.UUID, profilePicBytes, thumbnailBytes []byte) error {
	objKeyRoot := "users/" + userID.String()

	// first attempt to update the user's profile picture in the database - if this fails, we shouldn't upload the images to S3
	if err := r.dbClient.UpdateUserProfilePictureByID(ctx, postgres.UpdateUserProfilePictureByIDParams{
		UpdatedBy: pgtype.UUID{
			Bytes: userID,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  r.dateTimeTool.Now(),
			Valid: true,
		},
		ProfilePicture: pgtype.Text{
			String: r.awsS3.BucketNameMedia + "/" + objKeyRoot,
			Valid:  true,
		},
	}); err != nil {
		return dbErrorToServiceError(err)
	}

	// generate S3 object keys
	profilePicKey := objKeyRoot + "/profile_picture.webp"
	thumbnailKey := objKeyRoot + "/thumbnail.webp"

	// upload profile picture
	_, err := r.awsS3.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &r.awsS3.BucketNameMedia,
		Key:    &profilePicKey,
		Body:   bytes.NewReader(profilePicBytes),
	})
	if err != nil {
		return err
	}

	// upload thumbnail
	_, err = r.awsS3.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &r.awsS3.BucketNameMedia,
		Key:    &thumbnailKey,
		Body:   bytes.NewReader(thumbnailBytes),
	})
	if err != nil {
		return err
	}

	return nil
}
