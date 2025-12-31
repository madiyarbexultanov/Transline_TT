package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	customerpb "shipment-customer-test/api/proto/customerpb"
	customgrpc "shipment-customer-test/internal/customer/grpc"
	"shipment-customer-test/internal/customer/repo"
	"shipment-customer-test/internal/customer/service"
)

func main() {
    pool := NewPostgresPool()
    defer pool.Close()

    customerRepo := repo.NewCustomerRepository(pool)
    customerService := service.New(customerRepo)
    customerServer := customgrpc.NewServer(customerService)

    grpcServer := grpc.NewServer()
    customerpb.RegisterCustomerServiceServer(grpcServer, customerServer)

    lis, err := net.Listen("tcp", ":9090")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    log.Println("Customer gRPC server running on :9090")

    go func() {
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("failed to serve: %v", err)
        }
    }()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
    <-stop

    log.Println("Shutting down gRPC server...")
    grpcServer.GracefulStop()
    log.Println("Server stopped")
}


func NewPostgresPool() *pgxpool.Pool {
    dbURL := os.Getenv("DATABASE_URL")
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
