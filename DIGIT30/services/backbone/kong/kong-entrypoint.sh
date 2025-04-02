#!/bin/sh

# Debug environment variable
echo "KONG_NGINX_DAEMON=$KONG_NGINX_DAEMON"

echo "Entering Kong entrypoint script..."

# Wait for digit-postgres to be ready with a timeout
timeout=60
attempt=0
until pg_isready -h kong-db -p 5432 -U kong || [ $attempt -ge $timeout ]; do
  echo "Waiting for kong-db... ($attempt/$timeout)"
  sleep 2
  attempt=$((attempt + 2))
done

if [ $attempt -ge $timeout ]; then
  echo "Error: Timed out waiting for kong-db"
  exit 1
fi

# Bootstrap the database if not already done
kong migrations bootstrap
if [ $? -ne 0 ]; then
  echo "Error: Failed to bootstrap database"
  exit 1
fi

echo "Registering services with Consul..."

register_service() {
  SERVICE_ID=$1
  SERVICE_NAME=$2
  SERVICE_ADDRESS=$3
  SERVICE_PORT=$4

  echo "⏳ Registering $SERVICE_NAME ($SERVICE_ID) on $SERVICE_ADDRESS:$SERVICE_PORT..."

  curl --silent --fail --request PUT \
    --data "{
      \"ID\": \"$SERVICE_ID\",
      \"Name\": \"$SERVICE_NAME\",
      \"Address\": \"$SERVICE_ADDRESS\",
      \"Port\": $SERVICE_PORT
    }" \
    http://consul:8500/v1/agent/service/register

  if [ $? -eq 0 ]; then
    echo "✅ Registered $SERVICE_NAME"
  else
    echo "❌ Failed to register $SERVICE_NAME"
  fi
}

# Register each service
register_service jaeger jaeger jaeger 16686
register_service kafka kafka kafka 9092
register_service keycloak keycloak keycloak 8080
register_service redis redis redis 6379
register_service kong kong kong 8000
register_service kong-admin kong-admin kong 8001
register_service prometheus prometheus prometheus 9090
register_service grafana grafana grafana 3000

echo "✅ All registration requests sent"

# Start Kong with verbose output and explicit foreground
echo "Starting Kong..."
exec kong start --vv