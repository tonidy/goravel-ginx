package tonicadapter

import (
	"fmt"
	"strings"

	"github.com/TickLabVN/tonic/core"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gin-gonic/gin"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
)

// registerSwaggerRoutes registers Swagger JSON and UI routes on the given gin router
func registerSwaggerRoutes(router gin.IRoutes, spec *docs.OpenApi, path string) {
	swaggerPath := fmt.Sprintf("%s.json", path)
	router.GET(swaggerPath, gin.WrapH(core.JsonHttpHandler(spec)))
	router.GET(fmt.Sprintf("%s/*any", path), gin.WrapH(core.SwaggerUIHandler(swaggerPath)))
}

// registerSwaggerRoutesWithFacades registers Swagger JSON and UI routes using Goravel facades
func registerSwaggerRoutesWithFacades(router route.Route, spec *docs.OpenApi, path string) {
	swaggerPath := fmt.Sprintf("%s.json", path)

	// Wrap HTTP handlers for Goravel route.Route interface
	jsonHandler := core.JsonHttpHandler(spec)
	uiHandler := core.SwaggerUIHandler(swaggerPath)

	router.Get(swaggerPath, func(ctx http.Context) http.Response {
		w := ctx.Response().Writer()
		r := ctx.Request().Origin()
		jsonHandler.ServeHTTP(w, r)
		return nil
	})

	router.Get(fmt.Sprintf("%s/*any", path), func(ctx http.Context) http.Response {
		w := ctx.Response().Writer()
		r := ctx.Request().Origin()
		uiHandler.ServeHTTP(w, r)
		return nil
	})
}

// UIHandle wires Swagger UI routes exposing the generated OpenAPI spec.
func UIHandleWithEngine(e *gin.Engine, spec *docs.OpenApi, path string) {
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	} else {
		e.GET(path, func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s/index.html", path)) })
	}

	registerSwaggerRoutes(e, spec, path)
}

func UIHandleWithRouter(router gin.IRoutes, spec *docs.OpenApi, path string) {
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	} else {
		router.GET(path, func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s/index.html", path))
		})
	}

	registerSwaggerRoutes(router, spec, path)
}

func UIHandle(router route.Route, spec *docs.OpenApi, path string) {
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	} else {
		router.Get(path, func(ctx http.Context) http.Response {
			return ctx.Response().Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s/index.html", path))
		})
	}

	registerSwaggerRoutesWithFacades(router, spec, path)
}
