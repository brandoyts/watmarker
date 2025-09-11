package service

import (
	"context"

	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/core/port"
)

type WatermarkService struct {
	imageStorage port.ImageStorage
}

func NewWatermarkService(imageStorage port.ImageStorage) *WatermarkService {
	return &WatermarkService{
		imageStorage: imageStorage,
	}
}

func (ws *WatermarkService) ApplyWatermark(ctx context.Context, text string, size int32) (string, error) {
	return "http://test-link", nil
}
