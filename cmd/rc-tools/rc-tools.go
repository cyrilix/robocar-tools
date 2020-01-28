package main

import (
	"flag"
	"fmt"
	"github.com/cyrilix/robocar-base/cli"
	"github.com/cyrilix/robocar-tools/part"
	"github.com/cyrilix/robocar-tools/record"
	"github.com/cyrilix/robocar-tools/video"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	DefaultClientId = "robocar-tools"
)

func main() {
	var mqttBroker, username, password, clientId string
	var framePath string
	var fps int
	var frameTopic, objectsTopic, roadTopic, recordTopic string
	var withObjects, withRoad bool
	var jsonPath, imgPath string

	mqttQos := cli.InitIntFlag("MQTT_QOS", 0)
	_, mqttRetain := os.LookupEnv("MQTT_RETAIN")

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("  display\n  \tDisplay events on live frames\n")
		fmt.Printf("  record \n  \tRecord event for tensorflow training\n")
	}

	displayFlags := flag.NewFlagSet("display", flag.ExitOnError)
	cli.InitMqttFlagSet(displayFlags, DefaultClientId, &mqttBroker, &username, &password, &clientId, &mqttQos, &mqttRetain)
	displayFlags.StringVar(&frameTopic, "mqtt-topic-frame", os.Getenv("MQTT_TOPIC_FRAME"), "Mqtt topic that contains frame to display, use MQTT_TOPIC_FRAME if args not set")
	displayFlags.StringVar(&framePath, "frame-path", "", "Directory path where to read jpeg frame to inject in frame topic")
	displayFlags.IntVar(&fps, "frame-per-second", 25, "Video frame per second of frame to publish")

	displayFlags.StringVar(&objectsTopic, "mqtt-topic-objects", os.Getenv("MQTT_TOPIC_OBJECTS"), "Mqtt topic that contains detected objects, use MQTT_TOPIC_OBJECTS if args not set")
	displayFlags.BoolVar(&withObjects, "with-objects", false, "Display detected objects")

	displayFlags.StringVar(&roadTopic, "mqtt-topic-road", os.Getenv("MQTT_TOPIC_ROAD"), "Mqtt topic that contains road description, use MQTT_TOPIC_ROAD if args not set")
	displayFlags.BoolVar(&withRoad, "with-road", false, "Display detected road")

	recordFlags := flag.NewFlagSet("record", flag.ExitOnError)
	cli.InitMqttFlagSet(recordFlags, DefaultClientId, &mqttBroker, &username, &password, &clientId, &mqttQos, &mqttRetain)
	recordFlags.StringVar(&recordTopic, "mqtt-topic-records", os.Getenv("MQTT_TOPIC_RECORDS"), "Mqtt topic that contains record data for training, use MQTT_TOPIC_RECORDS if args not set")
	recordFlags.StringVar(&jsonPath, "record-json-path", os.Getenv("RECORD_JSON_PATH"), "Path where to write json files, use RECORD_JSON_PATH if args not set")
	recordFlags.StringVar(&imgPath, "record-image-path", os.Getenv("RECORD_IMAGE_PATH"), "Path where to write jpeg files, use RECORD_IMAGE_PATH if args not set")

	flag.Parse()

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
		client, err := cli.Connect(mqttBroker, username, password, clientId)
		if err != nil {
			log.Fatalf("unable to connect to mqtt bus: %v", err)
		}
		defer client.Disconnect(50)
		runDisplay(client, framePath, frameTopic, fps, objectsTopic, roadTopic, withObjects, withRoad)
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
		runRecord(client, jsonPath, imgPath, recordTopic)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	cmd := flag.Arg(1)
	switch cmd {
	case "display":
	case "record":
	default:
		log.Errorf("invalid command: %v", cmd)
	}

}

func runRecord(client mqtt.Client, jsonDir, imgDir string, recordTopic string) {

	r := record.New(client, jsonDir, imgDir, recordTopic)
	defer r.Stop()

	cli.HandleExit(r)

	err := r.Start()
	if err != nil {
		log.Fatalf("unable to start service: %v", err)
	}
}

func runDisplay(client mqtt.Client, framePath string, frameTopic string, fps int, objectsTopic string, roadTopic string, withObjects bool, withRoad bool) {

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
		objectsTopic, roadTopic,
		withObjects, withRoad)
	defer p.Stop()

	cli.HandleExit(p)

	err := p.Start()
	if err != nil {
		log.Fatalf("unable to start service: %v", err)
	}
}
