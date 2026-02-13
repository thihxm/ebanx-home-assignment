package domain

type EventRequest struct {
	Type        string `json:"type" validate:"required,oneof=deposit withdraw transfer"`
	Origin      string `json:"origin,omitempty" validate:"omitempty,required_if=Type withdraw,required_if=Type transfer,numeric"`
	Destination string `json:"destination,omitempty" validate:"omitempty,required_if=Type deposit,required_if=Type transfer,numeric"`
	Amount      int    `json:"amount" validate:"required,gt=0"`
}

type EventResponse struct {
	Origin      *Account `json:"origin,omitempty"`
	Destination *Account `json:"destination,omitempty"`
}

type EventService interface {
	ProcessEvent(event EventRequest) (*EventResponse, error)
}
