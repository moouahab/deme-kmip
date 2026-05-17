package audit

import "time"

type Event struct {
	Time      time.Time `json:"time"`
	Operation string    `json:"operation"`
	KeyID     string    `json:"key_id,omitempty"`
	Status    string    `json:"status,omitempty"`
	Result    string    `json:"result"`
	Error     string    `json:"error,omitempty"`
}
