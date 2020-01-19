package video

import (
	"encoding/base64"
	"fmt"
	"github.com/cyrilix/robocar-protobuf/go/events"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"sync"
	"testing"
	"time"
)

func TestNewCameraFake(t *testing.T) {
	oldPublish := publish
	defer func() {
		publish = oldPublish
	}()

	var muEventsPublished sync.Mutex
	eventsPublished := make([]*[]byte, 0, 10)
	publish = func(client mqtt.Client, topic string, payload *[]byte) {
		muEventsPublished.Lock()
		defer muEventsPublished.Unlock()
		eventsPublished = append(eventsPublished, payload)
	}

	videoDir := "testdata/video"
	fps := 250
	camera, err := NewCameraFake(nil, "topic/fake/video", videoDir, fps)
	if err != nil {
		t.Errorf("unable to load frame from directory %v: %v", videoDir, err)
	}

	// Video frame
	files, err := ioutil.ReadDir(videoDir)
	if err != nil {
		t.Errorf("unable to list files in directory %v: %v", videoDir, err)
	}

	if err := camera.Start(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer camera.Stop()

	time.Sleep(5 * time.Millisecond)
	for idx, frame := range files {
		muEventsPublished.Lock()
		if len(eventsPublished) < idx+1 {
			t.Errorf("frame %v has not been published", idx)
			t.FailNow()
		}

		msgPublished := eventsPublished[idx]
		var frameMsg events.FrameMessage
		err := proto.Unmarshal(*msgPublished, &frameMsg)
		if err != nil {
			t.Errorf("unable to unmarshal msg frame: %v", err)
		}

		frameContent, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", videoDir, frame.Name()))
		if err != nil {
			t.Errorf("unable to read frame file %v: %v", frame.Name(), err)
		}
		srcFrame := base64.StdEncoding.EncodeToString(frameContent)
		publishedFrame := base64.StdEncoding.EncodeToString(frameMsg.GetFrame())

		if srcFrame != publishedFrame {
			t.Errorf("frame signatures doesn't match: %v, wants %v", publishedFrame, srcFrame)
		}
		if frameMsg.GetId().GetName() != "camera" {
			t.Errorf("bad name frame: %v, wants %v", frameMsg.GetId().GetName(), "camera")
		}
		if len(frameMsg.GetId().GetId()) != 13 {
			t.Errorf("bad id length: %v, wants %v", len(frameMsg.GetId().GetId()), 13)
		}
		if frameMsg.GetId().GetCreatedAt() == nil {
			t.Errorf("missin CreatedAt field")
		}
		muEventsPublished.Unlock()
		time.Sleep(1 * time.Second / time.Duration(fps))
	}

}
