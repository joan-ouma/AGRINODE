package domain

import (
	"errors"
	"time"
)

// Define custom errors so we can handle them specifically in the use-case layer later
var (
	ErrInvalidTemp     = errors.New("temperature out of range")
	ErrInvalidHumidity = errors.New("humidity must be between 0 and 100")
	ErrInvalidMoisture = errors.New("moisture must be between 0 and 100")
	ErrMissingNodeID   = errors.New("node id is required")
)

// Telemetry represents the raw data coming from the edge nodes
// the struct itself is responsible for its own state
type Telemetry struct {
	NodeID       string
	Temperature  float64
	Humidity     float64
	SoilMoisture float64
	RecordedAt   time.Time
}

// Validate checks if the sensor readings make physical sense before we do anything with them
func (t *Telemetry) Validate() error {
	if t.NodeID == "" {
		return ErrMissingNodeID
	}

	// Setting practical environmental limits for an open farm
	if t.Temperature < -10 || t.Temperature > 60 {
		return ErrInvalidTemp
	}

	// Humidity is a percentage
	if t.Humidity < 0 || t.Humidity > 100 {
		return ErrInvalidHumidity
	}

	// Soil moisture is also a percentage (0 = completely dry, 100 = saturated)
	if t.SoilMoisture < 0 || t.SoilMoisture > 100 {
		return ErrInvalidMoisture
	}

	// If the hardware didn't send a timestamp, default to now
	if t.RecordedAt.IsZero() {
		t.RecordedAt = time.Now()
	}

	return nil
}
