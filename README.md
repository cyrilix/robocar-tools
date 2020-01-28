# robocar-display

Tools to manage robocar
     
## Display

Display events on live frames.

```
rc-tool display
   
   Usage of display:
     -frame-path string
           Directory path where to read jpeg frame to inject in frame topic
     -frame-per-second int
           Video frame per second of frame to publish (default 25)
     -mqtt-broker string
           Broker Uri, use MQTT_BROKER env if arg not set (default "tcp://127.0.0.1:1883")
     -mqtt-client-id string
           Mqtt client id, use MQTT_CLIENT_ID env if args not set (default "robocar-tools")
     -mqtt-password string
           Broker Password, MQTT_PASSWORD env if args not set
     -mqtt-qos int
           Qos to pusblish message, use MQTT_QOS env if arg not set
     -mqtt-retain
           Retain mqtt message, if not set, true if MQTT_RETAIN env variable is set
     -mqtt-topic-frame string
           Mqtt topic that contains frame to display, use MQTT_TOPIC_FRAME if args not set
     -mqtt-topic-objects string
           Mqtt topic that contains detected objects, use MQTT_TOPIC_OBJECTS if args not set
     -mqtt-topic-road string
           Mqtt topic that contains road description, use MQTT_TOPIC_ROAD if args not set
     -mqtt-username string
           Broker Username, use MQTT_USERNAME env if arg not set
     -with-objects
           Display detected objects
     -with-road
           Display detected road
```

## Record

Record event for tensorflow training

```
rc-tools record

Usage of record:
  -mqtt-broker string
        Broker Uri, use MQTT_BROKER env if arg not set (default "tcp://127.0.0.1:1883")
  -mqtt-client-id string
        Mqtt client id, use MQTT_CLIENT_ID env if args not set (default "robocar-tools")
  -mqtt-password string
        Broker Password, MQTT_PASSWORD env if args not set
  -mqtt-qos int
        Qos to pusblish message, use MQTT_QOS env if arg not set
  -mqtt-retain
        Retain mqtt message, if not set, true if MQTT_RETAIN env variable is set
  -mqtt-topic-records string
        Mqtt topic that contains record data for training, use MQTT_TOPIC_RECORDS if args not set
  -mqtt-username string
        Broker Username, use MQTT_USERNAME env if arg not set
  -record-image-path string
        Path where to write jpeg files, use RECORD_IMAGE_PATH if args not set
  -record-json-path string
        Path where to write json files, use RECORD_JSON_PATH if args not set
```
