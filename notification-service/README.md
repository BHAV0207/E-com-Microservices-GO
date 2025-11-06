## notification-service

Consumes events and stores/sends user notifications.

### Run

Exposed on port 2000 via docker-compose.

### Kafka Consumption

- Generic consumer (`GenericEvent`) reading a configured topic (set in code/compose) for events like:
  - `order.created`
  - `payment.success`
  - `payment.failed`
  - `order.shipped`, `order.delivered`
- Persists notifications to Mongo and logs a simulated send.

### Data

- Collection: `notifications` with fields: `userId`, `orderId`, `type`, `message`, `status`, `createdAt`.
