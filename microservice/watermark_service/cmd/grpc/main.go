package main

import (
	"log"
	"net"

	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/adapter/grpc/controller"
	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/core/service"
	pb "github.com/brandoyts/watmarker/proto/watermark"
	"google.golang.org/grpc"
)

func main() {
	port := ":6000"

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	watermarkService := service.NewWatermarkService()

	watermarkController := controller.NewWatermarkController(watermarkService)

	pb.RegisterWatermarkServiceServer(server, watermarkController)

	log.Println("watermark grpc server is listening on port", listener.Addr())

	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
