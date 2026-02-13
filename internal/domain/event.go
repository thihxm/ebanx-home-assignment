package domain

type EventRequest struct {
	Type        string `json:"type"`
	Origin      string `json:"origin,omitempty"`
	Destination string `json:"destination,omitempty"`
	Amount      int    `json:"amount"`
}

type EventResponse struct {
	Origin      *Account `json:"origin,omitempty"`
	Destination *Account `json:"destination,omitempty"`
}
