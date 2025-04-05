import json
import requests
import sys
import time

# Force stdout to flush immediately
sys.stdout.reconfigure(line_buffering=True)

class ModelContextClient:
    def __init__(self):
        print("Initializing ModelContextClient...", flush=True)
        # Use the Docker service name for identity service
        self.identity_url = "http://identity:8080"
        # Use the Docker service name for llama server
        self.llama_url = "http://llama-server:8082/v1/chat/completions"
        # Add retry configuration
        self.max_retries = 3
        self.retry_delay = 2  # seconds
        # Wait for model to load
        self.wait_for_model()

    def wait_for_model(self):
        print("Waiting for model to load...", flush=True)
        max_wait = 300  # 5 minutes
        start_time = time.time()
        while time.time() - start_time < max_wait:
            try:
                response = requests.get("http://llama-server:8082/health")
                if response.status_code == 200:
                    print("Model is ready!", flush=True)
                    return
            except Exception as e:
                print(f"Waiting for model... ({str(e)})", flush=True)
            time.sleep(5)
        print("Warning: Model loading timeout reached", flush=True)

    def get_available_services(self):
        try:
            # Check if llama server is healthy
            health_check = requests.get("http://llama-server:8082/health", timeout=5)
            if health_check.status_code != 200:
                print("Warning: Llama server health check failed")
            
            # Get available services from Consul
            consul_url = "http://consul:8500/v1/catalog/services"
            response = requests.get(consul_url, timeout=5)
            response.raise_for_status()
            services = list(response.json().keys())
            return services
        except Exception as e:
            print(f"Error getting available services: {str(e)}")
            return []

    def get_service_context(self, service_name):
        try:
            if service_name == "identity":
                swagger_url = f"{self.identity_url}/swagger/doc.json"
                print(f"Trying to connect to: {swagger_url}")  # Debug line
                response = requests.get(swagger_url)
                response.raise_for_status()
                return response.json()
            return None
        except requests.exceptions.RequestException as e:
            print(f"Error getting service context: {str(e)}")
            return None

    def call_llama(self, system_prompt, user_prompt, max_retries=None):
        if max_retries is None:
            max_retries = self.max_retries
            
        for attempt in range(max_retries):
            try:
                payload = {
                    "messages": [
                        {"role": "system", "content": system_prompt},
                        {"role": "user", "content": user_prompt}
                    ],
                    "temperature": 0.7,
                    "max_tokens": 500
                }
                
                response = requests.post(
                    self.llama_url,
                    json=payload,
                    timeout=120,  # Increased timeout for long responses
                    headers={"Content-Type": "application/json"}
                )
                response.raise_for_status()
                result = response.json()
                return result["choices"][0]["message"]["content"]
            except requests.exceptions.RequestException as e:
                print(f"Error calling Llama server (attempt {attempt + 1}/{max_retries}): {str(e)}")
                if attempt < max_retries - 1:
                    time.sleep(self.retry_delay * (attempt + 1))  # Exponential backoff
                continue
            except Exception as e:
                print(f"Unexpected error: {str(e)}")
                return None
        return None

    def generate_api_call(self, service_name, request_details):
        try:
            context = self.get_service_context(service_name)
            if not context:
                return "Failed to get service context"

            # Update the host in the context to use the Docker service name
            context["host"] = "identity:8080"
            context["basePath"] = ""  # Remove the basePath since we're using the Docker service directly

            system_prompt = "You are an API expert. Generate a curl command for the requested API endpoint."
            user_prompt = f"""Based on this API documentation, generate a curl command for: {request_details}

API Context:
{json.dumps(context, indent=2)}

Please provide only the curl command, no additional explanation. Use the exact host and port from the API context."""

            response = self.call_llama(system_prompt, user_prompt)
            if response:
                return response

            # Fallback to template-based response
            if "health check" in request_details.lower():
                return f"""curl -X GET "http://{context['host']}/ping" \\
  -H "Accept: application/json" \\
  -H "Content-Type: application/json\""""
            return "Failed to generate API call"

        except Exception as e:
            return f"Error generating API call: {str(e)}"

    def explain_api(self, service_name, aspect=None):
        try:
            context = self.get_service_context(service_name)
            if not context:
                return "Failed to get service context"

            # Update the host in the context to use the Docker service name
            context["host"] = "identity:8080"
            context["basePath"] = ""  # Remove the basePath since we're using the Docker service directly

            system_prompt = "You are an API documentation expert. Explain APIs clearly and concisely."
            user_prompt = f"""Explain the following API{f' {aspect}' if aspect else ''}:

API Context:
{json.dumps(context, indent=2)}

Focus on the key aspects like endpoints, methods, parameters, and responses. Use the exact host and port from the API context."""

            response = self.call_llama(system_prompt, user_prompt)
            if response:
                return response

            # Fallback to template-based response
            if aspect and "health check" in aspect.lower():
                paths = context.get("paths", {})
                ping_info = paths.get("/ping", {}).get("get", {})
                if ping_info:
                    return f"""Health Check Endpoint:
- Path: /ping
- Method: GET
- Description: {ping_info.get('description', 'No description available')}
- Response: {ping_info.get('responses', {}).get('200', {}).get('description', 'No response description available')}"""
                return "Health check endpoint information not found"

            info = context.get("info", {})
            return f"""API Information:
- Title: {info.get('title', 'No title')}
- Version: {info.get('version', 'No version')}
- Description: {info.get('description', 'No description')}"""

        except Exception as e:
            return f"Error explaining API: {str(e)}"

def main():
    try:
        client = ModelContextClient()
        
        # Test Llama server connection
        print("\nTesting Llama server connection...")
        test_response = client.call_llama(
            "You are a helpful assistant.",
            "Say 'Llama server is working!' if you can read this."
        )
        if not test_response:
            print("Failed to connect to Llama server. Please make sure it's running on port 8082.")
            sys.exit(1)
        print("Llama server connection successful!")

        # Get available services
        services = client.get_available_services()
        print(f"\nAvailable services: {services}")

        if not services:
            print("No services found")
            return

        # Generate API call for health check
        print("\nGenerating API call for health check:")
        request = "Generate a curl command for the health check endpoint"
        print(client.generate_api_call("identity", request))

        # Explain API
        print("\nExplaining API:")
        print(client.explain_api("identity"))

        # Explain specific endpoint
        print("\nExplaining health check endpoint:")
        print(client.explain_api("identity", "health check endpoint"))

    except Exception as e:
        print(f"An error occurred: {str(e)}")
        sys.exit(1)

if __name__ == "__main__":
    main() 