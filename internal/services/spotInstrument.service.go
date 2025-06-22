package services

import (
	"github.com/ewik2k21/grpcSpotInstrumentService/internal/mappers"
	"github.com/ewik2k21/grpcSpotInstrumentService/internal/repositories"
	spot "github.com/ewik2k21/grpcSpotInstrumentService/pkg/spot_instrument_v1"
	"log/slog"
)

type SpotInstrumentService struct {
	repo   repositories.SpotInstrumentRepository
	logger *slog.Logger
}

func NewSpotInstrumentService(repo repositories.SpotInstrumentRepository, logger *slog.Logger) *SpotInstrumentService {
	return &SpotInstrumentService{
		repo:   repo,
		logger: logger,
	}
}

func (s *SpotInstrumentService) GetAllMarkets(userRole string) ([]*spot.Market, error) {
	markets, err := s.repo.GetAllMarkets()
	if err != nil {
		return nil, err
	}
	var res []*spot.Market

	for _, market := range markets {
		if !market.Enabled || market.DeletedAt != nil {
			continue
		}
		protoMarket := mappers.MapMarketToProto(market)

		res = append(res, protoMarket)

	}
	return res, nil
}
