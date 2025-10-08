package service

import (
	"aviation-service/internal/dto"
	"context"
	"sync"

	"go.uber.org/zap"
)

type AirportWeatherService struct {
	logger         *zap.SugaredLogger
	airportService IAirportService
	weatherService IWeatherService
}

func NewAirportWeatherService(logger *zap.SugaredLogger, airportService IAirportService, weatherService IWeatherService) *AirportWeatherService {
	return &AirportWeatherService{
		logger:         logger,
		airportService: airportService,
		weatherService: weatherService,
	}
}

func (s *AirportWeatherService) SearchAirportWeather(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.AirportWeather, error) {
	airports, err := s.airportService.SearchAirport(ctx, icao, facilityName, limit, offset)
	if err != nil {
		s.logger.Errorw("Failed to search airports", "error", err)
		return nil, err
	}

	numWorkers := 10
	if len(airports) < numWorkers {
		numWorkers = len(airports)
	}
	airportChannel := make(chan dto.Airport, len(airports))
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	var airportWeathers []dto.AirportWeather

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for airport := range airportChannel {
				airportWeather := dto.AirportWeather{Airport: airport}
				if airport.City != nil {
					weather, weatherErr := s.weatherService.GetWeather(ctx, *airport.City)
					if weatherErr != nil {
						s.logger.Errorw("Weather not available for airport", "error", weatherErr, "icao", airport.ICAO)
					} else {
						airportWeather.Weather = weather
					}
				}
				mutex.Lock()
				airportWeathers = append(airportWeathers, airportWeather)
				mutex.Unlock()
			}
		}()
	}

	go func() {
		defer close(airportChannel)
		for _, airport := range airports {
			airportChannel <- airport
		}
	}()

	wg.Wait()
	return airportWeathers, nil
}
