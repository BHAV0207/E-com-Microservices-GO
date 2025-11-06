## payment-service

Validates user and order, processes payment, persists result, and publishes a `payment-events` message with outcome.

### Run

Exposed on port 3000 via docker-compose.

### Environment

- `PORT=3000`
- `KAFKA_BROKER=kafka:9092`
- `PAYMENT_TOPIC=payment-events`
- `MONGO_URI`

### REST API

- POST `/payment` — Create/Process payment
  - body: `{ "orderId": string, "userId": string, "reservationId": string, "amount": number, "method": string }`
  - writes a payment document and publishes Kafka event:
    - `status: success | failure`

### Kafka Production

- Topic: `payment-events`
- Event schema (simplified): `{ orderId, userId, reservationId, amount, method, status, timestamp }`
