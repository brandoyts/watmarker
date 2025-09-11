package watermark_grpc_client

import (
	"context"
	"testing"

	pb "github.com/brandoyts/watmarker/proto/watermark"
	"google.golang.org/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock gRPC Client ---

type MockWatermarkClient struct {
	mock.Mock
	pb.WatermarkServiceClient
}

func (m *MockWatermarkClient) ApplyWatermark(ctx context.Context, in *pb.ApplyWatermarkRequest, opts ...grpc.CallOption) (*pb.ApplyWatermarkResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.ApplyWatermarkResponse), args.Error(1)
}

// --- Tests ---

func TestApplyWatermark_Success(t *testing.T) {
	mockClient := new(MockWatermarkClient)
	adapter := &Client{client: mockClient}

	expectedReq := &pb.ApplyWatermarkRequest{
		Text: "Hello",
		Size: 42,
	}
	expectedResp := &pb.ApplyWatermarkResponse{
		ImageUrl: "http://test.com22",
	}

	mockClient.On("ApplyWatermark", mock.Anything, expectedReq).Return(expectedResp, nil)

	result, err := adapter.ApplyWatermark(context.Background(), "Hello", 42)

	assert.NoError(t, err)
	assert.Equal(t, "http://test.com22", result)
	mockClient.AssertExpectations(t)
}

func TestApplyWatermark_Error(t *testing.T) {
	mockClient := new(MockWatermarkClient)
	adapter := &Client{client: mockClient}

	expectedReq := &pb.ApplyWatermarkRequest{
		Text: "Fail",
		Size: 10,
	}
	expectedErr := assert.AnError

	mockClient.On("ApplyWatermark", mock.Anything, expectedReq).Return(&pb.ApplyWatermarkResponse{}, expectedErr)

	result, err := adapter.ApplyWatermark(context.Background(), "Fail", 10)

	assert.Error(t, err)
	assert.Empty(t, result)
	mockClient.AssertExpectations(t)
}
