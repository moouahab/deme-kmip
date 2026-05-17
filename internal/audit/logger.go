package audit

import "context"

type Logger interface {
	Log(ctx context.Context, event Event) error
}
