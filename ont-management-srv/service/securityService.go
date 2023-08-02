package service

import (
	"log"

	"github.com/vietbui1502/mqtt/ont-management-srv/domain"
	"github.com/vietbui1502/mqtt/ont-management-srv/dto"
)

type SecurityService interface {
	GetDomainCategory(dto.DomainRequestPayload) (*dto.DomainResponse, error)
}

type DefaultSecurityService struct {
	repo domain.DomainRepository
}

func (s DefaultSecurityService) GetDomainCategory(req dto.DomainRequestPayload) (*dto.DomainResponse, error) {
	// Todo: validate request

	category, err := s.repo.FindDomainCategory(req.Domain)

	if err != nil {
		log.Printf("Error on Find domain category : %v", err)
		return nil, err
	}

	responsePayload := dto.DomainResponsePayload{
		Domain:   req.Domain,
		Category: category,
		Verdict:  "block",
	}

	response := dto.DomainResponse{
		Event: "domain_response",
		Data:  responsePayload,
	}

	return &response, nil
}

func NewSecurityService(repo domain.DomainRepository) DefaultSecurityService {
	return DefaultSecurityService{repo}
}
