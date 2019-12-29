# robocar-display

Tool to display camera frame and metrics
     
## Usage
`rc-display <OPTIONS>`

```
  -mqtt-broker string
        Broker Uri, use MQTT_BROKER env if arg not set (default "tcp://127.0.0.1:1883")
  -mqtt-client-id string
        Mqtt client id, use MQTT_CLIENT_ID env if args not set (default "robocar-frame-display")
  -mqtt-password string
        Broker Password, MQTT_PASSWORD env if args not set
  -mqtt-qos int
        Qos to pusblish message, use MQTT_QOS env if arg not set
  -mqtt-retain
        Retain mqtt message, if not set, true if MQTT_RETAIN env variable is set
  -mqtt-topic-frame string
        Mqtt topic that contains frame to display, use MQTT_TOPIC_FRAME if args not set
  -mqtt-username string
        Broker Username, use MQTT_USERNAME env if arg not set
```

