package port

import "context"

type WatermarkService interface {
	ApplyWatermark(ctx context.Context, text string, size int32) (string, error)
}
