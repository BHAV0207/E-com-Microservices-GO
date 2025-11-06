## product-service

Manages products. On product creation/deletion it coordinates with `inventory-service` using an internal worker pool to create/delete inventory documents.

### Run

Exposed on port 4000 via docker-compose.

### REST API

- POST `/product` — Create product
- GET `/product` — List all products
- GET `/product/{id}` — Get product with inventory (parallel fetch)
- PATCH `/product/{id}` — Update fields
- DELETE `/product/{id}` — Delete product and its inventory

### Inventory Coupling

- After product insert, submits a job to create a corresponding inventory document in `inventory-service`.
- On delete, removes product and calls inventory removal.

### Notes

- Uses goroutines and `sync.WaitGroup` to parallelize product and inventory fetch in `GET /product/{id}`.
