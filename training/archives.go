package training

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var	sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.EuWest1RegionID),
	}))

func ListArchives() error {
	srv := s3.New(sess)
	input := s3.ListBucketsInput{}
	buckets, err := srv.ListBuckets(&input)
	if err != nil {
		return fmt.Errorf("unable to list buckets: %v", err)
	}

	for _, b := range buckets.Buckets {
		fmt.Printf("bucket: %v\n", b)
	}
	bucketName := aws.String("robocars-cyrilix-learning")
	prefix := aws.String("input/data/train/train.zip")
	listObjectInput := s3.ListObjectsV2Input{
		Bucket: bucketName,
		Prefix: prefix,
	}
	objects, err := srv.ListObjectsV2(&listObjectInput)
	if err != nil {
		return fmt.Errorf("unable to list objects in bucket %v: %v", *bucketName, err)
	}
	fmt.Printf("objects: %v\n", objects)

	return nil
}

func UploadArchive(archive *[]byte){
	bucketName := aws.String("robocars-cyrilix-learning")
	key := aws.String("input/data/train/train.zip")

	ctx := context.Background()

	svc := s3.New(sess)
	output, err := svc.PutObjectWithContext(ctx, s3.PutObjectInput{
		ACL:                       aws.String(s3.ObjectCannedACLPrivate),
		Body:                      bytes.NewReader(*archive),
		Bucket:                    bucketName,
		Key:                       key,
	}
	svc.PutObject(input)
}
