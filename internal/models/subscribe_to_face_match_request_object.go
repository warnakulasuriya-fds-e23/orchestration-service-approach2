package models

type SubscribeToFaceMatchRequest struct {
	EventTypes []int  `json:"eventTypes"`
	EventDest  string `json:"eventDest"`
}
