package video

import (
	"fmt"
	"github.com/cyrilix/robocar-protobuf/go/events"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

type CameraFake struct {
	client     mqtt.Client
	frameTopic string
	videoPath  string
	fps        int
	cancel     chan interface{}
}

func NewCameraFake(client mqtt.Client, frameTopic string, videoPath string, fps int) (*CameraFake, error) {

	files, err := ioutil.ReadDir(videoPath)
	if err != nil {
		return nil, fmt.Errorf("unable to found camera frame in directory %v: %v", videoPath, err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no files in directory %v", videoPath)
	}
	return &CameraFake{
		client:     client,
		frameTopic: frameTopic,
		videoPath:  videoPath,
		fps:        fps,
		cancel:     make(chan interface{}),
	}, nil
}

func (c CameraFake) Start() error {
	files, err := ioutil.ReadDir(c.videoPath)
	if err != nil {
		return fmt.Errorf("unable to found camera frame in directory %v: %v", c.videoPath, err)
	}

	go c.loop(files)
	return nil
}

func (c CameraFake) loop(files []os.FileInfo) {
	ticker := time.NewTicker(time.Second / time.Duration(c.fps))
	defer ticker.Stop()

	for {

		for _, file := range files {
			framePath := fmt.Sprintf("%s/%s", c.videoPath, file.Name())
			frameContent, err := ioutil.ReadFile(framePath)
			if err != nil {
				log.Errorf("unable to load frame content for %v: %v", framePath, err)
				continue
			}
			now := time.Now()
			msg := &events.FrameMessage{
				Id: &events.FrameRef{
					Name: "camera",
					Id:   fmt.Sprintf("%d%000d", now.Unix(), now.Nanosecond()/1000/1000),
					CreatedAt: &timestamp.Timestamp{
						Seconds: now.Unix(),
						Nanos:   int32(now.Nanosecond()),
					},
				},
				Frame: frameContent,
			}

			payload, err := proto.Marshal(msg)
			if err != nil {
				log.Errorf("unable to marshal protobuf message: %v", err)
			}

			publish(c.client, c.frameTopic, &payload)

			select {
			case <-ticker.C:
			case <-c.cancel:
				return
			}
		}
	}
}

func (c CameraFake) Stop() {
	close(c.cancel)
}

var publish = func(client mqtt.Client, topic string, payload *[]byte) {
	client.Publish(topic, 0, true, *payload)
}
