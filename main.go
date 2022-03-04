package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	"github.com/kevin-vargas/sidecar-log/pubsub"
)

const (
	TOPIC_LOGS          = "log"
	TOPIC_NOTIFICATIONS = "notification"
)

var toCompare = []byte("RESULT:")

func IsForNotification(payload []byte) bool {
	for i, elem := range payload {
		if i == (len(toCompare) - 1) {
			return true
		}
		if toCompare[i] != elem {
			break
		}
	}
	return false
}
func makeHandler(p pubsub.MQTTI, topic string) func(mqtt.Client, mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		payload := msg.Payload()
		if IsForNotification(payload) {
			fmt.Println(string(payload))
			p.Publish(topic, payload[len(toCompare):])
		}
	}
}

func main() {
	syncChan := make(chan bool, 0)
	fmt.Println("Init logs-filter")
	godotenv.Load()
	client := pubsub.New()
	handler := makeHandler(client, TOPIC_NOTIFICATIONS)
	client.SubscribeWithCB(TOPIC_LOGS, handler)
	<-syncChan
}
