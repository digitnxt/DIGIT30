FROM kong:latest
USER root
COPY plugins/api-costing /usr/local/share/lua/5.1/kong/plugins/api-costing
RUN apt-get update && apt-get install -y curl postgresql-client && rm -rf /var/lib/apt/lists/*
COPY kong-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/kong-entrypoint.sh
USER kong
ENTRYPOINT ["/usr/local/bin/kong-entrypoint.sh"]