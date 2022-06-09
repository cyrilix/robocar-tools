package display

import (
	"fmt"
	"github.com/cyrilix/robocar-protobuf/go/events"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"gocv.io/x/gocv"
	"image"
	"image/color"
)

func NewRecordDisplay(client mqtt.Client, recordTopic string) *Record {
	return &Record{
		client:      client,
		recordTopic: recordTopic,
		window:      gocv.NewWindow("recordTopic"),
		recordChan:  make(chan *events.RecordMessage),
		cancel:      make(chan interface{}),
	}

}

type Record struct {
	client      mqtt.Client
	recordTopic string

	window *gocv.Window

	recordChan chan *events.RecordMessage
	cancel     chan interface{}
}

func (r *Record) Start() error {
	if err := r.registerCallbacks(); err != nil {
		return fmt.Errorf("unable to start service: %v", err)
	}

	var rec *events.RecordMessage
	var objectsMsg events.ObjectsMessage
	var roadMsg events.RoadMessage

	for {
		select {
		case newRecord := <-r.recordChan:
			rec = newRecord
		case <-r.cancel:
			return nil
		}
		go r.drawRecord(rec, &objectsMsg, &roadMsg)
	}
}

func (r *Record) Stop() {
	defer r.window.Close()

	close(r.cancel)

	StopService("record-display", r.client, r.recordTopic)
}

func (r *Record) onRecord(_ mqtt.Client, message mqtt.Message) {
	var msg events.RecordMessage
	err := proto.Unmarshal(message.Payload(), &msg)
	if err != nil {
		zap.S().Errorf("unable to unmarshal protobuf FrameMessage: %v", err)
		return
	}
	message.Ack()
	r.recordChan <- &msg
}

func (r *Record) registerCallbacks() error {
	err := RegisterCallback(r.client, r.recordTopic, r.onRecord)
	if err != nil {
		return err
	}

	return nil
}

func (r *Record) drawRecord(rec *events.RecordMessage, objects *events.ObjectsMessage, road *events.RoadMessage) {

	img, err := gocv.IMDecode(rec.GetFrame().GetFrame(), gocv.IMReadUnchanged)
	if err != nil {
		zap.S().Errorf("unable to decode image: %v", err)
		return
	}
	defer img.Close()

	steering := rec.GetSteering().GetSteering()
	r.drawSteering(&img, steering)

	r.window.IMShow(img)
	r.window.WaitKey(1)
}


func (r *Record) drawSteering(img *gocv.Mat, steering float32) {
	gocv.PutText(
		img,
		fmt.Sprintf("Steering: %.3f", steering),
		image.Point{X: 20, Y: 20},
		gocv.FontHersheyPlain,
		1.,
		color.RGBA{R: 0, G: 255, B: 0, A: 255},
		1,
	)
}

func StopService(name string, client mqtt.Client, topics ...string) {
	zap.S().Infof("Stop %s service", name)
	token := client.Unsubscribe(topics...)
	token.Wait()
	if token.Error() != nil {
		zap.S().Errorf("unable to unsubscribe service: %v", token.Error())
	}
	client.Disconnect(50)
}

func RegisterCallback(client mqtt.Client, topic string, callback mqtt.MessageHandler) error {
	zap.S().Infof("Register callback on topic %v", topic)
	token := client.Subscribe(topic, 0, callback)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("unable to register callback on topic %s: %v", topic, token.Error())
	}
	return nil
}
