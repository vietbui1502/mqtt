package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/vietbui1502/mqtt/ont-management-srv/data"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Start() {
	// Initialize the domain data, to be update
	data.Init()

	// Configure MQTT client
	mqttBroker := "mqtt://127.0.0.1:1883"
	mqttClientID := "mqtt-client"
	mqttVcsIdentify := "VCS2023"
	mqttInitialTopic := "VCS2023/InitialTopic"

	// Create an MQTT client options
	opts := mqtt.NewClientOptions().AddBroker(mqttBroker)
	opts.SetClientID(mqttClientID)

	// Create an MQTT client
	client := mqtt.NewClient(opts)

	// Connect to the MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}

	// Define a callback function to handle incoming MQTT messages
	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		// Handle the incoming message as per your requirements
		fmt.Printf("Received message\nTopic: %s\nPayload: %s\n", msg.Topic(), msg.Payload())
		if msg.Topic() != mqttInitialTopic {
			securitySericesHandle(client, msg)
		} else {
			firstTimeConnection(client, msg)
		}
	}

	// Subscribe to all topics containing start with VCS2023
	topicFilter := fmt.Sprintf("%s/#", mqttVcsIdentify)
	if token := client.Subscribe(topicFilter, 0, messageHandler); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topics containing '%s': %v", mqttVcsIdentify, token.Error())
	} else {
		log.Printf("Subscribed to all topics of VCS Cloud Security")
	}

	// Wait for interrupt signal to gracefully stop the client
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Unsubscribe from the wildcard topics and disconnect from the broker
	if token := client.Unsubscribe(topicFilter); token.Wait() && token.Error() != nil {
		log.Printf("Failed to unsubscribe from topics containing '%s': %v", mqttVcsIdentify, token.Error())
	} else {
		log.Printf("Unsubscribed all topics of VCS Cloud Security")
	}

	client.Disconnect(250)
	log.Println("Disconnected from MQTT broker")
}
