package facades

import (
	"github.com/TickLabVN/tonic/core/docs"
	gin "github.com/tonidy/goravel-ginx"
)

func OpenApi() *docs.OpenApi {
	// Return the global instance if available
	if gin.OpenApiInstance != nil {
		return gin.OpenApiInstance
	}

	// Check if App is initialized
	if gin.App == nil {
		// Return a default instance if app not initialized
		return &docs.OpenApi{
			OpenAPI: "3.0.3",
			Info: docs.InfoObject{
				Version: "1.0.0",
				Title:   "API Documentation",
			},
		}
	}

	// Try to create from container
	instance, err := gin.App.Make(gin.BindingOpenApi)
	if err != nil {
		// Return a default instance if container not ready
		return &docs.OpenApi{
			OpenAPI: "3.0.1",
			Info: docs.InfoObject{
				Version: "1.0.0",
				Title:   "API Documentation",
			},
		}
	}

	return instance.(*docs.OpenApi)
}
