package models

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cyrilix/robocar-tools/pkg/awsutils"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

func ListModels(ctx context.Context, bucket string) error {

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(awsutils.MustLoadConfig())

	// Get the first page of results for ListObjectsV2 for a bucket
	outputs, err := client.ListObjectsV2(
		ctx,
		&s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String("output"),
		},
	)
	if err != nil {
		return fmt.Errorf("unable to list models in bucket %v: %w", bucket, err)
	}
	for _, output := range outputs.Contents {
		fmt.Printf("model: %s\n", *output.Key)
	}

	return nil
}

func DownloadArchive(ctx context.Context, bucketName, modelPath string) ([]byte, error) {
	l := zap.S().With(
		"bucket", bucketName,
		"model", modelPath,
	)
	client := s3.NewFromConfig(awsutils.MustLoadConfig())

	l.Debug("download model")
	archive, err := client.GetObject(
		ctx,

		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(modelPath),
		})
	if err != nil {
		return nil, fmt.Errorf("unable to download model: %w", err)
	}
	l.Debug("model downloaded")
	resp, err := ioutil.ReadAll(archive.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read model archive content: %w", err)
	}
	return resp, nil
}

func DownloadArchiveToFile(ctx context.Context, bucketName, modelPath, outputFile string) error {
	arch, err := DownloadArchive(ctx, bucketName, modelPath)
	if err != nil {
		return fmt.Errorf("unable to download model '%v/%v': %w", bucketName, modelPath, err)
	}
	err = ioutil.WriteFile(outputFile, arch, os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("unable to write model '%s' to file '%s': %v", modelPath, outputFile, err)
	}
	return nil
}
