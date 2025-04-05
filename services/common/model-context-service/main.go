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

// ServiceConfig holds service-specific configuration
type ServiceConfig struct {
	SwaggerPath string
	Port        int
}

var serviceConfigs = map[string]ServiceConfig{
	"identity": {SwaggerPath: "/swagger/doc.json", Port: 8080},
	// Add more services here as they become available
}

// ModelContext represents the structured format for LLM consumption
type ModelContext struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Endpoints   []EndpointContext      `json:"endpoints"`
	Schemas     map[string]SchemaInfo  `json:"schemas"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// EndpointContext represents an API endpoint context
type EndpointContext struct {
	Path        string                 `json:"path"`
	Method      string                 `json:"method"`
	Summary     string                 `json:"summary"`
	Description string                 `json:"description"`
	Parameters  []ParameterContext     `json:"parameters,omitempty"`
	RequestBody map[string]interface{} `json:"requestBody,omitempty"`
	Responses   map[string]interface{} `json:"responses"`
	Tags        []string               `json:"tags,omitempty"`
}

// ParameterContext represents an API parameter context
type ParameterContext struct {
	Name        string     `json:"name"`
	In          string     `json:"in"` // query, path, header, cookie
	Required    bool       `json:"required"`
	Schema      SchemaInfo `json:"schema"`
	Description string     `json:"description"`
}

// SchemaInfo represents type information
type SchemaInfo struct {
	Type       string                 `json:"type"`
	Format     string                 `json:"format,omitempty"`
	Properties map[string]SchemaInfo  `json:"properties,omitempty"`
	Items      *SchemaInfo            `json:"items,omitempty"`
	Required   []string               `json:"required,omitempty"`
	Example    interface{}            `json:"example,omitempty"`
	Additional map[string]interface{} `json:"additional,omitempty"`
}

func main() {
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

	// Endpoint to get model context for a specific service
	r.GET("/context/:service", getModelContext)

	// Endpoint to list all available service contexts
	r.GET("/contexts", listAvailableContexts)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	log.Fatal(r.Run(":8085"))
}

func getModelContext(c *gin.Context) {
	service := c.Param("service")

	// Get service configuration
	config, exists := serviceConfigs[service]
	if !exists {
		c.JSON(404, gin.H{"error": "service not configured"})
		return
	}

	// Get service address from Consul
	serviceAddr, err := getServiceAddress(service)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to get service address: %v", err)})
		return
	}

	// Build Swagger doc URL
	swaggerURL := fmt.Sprintf("http://%s:%d%s", serviceAddr, config.Port, config.SwaggerPath)

	// Fetch Swagger documentation
	client := resty.New()
	resp, err := client.R().Get(swaggerURL)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to fetch swagger doc: %v", err)})
		return
	}

	var swaggerDoc map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &swaggerDoc); err != nil {
		c.JSON(500, gin.H{"error": "failed to parse swagger doc"})
		return
	}

	// Convert Swagger to ModelContext
	modelContext := convertSwaggerToModelContext(swaggerDoc)
	c.JSON(200, modelContext)
}

func listAvailableContexts(c *gin.Context) {
	// Get list of services from Consul
	services, err := getConsulServices()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to get services: %v", err)})
		return
	}

	// Filter services that have Swagger documentation configured
	var availableContexts []string
	for service := range services {
		if _, exists := serviceConfigs[service]; exists {
			availableContexts = append(availableContexts, service)
		}
	}

	c.JSON(200, gin.H{
		"available_contexts": availableContexts,
	})
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

func getConsulServices() (map[string][]string, error) {
	config := consulapi.DefaultConfig()
	config.Address = "consul:8500"

	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %v", err)
	}

	services, _, err := client.Catalog().Services(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get services from consul: %v", err)
	}

	return services, nil
}

func convertSwaggerToModelContext(swagger map[string]interface{}) ModelContext {
	mc := ModelContext{
		Name:        swagger["info"].(map[string]interface{})["title"].(string),
		Description: swagger["info"].(map[string]interface{})["description"].(string),
		Endpoints:   make([]EndpointContext, 0),
		Schemas:     make(map[string]SchemaInfo),
		Metadata: map[string]interface{}{
			"version":  swagger["info"].(map[string]interface{})["version"],
			"basePath": swagger["basePath"],
			"host":     swagger["host"],
		},
	}

	// Convert paths to endpoints
	paths := swagger["paths"].(map[string]interface{})
	for path, methods := range paths {
		for method, details := range methods.(map[string]interface{}) {
			endpoint := convertToEndpoint(path, method, details.(map[string]interface{}))
			mc.Endpoints = append(mc.Endpoints, endpoint)
		}
	}

	// Convert schemas/definitions
	if definitions, ok := swagger["definitions"].(map[string]interface{}); ok {
		for name, schema := range definitions {
			mc.Schemas[name] = convertToSchema(schema.(map[string]interface{}))
		}
	}

	return mc
}

func convertToEndpoint(path, method string, details map[string]interface{}) EndpointContext {
	endpoint := EndpointContext{
		Path:        path,
		Method:      strings.ToUpper(method),
		Summary:     getStringValue(details, "summary"),
		Description: getStringValue(details, "description"),
		Parameters:  make([]ParameterContext, 0),
		Responses:   make(map[string]interface{}),
		Tags:        make([]string, 0),
	}

	// Convert parameters
	if params, ok := details["parameters"].([]interface{}); ok {
		for _, param := range params {
			paramMap := param.(map[string]interface{})
			parameter := ParameterContext{
				Name:        paramMap["name"].(string),
				In:          paramMap["in"].(string),
				Required:    getBoolValue(paramMap, "required"),
				Description: getStringValue(paramMap, "description"),
			}
			if schema, ok := paramMap["schema"].(map[string]interface{}); ok {
				parameter.Schema = convertToSchema(schema)
			}
			endpoint.Parameters = append(endpoint.Parameters, parameter)
		}
	}

	// Convert responses
	if responses, ok := details["responses"].(map[string]interface{}); ok {
		endpoint.Responses = responses
	}

	// Convert tags
	if tags, ok := details["tags"].([]interface{}); ok {
		for _, tag := range tags {
			endpoint.Tags = append(endpoint.Tags, tag.(string))
		}
	}

	return endpoint
}

func convertToSchema(schema map[string]interface{}) SchemaInfo {
	schemaInfo := SchemaInfo{
		Type:       getStringValue(schema, "type"),
		Format:     getStringValue(schema, "format"),
		Properties: make(map[string]SchemaInfo),
		Additional: make(map[string]interface{}),
	}

	if properties, ok := schema["properties"].(map[string]interface{}); ok {
		for name, prop := range properties {
			schemaInfo.Properties[name] = convertToSchema(prop.(map[string]interface{}))
		}
	}

	if required, ok := schema["required"].([]interface{}); ok {
		for _, req := range required {
			schemaInfo.Required = append(schemaInfo.Required, req.(string))
		}
	}

	if items, ok := schema["items"].(map[string]interface{}); ok {
		itemSchema := convertToSchema(items)
		schemaInfo.Items = &itemSchema
	}

	// Store any additional fields
	for k, v := range schema {
		if !isStandardField(k) {
			schemaInfo.Additional[k] = v
		}
	}

	return schemaInfo
}

func isStandardField(field string) bool {
	standardFields := []string{"type", "format", "properties", "required", "items"}
	for _, f := range standardFields {
		if f == field {
			return true
		}
	}
	return false
}

func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getBoolValue(m map[string]interface{}, key string) bool {
	if val, ok := m[key].(bool); ok {
		return val
	}
	return false
}
