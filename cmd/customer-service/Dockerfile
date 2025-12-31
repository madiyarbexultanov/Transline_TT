FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o customer-service ./cmd/customer-service

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/customer-service .
CMD ["./customer-service"]
