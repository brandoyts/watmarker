package controller

import (
	"context"

	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/core/port"
	pb "github.com/brandoyts/watmarker/proto/watermark"
)

type WatermarkController struct {
	pb.UnimplementedWatermarkServiceServer
	service port.WatermarkService
}

func NewWatermarkController(service port.WatermarkService) *WatermarkController {
	return &WatermarkController{
		UnimplementedWatermarkServiceServer: pb.UnimplementedWatermarkServiceServer{},
		service:                             service,
	}
}

func (wc *WatermarkController) ApplyWatermark(ctx context.Context, request *pb.ApplyWatermarkRequest) (*pb.ApplyWatermarkResponse, error) {
	res, err := wc.service.ApplyWatermark(ctx, request.Text, request.Size)
	if err != nil {
		return nil, err
	}

	return &pb.ApplyWatermarkResponse{ImageUrl: res}, nil
}
