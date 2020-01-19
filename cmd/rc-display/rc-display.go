package main

import (
	"flag"
	"github.com/cyrilix/robocar-base/cli"
	"github.com/cyrilix/robocar-display/part"
	"github.com/cyrilix/robocar-display/video"
	"log"
	"os"
)

const (
	DefaultClientId = "robocar-display"
)

func main() {
	var mqttBroker, username, password, clientId string
	var framePath string
	var fps int
	var frameTopic, objectsTopic, roadTopic string
	var withObjects, withRoad bool

	mqttQos := cli.InitIntFlag("MQTT_QOS", 0)
	_, mqttRetain := os.LookupEnv("MQTT_RETAIN")

	cli.InitMqttFlags(DefaultClientId, &mqttBroker, &username, &password, &clientId, &mqttQos, &mqttRetain)

	flag.StringVar(&frameTopic, "mqtt-topic-frame", os.Getenv("MQTT_TOPIC_FRAME"), "Mqtt topic that contains frame to display, use MQTT_TOPIC_FRAME if args not set")
	flag.StringVar(&framePath, "frame-path", "", "Directory path where to read jpeg frame to inject in frame topic")
	flag.IntVar(&fps, "frame-per-second", 25, "Video frame per second of frame to publish")

	flag.StringVar(&objectsTopic, "mqtt-topic-objects", os.Getenv("MQTT_TOPIC_OBJECTS"), "Mqtt topic that contains detected objects, use MQTT_TOPIC_OBJECTS if args not set")
	flag.BoolVar(&withObjects, "with-objects", false, "Display detected objects")

	flag.StringVar(&roadTopic, "mqtt-topic-road", os.Getenv("MQTT_TOPIC_ROAD"), "Mqtt topic that contains road description, use MQTT_TOPIC_ROAD if args not set")
	flag.BoolVar(&withRoad, "with-road", false, "Display detected road")

	flag.Parse()
	if len(os.Args) <= 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	client, err := cli.Connect(mqttBroker, username, password, clientId)
	if err != nil {
		log.Fatalf("unable to connect to mqtt bus: %v", err)
	}
	defer client.Disconnect(50)

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
		withObjects, withRoad )
	defer p.Stop()

	cli.HandleExit(p)

	err = p.Start()
	if err != nil {
		log.Fatalf("unable to start service: %v", err)
	}
}
