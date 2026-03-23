package usecases

import (
	"context"
	"fmt"

	"agrinode/internal/core/domain"
	"agrinode/internal/core/ports"
)

// IngestionUseCase is the struct that holds our dependencies (like kafka)
type IngestionUseCase struct {
	publisher ports.TelemetryPublisher
}

// NewIngestionUseCase is a constructor function. It forces anyone creating
// this use case to provide a valid publisher. This is called Dependency Injection
func NewIngestionUseCase(pub ports.TelemetryPublisher) *IngestionUseCase {
	return &IngestionUseCase{
		publisher: pub,
	}
}

// ProcessRawTelemetry is the actual business flow. It takes raw data,
// validates it, and publishes it to the kafka topic.
func (uc *IngestionUseCase) ProcessRawTelemetry(ctx context.Context, t domain.Telemetry) error {
	// Step 1: Validate the rules using the Domain entity we wrote earlier.
	// Notice the capital 'V' in Validate()!
	err := t.Validate()
	if err != nil {
		// If data is physically impossible we reject it immediately
		return fmt.Errorf("invalid telemetry data: %w", err)
	}

	// Step 2: If valid, drop it off at the kafka "Post Office".
	// We are hardcoding the topic name "raw-telemetry-stream" for now.
	err = uc.publisher.Publish(ctx, "raw-telemetry-stream", t)
	if err != nil {
		return fmt.Errorf("failed to publish to kafka: %w", err)
	}

	return nil
}
