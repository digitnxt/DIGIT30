#!/bin/sh

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
register_service model-context model-context model-context-service 8085

echo "✅ All registration requests sent"