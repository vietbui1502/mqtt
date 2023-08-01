package dto

import "encoding/json"

//Define inital query format on mqtt message payload
type GeneralMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:data`
}
