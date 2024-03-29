package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cyrilix/robocar-base/cli"
	"github.com/cyrilix/robocar-tools/dkimpt"
	"github.com/cyrilix/robocar-tools/part"
	"github.com/cyrilix/robocar-tools/pkg/data"
	"github.com/cyrilix/robocar-tools/pkg/display"
	"github.com/cyrilix/robocar-tools/pkg/models"
	"github.com/cyrilix/robocar-tools/pkg/train"
	"github.com/cyrilix/robocar-tools/record"
	"github.com/cyrilix/robocar-tools/video"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"log"
	"os"
)

const (
	DefaultClientId       = "robocar-tools"
	DefaultTrainSliceSize = 0
)

func main() {
	var mqttBroker, username, password, clientId string
	var framePath string
	var fps int
	var frameTopic, objectsTopic, roadTopic, recordTopic, throttleFeedbackTopic string
	var withObjects, withRoad, withThrottleFeedback bool
	var recordsPath string
	var trainArchiveName string
	var trainSliceSize int
	var bucket, ociImage string
	var debug bool

	mqttQos := cli.InitIntFlag("MQTT_QOS", 0)
	_, mqttRetain := os.LookupEnv("MQTT_RETAIN")

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("  display\n  \tDisplay events on live frames\n")
		fmt.Printf("  record \n  \tRecord event for tensorflow training\n")
		fmt.Printf("  training  \n  \tManage training\n")
		fmt.Printf("  models  \n  \tManage models\n")
		fmt.Printf("  import-donkey-records \n  \tCopy donkeycar records to new format\n")
	}

	err := cli.SetIntDefaultValueFromEnv(&trainSliceSize, "RC_TRAIN_SLICE_SIZE", DefaultTrainSliceSize)
	if err != nil {
		log.Printf("unable to init TRAIN_SLICE_SIZE: %v", err)
	}
	cli.SetDefaultValueFromEnv(&ociImage, "TRAIN_OCI_IMAGE", "")
	cli.SetDefaultValueFromEnv(&bucket, "TRAIN_BUCKET", "")

	flag.BoolVar(&debug, "debug", false, "Display debug logs")

	displayFlags := flag.NewFlagSet("display", flag.ExitOnError)

	displayFlags.Usage = func() {
		fmt.Printf("Usage of %s %s:\n", os.Args[0], displayFlags.Name())
		fmt.Printf("  camera\n  \tLive from car camera\n")
		fmt.Printf("  record\n  \tLive from published records\n")
	}

	displayRecordFlags := flag.NewFlagSet("record", flag.ExitOnError)
	cli.InitMqttFlagSet(displayRecordFlags, DefaultClientId, &mqttBroker, &username, &password, &clientId, &mqttQos, &mqttRetain)
	displayRecordFlags.StringVar(&recordTopic, "mqtt-topic-records", os.Getenv("MQTT_TOPIC_RECORDS"), "Mqtt topic that contains record data for training, use MQTT_TOPIC_RECORDS if args not set")

	displayCameraFlags := flag.NewFlagSet("camera", flag.ExitOnError)
	cli.InitMqttFlagSet(displayCameraFlags, DefaultClientId, &mqttBroker, &username, &password, &clientId, &mqttQos, &mqttRetain)
	displayCameraFlags.StringVar(&frameTopic, "mqtt-topic-frame", os.Getenv("MQTT_TOPIC_FRAME"), "Mqtt topic that contains frame to display, use MQTT_TOPIC_FRAME if args not set")
	displayCameraFlags.StringVar(&framePath, "frame-path", "", "Directory path where to read jpeg frame to inject in frame topic")
	displayCameraFlags.IntVar(&fps, "frame-per-second", 25, "Video frame per second of frame to publish")

	displayCameraFlags.StringVar(&objectsTopic, "mqtt-topic-objects", os.Getenv("MQTT_TOPIC_OBJECTS"), "Mqtt topic that contains detected objects, use MQTT_TOPIC_OBJECTS if args not set")
	displayCameraFlags.BoolVar(&withObjects, "with-objects", false, "Display detected objects")

	displayCameraFlags.StringVar(&roadTopic, "mqtt-topic-road", os.Getenv("MQTT_TOPIC_ROAD"), "Mqtt topic that contains road description, use MQTT_TOPIC_ROAD if args not set")
	displayCameraFlags.BoolVar(&withRoad, "with-road", false, "Display detected road")

	displayCameraFlags.StringVar(&throttleFeedbackTopic, "mqtt-topic-throttle-feedback", os.Getenv("MQTT_TOPIC_THROTTLE_FEEDBACK"), "Mqtt topic where to publish throttle feedback, use MQTT_TOPIC_THROTTLE_FEEDBACK if args not set")
	displayCameraFlags.BoolVar(&withThrottleFeedback, "with-throttle-feedback", false, "Display throttle feedback")

	recordFlags := flag.NewFlagSet("record", flag.ExitOnError)
	cli.InitMqttFlagSet(recordFlags, DefaultClientId, &mqttBroker, &username, &password, &clientId, &mqttQos, &mqttRetain)
	recordFlags.StringVar(&recordTopic, "mqtt-topic-records", os.Getenv("MQTT_TOPIC_RECORDS"), "Mqtt topic that contains record data for training, use MQTT_TOPIC_RECORDS if args not set")
	recordFlags.StringVar(&recordsPath, "record-path", os.Getenv("RECORD_PATH"), "Path where to write records files, use RECORD_PATH if args not set")

	var basedir, destdir string
	impdkFlags := flag.NewFlagSet("import-donkey-records", flag.ExitOnError)
	impdkFlags.StringVar(&basedir, "from", "", "source directory")
	impdkFlags.StringVar(&destdir, "to", "", "destination directory")

	trainingFlags := flag.NewFlagSet("training", flag.ExitOnError)
	trainingFlags.Usage = func() {
		fmt.Printf("Usage of %s %s:\n", os.Args[0], trainingFlags.Name())
		fmt.Printf("  list\n  \tList existing training jobs\n")
		fmt.Printf("  archive\n  \tBuild tar.gz archive for training\n")
		fmt.Printf("  run\n  \tRun training job\n")
	}

	var modelPath, roleArn, trainJobName, modelType string
	var horizon int
	var withFlipImage bool
	var trainImageHeight, trainImageWidth int
	var enableSpotTraining bool
	trainingRunFlags := flag.NewFlagSet("run", flag.ExitOnError)
	trainingRunFlags.StringVar(&bucket, "bucket", os.Getenv("RC_TRAIN_BUCKET"), "AWS bucket where store data required, use RC_TRAIN_BUCKET if arg not set")
	trainingRunFlags.StringVar(&recordsPath, "record-path", os.Getenv("RECORD_PATH"), "Input data path where records and img files are stored, use RECORD_PATH if arg not set")
	trainingRunFlags.StringVar(&modelPath, "output-model-path", "", "Path where to write output model archive")
	trainingRunFlags.IntVar(&trainSliceSize, "slice-size", trainSliceSize, "Number of record to shift with image, use RC_TRAIN_SLICE_SIZE if args not set")
	trainingRunFlags.StringVar(&ociImage, "oci-image", os.Getenv("RC_TRAIN_OCI_IMAGE"), "OCI image to run (required), use RC_TRAIN_OCI_IMAGE if args not set")
	trainingRunFlags.StringVar(&roleArn, "role-arn", os.Getenv("RC_TRAIN_ROLE"), "AWS ARN role to use to run training (required), use RC_TRAIN_ROLE if arg not set")
	trainingRunFlags.StringVar(&trainJobName, "job-name", "", "Training job name (required)")
	trainingRunFlags.BoolVar(&withFlipImage, "with-flip-image", withFlipImage, "Flip horiontal image and reverse steering to increase data into training archive")

	trainingRunFlags.IntVar(&trainImageHeight, "image-height", 128, "Pixels image height")
	trainingRunFlags.IntVar(&trainImageWidth, "image-width", 160, "Pixels image width")
	trainingRunFlags.IntVar(&horizon, "horizon", 0, "Upper zone image to crop (in pixels)")
	trainingRunFlags.StringVar(&modelType, "model-type", train.ModelTypeCategorical.String(), "Type model to build")

	trainingRunFlags.BoolVar(&enableSpotTraining, "enable-spot-training", true, "Train models using managed spot training")
	trainingListJobFlags := flag.NewFlagSet("list", flag.ExitOnError)

	trainArchiveFlags := flag.NewFlagSet("archive", flag.ExitOnError)
	trainArchiveFlags.StringVar(&recordsPath, "record-path", os.Getenv("RECORD_PATH"), "Path where records files are stored, use RECORD_PATH if args not set")
	trainArchiveFlags.StringVar(&trainArchiveName, "output", os.Getenv("TRAIN_ARCHIVE_NAME"), "Zip archive file name, use TRAIN_ARCHIVE_NAME if args not set")
	trainArchiveFlags.IntVar(&trainSliceSize, "slice-size", trainSliceSize, "Number of record to shift with image, use TRAIN_SLICE_SIZE if args not set")
	trainArchiveFlags.IntVar(&trainImageWidth, "image-width", 0, "Resize image width")
	trainArchiveFlags.IntVar(&trainImageHeight, "image-height", 0, "Resize image height")
	trainArchiveFlags.IntVar(&horizon, "horizon", 0, "Upper zone image to crop (in pixels)")
	trainArchiveFlags.BoolVar(&withFlipImage, "with-flip-image", withFlipImage, "Flip horiontal image and reverse steering to increase data into training archive")

	modelsFlags := flag.NewFlagSet("models", flag.ExitOnError)
	modelsFlags.Usage = func() {
		fmt.Printf("Usage of %s %s:\n", os.Args[0], modelsFlags.Name())
		fmt.Printf("  list\n  \tList existing models\n")
		fmt.Printf("  download\n  \tDownload existing models\n")
	}

	modelsListFlags := flag.NewFlagSet("list", flag.ExitOnError)
	modelsListFlags.StringVar(&bucket, "bucket", os.Getenv("RC_TRAIN_BUCKET"), "AWS bucket where store data required, use RC_TRAIN_BUCKET if arg not set")

	var modelPathBucket string
	modelsDownloadFlags := flag.NewFlagSet("download", flag.ExitOnError)
	modelsDownloadFlags.StringVar(&bucket, "bucket", os.Getenv("RC_TRAIN_BUCKET"), "AWS bucket where store data required, use RC_TRAIN_BUCKET if arg not set")
	modelsDownloadFlags.StringVar(&modelPathBucket, "model", "", "S3 Model key into bucket (mandatory)")
	modelsDownloadFlags.StringVar(&trainArchiveName, "output", os.Getenv("TRAIN_ARCHIVE_NAME"), "Zip archive file name, use TRAIN_ARCHIVE_NAME if args not set")
	flag.Parse()

	config := zap.NewDevelopmentConfig()
	if debug {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	lgr, err := config.Build()
	if err != nil {
		log.Fatalf("unable to init logger: %v", err)
	}
	defer func() {
		if err := lgr.Sync(); err != nil {
			log.Printf("unable to Sync logger: %v\n", err)
		}
	}()
	zap.ReplaceGlobals(lgr)

	// Switch on the subcommand
	// Parse the flags for appropriate FlagSet
	// FlagSet.Parse() requires a set of arguments to parse as input
	// os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
	switch flag.Arg(0) {
	case displayFlags.Name():
		if err := displayFlags.Parse(os.Args[2:]); err == flag.ErrHelp {
			displayFlags.PrintDefaults()
			os.Exit(0)
		}
		switch displayFlags.Arg(0) {
		case displayRecordFlags.Name():
			if err := displayRecordFlags.Parse(os.Args[3:]); err == flag.ErrHelp {
				displayRecordFlags.PrintDefaults()
				os.Exit(0)
			}
			client, err := cli.Connect(mqttBroker, username, password, clientId)
			if err != nil {
				zap.S().Fatalf("unable to connect to mqtt bus: %v", err)
			}
			runDisplayRecord(client, recordTopic)
		case displayCameraFlags.Name():
			if err := displayCameraFlags.Parse(os.Args[3:]); err == flag.ErrHelp {
				displayCameraFlags.PrintDefaults()
				os.Exit(0)
			}
			client, err := cli.Connect(mqttBroker, username, password, clientId)
			if err != nil {
				zap.S().Fatalf("unable to connect to mqtt bus: %v", err)
			}
			defer client.Disconnect(50)
			runDisplay(client, framePath, frameTopic, fps, objectsTopic, roadTopic, throttleFeedbackTopic, withObjects, withRoad, withThrottleFeedback)
		default:
			displayFlags.PrintDefaults()
			os.Exit(0)
		}
	case recordFlags.Name():
		if err := recordFlags.Parse(os.Args[2:]); err == flag.ErrHelp {
			recordFlags.PrintDefaults()
			os.Exit(0)
		}
		client, err := cli.Connect(mqttBroker, username, password, clientId)
		if err != nil {
			log.Fatalf("unable to connect to mqtt bus: %v", err)
		}
		defer client.Disconnect(50)
		runRecord(client, recordsPath, recordTopic)
	case impdkFlags.Name():
		if err := impdkFlags.Parse(os.Args[2:]); err == flag.ErrHelp {
			impdkFlags.PrintDefaults()
			os.Exit(0)
		}
		runImportDonkeyRecords(basedir, destdir)
	case trainingFlags.Name():
		if err := trainingFlags.Parse(os.Args[2:]); err == flag.ErrHelp {
			trainingFlags.PrintDefaults()
			os.Exit(0)
		}
		switch trainingFlags.Arg(0) {
		case trainingListJobFlags.Name():
			if err := trainingListJobFlags.Parse(os.Args[3:]); err == flag.ErrHelp {
				trainingListJobFlags.PrintDefaults()
				os.Exit(0)
			}
			runTrainList()
		case trainingRunFlags.Name():
			if err := trainingRunFlags.Parse(os.Args[3:]); err == flag.ErrHelp {
				trainingRunFlags.PrintDefaults()
				os.Exit(0)
			}
			runTraining(bucket, ociImage, roleArn, trainJobName, recordsPath, train.ParseModelType(modelType), trainSliceSize, trainImageWidth, trainImageHeight, horizon, withFlipImage, modelPath, enableSpotTraining)
		case trainArchiveFlags.Name():
			if err := trainArchiveFlags.Parse(os.Args[3:]); err == flag.ErrHelp {
				trainArchiveFlags.PrintDefaults()
				os.Exit(0)
			}
			runTrainArchive(recordsPath, trainArchiveName, trainSliceSize, trainImageWidth, trainImageHeight, horizon, withFlipImage)
		default:
			trainingFlags.PrintDefaults()
			os.Exit(0)

		}
	case modelsFlags.Name():

		if err := modelsFlags.Parse(os.Args[2:]); err == flag.ErrHelp {
			modelsFlags.PrintDefaults()
			os.Exit(0)
		}
		switch modelsFlags.Arg(0) {
		case modelsListFlags.Name():
			if err := modelsListFlags.Parse(os.Args[3:]); err == flag.ErrHelp {
				modelsListFlags.PrintDefaults()
				os.Exit(0)
			}
			runModelsList(bucket)
		case modelsDownloadFlags.Name():
			if err := modelsDownloadFlags.Parse(os.Args[3:]); err == flag.ErrHelp {
				modelsDownloadFlags.PrintDefaults()
				os.Exit(0)
			}
			if trainArchiveName == "" {
				zap.S().Error("output model file is mandatory")
				modelsDownloadFlags.PrintDefaults()
				os.Exit(1)
			}
			runModelsDownload(bucket, modelPathBucket, trainArchiveName)
		default:
			modelsFlags.PrintDefaults()
			os.Exit(0)
		}

	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func runRecord(client mqtt.Client, recordsDir, recordTopic string) {

	r, err := record.New(client, recordsDir, recordTopic)
	if err != nil {
		zap.S().Fatalf("unable to init record part: %v", err)
	}
	defer r.Stop()

	cli.HandleExit(r)

	err = r.Start()
	if err != nil {
		zap.S().Fatalf("unable to start service: %v", err)
	}
}

func runTrainArchive(basedir, archiveName string, sliceSize int, imgWidth, imgHeight int, horizon int, withFlipImage bool) {

	err := data.WriteArchive(basedir, archiveName, sliceSize, imgWidth, imgHeight, horizon, withFlipImage)
	if err != nil {
		zap.S().Fatalf("unable to build archive file %v: %v", archiveName, err)
	}
}

func runImportDonkeyRecords(basedir, destdir string) {
	if destdir == "" || basedir == "" {
		zap.S().Fatal("invalid arg")
	}
	err := dkimpt.ImportDonkeyRecords(basedir, destdir)
	if err != nil {
		zap.S().Fatalf("unable to import files from %v to %v: %v", basedir, destdir, err)
	}
}
func runDisplayRecord(client mqtt.Client, recordTopic string) {
	r := display.NewRecordDisplay(client, recordTopic)
	defer r.Stop()

	cli.HandleExit(r)
	err := r.Start()
	if err != nil {
		zap.S().Fatalf("unable to start service: %v", err)
	}
}
func runDisplay(client mqtt.Client, framePath string, frameTopic string, fps int, objectsTopic, roadTopic, throttleFeedbackTopic string,
	withObjects, withRoad, withThrottleFeedback bool) {

	if framePath != "" {
		camera, err := video.NewCameraFake(client, frameTopic, framePath, fps)
		if err != nil {
			log.Fatalf("unable to load fake camera: %v", err)
		}
		if err = camera.Start(); err != nil {
			log.Fatalf("unable to start fake camera: %v", err)
		}
		defer camera.Stop()
	}

	p := part.NewPart(client, frameTopic,
		objectsTopic, roadTopic, throttleFeedbackTopic,
		withObjects, withRoad, withThrottleFeedback)
	defer p.Stop()

	cli.HandleExit(p)

	err := p.Start()
	if err != nil {
		zap.S().Fatalf("unable to start service: %v", err)
	}
}

func runTraining(bucketName, ociImage, roleArn, jobName, dataDir string, modelType train.ModelType, sliceSize, imgWidth, imgHeight int, horizon int, withFlipImage bool, outputModel string, enableSpotTraining bool) {

	l := zap.S()
	if bucketName == "" {
		l.Fatalf("no bucket define, see help")
	}
	if ociImage == "" {
		l.Fatalf("no oci image define, see help")
	}
	if jobName == "" {
		l.Fatalf("no job name define, see help")
	}
	if dataDir == "" {
		l.Fatalf("no training data define, see help")
	}
	if outputModel == "" {
		l.Fatalf("no output model path define, see help")
	}

	if sliceSize != 0 && sliceSize != 2 {
		l.Fatalf("invalid value for sie-slice, only '0' or '2' are allowed")
	}

	if modelType == train.ModelTypeUnknown {
		l.Fatalf("invalid model type: %v", modelType)
	}

	training := train.New(bucketName, ociImage, roleArn)
	err := training.TrainDir(context.Background(), jobName, dataDir, modelType, imgWidth, imgHeight, sliceSize, horizon, withFlipImage, outputModel, enableSpotTraining)

	if err != nil {
		l.Fatalf("unable to run training: %v", err)
	}
}

func runTrainList() {
	err := train.ListJob(context.Background())
	if err != nil {
		zap.S().Fatalf("unable to list training jobs: %w", err)
	}
}

func runModelsList(bucketName string) {
	err := models.ListModels(context.Background(), bucketName)
	if err != nil {
		zap.S().Fatalf("unable to list models: %s", err)
	}

}

func runModelsDownload(bucketName, modelPath, output string) {
	err := models.DownloadArchiveToFile(context.Background(), bucketName, modelPath, output)
	if err != nil {
		zap.S().Fatalf("unable to download model: %s", err)
	}
}
