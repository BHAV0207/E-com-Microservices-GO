## E-com-Microservices-GO

A Go-based, event-driven e-commerce system built with microservices. Services communicate via REST and Kafka, with MongoDB as the primary datastore. The stack is containerized via Docker and orchestrated with docker-compose.

### Services

- **user-service (8080)**: User registration, login, CRUD. Publishes `user-created` and `user-deleted` events.
- **product-service (4000)**: Product CRUD. Triggers inventory creation/deletion via internal worker pool and inventory service calls.
- **inventory-service (6000)**: Inventory CRUD, reservation API, and Kafka consumer for payment events to commit/cancel reservations.
- **cart-service (9000)**: Add/get/delete cart. Consumes user events to create/delete carts.
- **order-service (7000)**: Validates user/cart/product/inventory, creates orders, reserves stock via inventory, and invokes payment.
- **payment-service (3000)**: Validates user/order, processes payments, publishes payment outcome events.
- **notification-service (2000)**: Consumes generic events and stores/sends user notifications.
- **Kafka + Zookeeper**: Inter-service event bus.

### High-Level Flow

1. User registers via `user-service` → publishes `user-created`.
2. `cart-service` consumes `user-created` → creates an empty cart.
3. User browses products from `product-service` and adds to cart via `cart-service`.
4. User places order via `order-service`:
   - Validates user, cart, product, and inventory.
   - Calls `inventory-service` to reserve stock → returns `reservationId`.
   - Persists order with status `pending`.
   - Calls `payment-service` with `orderId`, `userId`, `amount`, `reservationId`.
5. `payment-service` processes payment and publishes `payment-events` with status `success`/`failure`.
6. `inventory-service` consumes `payment-events`:
   - On `success`: commits inventory (decrement stock) and marks reservation `COMMITTED`.
   - On `failure`: cancels reservation (`CANCELLED`).
7. `order-service` consumes payment events (updates order status to `confirmed`/`failed`).
8. `notification-service` consumes events and persists/sends user notifications.

### Kafka Topics

- `user-created`, `user-deleted` (created in compose via KAFKA_CREATE_TOPICS)
- `payment-events` (emitted by `payment-service`, consumed by `inventory-service` and `order-service`)

### Ports

- 8080 user-service
- 4000 product-service
- 6000 inventory-service
- 9000 cart-service
- 7000 order-service
- 3000 payment-service
- 2000 notification-service
- 9092 Kafka, 2181 Zookeeper

### Running Locally

1. Prerequisites: Docker, docker-compose.
2. From the project root:

```bash
docker-compose up -d --build
```

3. Healthchecks (as defined in compose):
   - user: GET http://localhost:8080/register (route available for probe)
   - product: GET http://localhost:8080/get (probe target, service runs on 4000)
   - inventory: GET http://localhost:6000/get/health
   - cart: GET http://localhost:5000/cart (probe target, service runs on 9000)
   - order: GET http://localhost:7000/order
   - payment: GET http://localhost:3000/payment
   - notification: GET http://localhost:2000/notifications

Note: Some healthcheck endpoints are placeholders; refer to each service README for actual routes.

### Environment

Provide a `.env` in the project root with the necessary variables (Mongo URIs, Kafka broker, JWT secret, etc.). Services also use PORT from compose. Example (adjust as needed):

```bash
MONGO_URI=mongodb://mongo:27017
KAFKA_BROKER=kafka:9092
PAYMENT_TOPIC=payment-events
JWT_SECRET=your_secret_key
```

### Data Stores

- MongoDB collections per service (connection details configured within each service).
- `inventory-service` additional `reservations` collection for soft reservations.

### Service READMEs

- user-service: see `user-service/README.md`
- product-service: see `product-service/README.md`
- inventory-service: see `inventory-service/README.md`
- cart-service: see `cart-service/README.md`
- order-service: see `order-service/README.md`
- payment-service: see `Payment-Service/README.md`
- notification-service: see `notification-service/README.md`

### Request Samples

- Create User: POST http://localhost:8080/register
- Login: POST http://localhost:8080/login
- Create Product: POST http://localhost:4000/product
- Add to Cart: POST http://localhost:9000/cart
- Create Order: POST http://localhost:7000/order
- Create Reservation: POST http://localhost:6000/inventory/reserve
- Create Payment: POST http://localhost:3000/payment

### Development Notes

- Use Kafka for asynchronous boundaries and side-effects (carts, inventory commit, notifications).
- Use REST for synchronous validations and orchestrations (orders → payments, orders → inventory reservation, products → inventory reads).
