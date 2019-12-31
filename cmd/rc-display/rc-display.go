package main

import (
	"flag"
	"github.com/cyrilix/robocar-base/cli"
	"github.com/cyrilix/robocar-display/part"
	"log"
	"os"
)

const (
	DefaultClientId = "robocar-display"
)

func main() {
	var mqttBroker, username, password, clientId string
	var frameTopic, objectsTopic string
	var withObjects bool

	mqttQos := cli.InitIntFlag("MQTT_QOS", 0)
	_, mqttRetain := os.LookupEnv("MQTT_RETAIN")

	cli.InitMqttFlags(DefaultClientId, &mqttBroker, &username, &password, &clientId, &mqttQos, &mqttRetain)

	flag.StringVar(&frameTopic, "mqtt-topic-frame", os.Getenv("MQTT_TOPIC_FRAME"), "Mqtt topic that contains frame to display, use MQTT_TOPIC_FRAME if args not set")
	flag.StringVar(&objectsTopic, "mqtt-topic-objects", os.Getenv("MQTT_TOPIC_OBJECTS"), "Mqtt topic that contains detected objects, use MQTT_TOPIC_OBJECTS if args not set")
	flag.BoolVar(&withObjects, "with-objects", false, "Display detected objects")

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

	p := part.NewPart(client, frameTopic, objectsTopic, withObjects)
	defer p.Stop()

	cli.HandleExit(p)

	err = p.Start()
	if err != nil {
		log.Fatalf("unable to start service: %v", err)
	}
}
