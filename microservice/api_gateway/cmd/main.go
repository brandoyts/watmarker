package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brandoyts/watmarker/microservice/api_gateway/config"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/adapter/cache"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/adapter/grpc/watermark_grpc_client"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/adapter/http"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/adapter/http/controller"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/adapter/http/middleware"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/service"
	"github.com/brandoyts/watmarker/pkg/logger/v1"
	"github.com/joho/godotenv"
)

func init() {
	_, exits := os.LookupEnv("APP_ENV")
	if exits {
		return
	}

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
		log.Fatal(errors.New("missing .env file"))
	}
}

func main() {
	// load gateway configuration
	gatewayConfig := config.LoadGatewayConfig("../config/")

	// Setup logger
	appLoger, err := logger.New()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer appLoger.Sync()

	// initialize cache provider
	redisCache, err := cache.NewRedisClient(cache.RedisClientConfig{
		Addr:     os.Getenv("REDIS_ADDRESS"), // address + port
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	if err != nil {
		appLoger.Fatalw("failed to establish a connection to redis", err)
	}

	// initialize watermark grpc client
	watermarkGrpcClient, err := watermark_grpc_client.New(gatewayConfig.Services[0].Url)
	if err != nil {
		appLoger.Fatalf("failed connecting to watermark grpc client:", err)
	}

	// initialize watermark service
	watermarkService := service.NewWatermarkService(watermarkGrpcClient)

	// initialize watermark http controller
	watermarkController := controller.NewWatermarkController(watermarkService)

	srv := http.NewServer(gatewayConfig.Address)

	// middleware
	srv.Use(middleware.Log(appLoger))
	srv.Use(middleware.RateLimit(redisCache, 10, time.Minute*5))

	// register health check route
	srv.RegisterHandler("/health", controller.HealthCheck)
	srv.RegisterHandler("/watermark", watermarkController.ServeHTTP)

	go func() {
		appLoger.Info("Starting API Gateway server on ", gatewayConfig.Address)
		err := srv.Run()
		if err != nil {
			appLoger.Error("Server failed to start: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	appLoger.Info("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		appLoger.Error(err)
		os.Exit(1)
	}

	appLoger.Info("Server exited successfully.")
}
