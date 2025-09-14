package port

import (
	"context"

	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/core/model/request"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/core/model/response"
)

type WatermarkService interface {
	ApplyWatermark(ctx context.Context, req request.ApplyWatermarkRequest) (response.ApplyWatermarkResponse, error)
}

type WatermarkClient interface {
	ApplyWatermark(ctx context.Context, text string, size []byte) (string, error)
}
