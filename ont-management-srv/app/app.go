package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/vietbui1502/mqtt/ont-management-srv/domain"
	"github.com/vietbui1502/mqtt/ont-management-srv/dto"
	"github.com/vietbui1502/mqtt/ont-management-srv/service"
)

func Start() {
	// Initialize the domain data, to be update
	//data.Init()
	domainRepositoryDemo := domain.NewDomainRepositoryDemo()

	// Register service handler
	rh := RegisterHandlers{service: service.NewRegisterService()}

	// Secrity service handler, include domain query and further query
	sh := SecurityHandlers{service: service.NewSecurityService(domainRepositoryDemo)}

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
			//securitySericesHandle(client, msg)
			fmt.Printf("Received message for topic\nTopic: %s\nPayload: %s\n", msg.Topic(), msg.Payload())

			payload := msg.Payload()

			// Unmarshal the JSON payload into the General Message struct
			var request dto.GeneralMessage
			if err := json.Unmarshal(payload, &request); err != nil {
				log.Printf("Error unmarshaling JSON message Payload: %v", err)
			}

			// For debugging
			log.Printf("Received JSON data:\nEvent: %s\nData: %s\n", request.Event, request.Data)

			switch request.Event {
			case "domain_query":
				log.Printf("Domain query activity\n")
				sh.domainQuery(client, msg)
			case "domain_response":
				log.Printf("Enter Domain response activity\n")
			case "client_connected":
				log.Printf("Client connected event\n")
			case "client_disconnected":
				log.Printf("Client disconnected event\n")
			default:
				log.Printf("Unknown event\n")
			}
		} else {
			rh.firstTimeConnection(client, msg)
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
