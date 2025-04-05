package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	consulapi "github.com/hashicorp/consul/api"
)

// MCPConfig holds MCP server configuration
type MCPConfig struct {
	Port            int
	ModelContextURL string
	LlamaServerURL  string
}

// Service represents a discovered service
type Service struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Endpoints   []Endpoint `json:"endpoints"`
}

// Endpoint represents an API endpoint
type Endpoint struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

func main() {
	config := MCPConfig{
		Port:            8086,
		ModelContextURL: "http://model-context-service:8085",
		LlamaServerURL:  "http://llama-server:8082",
	}

	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Register with Consul
	err := registerWithConsul(config.Port)
	if err != nil {
		log.Fatalf("Failed to register with Consul: %v", err)
	}

	// MCP endpoints
	r.GET("/services", listServices)
	r.GET("/services/:name", getServiceDetails)
	r.POST("/chat", handleChat)
	r.GET("/health", healthCheck)

	log.Printf("MCP server starting on port %d", config.Port)
	log.Fatal(r.Run(fmt.Sprintf(":%d", config.Port)))
}

func registerWithConsul(port int) error {
	config := consulapi.DefaultConfig()
	config.Address = "consul:8500"

	client, err := consulapi.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create consul client: %v", err)
	}

	service := &consulapi.AgentServiceRegistration{
		ID:   "mcp-service",
		Name: "mcp-service",
		Port: port,
		Check: &consulapi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://mcp-service:%d/health", port),
			Interval: "10s",
			Timeout:  "5s",
		},
	}

	return client.Agent().ServiceRegister(service)
}

func listServices(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().Get("http://model-context-service:8085/contexts")
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to get services: %v", err)})
		return
	}

	var result map[string][]string
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		c.JSON(500, gin.H{"error": "failed to parse response"})
		return
	}

	c.JSON(200, result)
}

func getServiceDetails(c *gin.Context) {
	serviceName := c.Param("name")
	client := resty.New()
	resp, err := client.R().Get(fmt.Sprintf("http://model-context-service:8085/context/%s", serviceName))
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to get service details: %v", err)})
		return
	}

	c.Data(200, "application/json", resp.Body())
}

func getServiceAddress(serviceName string) (string, error) {
	config := consulapi.DefaultConfig()
	config.Address = "consul:8500"

	client, err := consulapi.NewClient(config)
	if err != nil {
		return "", fmt.Errorf("failed to create consul client: %v", err)
	}

	services, _, err := client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get service from consul: %v", err)
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no healthy instances found for service: %s", serviceName)
	}

	return services[0].Service.Address, nil
}

func handleChat(c *gin.Context) {
	var request struct {
		Message string `json:"message"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// Get available services
	client := resty.New()
	servicesResp, err := client.R().Get("http://model-context-service:8085/contexts")
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to get services: %v", err)})
		return
	}

	var services map[string][]string
	if err := json.Unmarshal(servicesResp.Body(), &services); err != nil {
		c.JSON(500, gin.H{"error": "failed to parse services"})
		return
	}

	// Create context for LLM with service details
	context := fmt.Sprintf("Available services: %v\n", services["available_contexts"])

	// Get details for each service
	for _, service := range services["available_contexts"] {
		detailsResp, err := client.R().Get(fmt.Sprintf("http://model-context-service:8085/context/%s", service))
		if err != nil {
			continue
		}
		context += fmt.Sprintf("\nService %s details:\n%s\n", service, string(detailsResp.Body()))
	}

	// Call llama-server to understand the request and get API call details
	llamaResp, err := client.R().
		SetBody(map[string]interface{}{
			"messages": []map[string]string{
				{"role": "system", "content": fmt.Sprintf(`You are an API assistant. Given a user request and available services, determine:
1. Which service to call
2. Which endpoint to use
3. What parameters are needed
4. Return ONLY a JSON object in this exact format:
{
    "service": "service_name",
    "endpoint": "endpoint_path",
    "method": "HTTP_METHOD",
    "parameters": {}
}

Available services and their details:
%s`, context)},
				{"role": "user", "content": request.Message},
			},
		}).
		Post("http://llama-server:8082/v1/chat/completions")

	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to get LLM response: %v", err)})
		return
	}

	// Parse the LLM response to get API call details
	var llamaResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(llamaResp.Body(), &llamaResponse); err != nil {
		c.JSON(500, gin.H{"error": "failed to parse LLM response"})
		return
	}

	// Extract the JSON object from the LLM's response
	content := llamaResponse.Choices[0].Message.Content
	// Find the JSON object in the response
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start == -1 || end == -1 {
		c.JSON(500, gin.H{"error": "no API call details found in LLM response"})
		return
	}
	jsonStr := content[start : end+1]

	var apiCall struct {
		Service    string                 `json:"service"`
		Endpoint   string                 `json:"endpoint"`
		Method     string                 `json:"method"`
		Parameters map[string]interface{} `json:"parameters"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &apiCall); err != nil {
		c.JSON(500, gin.H{"error": "failed to parse API call details"})
		return
	}

	// Get service address from Consul
	serviceAddr, err := getServiceAddress(apiCall.Service)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to get service address: %v", err)})
		return
	}

	// Convert parameters to string map for query params
	queryParams := make(map[string]string)
	for k, v := range apiCall.Parameters {
		if str, ok := v.(string); ok {
			queryParams[k] = str
		}
	}

	// Make the actual API call
	var apiResp *resty.Response
	url := fmt.Sprintf("http://%s:8080%s", serviceAddr, apiCall.Endpoint)

	switch apiCall.Method {
	case "GET":
		apiResp, err = client.R().SetQueryParams(queryParams).Get(url)
	case "POST":
		apiResp, err = client.R().SetBody(apiCall.Parameters).Post(url)
	default:
		c.JSON(400, gin.H{"error": "unsupported HTTP method"})
		return
	}

	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to call API: %v", err)})
		return
	}

	// Return both the API response and the LLM's explanation
	c.JSON(200, gin.H{
		"api_response": string(apiResp.Body()),
		"explanation":  content,
	})
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "healthy"})
}
