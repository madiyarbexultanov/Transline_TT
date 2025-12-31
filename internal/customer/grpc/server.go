package grpc

import (
	"context"
	"time"

	customerpb "shipment-customer-test/api/proto/customerpb"
	"shipment-customer-test/internal/customer/service"

	"go.opentelemetry.io/otel"
)

type Server struct {
    customerpb.UnimplementedCustomerServiceServer
    svc *service.CustomerService
}

func NewServer(svc *service.CustomerService) *Server {
    return &Server{svc: svc}
}

func (s *Server) UpsertCustomer(
	ctx context.Context,
	req *customerpb.UpsertCustomerRequest,
) (*customerpb.CustomerResponse, error) {
	tr := otel.Tracer("customer-service")
	ctx, span := tr.Start(ctx, "UpsertCustomer")
	defer span.End()

	customer, err := s.svc.UpsertCustomer(ctx, req.Idn)
	if err != nil {
		return nil, err
	}

	return &customerpb.CustomerResponse{
		Id:        customer.ID.String(),
		Idn:       customer.IDN,
		CreatedAt: customer.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Server) GetCustomer(
	ctx context.Context,
	req *customerpb.GetCustomerRequest,
) (*customerpb.CustomerResponse, error) {
	tr := otel.Tracer("customer-service")
	ctx, span := tr.Start(ctx, "GetCustomer")
	defer span.End()

	customer, err := s.svc.GetCustomer(ctx, req.Idn)
	if err != nil {
		return nil, err
	}

	return &customerpb.CustomerResponse{
		Id:        customer.ID.String(),
		Idn:       customer.IDN,
		CreatedAt: customer.CreatedAt.Format(time.RFC3339),
	}, nil
}
