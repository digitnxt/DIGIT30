#!/bin/bash
# Start Keycloak in the background
/opt/keycloak/bin/kc.sh start-dev &

# Give Keycloak time to start
sleep 10

# Register the Keycloak service with Consul
curl --request PUT \
  --data '{
    "ID": "keycloak",
    "Name": "keycloak",
    "Address": "keycloak",
    "Port": 8080,
    "Tags": ["idp", "keycloak"]
  }' \
  http://consul:8500/v1/agent/service/register

# Wait for Keycloak process to end
wait %1