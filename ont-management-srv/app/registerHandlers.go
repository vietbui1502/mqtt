package app

import (
	"encoding/json"
	"log"

	"github.com/vietbui1502/mqtt/ont-management-srv/dto"
	"github.com/vietbui1502/mqtt/ont-management-srv/service"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type RegisterHandlers struct {
	service service.RegisterService
}

func (rh *RegisterHandlers) firstTimeConnection(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Topic: %s\nPayload: %s\n", msg.Topic(), msg.Payload())
	// Extract the message payload as a string
	payload := msg.Payload()

	// Unmarshal the JSON payload into the General struct
	var request dto.GeneralMessage
	if err := json.Unmarshal(payload, &request); err != nil {
		log.Printf("Error unmarshaling JSON message Payload: %v", err)
		return
	}

	// For debugging
	log.Printf("Received JSON data:\nEvent: %s\nData: %s\n", request.Event, request.Data)

	if request.Event == "hello_vcs" {
		log.Printf("Security Agent first time connect to VCS Security Cloud\n")

		response, err := rh.service.NewOntRegister(request)

		if err != nil {
			log.Printf("Error New Ont register service: %v", err)
			return
		}

		// Marshal the struct to JSON format
		jsonData, err := json.Marshal(response)

		if err != nil {
			// Handle the error, if necessary
			log.Println("Error marshaling JSON:", err)
			return
		}

		responseJson := string(jsonData)

		responseTopic := msg.Topic()

		if token := client.Publish(responseTopic, 0, false, responseJson); token.Wait() && token.Error() != nil {
			log.Printf("Failed to publish new topic to ONT %v", token.Error())
		} else {
			log.Printf("Published new topic for the ONT : '%s'", responseJson)
		}

	} else if request.Event == "topic_offer" {
		log.Printf("ONT will handle this\n")
		return
	} else {
		log.Printf("Unknow event on Initial Topic")
		return
	}
}
