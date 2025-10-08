package dto

type AirportWeather struct {
	Airport Airport  `json:"airport"`
	Weather *Weather `json:"weather,omitempty"`
}

