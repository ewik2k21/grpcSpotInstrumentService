package repositories

import (
	"fmt"
	"github.com/ewik2k21/grpcSpotInstrumentService/internal/models"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type ISpotInstrumentRepository interface {
	GetAllMarkets() (map[string]models.Market, error)
}

type SpotInstrumentRepository struct {
	markets map[string]models.Market
	logger  *slog.Logger
}

const (
	UUID1 = "758268a7-8412-457a-971a-877189971093"
	UUID2 = "2a9e394b-304e-46a5-808c-61d798bbd714"
	UUID3 = "9c5d898f-56c7-49d7-a69b-2c73b711c4f7"
	UUID4 = "a8f91932-a0c6-4001-92b5-066a5d585250"
	UUID5 = "b2135e53-1d84-4386-b357-940861698b28"
)

func NewSpotInstrumentRepository(logger *slog.Logger) *SpotInstrumentRepository {
	markets := make(map[string]models.Market)
	uuid1, _ := uuid.Parse(UUID1)
	uuid2, _ := uuid.Parse(UUID2)
	uuid3, _ := uuid.Parse(UUID3)
	uuid4, _ := uuid.Parse(UUID4)
	uuid5, _ := uuid.Parse(UUID5)

	markets["1"] = models.Market{
		ID:        uuid1,
		Name:      "1",
		Enabled:   true,
		DeletedAt: nil,
	}

	markets["2"] = models.Market{
		ID:        uuid2,
		Name:      "2",
		Enabled:   true,
		DeletedAt: nil,
	}

	markets["3"] = models.Market{
		ID:        uuid3,
		Name:      "3",
		Enabled:   false,
		DeletedAt: nil,
	}

	del := time.Now()
	markets["4"] = models.Market{
		ID:        uuid4,
		Name:      "4",
		Enabled:   true,
		DeletedAt: &del,
	}

	markets["5"] = models.Market{
		ID:        uuid5,
		Name:      "5",
		Enabled:   true,
		DeletedAt: nil,
	}

	return &SpotInstrumentRepository{
		markets: markets,
		logger:  logger,
	}
}

func (r *SpotInstrumentRepository) GetAllMarkets() (map[string]models.Market, error) {
	if len(r.markets) == 0 {
		r.logger.Error("no markets in memory")
		return nil, fmt.Errorf("zero markets")
	}
	return r.markets, nil
}
