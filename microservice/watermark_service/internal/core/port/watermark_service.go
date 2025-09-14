package port

import "context"

//go:generate mockgen -package mock -source=../port/watermark_service.go -destination=../../mock/watermark_service_mock.go
type WatermarkService interface {
	ProcessImage(ctx context.Context, in ProcessImageInput) (string, error)
}

type ProcessImageInput struct {
	ImageData         []byte
	WatermarkText     string
	WatermarkFontSize int16
}
