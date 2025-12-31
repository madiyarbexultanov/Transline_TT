package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Shipment struct {
    ID         uuid.UUID
    Route      string
    Price      float64
    Status     string
    CustomerID uuid.UUID
    CreatedAt  time.Time
}

type ShipmentRepository struct {
    pool *pgxpool.Pool
}

func NewShipmentRepository(pool *pgxpool.Pool) *ShipmentRepository {
    return &ShipmentRepository{pool: pool}
}



func (r *ShipmentRepository) Create(ctx context.Context, s *Shipment) error {
	tr := otel.Tracer("shipment-service")         // имя tracer'а можно выбрать для сервиса
	ctx, span := tr.Start(ctx, "ShipmentRepository.Create") // создаём span
	defer span.End()

	span.SetAttributes(
		attribute.String("shipment.route", s.Route),
		attribute.Float64("shipment.price", s.Price),
	)

	err := r.pool.QueryRow(
		ctx,
		`INSERT INTO shipments (route, price, status, customer_id) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id, status, created_at`,
		s.Route, s.Price, s.Status, s.CustomerID,
	).Scan(&s.ID, &s.Status, &s.CreatedAt)

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *ShipmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*Shipment, error) {
	tr := otel.Tracer("shipment-service")
	ctx, span := tr.Start(ctx, "ShipmentRepository.GetByID")
	defer span.End()

	s := &Shipment{}
	err := r.pool.QueryRow(
		ctx,
		`SELECT id, route, price, status, customer_id, created_at
		 FROM shipments
		 WHERE id = $1`,
		id,
	).Scan(&s.ID, &s.Route, &s.Price, &s.Status, &s.CustomerID, &s.CreatedAt)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return s, nil
}
