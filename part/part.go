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

func NewPart(client mqtt.Client, frameTopic, objectsTopic, roadTopic string, withObjects, withRoad bool) *FramePart {
	return &FramePart{
		client:       client,
		frameTopic:   frameTopic,
		objectsTopic: objectsTopic,
		roadTopic:    roadTopic,
		window:       gocv.NewWindow(frameTopic),
		withObjects:  withObjects,
		withRoad:     withRoad,
		imgChan:      make(chan gocv.Mat),
		objectsChan:  make(chan events.ObjectsMessage),
		roadChan:     make(chan events.RoadMessage),
		cancel:       make(chan interface{}),
	}

}

type FramePart struct {
	client                              mqtt.Client
	frameTopic, objectsTopic, roadTopic string

	window      *gocv.Window
	withObjects bool
	withRoad    bool

	imgChan     chan gocv.Mat
	objectsChan chan events.ObjectsMessage
	roadChan    chan events.RoadMessage
	cancel      chan interface{}
}

func (p *FramePart) Start() error {
	if err := p.registerCallbacks(); err != nil {
		return fmt.Errorf("unable to start service: %v", err)
	}

	var img = gocv.NewMat()
	var objectsMsg events.ObjectsMessage
	var roadMsg events.RoadMessage

	for {
		select {
		case newImg := <-p.imgChan:
			img.Close()
			img = newImg
		case objects := <-p.objectsChan:
			objectsMsg = objects
		case road := <-p.roadChan:
			roadMsg = road
		case <-p.cancel:
			img.Close()
			return nil
		}
		p.drawFrame(&img, &objectsMsg, &roadMsg)
	}
}

func (p *FramePart) Stop() {
	defer p.window.Close()

	close(p.cancel)

	StopService("frame-display", p.client, p.frameTopic, p.roadTopic)
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
	var msg events.ObjectsMessage

	err := proto.Unmarshal(message.Payload(), &msg)
	if err != nil {
		log.Errorf("unable to unmarshal msg %T: %v", msg, err)
		return
	}

	p.objectsChan <- msg
}

func (p *FramePart) onRoad(_ mqtt.Client, message mqtt.Message) {
	var msg events.RoadMessage

	err := proto.Unmarshal(message.Payload(), &msg)
	if err != nil {
		log.Errorf("unable to unmarshal msg %T: %v", msg, err)
		return
	}

	p.roadChan <- msg
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

	if p.withRoad {
		err := service.RegisterCallback(p.client, p.roadTopic, p.onRoad)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *FramePart) drawFrame(img *gocv.Mat, objects *events.ObjectsMessage, road *events.RoadMessage) {

	if p.withObjects {
		p.drawObjects(img, objects)
	}
	if p.withRoad {
		p.drawRoad(img, road)
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

func (p *FramePart) drawRoad(img *gocv.Mat, road *events.RoadMessage) {
	cntr := make([]image.Point, 0, len(road.GetContour()))
	if road.GetContour() == nil || len(road.GetContour()) < 3 {
		return
	}
	for _, pt := range road.GetContour() {
		cntr = append(cntr, image.Point{X: int(pt.GetX()), Y: int(pt.GetY())})
	}

	gocv.DrawContours(
		img,
		[][]image.Point{cntr},
		0,
		color.RGBA{R: 255, G: 0, B: 0, A: 128,},
		-1)

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
