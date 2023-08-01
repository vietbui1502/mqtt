package dto

type DomainRequestPayload struct {
	Domain string `json: "domain"`
}

type DomainResponsePayload struct {
	Domain   string `json: "domain"`
	Category string `json: "category"`
}

type DomainResponse struct {
	Event string                `json: "event"`
	Data  DomainResponsePayload `json: "data"`
}
