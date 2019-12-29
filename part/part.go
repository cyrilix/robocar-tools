package part

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gocv.io/x/gocv"
	"log"
	"time"
)

func NewPart(client mqtt.Client, frameTopic string) *FramePart {
	return &FramePart{
		client:     client,
		frameTopic: frameTopic,
		window:     gocv.NewWindow(frameTopic),
		img:        gocv.NewMat(),
	}

}

type FramePart struct {
	client     mqtt.Client
	frameTopic string
	window     *gocv.Window
	img        gocv.Mat
}

func (p *FramePart) Start() error {
	if err := p.registerCallbacks(); err != nil {
		return fmt.Errorf("unable to start service: %v", err)
	}
	for {
		time.Sleep(1 * time.Hour)
	}
}

func (p *FramePart) Stop() {
	defer func (){
		err := p.img.Close()
		if err != nil {
			log.Printf("unable to close resource: %v", err)
		}
	}()
	defer p.window.Close()
	StopService("frame-display", p.client, p.frameTopic)
}

func (p *FramePart) onFrame(_ mqtt.Client, message mqtt.Message) {
		img, err := gocv.IMDecode(message.Payload(), gocv.IMReadUnchanged)
		if err != nil {
			log.Printf("unable to decode image: %v", err)
			return
		}
		defer img.Close()
		img.CopyTo(&p.img)
		p.window.IMShow(p.img)
		p.window.WaitKey(1)
}

func (p *FramePart) registerCallbacks() error {
	err := RegisterCallback(p.client, p.frameTopic, p.onFrame)
	if err != nil {
		return err
	}

	return nil
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
