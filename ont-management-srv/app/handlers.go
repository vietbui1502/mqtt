package app

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/vietbui1502/mqtt/ont-management-srv/data"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// ONT information send to the VCS Cloud first time
type OntInfo struct {
	Serial string `json:"sn"`
	Vendor string `json:"vendor"`
}

//Define inital query format on mqtt message payload
type InitialQuery struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:data`
}

type OntTopicResp struct {
	Serial string `json:"sn"`
	Vendor string `json:"vendor"`
	Topic  string `json:"topic"`
}

type InitialResponse struct {
	Event string       `json:"event"`
	Data  OntTopicResp `json:data`
}

type DomainQuery struct {
	Domain string `json:"domain"`
}

type DomainResponse struct {
	Domain   string `json:domain`
	Category string `json:category`
}

//Function handle client first time connect to the VCS Cloud on default topic
func FirstTimeConnection(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Topic: %s\nPayload: %s\n", msg.Topic(), msg.Payload())
	// Extract the message payload as a string
	payload := msg.Payload()

	// Unmarshal the JSON payload into the InitialQuery struct
	var query InitialQuery
	if err := json.Unmarshal(payload, &query); err != nil {
		log.Printf("Error unmarshaling JSON message Payload: %v", err)
		return
	}

	// For debugging
	fmt.Printf("Received JSON data:\nEvent: %s\nData: %s\n", query.Event, query.Data)

	fmt.Printf("zzzzzzzz\n")

	if query.Event == "Hello_VCS" {
		fmt.Printf("Security Agent first time connect to VCS Security Cloud\n")
		// Unmarshal Payload as OntInfo struct
		var ontInfo OntInfo
		if err := json.Unmarshal([]byte(query.Data), &ontInfo); err != nil {
			log.Printf("Error unmarshaling data JSON: %v", err)
			return
		}

		// For debugging ontInfo
		fmt.Printf("ontInfo:\nSerial: %s\nVendor: %s\n", ontInfo.Serial, ontInfo.Vendor)

		//Description: Saving ONT information to the database and send new topic name back to the ONT
		//After that, ONT and Cloud will commnunicate on new topic

		responseTopic := msg.Topic()
		// Process the message and generate the response
		responsePayload, err := generateNewTopic(ontInfo)

		if err != nil {
			log.Printf("Error generateNewTopic: %v", err)
			return
		}

		if token := client.Publish(responseTopic, 0, false, responsePayload); token.Wait() && token.Error() != nil {
			log.Printf("Failed to publish new topic to ONT %v", token.Error())
		} else {
			log.Printf("Published new topic for the ONT : '%s'", responsePayload)
		}

	} else if query.Event == "Topic_offer" {
		fmt.Printf("ONT will handle this\n")
		return
	} else {
		return
	}
}

func SecuritySericesHandle(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message for topic\nTopic: %s\nPayload: %s\n", msg.Topic(), msg.Payload())

	payload := msg.Payload()

	// Unmarshal the JSON payload into the InitialQuery struct
	var query InitialQuery
	if err := json.Unmarshal(payload, &query); err != nil {
		log.Printf("Error unmarshaling JSON message Payload: %v", err)
		return
	}

	// For debugging
	fmt.Printf("Received JSON data:\nEvent: %s\nData: %s\n", query.Event, query.Data)

	switch query.Event {
	case "domain_query":
		fmt.Printf("Enter Domain query activity\n")
		responsePayload, err := processDomainQuery(query)
		if err != nil {
			log.Printf("Error processDomainQuery: %v", err)
			return
		}
		if token := client.Publish(msg.Topic(), 0, false, responsePayload); token.Wait() && token.Error() != nil {
			log.Printf("Failed to publish new topic to ONT %v", token.Error())
		} else {
			log.Printf("Published new topic for the ONT : '%s'", responsePayload)
		}
		return
	case "domain_response":
		fmt.Printf("Enter Domain response activity\n")
		return
	case "client_connected":
		fmt.Printf("Client connected event\n")
		return
	case "client_disconnected":
		fmt.Printf("Client disconnected event\n")
		return
	default:
		fmt.Printf("Unknown event\n")
		return
	}

	// // Process the message and generate the response
	// //responsePayload := processMessage(msg.Payload())

	// // Publish the response back to the client on the same topic
	// //responseTopic := msg.Topic()
	// if token := client.Publish(responseTopic, 0, false, responsePayload); token.Wait() && token.Error() != nil {
	// 	log.Printf("Failed to publish response to topic '%s': %v", responseTopic, token.Error())
	// } else {
	// 	log.Printf("Published response '%s' to topic '%s'", responsePayload, responseTopic)
	// }
}

// Process the query and generate new topic
func generateNewTopic(ontInfo OntInfo) (response string, err error) {

	newTopic := "VCS2023/" + ontInfo.Vendor + "_" + ontInfo.Serial

	topicResp := OntTopicResp{
		Serial: ontInfo.Serial,
		Vendor: ontInfo.Vendor,
		Topic:  newTopic,
	}

	responseData := InitialResponse{
		Event: "Topic_offer",
		Data:  topicResp,
	}

	// Marshal the struct to JSON format
	jsonData, err := json.Marshal(responseData)

	if err != nil {
		// Handle the error, if necessary
		fmt.Println("Error marshaling JSON:", err)
		return "", err
	}

	return string(jsonData), nil
}

func processDomainQuery(query InitialQuery) (response string, er error) {
	var domainQuery DomainQuery
	if err := json.Unmarshal([]byte(query.Data), &domainQuery); err != nil {
		log.Printf("Error unmarshaling data JSON: %v", err)
		return "", err
	}

	category := "unknown"

	// For debugging ontInfo
	fmt.Printf("Domain query:%s\n", domainQuery.Domain)

	for _, item := range data.SexualDomain {
		if item == domainQuery.Domain {
			category = "sexual"
		}
	}

	outerJson := map[string]interface{}{
		"event": "domain_response",
		"data": DomainResponse{
			Domain:   domainQuery.Domain,
			Category: category,
		},
	}

	// Marshal the JSON object to JSON format
	jsonBytes, err := json.Marshal(outerJson)
	if err != nil {
		fmt.Println("Error:", err)
		return "", nil
	}

	// Convert the JSON byte slice to a string for printing
	jsonString := string(jsonBytes)
	return jsonString, nil
}
