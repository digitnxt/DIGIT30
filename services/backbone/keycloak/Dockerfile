FROM quay.io/keycloak/keycloak:latest

# Copy the entrypoint script to a writable directory
COPY keycloak-entrypoint.sh /usr/local/bin/keycloak-entrypoint.sh

# Ensure the script is executable
RUN chmod +x /usr/local/bin/keycloak-entrypoint.sh

ENTRYPOINT ["/opt/keycloak-entrypoint.sh"]