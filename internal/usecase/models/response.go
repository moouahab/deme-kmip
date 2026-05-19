package models

type OperationResponse struct {
	KeyID      string         `json:"key_id,omitempty"`
	Status     string         `json:"status,omitempty"`
	Keys       []KeySummary   `json:"keys,omitempty"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

type KeySummary struct {
	ID         string `json:"id"`
	ObjectType uint32 `json:"object_type"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}
