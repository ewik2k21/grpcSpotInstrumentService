package handlers

import (
	"context"
	"github.com/ewik2k21/grpcSpotInstrumentService/internal/services"
	spotInstrument "github.com/ewik2k21/grpcSpotInstrumentService/pkg/spot_instrument_v1"
)

type SpotInstrumentHandler struct {
	spotInstrument.UnimplementedSpotInstrumentServiceServer
	service services.SpotInstrumentService
}

func NewSpotInstrumentHandler(service services.SpotInstrumentService) *SpotInstrumentHandler {
	return &SpotInstrumentHandler{service: service}
}

func (h *SpotInstrumentHandler) ViewMarkets(ctx context.Context, req *spotInstrument.ViewMarketsRequest) (res *spotInstrument.ViewMarketsResponse, err error) {
	userRole := req.GetUserRole()

	markets, err := h.service.GetAllMarkets(userRole.String())
	if err != nil {
		return nil, err
	}

	return &spotInstrument.ViewMarketsResponse{
		Markets: markets,
	}, nil
}
