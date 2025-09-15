package cache

import (
	"context"
	"testing"
	"time"

	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/mock"
	"github.com/golang/mock/gomock"
)

func TestRedisCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock.NewMockRateLimit(ctrl)

	mockCache.EXPECT().Increment(gomock.Any(), "counter").Return(int64(1), nil).Times(1)

	mockCache.EXPECT().Expire(gomock.Any(), "counter", time.Second*5).Return(nil).Times(1)

	mockCache.EXPECT().Ping(gomock.Any()).Return(nil).Times(1)

	ctx := context.Background()

	val, err := mockCache.Increment(ctx, "counter")
	if err != nil || val != 1 {
		t.Fatalf("expected 1, got %d, err: %v", val, err)
	}

	if err := mockCache.Expire(ctx, "counter", 5*time.Second); err != nil {
		t.Fatal(err)
	}

	if err := mockCache.Ping(ctx); err != nil {
		t.Fatal(err)
	}
}
