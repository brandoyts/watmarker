package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brandoyts/watmarker/microservice/api_gateway/config"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/adapter/grpc/watermark_grpc_client"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/adapter/http/controller"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/core/server"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/core/service"
	"github.com/brandoyts/watmarker/pkg/logger/v1"
)

func main() {
	// load gateway configuration
	gatewayConfig := config.LoadGatewayConfig("../../config")

	// initialize logger
	appLogger, err := logger.NewLogger(&logger.Config{LogLevel: "info"})
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
		os.Exit(1)
	}
	defer appLogger.Sync()

	// initialize watermark grpc client
	watermarkGrpcClient, err := watermark_grpc_client.New(gatewayConfig.Services[0].Url)
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

	go func() {
		appLogger.Info("Starting API Gateway server on ", gatewayConfig.Address)
		err := srv.Run()
		if err != nil {
			appLogger.Error("Server failed to start: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	appLogger.Info("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		appLogger.Error(err)
		os.Exit(1)
	}

	appLogger.Info("Server exited successfully.")
}
