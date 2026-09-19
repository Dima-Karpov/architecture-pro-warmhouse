#!/bin/bash

# Exit on any error
set -e

if [ ! -f .env ]; then
  cp .env.example .env
  echo "Created .env from .env.example"
fi

echo "Starting the Smart Home Sensor API..."
echo "Building and starting containers..."
docker-compose up --build -d

echo "Waiting for services to be ready..."
# Wait for PostgreSQL to be ready
for i in {1..30}; do
  if docker exec db sh -c 'pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB"' > /dev/null 2>&1; then
    echo "PostgreSQL is ready!"
    break
  fi
  echo "Waiting for PostgreSQL to start... ($i/30)"
  sleep 1
done

# Check if PostgreSQL is ready
if ! docker exec db sh -c 'pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB"' > /dev/null 2>&1; then
  echo "Error: PostgreSQL did not start within the expected time."
  exit 1
fi

echo "All services are up and running!"
echo "Monolith:        http://localhost:8080"
echo "temperature-api: http://localhost:8081/swagger"
echo "device:          http://localhost:8082/swagger"
echo "telemetry:       http://localhost:8083/swagger"
echo "scenario:        http://localhost:8084/health"
echo "RabbitMQ UI:     http://localhost:15672  (guest/guest)"
echo ""
echo "To view logs, run: docker-compose logs -f"
echo "To stop the services, run: docker-compose down"