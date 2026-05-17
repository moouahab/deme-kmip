package metrics

import "context"

type Snapshot struct {
	HTTPRequestsTotal int `json:"http_requests_total"`
	HTTPErrorsTotal   int `json:"http_errors_total"`

	CreateKeyTotal  int `json:"create_key_total"`
	GetKeyTotal     int `json:"get_key_total"`
	DestroyKeyTotal int `json:"destroy_key_total"`

	SuccessTotal  int `json:"success_total"`
	NotFoundTotal int `json:"not_found_total"`
}

type Collector interface {
	IncHTTPRequests(ctx context.Context)
	IncHTTPErrors(ctx context.Context)

	IncCreateKey(ctx context.Context)
	IncGetKey(ctx context.Context)
	IncDestroyKey(ctx context.Context)

	IncSuccess(ctx context.Context)
	IncNotFound(ctx context.Context)

	Snapshot(ctx context.Context) Snapshot
}
