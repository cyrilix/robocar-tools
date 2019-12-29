package cli

import (
	"flag"
	"fmt"
	"github.com/cyrilix/robocar-base/service"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func SetDefaultValueFromEnv(value *string, key string, defaultValue string) {
	if os.Getenv(key) != "" {
		*value = os.Getenv(key)
	} else {
		*value = defaultValue
	}
}
func SetIntDefaultValueFromEnv(value *int, key string, defaultValue int) error {
	var sVal string
	if os.Getenv(key) != "" {
		sVal = os.Getenv(key)
		val, err := strconv.Atoi(sVal)
		if err != nil {
			log.Printf("unable to convert string to int: %v", err)
			return err
		}
		*value = val
	} else {
		*value = defaultValue
	}
	return nil
}
func SetFloat64DefaultValueFromEnv(value *float64, key string, defaultValue float64) error {
	var sVal string
	if os.Getenv(key) != "" {
		sVal = os.Getenv(key)
		val, err := strconv.ParseFloat(sVal, 64)
		if err != nil {
			log.Printf("unable to convert string to float: %v", err)
			return err
		}
		*value = val
	} else {
		*value = defaultValue
	}
	return nil
}

func HandleExit(p service.Part) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals
		p.Stop()
		os.Exit(0)
	}()
}

func InitMqttFlags(defaultClientId string, mqttBroker, username, password, clientId *string, mqttQos *int, mqttRetain *bool) {
	SetDefaultValueFromEnv(clientId, "MQTT_CLIENT_ID", defaultClientId)
	SetDefaultValueFromEnv(mqttBroker, "MQTT_BROKER", "tcp://127.0.0.1:1883")

	flag.StringVar(mqttBroker, "mqtt-broker", *mqttBroker, "Broker Uri, use MQTT_BROKER env if arg not set")
	flag.StringVar(username, "mqtt-username", os.Getenv("MQTT_USERNAME"), "Broker Username, use MQTT_USERNAME env if arg not set")
	flag.StringVar(password, "mqtt-password", os.Getenv("MQTT_PASSWORD"), "Broker Password, MQTT_PASSWORD env if args not set")
	flag.StringVar(clientId, "mqtt-client-id", *clientId, "Mqtt client id, use MQTT_CLIENT_ID env if args not set")
	flag.IntVar(mqttQos, "mqtt-qos", *mqttQos, "Qos to pusblish message, use MQTT_QOS env if arg not set")
	flag.BoolVar(mqttRetain, "mqtt-retain", *mqttRetain, "Retain mqtt message, if not set, true if MQTT_RETAIN env variable is set")
}

func InitIntFlag(key string, defValue int) int {
	var value int
	err := SetIntDefaultValueFromEnv(&value, key, defValue)
	if err != nil {
		log.Panicf("invalid int value: %v", err)
	}
	return value
}

func InitFloat64Flag(key string, defValue float64) float64 {
	var value float64
	err := SetFloat64DefaultValueFromEnv(&value, key, defValue)
	if err != nil {
		log.Panicf("invalid value: %v", err)
	}
	return value
}

func Connect(uri, username, password, clientId string) (MQTT.Client, error) {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(uri)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	opts.SetAutoReconnect(true)
	opts.SetDefaultPublishHandler(
		//define a function for the default message handler
		func(client MQTT.Client, msg MQTT.Message) {
			fmt.Printf("TOPIC: %s\n", msg.Topic())
			fmt.Printf("MSG: %s\n", msg.Payload())
		})

	//create and start a client using the above ClientOptions
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("unable to connect to mqtt bus: %v", token.Error())
	}
	return client, nil
}
