## user-service

Handles user registration, login (JWT), and CRUD. Publishes user lifecycle events to Kafka for other services (like `cart-service`).

### Run

Service is wired in docker-compose (port 8080). To run standalone, export env and run the service binary.

### Environment

- `PORT` (default 8080 via compose)
- `MONGO_URI`
- `KAFKA_BROKER` (e.g., `kafka:9092`)
- Kafka Topics: `user-created`, `user-deleted`

### REST API

- POST `/register` — Create a user
  - body: `{ "name", "email", "password" }`
  - emits: `user-created` with `{ userId, email }`
- POST `/login` — Issue JWT token
  - body: `{ "email", "password" }`
  - returns: `{ token }`
- GET `/users` — List users
- GET `/users/{id}` — Get user by id
- PATCH `/users/{id}` — Update fields
- DELETE `/users/{id}` — Delete user
  - emits: `user-deleted` with `{ userId }`

### Events

- Produced
  - `user-created` on successful registration
  - `user-deleted` on delete

### Notes

- Passwords are hashed before persistence.
- JWT uses HS256; secret configured in code/env.
