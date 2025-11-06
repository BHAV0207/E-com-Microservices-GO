## order-service

Coordinates checkout. Validates user, cart, product pricing, and inventory; creates orders; reserves inventory; invokes payment; updates order status on payment events.

### Run

Exposed on port 7000 via docker-compose.

### REST API

- POST `/order` — Create an order
  - body: `{ "userId": string, "cartId": string, "address": string, "paymentInfo": string }`
  - flow:
    1. Validate user, cart items, product prices
    2. Validate inventory for each item
    3. Reserve inventory via `inventory-service` → `reservationId`
    4. Create order with status `pending`
    5. Call `payment-service` with `{ orderId, userId, amount, reservationId }`
- GET `/order/{id}` — Get order by id
- GET `/orders/user/{id}` — Get all orders for a user

### Kafka Consumption

- Topic: `payment-events`
- Updates order status to `confirmed` on success; `failed` on failure.
