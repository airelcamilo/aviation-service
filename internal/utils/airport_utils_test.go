package utils_test

import (
	"fmt"
	"testing"

	"aviation-service/internal/dto"
	. "aviation-service/internal/utils"
)

var (
	atype        = "Airport"
	facilityName = "Lorem ipsum"
	faa          = "LAX"
	region       = "LA"
	state        = "CALIFORNIA"
	county       = "LOS ANGELES"
	city         = "LOS ANGELES"
	ownership    = "PU"
	use          = "PU"
	manager      = "Manager1"
	managerPhone = "317-487-9594"
	latitude     = "39-43-02.3000N"
	longitude    = "086-17-40.7000W"
)

func TestAirportValidator_IsComplete(t *testing.T) {
	tests := []struct {
		name           string
		airport        *dto.Airport
		expectedResult bool
	}{
		{
			name: "Success airport data is complete",
			airport: &dto.Airport{ID: 0, Type: &atype, FacilityName: &facilityName, FAA: &faa, ICAO: "KLAX", Region: &region, State: &state, County: &county, City: &city, Ownership: &ownership,
				Use: &use, Manager: &manager, ManagerPhone: &managerPhone, Latitude: &latitude, Longitude: &longitude, Status: "PENDING"},
			expectedResult: true,
		},
		{
			name:    "Success airport data is not complete",
			airport: &dto.Airport{ICAO: "KLAX"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := NewAirportValidator()
			got := av.IsComplete(tt.airport)
			if got != tt.expectedResult {
				t.Errorf("Expected result %v, got %v", tt.expectedResult, got)
			}
		})
	}
}

func TestAirportValidator_Validate(t *testing.T) {
	falseOwnership := "A"
	falseUse := "A"
	falseLatitude := "1"
	falseLongitude := "1"
	falseManagerPhone := "62"
	tests := []struct {
		name        string
		airport     *dto.Airport
		expectedErr error
	}{
		{
			name: "Success validate airport",
			airport: &dto.Airport{ID: 0, Type: &atype, FacilityName: &facilityName, FAA: &faa, ICAO: "KLAX", Region: &region, State: &state, County: &county, City: &city, Ownership: &ownership,
				Use: &use, Manager: &manager, ManagerPhone: &managerPhone, Latitude: &latitude, Longitude: &longitude, Status: "PENDING"},
			expectedErr: nil,
		},
		{
			name:        "Error ICAO is required",
			airport:     &dto.Airport{},
			expectedErr: fmt.Errorf("ICAO is required"),
		},
		{
			name:        "Error Ownership must be PU or PR",
			airport:     &dto.Airport{ICAO: "KAVL", Ownership: &falseOwnership},
			expectedErr: fmt.Errorf("Ownership must be PU or PR"),
		},
		{
			name:        "Error Use must be PU or PR",
			airport:     &dto.Airport{ICAO: "KAVL", Ownership: &ownership, Use: &falseUse},
			expectedErr: fmt.Errorf("Use must be PU or PR"),
		},
		{
			name:        "Error invalid latitude format",
			airport:     &dto.Airport{ICAO: "KAVL", Ownership: &ownership, Use: &use, Latitude: &falseLatitude},
			expectedErr: fmt.Errorf("Invalid latitude format (expected DD-MM-SS.sssN/S)"),
		},
		{
			name:        "Error invalid longitude format",
			airport:     &dto.Airport{ICAO: "KAVL", Ownership: &ownership, Use: &use, Latitude: &latitude, Longitude: &falseLongitude},
			expectedErr: fmt.Errorf("Invalid longitude format (expected DDD-MM-SS.sssE/W)"),
		},
		{
			name:        "Error invalid manager phone format",
			airport:     &dto.Airport{ICAO: "KAVL", Ownership: &ownership, Use: &use, Latitude: &latitude, Longitude: &longitude, ManagerPhone: &falseManagerPhone},
			expectedErr: fmt.Errorf("Invalid manager phone format"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := NewAirportValidator()
			got := av.Validate(tt.airport)
			if got != nil && got.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, got)
			}
		})
	}
}
