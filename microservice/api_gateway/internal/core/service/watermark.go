package service

import (
	"context"

	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/core/model/request"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/core/model/response"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/core/port"
)

type WatermarkService struct {
	client port.WatermarkClient
}

func NewWatermarkService(client port.WatermarkClient) *WatermarkService {
	return &WatermarkService{
		client: client,
	}
}

func (ws *WatermarkService) ApplyWatermark(ctx context.Context, in request.ApplyWatermarkRequest) (response.ApplyWatermarkResponse, error) {

	res, err := ws.client.ApplyWatermark(context.Background(), in.Text, in.Size)
	if err != nil {
		return response.ApplyWatermarkResponse{}, err
	}

	return response.ApplyWatermarkResponse{
		ImageUrl: res,
	}, nil
}
