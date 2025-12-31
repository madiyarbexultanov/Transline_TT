# Shipment & Customer Service

Этот проект реализует систему управления заявками на доставку (Shipment-Service) и сервис клиентов (Customer-Service) с использованием gRPC и REST через Envoy.  

REST API доступно снаружи, а gRPC используется только внутри сети между сервисами. Все данные хранятся в PostgreSQL, а трассировка осуществляется через OpenTelemetry и Jaeger.

## Быстрый старт

1. Клонируем репозиторий:

```
git clone https://github.com/madiyarbexultanov/Transline_TT.git
cd Transline_TT
```

Поднимаем все сервисы через Docker Compose:

```
docker-compose up --build
```

Сервисы:

Envoy — REST на :8080, gRPC проксирование на :9090 внутри сети.
Shipment-Service — REST API.
Customer-Service — gRPC сервер.
Postgres — база данных.
Jaeger — для просмотра трассировки (:16686).

## REST API
Создание заявки (Shipment)

```
curl -X POST http://localhost:8080/api/v1/shipments \
  -H "Content-Type: application/json" \
  -d '{
    "route":"ALMATY→ASTANA",
    "price":120000,
    "customer":{"idn":"990101123456"}
  }'
```

Пример ответа
```
{
  "id": "shp-uuid",
  "status": "CREATED",
  "customerId": "cus-uuid"
}
```

Получение заявки по ID
```
curl http://localhost:8080/api/v1/shipments/<id>
```

```
{
  "id": "shp-uuid",
  "route": "ALMATY→ASTANA",
  "price": 120000,
  "status": "CREATED",
  "customerId": "cus-uuid",
  "created_at": "2025-12-31T13:30:27.826539Z"
}
```

## gRPC / Внутренние вызовы

Shipment-Service вызывает Customer-Service через gRPC (порт :9090)
gRPC порт не пробрасывается наружу
Customer-Service создаёт или получает клиента и возвращает customerId

## OpenTelemetry / Tracing

Все REST-запросы и внутренние gRPC вызовы трассируются.
Заходим в Jaeger UI: http://localhost:16686
Там видна цепочка:
```
REST → shipment-service → gRPC → customer-service → PostgreSQL
```

## Кратко о работе системы

Клиент делает REST-запрос через Envoy.
Envoy проксирует REST в Shipment-Service.
Shipment-Service через gRPC вызывает Customer-Service для создания или получения клиента.
Customer-Service сохраняет данные в PostgreSQL и возвращает ID.
Shipment-Service создаёт запись shipment с полученным customerId.
Ответ возвращается клиенту через Envoy.
Все шаги трассируются и видны в Jaeger.

gRPC порт (:9090) не доступен снаружи, только внутри Docker-сети.
Envoy применяет локальное ограничение 10 rps на IP.
Таблицы customers и shipments создаются автоматически через миграции при старте.
