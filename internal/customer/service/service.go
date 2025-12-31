package service

import (
	"context"
	"errors"
	"time"

	"shipment-customer-test/internal/customer/repo"

	"github.com/google/uuid"
)

type Customer struct {
    ID        uuid.UUID
    IDN       string
    CreatedAt time.Time
}

type CustomerService struct {
    repo *repo.CustomerRepository
}

func New(repo *repo.CustomerRepository) *CustomerService {
    return &CustomerService{repo: repo}
}

func (s *CustomerService) UpsertCustomer(ctx context.Context, idn string) (*Customer, error) {
    if len(idn) != 12 {
        return nil, errors.New("invalid idn")
    }

    r, err := s.repo.Upsert(ctx, idn)
    if err != nil {
        return nil, err
    }

    return &Customer{
        ID:        r.ID,
        IDN:       r.IDN,
        CreatedAt: r.CreatedAt,
    }, nil
}

func (s *CustomerService) GetCustomer(ctx context.Context, idn string) (*Customer, error) {
	if len(idn) != 12 {
		return nil, errors.New("invalid idn")
	}

	r, err := s.repo.GetByIDN(ctx, idn)
	if err != nil {
		return nil, err
	}

	return &Customer{
		ID:        r.ID,
		IDN:       r.IDN,
		CreatedAt: r.CreatedAt,
	}, nil
}
