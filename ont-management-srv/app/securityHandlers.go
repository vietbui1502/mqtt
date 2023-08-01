package app

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/vietbui1502/mqtt/ont-management-srv/dto"
	"github.com/vietbui1502/mqtt/ont-management-srv/service"
)

type SecurityHandlers struct {
	service service.SecurityService
}

func (sh *SecurityHandlers) domainQuery(client mqtt.Client, msg mqtt.Message) {
	//For debug
	log.Printf("Topic: %s\nPayload: %s\n", msg.Topic(), msg.Payload())
	// Extract the message payload as a string
	payload := msg.Payload()

	// Unmarshal the JSON payload into the General struct
	var raw_request dto.GeneralMessage
	if err := json.Unmarshal(payload, &raw_request); err != nil {
		log.Printf("Error unmarshaling JSON message Payload: %v", err)
		return
	}

	// Unmarshall the JSON raw request to domain request struct
	var domainRequest dto.DomainRequestPayload
	if err := json.Unmarshal([]byte(raw_request.Data), &domainRequest); err != nil {
		log.Printf("Error unmarshaling data JSON: %v", err)
		return
	}

	response, err := sh.service.GetDomainCategory(domainRequest)

	if err != nil {
		log.Printf("Error Get domain category service: %v", err)
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

}
