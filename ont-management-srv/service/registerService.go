package service

import (
	"encoding/json"
	"log"

	"github.com/vietbui1502/mqtt/ont-management-srv/dto"
)

type RegisterService interface {
	NewOntRegister(dto.GeneralMessage) (*dto.NewOntResponse, error)
}

type DefaultRegisterService struct {
}

func (s DefaultRegisterService) NewOntRegister(req dto.GeneralMessage) (*dto.NewOntResponse, error) {
	//Todo: Validate request
	//req.Validate()

	var registerPayload dto.NewOntRegisterPayload
	if err := json.Unmarshal([]byte(req.Data), &registerPayload); err != nil {
		log.Printf("Error unmarshaling data JSON: %v", err)
		return nil, err
	}

	// For debugging ontInfo
	log.Printf("registerPayload:\nSerial: %s\nVendor: %s\n", registerPayload.Serial, registerPayload.Vendor)

	//Description: Saving ONT information to the database and send new topic name back to the ONT
	//After that, ONT and Cloud will commnunicate on new topic

	response := generateNewTopic(registerPayload)

	//Todo: Save ONT information to database

	return response, nil
}

// Process the query and generate new topic
func generateNewTopic(registerPayload dto.NewOntRegisterPayload) *dto.NewOntResponse {

	newTopic := "VCS2023/" + registerPayload.Vendor + "_" + registerPayload.Serial

	responsePayload := dto.NewOntResponsePayload{
		Serial: registerPayload.Serial,
		Vendor: registerPayload.Vendor,
		Topic:  newTopic,
	}

	response := dto.NewOntResponse{
		Event: "topic_offer",
		Data:  responsePayload,
	}

	return &response
}

func NewRegisterService() DefaultRegisterService {
	return DefaultRegisterService{}
}
