package train

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/cyrilix/robocar-tools/pkg/awsutils"
	"github.com/cyrilix/robocar-tools/pkg/data"
	"github.com/cyrilix/robocar-tools/pkg/models"
	"log"
	"strconv"
	"time"
)

func New(bucketName string, ociImage, roleArn string) *Training {
	return &Training{
		config:       awsutils.MustLoadConfig(),
		bucketName:   bucketName,
		ociImage:     ociImage,
		roleArn:      roleArn,
		prefixInput:  prefixInput,
		outputBucket: fmt.Sprintf("s3://%s/output", bucketName),
	}
}

type Training struct {
	config       aws.Config
	bucketName   string
	ociImage     string
	roleArn      string
	prefixInput  string
	outputBucket string
}

func (t *Training) TrainDir(ctx context.Context, jobName, basedir string, imgHeight, imgWidth int, sliceSize int, withFlipImage bool, outputModelFile string, enableSpotTraining bool) error {
	log.Printf("run training with data from %s\n", basedir)
	archive, err := data.BuildArchive(basedir, sliceSize, withFlipImage)
	if err != nil {
		return fmt.Errorf("unable to build data archive: %w", err)
	}
	log.Println("")

	err = t.UploadArchive(ctx, archive)
	if err != nil {
		return fmt.Errorf("unable to upload data arrchive: %w", err)
	}
	log.Println("")

	err = t.runTraining(
		ctx,
		jobName,
		sliceSize,
		imgHeight,
		imgWidth,
		enableSpotTraining,
	)
	if err != nil {
		return fmt.Errorf("unable to run training: %w", err)
	}

	err = t.GetTrainingOutput(ctx, jobName, outputModelFile)
	if err != nil {
		return fmt.Errorf("unable to get output model file '%s': %w", outputModelFile, err)
	}
	return nil
}

func List(bucketName string) error {
	l := zap.S()
	pfxInput := prefixInput

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(awsutils.MustLoadConfig())

	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: &pfxInput,
	})
	if err != nil {
		l.Fatal(err)
	}

	l.Info("first page results:")
	for _, object := range output.Contents {
		if *object.Key == pfxInput {
			continue
		}
		l.Infof("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}
	return nil
}

func (t *Training) runTraining(ctx context.Context, jobName string, slideSize int, imgHeight, imgWidth int, enableSpotTraining bool) error {
	l := zap.S()
	client := sagemaker.NewFromConfig(awsutils.MustLoadConfig())
	l.Infof("Start training job '%s'", jobName)

	trainingJobInput := sagemaker.CreateTrainingJobInput{
		EnableManagedSpotTraining: enableSpotTraining,
		AlgorithmSpecification: &types.AlgorithmSpecification{
			TrainingInputMode: types.TrainingInputModeFile,
			TrainingImage:     aws.String(t.ociImage),
		},
		OutputDataConfig: &types.OutputDataConfig{
			S3OutputPath: aws.String(t.outputBucket),
		},
		ResourceConfig: &types.ResourceConfig{
			InstanceCount:  1,
			//InstanceType:   types.TrainingInstanceTypeMlP2Xlarge,
			InstanceType:   types.TrainingInstanceTypeMlG4dnXlarge,
			VolumeSizeInGB: 1,
		},
		RoleArn: aws.String(t.roleArn),
		StoppingCondition: &types.StoppingCondition{
			MaxRuntimeInSeconds: 1800,
		},
		TrainingJobName: aws.String(jobName),
		HyperParameters: map[string]string{
			"sagemaker_region": "eu-west-1",
			"slide_size":       strconv.Itoa(slideSize),
			"img_height":       strconv.Itoa(imgHeight),
			"img_width":        strconv.Itoa(imgWidth),
		},
		InputDataConfig: []types.Channel{
			{
				ChannelName: aws.String("train"),
				DataSource: &types.DataSource{
					S3DataSource: &types.S3DataSource{
						S3DataType:             types.S3DataTypeS3Prefix,
						S3Uri:                  aws.String(fmt.Sprintf("s3://%s/%s", t.bucketName, t.prefixInput)),
						S3DataDistributionType: types.S3DataDistributionFullyReplicated,
					},
				},
			},
		},
	}
	if enableSpotTraining {
		trainingJobInput.StoppingCondition.MaxWaitTimeInSeconds = aws.Int32(3600)
	}

	// TODO: check train data exist
	jobOutput, err := client.CreateTrainingJob(
		ctx,
		&trainingJobInput,
	)

	if err != nil {
		return fmt.Errorf("unable to run sagemeker job: %w", err)
	}

	for {
		time.Sleep(30 * time.Second)

		status, err := client.DescribeTrainingJob(
			ctx,
			&sagemaker.DescribeTrainingJobInput{
				TrainingJobName: aws.String(jobName),
			},
		)
		if err != nil {
			log.Printf("unable to get status from ob %v: %v\n", jobOutput.TrainingJobArn, err)
			continue
		}
		switch status.TrainingJobStatus {
		case types.TrainingJobStatusInProgress:
			log.Printf("job in progress: %v - %v - %v\n", status.TrainingJobStatus, status.SecondaryStatus, *status.SecondaryStatusTransitions[len(status.SecondaryStatusTransitions) - 1].StatusMessage)
			continue
		case types.TrainingJobStatusFailed:
			return fmt.Errorf("job %s finished with status %v\n", jobName, status.TrainingJobStatus)
		default:
			log.Printf("job %s finished with status %v\n", jobName, status.TrainingJobStatus)
			return nil
		}
	}
}

func (t *Training) GetTrainingOutput(ctx context.Context, jobName, outputFile string) error {
	modelPath := fmt.Sprintf("output/%s/output/model.tar.gz", jobName)
	err := models.DownloadArchiveToFile(ctx, t.bucketName, modelPath, outputFile)
	if err != nil {
		return fmt.Errorf("unable to download training model '%s' to '%s' file: %w", modelPath, outputFile, err)
	}
	return nil
}

func ListJob(ctx context.Context) error {

	client := sagemaker.NewFromConfig(awsutils.MustLoadConfig())
	jobs, err := client.ListTrainingJobs(ctx, &sagemaker.ListTrainingJobsInput{})
	if err != nil {
		return fmt.Errorf("unable to list trainings jobs: %w", err)
	}
	for _, job := range jobs.TrainingJobSummaries {
		fmt.Printf("%s\t\t%s\n", *job.TrainingJobName, job.TrainingJobStatus)
	}
	return  nil
}