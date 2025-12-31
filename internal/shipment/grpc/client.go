package grpc

import (
	"context"
	"os"

	customerpb "shipment-customer-test/api/proto/customerpb"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client customerpb.CustomerServiceClient
	conn   *grpc.ClientConn
}

func NewClient(addr string) (*Client, error) {
    if addr == "" {
        addr = os.Getenv("CUSTOMER_GRPC_ADDR")
    }

    conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, err
    }

    c := customerpb.NewCustomerServiceClient(conn)
    return &Client{client: c, conn: conn}, nil
}


func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) UpsertCustomer(ctx context.Context, idn string) (*customerpb.CustomerResponse, error) {
	tr := otel.Tracer("shipment-service")
	ctx, span := tr.Start(ctx, "CustomerGRPC.UpsertCustomer")
	defer span.End()

	span.SetAttributes(attribute.String("customer.idn", idn))

	resp, err := c.client.UpsertCustomer(ctx, &customerpb.UpsertCustomerRequest{Idn: idn})
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return resp, nil
}
