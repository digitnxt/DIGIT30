# Model Context LLM Client Example

This example demonstrates how to use the model-context-service with an LLM (Large Language Model) to:
1. Generate API calls based on natural language requests
2. Get explanations of APIs
3. Query specific aspects of services

## Setup

1. Install dependencies:
```bash
pip install -r requirements.txt
```

2. Set your OpenAI API key:
```bash
export OPENAI_API_KEY=your_api_key_here
```

3. Ensure the model-context-service is running:
```bash
# From the root directory
cd deployment/docker-compose
docker-compose up -d
```

## Usage

Run the example script:
```bash
python llm_client.py
```

### Using in Your Code

```python
from llm_client import ModelContextLLMClient

# Initialize the client
client = ModelContextLLMClient()

# Get available services
services = client.get_available_services()
print(f"Available services: {services}")

# Generate an API call
request = "I want to check if the identity service is healthy"
curl_command = client.generate_api_call("identity", request)
print(f"Generated curl command: {curl_command}")

# Get API explanation
explanation = client.explain_api("identity")
print(f"API explanation: {explanation}")

# Get specific aspect explanation
endpoints = client.explain_api("identity", "available endpoints")
print(f"Endpoints explanation: {endpoints}")
```

## Features

1. **Service Discovery**: Automatically discovers available services through the model-context-service.

2. **API Call Generation**: Uses GPT-4 to generate curl commands based on:
   - Service API documentation
   - Natural language requests
   - Service metadata (host, basePath, etc.)

3. **API Explanation**: Get natural language explanations of:
   - Overall API functionality
   - Specific endpoints
   - Authentication methods
   - Data models
   - And more...

4. **Extensible**: Easy to add support for new services - just add them to the model-context-service.

## Example Outputs

1. **API Call Generation**:
   ```
   User: "I want to check if the identity service is healthy"
   Output: "curl http://localhost:8000/identity/ping"
   ```

2. **API Explanation**:
   ```
   User: "Explain this API"
   Output: "This is the Identity Service API (v1.0) which handles identity and authentication operations. 
   It provides endpoints for health monitoring and user authentication..."
   ```

## Error Handling

The client includes error handling for:
- Connection issues with model-context-service
- Invalid service names
- OpenAI API errors
- Missing API keys

## Contributing

Feel free to contribute by:
1. Adding new features
2. Improving prompts
3. Adding support for more LLM providers
4. Enhancing documentation 