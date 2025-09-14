package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/adapter/grpc/controller"
	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/adapter/storage/s3"
	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/config"
	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/core/service"
	pb "github.com/brandoyts/watmarker/proto/watermark"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func init() {
	_, exits := os.LookupEnv("APP_ENV")
	if exits {
		return
	}

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(errors.New("missing .env file"))
	}
}

func main() {
	// load configuration
	cnf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v\n", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", cnf.AppPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	// initialize s3
	s3Adapter := s3.New(s3.Configuration{
		Bucket:          cnf.AwsBucket,
		Region:          cnf.AwsRegion,
		AccessKeyId:     cnf.AwsAccessKeyId,
		SecretAccessKey: cnf.AwsSecretAccessKey,
		BaseEndpoint:    cnf.LeapcellBaseEndpoint,
	})

	watermarkService := service.NewWatermarkService(s3Adapter, cnf.LeapcellCdn)

	watermarkController := controller.NewWatermarkController(watermarkService)

	pb.RegisterWatermarkServiceServer(server, watermarkController)

	log.Println("watermark grpc server is listening on", listener.Addr())

	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
