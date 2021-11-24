package train

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/cyrilix/robocar-tools/pkg/awsutils"
	"go.uber.org/zap"
)

func ListArchives(ctx context.Context, bucket string) error {
	client := s3.NewFromConfig(awsutils.MustLoadConfig())

	prefix := aws.String("input/data/train/train.zip")
	objects, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: prefix,
	})
	if err != nil {
		return fmt.Errorf("unable to list objects in bucket %v: %w", bucket, err)
	}
	fmt.Printf("objects: %v\n", objects)

	return nil
}

func (t Training) UploadArchive(ctx context.Context, archive []byte) error {
	client := s3.NewFromConfig(t.config)
	key := aws.String("input/data/train/train.zip")

	zap.S().Infof("upload archive to bucket '%s/%s'", t.bucketName, *key)
	_, err := client.PutObject(
		ctx,
		&s3.PutObjectInput{
			ACL:    types.ObjectCannedACLPrivate,
			Body:   bytes.NewReader(archive),
			Bucket: aws.String(t.bucketName),
			Key:    key,
		})
	if err != nil {
		return fmt.Errorf("unable to upload archive: %w", err)
	}
	zap.S().Info("archive uploaded")
	return nil
}
