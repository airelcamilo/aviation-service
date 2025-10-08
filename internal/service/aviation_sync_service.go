package service

import (
	"aviation-service/internal/dto"
	"aviation-service/internal/repository"
	"context"
	"strings"
	"sync"

	"go.uber.org/zap"
)

type AviationSyncService struct {
	logger         *zap.SugaredLogger
	airportRepo    repository.IAirportRepository
	airportService IAirportService
}

type SyncStats struct {
	mu      sync.Mutex
	success int
	failed  int
	err     int
}

func NewAviationSyncService(logger *zap.SugaredLogger, airportRepo repository.IAirportRepository, airportService IAirportService) *AviationSyncService {
	return &AviationSyncService{
		logger:         logger,
		airportRepo:    airportRepo,
		airportService: airportService,
	}
}

func (s *AviationSyncService) Sync(ctx context.Context) (*dto.SyncResponse, error) {
	var syncStats SyncStats
	airports, err := s.airportRepo.GetAllPending(ctx)
	if err != nil {
		s.logger.Errorw("Failed get all pending airports", "error", err)
		return nil, err
	}

	if len(airports) == 0 {
		return nil, nil
	}
	s.logger.Infow("Syncing pending airports", "count", len(airports))

	batchLength := 30
	numWorkers := 10
	icaoChannel := make(chan string, len(airports)/batchLength)
	wg := sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go s.getAirportData(ctx, &wg, &syncStats, icaoChannel)
	}

	go func() {
		batch := []string{}
		for _, apt := range airports {
			batch = append(batch, apt.ICAO)

			if len(batch) == batchLength {
				batchString := strings.Join(batch, ",")
				icaoChannel <- batchString

				batch = []string{}
			}
		}

		if len(batch) > 0 {
			batchString := strings.Join(batch, ",")
			icaoChannel <- batchString
		}
		close(icaoChannel)
	}()
	wg.Wait()

	syncResponse := dto.SyncResponse{
		Total:   len(airports),
		Success: syncStats.success,
		Failed:  syncStats.failed,
		Error:   syncStats.err,
	}
	return &syncResponse, nil
}

func (s *AviationSyncService) getAirportData(ctx context.Context, wg *sync.WaitGroup, syncStats *SyncStats, icaoChannel chan string) {
	defer wg.Done()

	for icaos := range icaoChannel {
		airports, err := s.airportService.FetchAirportData(icaos)
		if err != nil {
			s.logger.Errorw("Failed to fetch airports from API", "error", err)
			syncStats.mu.Lock()
			syncStats.err += 1
			syncStats.mu.Unlock()
			continue
		}

		err = s.updateAirportData(ctx, syncStats, airports)
		if err != nil {
			s.logger.Errorw("Failed to update airports from API", "error", err)
			syncStats.mu.Lock()
			syncStats.err += 1
			syncStats.mu.Unlock()
			continue
		}
	}
}

func (s *AviationSyncService) updateAirportData(ctx context.Context, syncStats *SyncStats, airports *dto.AirportDataResponse) error {
	s.logger.Infow("Updating airports data", "count", len(*airports))
	var success, failed int
	var toUpdate []dto.Airport

	for icao, apt := range *airports {
		if len(apt) == 0 {
			toUpdate = append(toUpdate, dto.Airport{
				ICAO:   icao,
				Status: "FAILED",
			})
			failed++
			continue
		}

		apt[0].Status = "DONE"
		toUpdate = append(toUpdate, apt[0])
		success++
	}

	if len(toUpdate) > 0 {
		if err := s.airportRepo.UpdateByICAO(ctx, toUpdate); err != nil {
			s.logger.Errorw("Batch update failed", "error", err)
			return err
		}
	}
	syncStats.mu.Lock()
	syncStats.success += success
	syncStats.failed += failed
	syncStats.mu.Unlock()
	return nil
}
