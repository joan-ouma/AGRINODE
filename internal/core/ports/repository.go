package ports

import (
	"agrinode/internal/core/domain"
	"context"
)

// TelemetryRepository defines the contract for saving telemetry data
// it says noting about postgresql
type TelemetryRepository interface {
	save(ctx context.Context, t domain.Telemetry) error
}
