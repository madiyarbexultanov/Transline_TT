package http

import (
	"errors"
	"log"
	"net/http"

	"shipment-customer-test/internal/shipment/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

type CreateShipmentRequest struct {
	Route    string `json:"route"`
	Price    float64 `json:"price"`
	Customer struct {
		IDN string `json:"idn"`
	} `json:"customer"`
}

type Handler struct {
	service *service.ShipmentService
}

func NewHandler(s *service.ShipmentService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateShipment(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shipment, err := h.service.CreateShipment(ctx, req.Route, req.Price, req.Customer.IDN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         shipment.ID.String(),
		"status":     shipment.Status,
		"customerId": shipment.CustomerID.String(),
	})
}

func (h *Handler) GetShipment(c *gin.Context) {
    ctx := c.Request.Context()
    span := trace.SpanFromContext(ctx) 
    log.Printf("GET shipment, trace_id=%s", span.SpanContext().TraceID())


    idParam := c.Param("id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
        return
    }

    shipment, err := h.service.GetShipment(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, shipment)
}

