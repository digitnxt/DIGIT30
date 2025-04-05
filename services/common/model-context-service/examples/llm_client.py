import os
import json
import requests
from openai import OpenAI

class ModelContextLLMClient:
    def __init__(self, model_context_url="http://localhost:8085", openai_api_key=None):
        self.model_context_url = model_context_url
        self.client = OpenAI(api_key=openai_api_key or os.getenv("OPENAI_API_KEY"))
        
    def get_available_services(self):
        """Get list of available services."""
        response = requests.get(f"{self.model_context_url}/contexts")
        return response.json()["available_contexts"]
    
    def get_service_context(self, service_name):
        """Get the model context for a specific service."""
        response = requests.get(f"{self.model_context_url}/context/{service_name}")
        return response.json()
    
    def generate_api_call(self, service_name, user_request):
        """Generate an API call based on user request and service context."""
        # Get service context
        context = self.get_service_context(service_name)
        
        # Create prompt for the LLM
        system_prompt = f"""You are an API expert. Given a service's API documentation and a user request, 
        generate the appropriate curl command to call the API. Here's the API documentation:
        {json.dumps(context, indent=2)}
        
        Rules:
        1. Use the host and basePath from the metadata
        2. Include all required parameters
        3. Format the curl command properly
        4. If the API doesn't support the requested functionality, say so
        """
        
        # Call the LLM
        response = self.client.chat.completions.create(
            model="gpt-4",  # or "gpt-3.5-turbo"
            messages=[
                {"role": "system", "content": system_prompt},
                {"role": "user", "content": user_request}
            ]
        )
        
        return response.choices[0].message.content

    def explain_api(self, service_name, aspect=None):
        """Explain the API or a specific aspect of it."""
        context = self.get_service_context(service_name)
        
        system_prompt = f"""You are an API documentation expert. Given a service's API documentation,
        provide a clear explanation of the API. Here's the API documentation:
        {json.dumps(context, indent=2)}
        """
        
        user_prompt = "Explain this API" if not aspect else f"Explain {aspect} of this API"
        
        response = self.client.chat.completions.create(
            model="gpt-4",  # or "gpt-3.5-turbo"
            messages=[
                {"role": "system", "content": system_prompt},
                {"role": "user", "content": user_prompt}
            ]
        )
        
        return response.choices[0].message.content

def main():
    # Initialize the client
    client = ModelContextLLMClient()
    
    # Example 1: List available services
    print("Available services:", client.get_available_services())
    
    # Example 2: Generate an API call for the health check
    request = "I want to check if the identity service is healthy"
    print("\nGenerating API call for health check:")
    print(client.generate_api_call("identity", request))
    
    # Example 3: Explain the API
    print("\nAPI Explanation:")
    print(client.explain_api("identity"))
    
    # Example 4: Explain specific aspects
    print("\nEndpoints Explanation:")
    print(client.explain_api("identity", "available endpoints"))

if __name__ == "__main__":
    main() 