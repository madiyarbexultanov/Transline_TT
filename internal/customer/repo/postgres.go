package repo

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/google/uuid"
)

type Customer struct {
    ID        uuid.UUID
    IDN       string
    CreatedAt time.Time
}

type CustomerRepository struct {
    pool *pgxpool.Pool
}

func NewCustomerRepository(pool *pgxpool.Pool) *CustomerRepository {
    return &CustomerRepository{pool: pool}
}

func (r *CustomerRepository) Upsert(ctx context.Context, idn string) (*Customer, error) {
    customer := &Customer{}
    err := r.pool.QueryRow(
        ctx,
        `INSERT INTO customers (idn) VALUES ($1)
         ON CONFLICT (idn) DO UPDATE SET idn = EXCLUDED.idn
         RETURNING id, idn, created_at`,
        idn,
    ).Scan(&customer.ID, &customer.IDN, &customer.CreatedAt)

    if err != nil {
        return nil, err
    }
    return customer, nil
}

func (r *CustomerRepository) GetByIDN(ctx context.Context, idn string) (*Customer, error) {
	customer := &Customer{}
	err := r.pool.QueryRow(
		ctx,
		`SELECT id, idn, created_at FROM customers WHERE idn = $1`,
		idn,
	).Scan(&customer.ID, &customer.IDN, &customer.CreatedAt)

	if err != nil {
		return nil, err
	}
	return customer, nil
}
