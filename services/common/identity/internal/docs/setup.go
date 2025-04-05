package docs

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupDocumentation registers Swagger UI and OpenAPI endpoints on the given Gin router.
func SetupDocumentation(router *gin.Engine) {
	// Serve Swagger UI and OpenAPI JSON at /swagger/*any
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
