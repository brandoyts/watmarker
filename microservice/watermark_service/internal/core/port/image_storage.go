package port

import "context"

type ImageStorage interface {
	Upload(ctx context.Context, fileContent []byte) error
}
