package domain //means the file belongs to domain package since it is in the same directory as domain
//it has direct access to the telemetry struct and the custom errors in there without the need to import them

import (
	"testing" //this library provides tools needed to run automated tests
)

func TestTelemetryValidation(t *testing.T) {
	tests := []struct {
		name    string
		data    Telemetry
		wantErr error
	}{
		{
			name: "valid telemetry",
			data: Telemetry{
				NodeID:       "node-1",
				Temperature:  25.5,
				Humidity:     60.0,
				SoilMoisture: 45.0,
			},
			wantErr: nil,
		},
		{
			name: "missing node id",
			data: Telemetry{
				Temperature:  25.5,
				Humidity:     60.0,
				SoilMoisture: 45.0,
			},
			wantErr: ErrMissingNodeID,
		},
		{
			name: "invalid temperature",
			data: Telemetry{
				NodeID:       "node-2",
				Temperature:  100.0,
				Humidity:     50.0,
				SoilMoisture: 50.0,
			},
			wantErr: ErrInvalidTemp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate()
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}
