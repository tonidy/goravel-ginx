package tonicadapter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TickLabVN/tonic/core"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gin-gonic/gin"
)

// UIHandle wires Swagger UI routes exposing the generated OpenAPI spec.
func UIHandle(e *gin.Engine, spec *docs.OpenApi, path string) {
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	} else {
		e.GET(path, func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s/index.html", path)) })
	}

	swaggerPath := fmt.Sprintf("%s.json", path)
	e.GET(swaggerPath, gin.WrapH(core.JsonHttpHandler(spec)))
	e.GET(fmt.Sprintf("%s/*subpaths", path), gin.WrapH(core.SwaggerUIHandler(swaggerPath)))
}
