// Package docs provides utilities for setting up API documentation using Swagger.
package docs

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupDocumentation registers the Swagger and OpenAPI endpoints on the given Gin router.
// swaggerJSONPath is the path to your swagger.json file.
func SetupDocumentation(router *gin.Engine, swaggerJSONPath string) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/openapi.json", func(c *gin.Context) {
		c.File(swaggerJSONPath)
	})
}
