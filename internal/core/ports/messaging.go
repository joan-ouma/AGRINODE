package ports

import (
	"agrinode/internal/core/domain"
	"context"
)

// TelemetryPublisher defines the contract for sending data to a message queue
// it says noting about kafka
type TelemetryPublisher interface {
	Publish(ctx context.Context, topic string, t domain.Telemetry) error
}
