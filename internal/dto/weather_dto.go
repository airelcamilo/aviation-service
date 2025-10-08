package dto

type Weather struct {
	LastUpdated string  `json:"last_updated"`
	TempC       float64 `json:"temp_c"`
	IsDay       int     `json:"is_day"`
	FeelslikeC  float64 `json:"feelslike_c"`
	WindchillC  float64 `json:"windchill_c"`
	HeatindexC  float64 `json:"heatindex_c"`
	DewpointC   float64 `json:"dewpoint_c"`
	Condition   struct {
		Text string `json:"text"`
		Icon string `json:"icon"`
		Code int    `json:"code"`
	} `json:"condition"`
	WindKph    float64 `json:"wind_kph"`
	WindDegree int     `json:"wind_degree"`
	WindDir    string  `json:"wind_dir"`
	PressureMb float64 `json:"pressure_mb"`
	PrecipMm   float64 `json:"precip_mm"`
	Humidity   int     `json:"humidity"`
	Cloud      int     `json:"cloud"`
	VisKm      float64 `json:"vis_km"`
	UV         float64 `json:"uv"`
	GustKph    float64 `json:"gust_kph"`
}

type WeatherDataResponse struct {
	Current Weather `json:"current"`
}
