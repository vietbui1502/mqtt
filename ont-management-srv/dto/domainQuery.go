package dto

const (
	Unknown                = 1
	Education              = 2
	News                   = 3
	Game                   = 4
	Pornography            = 5
	Online_Dating          = 6
	Financial              = 7
	Gambling_and_Phishing  = 8
	Force                  = 9
	Shopping               = 10
	Sports                 = 11
	Social_Networks        = 12
	Chat                   = 13
	Entertaiment           = 14
	Heath                  = 15
	Computers_and_Software = 16
	Social_Science         = 17
	Travel                 = 18
	Government             = 19
	Food_and_Drink         = 20
	Job_search             = 21
	Video_and_Audio        = 22
	Superstition           = 23
	Natural_World          = 24
	Take_Time              = 25
)

type DomainRequestPayload struct {
	Domain string `json:"domain"`
}

type DomainResponsePayload struct {
	Domain   string `json:"domain"`
	Category int    `json:"category"`
	Verdict  string `json:"verdict"`
}

type DomainResponse struct {
	Event string                `json:"event"`
	Data  DomainResponsePayload `json:"data"`
}
