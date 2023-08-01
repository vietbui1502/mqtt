package dto

// ONT information send to the VCS Cloud first time
type NewOntRegisterPayload struct {
	Serial string `json:"sn"`
	Vendor string `json:"vendor"`
}

type NewOntRegister struct {
	Event string                `json:"event"`
	Data  NewOntRegisterPayload `json:"data"`
}

type NewOntResponsePayload struct {
	Serial string `json:"sn"`
	Vendor string `json:"vendor"`
	Topic  string `json:"topic"`
}

type NewOntResponse struct {
	Event string                `json:"event"`
	Data  NewOntResponsePayload `json:data`
}
