#!/bin/bash

PROJECT_NAME=digit30  # change if you're using a different project name
SERVICES=(
  postgres
  consul
  kafka
  kong-db
  kong
  kong-synch
  keycloak
  redis
)

echo "Checking status of Docker Compose services under project: $PROJECT_NAME"
echo "---------------------------------------------------------------"

all_up=true

for service in "${SERVICES[@]}"; do
  status=$(docker inspect -f '{{.State.Status}}' "${PROJECT_NAME}_${service}_1" 2>/dev/null)
  
  if [ "$status" == "running" ]; then
    echo "✅ $service is running"
  else
    echo "❌ $service is NOT running (status: $status)"
    all_up=false
  fi
done

echo "---------------------------------------------------------------"
if [ "$all_up" = true ]; then
  echo "🎉 All services are up and running."
else
  echo "⚠️  Some services are not running."
fi