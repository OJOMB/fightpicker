package users

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// GetS3Object retrieves an object in the given bucket from S3 and returns the fetched bytes.
func (r *Repo) GetS3Object(ctx context.Context, bucketName, objectKey string) ([]byte, error) {
	logger := r.logger.With("bucket", bucketName, "object_key", objectKey)

	objOutput, err := r.awsS3.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get S3 object", "error", err)
		return nil, err
	}

	data, err := io.ReadAll(objOutput.Body)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to read S3 object body", "error", err)
		return nil, err
	}
	defer objOutput.Body.Close()

	if len(data) == 0 {
		logger.ErrorContext(ctx, "S3 object has no content")
		return nil, nil
	}

	return data, nil
}
