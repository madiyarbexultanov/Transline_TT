package service

import (
	"context"

	"shipment-customer-test/internal/shipment/grpc"
	"shipment-customer-test/internal/shipment/repo"

	"github.com/google/uuid"
)

type ShipmentService struct {
	repo       *repo.ShipmentRepository
	grpcClient *grpc.Client
}

func NewShipmentService(r *repo.ShipmentRepository, g *grpc.Client) *ShipmentService {
	return &ShipmentService{repo: r, grpcClient: g}
}

func (s *ShipmentService) CreateShipment(ctx context.Context, route string, price float64, idn string) (*repo.Shipment, error) {
	customer, err := s.grpcClient.UpsertCustomer(ctx, idn)
	if err != nil {
		return nil, err
	}

	shipment := &repo.Shipment{
		Route:      route,
		Price:      price,
		CustomerID: uuid.MustParse(customer.Id),
		Status:     "CREATED",
	}
	
	if err := s.repo.Create(ctx, shipment); err != nil {
		return nil, err
	}

	return shipment, nil
}

func (s *ShipmentService) GetShipment(ctx context.Context, id uuid.UUID) (*repo.Shipment, error) {
	return s.repo.GetByID(ctx, id)
}