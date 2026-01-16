package users

import (
	"context"
	"net/http"

	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (r *Repo) GeneratePresignedPutURL(ctx context.Context, userID id.UUID7, contentType string) (string, http.Header, error) {
	objectKey := "users/" + userID.String() + "/media/profile_image.jpeg"

	// tags := "Key1=Value1&Key2=Value2" // Example tags; modify as needed
	request, err := r.awsS3.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.awsS3.BucketNameMediaRaw),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
		// Tagging:     aws.String(tags),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = r.awsS3.presignedPutURLTTL
	})
	if err != nil {
		r.logger.ErrorContext(ctx, "Couldn't get a presigned request to put %s to %s: %v", objectKey, r.awsS3.BucketNameMedia, err)
		return "", nil, err
	}

	r.logger.DebugContext(ctx, "Generated presigned URL for user profile picture upload", "url", request.URL)

	return request.URL, request.SignedHeader, err
}

func (r *Repo) GeneratePresignedGetURL(ctx context.Context, userID id.UUID7) (string, http.Header, error) {
	objectKey := "users/" + userID.String() + "/media/profile_image.jpeg"

	// tags := "Key1=Value1&Key2=Value2" // Example tags; modify as needed
	request, err := r.awsS3.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.awsS3.BucketNameMedia),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = r.awsS3.presignedGetURLTTL
	})
	if err != nil {
		r.logger.ErrorContext(ctx, "Couldn't get a presigned request to get %s from %s: %v", objectKey, r.awsS3.BucketNameMedia, err)
		return "", nil, err
	}

	r.logger.DebugContext(ctx, "Generated presigned URL for user profile picture download", "url", request.URL)

	return request.URL, request.SignedHeader, err
}
