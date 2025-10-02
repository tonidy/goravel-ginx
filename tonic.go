package gin

import (
	"net/http"

	ginAdapter "github.com/TickLabVN/tonic/adapters/gin"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gin-gonic/gin"
	contractshttp "github.com/goravel/framework/contracts/http"
	contractsroute "github.com/goravel/framework/contracts/route"
)

// TonicRoute represents a Tonic-enabled route configuration
type TonicRoute struct {
	Method   string
	Path     string
	Handlers []contractshttp.HandlerFunc
}

// WithTonic enables Tonic documentation for this route group
// Returns a TonicGroup that wraps the existing Group with Tonic support
func (r *Group) WithTonic(openapi *docs.OpenApi) *TonicGroup {
	if openapi == nil {
		openapi = OpenApiInstance
	}
	return &TonicGroup{
		Group:   r,
		openapi: openapi,
	}
}

// TonicGroup wraps a Group with Tonic documentation support
type TonicGroup struct {
	*Group
	openapi *docs.OpenApi
}

// Get adds a GET route with Tonic documentation
func (t *TonicGroup) Get(path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return t.addTonicRoute(http.MethodGet, path, handlers...)
}

// Post adds a POST route with Tonic documentation
func (t *TonicGroup) Post(path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return t.addTonicRoute(http.MethodPost, path, handlers...)
}

// Put adds a PUT route with Tonic documentation
func (t *TonicGroup) Put(path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return t.addTonicRoute(http.MethodPut, path, handlers...)
}

// Delete adds a DELETE route with Tonic documentation
func (t *TonicGroup) Delete(path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return t.addTonicRoute(http.MethodDelete, path, handlers...)
}

// Patch adds a PATCH route with Tonic documentation
func (t *TonicGroup) Patch(path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return t.addTonicRoute(http.MethodPatch, path, handlers...)
}

// Options adds an OPTIONS route with Tonic documentation
func (t *TonicGroup) Options(path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return t.addTonicRoute(http.MethodOptions, path, handlers...)
}

// addTonicRoute is the internal method that handles Tonic route registration
func (t *TonicGroup) addTonicRoute(method, path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		ginHandlers[i] = handlerToGinHandler(h)
	}

	// Add route to Tonic OpenAPI documentation
	ginAdapter.AddRoute[any, any](t.openapi, t.Group.WithMiddlewares(), ginAdapter.Route{
		Method:   method,
		Path:     pathToGinPath(path),
		Handlers: ginHandlers,
	})

	return NewAction(method, t.Group.getFullPath(path), t.Group.getHandlerName(handlers[len(handlers)-1]))
}

// GetTyped adds a typed GET route with Tonic documentation
func GetTyped[Req any, Res any](t *TonicGroup, path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return addTonicRouteTyped[Req, Res](t, http.MethodGet, path, handlers...)
}

// PostTyped adds a typed POST route with Tonic documentation
func PostTyped[Req any, Res any](t *TonicGroup, path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return addTonicRouteTyped[Req, Res](t, http.MethodPost, path, handlers...)
}

// PutTyped adds a typed PUT route with Tonic documentation
func PutTyped[Req any, Res any](t *TonicGroup, path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return addTonicRouteTyped[Req, Res](t, http.MethodPut, path, handlers...)
}

// DeleteTyped adds a typed DELETE route with Tonic documentation
func DeleteTyped[Req any, Res any](t *TonicGroup, path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return addTonicRouteTyped[Req, Res](t, http.MethodDelete, path, handlers...)
}

// PatchTyped adds a typed PATCH route with Tonic documentation
func PatchTyped[Req any, Res any](t *TonicGroup, path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return addTonicRouteTyped[Req, Res](t, http.MethodPatch, path, handlers...)
}

// addTonicRouteTyped is the internal method for typed route registration
func addTonicRouteTyped[Req any, Res any](t *TonicGroup, method, path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		ginHandlers[i] = handlerToGinHandler(h)
	}

	// Add typed route to Tonic OpenAPI documentation
	ginAdapter.AddRoute[Req, Res](t.openapi, t.Group.WithMiddlewares(), ginAdapter.Route{
		Method:   method,
		Path:     pathToGinPath(path),
		Handlers: ginHandlers,
	})

	return NewAction(method, t.Group.getFullPath(path), t.Group.getHandlerName(handlers[len(handlers)-1]))
}

// AddTonicRoute adds a route with Tonic documentation support (legacy method)
func (r *Group) AddTonicRoute(openapi *docs.OpenApi, method, path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	if openapi == nil {
		openapi = OpenApiInstance
	}
	return r.WithTonic(openapi).addTonicRoute(method, path, handlers...)
}

// AddTonicRouteTyped adds a route with typed request/response for Tonic documentation
func AddTonicRouteTyped[Req any, Res any](group contractsroute.Router, openapi *docs.OpenApi, method, path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	if openapi == nil {
		openapi = OpenApiInstance
	}

	g, ok := group.(*Group)
	if !ok {
		panic("group must be *gin.Group")
	}

	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		ginHandlers[i] = handlerToGinHandler(h)
	}

	// Add typed route to Tonic OpenAPI documentation
	ginAdapter.AddRoute[Req, Res](openapi, g.instance.(*gin.RouterGroup), ginAdapter.Route{
		Method:   method,
		Path:     pathToGinPath(path),
		Handlers: ginHandlers,
	})

	return NewAction(method, g.getFullPath(path), g.getHandlerName(handlers[len(handlers)-1]))
}

// BindMiddleware creates a middleware that binds and validates request data using Gin's binding
func BindMiddleware[T any]() contractshttp.HandlerFunc {
	return func(ctx contractshttp.Context) contractshttp.Response {
		var req T
		ginCtx := ctx.(*Context).instance

		// Try binding URI parameters
		if err := ginCtx.ShouldBindUri(&req); err == nil {
			ctx.WithValue("request", req)
			ctx.Request().Next()
			return nil
		}

		// Try binding JSON body
		if err := ginCtx.ShouldBindJSON(&req); err == nil {
			ctx.WithValue("request", req)
			ctx.Request().Next()
			return nil
		}

		// Try binding query parameters
		if err := ginCtx.ShouldBindQuery(&req); err == nil {
			ctx.WithValue("request", req)
			ctx.Request().Next()
			return nil
		}

		// Try general binding
		if err := ginCtx.ShouldBind(&req); err != nil {
			return ctx.Response().Json(http.StatusBadRequest, contractshttp.Json{
				"error": err.Error(),
			})
		}

		ctx.WithValue("request", req)
		ctx.Request().Next()
		return nil
	}
}

// GetRequest retrieves the bound request data from context
func GetRequest[T any](ctx contractshttp.Context) (T, bool) {
	val := ctx.Value("request")
	if val == nil {
		var zero T
		return zero, false
	}

	req, ok := val.(T)
	return req, ok
}

// TonicGet adds a GET route with Tonic documentation
func (r *Group) TonicGet(openapi *docs.OpenApi, path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return r.AddTonicRoute(openapi, http.MethodGet, path, handlers...)
}

// TonicPost adds a POST route with Tonic documentation
func (r *Group) TonicPost(openapi *docs.OpenApi, path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return r.AddTonicRoute(openapi, http.MethodPost, path, handlers...)
}

// TonicPut adds a PUT route with Tonic documentation
func (r *Group) TonicPut(openapi *docs.OpenApi, path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return r.AddTonicRoute(openapi, http.MethodPut, path, handlers...)
}

// TonicDelete adds a DELETE route with Tonic documentation
func (r *Group) TonicDelete(openapi *docs.OpenApi, path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return r.AddTonicRoute(openapi, http.MethodDelete, path, handlers...)
}

// TonicPatch adds a PATCH route with Tonic documentation
func (r *Group) TonicPatch(openapi *docs.OpenApi, path string, handlers ...contractshttp.HandlerFunc) contractsroute.Action {
	return r.AddTonicRoute(openapi, http.MethodPatch, path, handlers...)
}
