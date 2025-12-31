package main

import (
	"context"
	"log"
	"os"
	"time"

	"shipment-customer-test/internal/shipment/grpc"
	"shipment-customer-test/internal/shipment/http"
	"shipment-customer-test/internal/shipment/http/middleware"
	"shipment-customer-test/internal/shipment/repo"
	"shipment-customer-test/internal/shipment/service"
	"shipment-customer-test/internal/tracing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	shutdown, err := tracing.Init("shipment-service")
	if err != nil {
		log.Fatalf("failed to init tracing: %v", err)
	}
	defer shutdown(context.Background())

	pool := NewPostgresPool()
	defer pool.Close()

	shipmentRepo := repo.NewShipmentRepository(pool)

	addr := os.Getenv("CUSTOMER_GRPC_ADDR")
	grpcClient, err := grpc.NewClient(addr)
	if err != nil {
		log.Fatalf("failed to connect to customer-service: %v", err)
	}


	shipmentService := service.NewShipmentService(
		shipmentRepo,
		grpcClient,
	)
	handler := http.NewHandler(shipmentService)


	r := gin.New()

	r.Use(
		gin.Recovery(),
		gin.Logger(),
		middleware.Tracing("shipment-service"),
	)

	api := r.Group("/api/v1")
	api.POST("/shipments", handler.CreateShipment)
	api.GET("/shipments/:id", handler.GetShipment)

	log.Println("Shipment REST server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}


func NewPostgresPool() *pgxpool.Pool {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL is not set")
    }
    config, err := pgxpool.ParseConfig(dbURL)
    if err != nil {
        log.Fatalf("failed to parse database config: %v", err)
    }

    config.MaxConns = 10
    config.MinConns = 2
    config.MaxConnLifetime = time.Hour

    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        log.Fatalf("failed to create pgxpool: %v", err)
    }

    return pool
}

