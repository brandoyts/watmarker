package port

import "context"

//go:generate mockgen -package mock -source=../port/image_storage.go -destination=../../mock/image_storage_mock.go
type ImageStorage interface {
	Upload(ctx context.Context, filename string, imageData []byte) error
}
