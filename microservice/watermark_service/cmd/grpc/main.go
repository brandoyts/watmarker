package main

import (
	"log"
	"net"
	"os"

	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/adapter/grpc/controller"
	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/adapter/storage/s3"
	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/config"
	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/core/service"
	pb "github.com/brandoyts/watmarker/proto/watermark"
	"google.golang.org/grpc"
)

func main() {
	cnf, err := config.LoadConfig("../../")
	if err != nil {
		log.Fatalf("failed to read configuration file: %v", err)
		os.Exit(1)
	}

	listener, err := net.Listen("tcp", cnf.AppUrl)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		os.Exit(1)
	}

	server := grpc.NewServer()

	// initialize s3
	s3Adapter := s3.New(s3.Configuration{
		Bucket:          cnf.AwsBucket,
		Region:          cnf.AwsRegion,
		AccessKeyId:     cnf.AwsAccessKeyId,
		SecretAccessKey: cnf.AwsSecretAccessKey,
		BaseEndpoint:    cnf.AwsEndpoint,
	})

	watermarkService := service.NewWatermarkService(s3Adapter)

	watermarkController := controller.NewWatermarkController(watermarkService)

	pb.RegisterWatermarkServiceServer(server, watermarkController)

	log.Println("watermark grpc server is listening on port", listener.Addr())

	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
