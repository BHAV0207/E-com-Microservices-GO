## cart-service

Manages user carts. Reacts to user lifecycle events to provision and cleanup carts.

### Run

Exposed on port 9000 via docker-compose.

### REST API

- POST `/cart` — Add item to cart
  - body: `{ "userId": string, "productId": string, "quantity": number }`
- GET `/cart/{userId}` — Get a user's cart
- DELETE `/cart/{userId}` — Delete a user's cart

### Kafka Consumption

- `user-created` → create an empty cart
- `user-deleted` → delete cart for that user

### Validations

- Validates user existence via `user-service` and product availability via `product-service` before adding items.
