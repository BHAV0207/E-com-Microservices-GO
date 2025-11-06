## inventory-service

Provides inventory CRUD and reservation workflow. Consumes `payment-events` from Kafka to commit or cancel reservations.

### Run

Exposed on port 6000 via docker-compose.

### Environment

- `PORT=6000`
- `KAFKA_BROKER=kafka:9092`
- `PAYMENT_TOPIC=payment-events`
- `MONGO_URI`

### REST API

- GET `/inventory/{productId}` — Get inventory for product
- POST `/inventory` — Create inventory document
- PATCH `/inventory/{id}` — Update inventory fields
- DELETE `/inventory/{id}` — Delete inventory document
- POST `/inventory/reserve` — Reserve stock for an order
  - body: `{ "orderId": string, "items": [{"productId": string, "quantity": number}] }`
  - returns: `{ reservationId, expiresAt }`

### Kafka Consumption

- Topic: `payment-events`
- On `success`: decrements inventory and marks reservation `COMMITTED`.
- On `failed`: marks reservation `CANCELLED`.

### Data

- Collection: `inventory`
- Collection: `reservations` (reservation state with TTL semantics enforced by application logic)
