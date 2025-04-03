-- handler.lua
local prometheus = require "kong.plugins.prometheus.exporter"  -- Ensure the prometheus library is available

local ApiCosting = {}
ApiCosting.PRIORITY = 900
ApiCosting.VERSION = "1.0.0"

-- We'll lazily initialize the counter.
local api_usage_counter = nil

function ApiCosting:access(conf)
  -- Initialize the counter if it hasn't been already.
  if not api_usage_counter then
    api_usage_counter = prometheus.counter(
      "api_usage_total",
      "Total number of API calls by account, service, and API endpoint",
      {"account", "service", "api", "cost"}
    )
  end

  -- Extract account information (could also come from the consumer object if authenticated)
  local account = kong.request.get_header("X-Account") or "unknown"
  -- Get service name (or assign a default if unavailable)
  local service = kong.router.get_service() and kong.router.get_service().name or "unknown"
  -- Use the request path as the API endpoint
  local api = kong.request.get_path()
  -- Set a cost value (could be static or dynamic based on your pricing logic)
  local cost = "1"

  -- Increment the counter with the labels
  api_usage_counter:inc(1, {account, service, api, cost})
end

return ApiCosting