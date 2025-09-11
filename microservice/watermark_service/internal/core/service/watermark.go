package service

import "context"

type WatermarkService struct {
}

func NewWatermarkService() *WatermarkService {
	return &WatermarkService{}
}

func (ws *WatermarkService) ApplyWatermark(ctx context.Context, text string, size int32) (string, error) {
	return "http://test-link", nil
}
