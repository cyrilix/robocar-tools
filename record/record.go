package record

import (
	"encoding/json"
	"fmt"
	"github.com/cyrilix/robocar-base/service"
	"github.com/cyrilix/robocar-protobuf/go/events"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func New(client mqtt.Client, jsonDir, imgDir string, recordTopic string) *Recorder {
	return &Recorder{
		client:      client,
		jsonDir:     jsonDir,
		imgDir:      imgDir,
		recordTopic: recordTopic,
		cancel:      make(chan interface{}),
	}

}

type Recorder struct {
	client          mqtt.Client
	jsonDir, imgDir string
	recordTopic     string
	cancel          chan interface{}
}

func (r *Recorder) Start() error {
	err := service.RegisterCallback(r.client, r.recordTopic, r.onRecordMsg)
	if err != nil {
		return fmt.Errorf("unable to start recorder part: %v", err)
	}
	<-r.cancel
	return nil
}

func (r *Recorder) Stop() {
	service.StopService("record", r.client, r.recordTopic)
	close(r.cancel)
}

func (r *Recorder) onRecordMsg(_ mqtt.Client, message mqtt.Message) {
	var msg events.RecordMessage
	err := proto.Unmarshal(message.Payload(), &msg)
	if err != nil {
		log.Errorf("unable to unmarshal protobuf %T: %v", msg, err)
		return
	}

	imgName := fmt.Sprintf("cam-image_array_%s.jpg", msg.GetFrame().GetId().GetId())
	err = ioutil.WriteFile(imgName, msg.GetFrame().GetFrame(), 0755)
	if err != nil {
		log.Errorf("unable to write json file %v: %v", imgName, err)
		return
	}

	recordName := fmt.Sprintf("record_%s.jpg", msg.GetFrame().GetId().GetId())
	record := Record{
		UserAngle:     msg.GetSteering().GetSteering(),
		CamImageArray: imgName,
	}
	jsonBytes, err := json.Marshal(&record)
	if err != nil {
		log.Errorf("unable to marshal json content: %v", err)
		return
	}

	err = ioutil.WriteFile(recordName, jsonBytes, 0755)
	if err != nil {
		log.Errorf("unable to write json file %v: %v", recordName, err)
	}

}

type Record struct {
	UserAngle     float32 `json:"user/angle,"`
	CamImageArray string  `json:"cam/image_array,"`
}
