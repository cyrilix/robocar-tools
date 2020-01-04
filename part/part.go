package part

import (
	"fmt"
	"github.com/cyrilix/robocar-base/service"
	"github.com/cyrilix/robocar-protobuf/go/events"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
	"image"
	"image/color"
)

func NewPart(client mqtt.Client, frameTopic, objectsTopic string, withObjects bool) *FramePart {
	return &FramePart{
		client:       client,
		frameTopic:   frameTopic,
		objectsTopic: objectsTopic,
		window:       gocv.NewWindow(frameTopic),
		withObjects:  withObjects,
		imgChan:      make(chan gocv.Mat),
		objectsChan:  make(chan events.ObjectsMessage),
		cancel:       make(chan interface{}),
	}

}

type FramePart struct {
	client                   mqtt.Client
	frameTopic, objectsTopic string

	window      *gocv.Window
	withObjects bool

	imgChan     chan gocv.Mat
	objectsChan chan events.ObjectsMessage
	cancel      chan interface{}
}

func (p *FramePart) Start() error {
	if err := p.registerCallbacks(); err != nil {
		return fmt.Errorf("unable to start service: %v", err)
	}

	var img = gocv.NewMat()
	var objectsMsg events.ObjectsMessage
	for {
		select {
		case newImg := <-p.imgChan:
			img.Close()
			img = newImg
		case objects := <-p.objectsChan:
			objectsMsg = objects
		case <-p.cancel:
			img.Close()
			return nil
		}
		p.drawFrame(&img, &objectsMsg)
	}
}

func (p *FramePart) Stop() {
	defer p.window.Close()

	close(p.cancel)

	StopService("frame-display", p.client, p.frameTopic)
}

func (p *FramePart) onFrame(_ mqtt.Client, message mqtt.Message) {
	var msg events.FrameMessage
	err := proto.Unmarshal(message.Payload(), &msg)
	if err != nil {
		log.Errorf("unable to unmarshal protobuf FrameMessage: %v", err)
		return
	}

	img, err := gocv.IMDecode(msg.Frame, gocv.IMReadUnchanged)
	if err != nil {
		log.Errorf("unable to decode image: %v", err)
		return
	}
	log.Infof("[%v] frame %v", message.Topic(), msg.GetId())
	p.imgChan <- img
}

func (p *FramePart) onObjects(_ mqtt.Client, message mqtt.Message) {
	var objectsMsg events.ObjectsMessage

	err := proto.Unmarshal(message.Payload(), &objectsMsg)
	if err != nil {
		log.Errorf("unable to unmarshal detected objects: %v", err)
		return
	}

	p.objectsChan <- objectsMsg
}

func (p *FramePart) registerCallbacks() error {
	err := RegisterCallback(p.client, p.frameTopic, p.onFrame)
	if err != nil {
		return err
	}

	if p.withObjects {
		err := service.RegisterCallback(p.client, p.objectsTopic, p.onObjects)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *FramePart) drawFrame(img *gocv.Mat, objects *events.ObjectsMessage) {

	if p.withObjects {
		p.drawObjects(img, objects)
	}

	p.window.IMShow(*img)
	p.window.WaitKey(1)
}

func (p *FramePart) drawObjects(img *gocv.Mat, objects *events.ObjectsMessage) {
	for _, obj := range objects.GetObjects() {
		gocv.Rectangle(
			img,
			image.Rect(int(obj.GetLeft()), int(obj.GetTop()), int(obj.GetRight()), int(obj.GetBottom())),
			color.RGBA{0, 255, 0, 0},
			2)
	}
}

func StopService(name string, client mqtt.Client, topics ...string) {
	log.Printf("Stop %s service", name)
	token := client.Unsubscribe(topics...)
	token.Wait()
	if token.Error() != nil {
		log.Printf("unable to unsubscribe service: %v", token.Error())
	}
	client.Disconnect(50)
}

func RegisterCallback(client mqtt.Client, topic string, callback mqtt.MessageHandler) error {
	log.Printf("Register callback on topic %v", topic)
	token := client.Subscribe(topic, 0, callback)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("unable to register callback on topic %s: %v", topic, token.Error())
	}
	return nil
}
