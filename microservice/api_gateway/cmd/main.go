package main

import (
	"log"
	"os"

	"github.com/brandoyts/watmarker/microservice/api_gateway/config"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/adapter/grpc/watermark_grpc_client"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/adapter/http/controller"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/core/server"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/core/service"
	"github.com/brandoyts/watmarker/pkg/logger/v1"
)

func main() {
	// load gateway configuration
	gatewayConfig := config.LoadGatewayConfig()

	// initialize logger
	appLogger, err := logger.NewLogger(&logger.Config{LogLevel: "info"})
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
		os.Exit(1)
	}
	defer appLogger.Sync()

	// initialize watermark grpc client
	watermarkGrpcClient, err := watermark_grpc_client.New(":6000")
	if err != nil {
		log.Fatalf("failed connecting to watermark grpc client: %v", err)
	}

	// initialize watermark service
	watermarkService := service.NewWatermarkService(watermarkGrpcClient)

	// initialize watermark http controller
	watermarkController := controller.NewWatermarkController(watermarkService)

	srv := server.NewServer(gatewayConfig, appLogger)

	srv.Use(server.Log)

	// register health check route
	srv.RegisterHandler("/health", controller.HealthCheck)
	srv.RegisterHandler("/watermark", watermarkController.ServeHTTP)

	err = srv.Run()
	if err != nil {
		appLogger.Error("Failed to start api gateway: ", err)
	}
}
