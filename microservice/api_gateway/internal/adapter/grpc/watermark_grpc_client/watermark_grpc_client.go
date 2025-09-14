package watermark_grpc_client

import (
	"context"

	pb "github.com/brandoyts/watmarker/proto/watermark"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client pb.WatermarkServiceClient
}

func New(addr string) (*Client, error) {
	// insecure for local development
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{client: pb.NewWatermarkServiceClient(conn)}, nil
}

func (a *Client) ApplyWatermark(ctx context.Context, text string, chunks []byte) (string, error) {
	resp, err := a.client.ApplyWatermark(ctx, &pb.ApplyWatermarkRequest{
		WatermarkText: text,
		ImageData:     chunks,
	})
	if err != nil {
		return "", err
	}
	return resp.ImageUrl, nil
}
