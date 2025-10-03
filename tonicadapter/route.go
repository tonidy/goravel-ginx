package tonicadapter

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/TickLabVN/tonic/core/docs"
	"github.com/TickLabVN/tonic/core/utils"
	"github.com/gin-gonic/gin"
)

// Route describes an HTTP route to register with gin along with optional Tonic documentation tweaks.
type Route struct {
	Method   string
	Path     string
	Handlers []gin.HandlerFunc
	opts     []docs.OperationObject
}

// WithOptions attaches extra OpenAPI operation overrides for this route.
func (r Route) WithOptions(opts ...docs.OperationObject) Route {
	r.opts = append(r.opts, opts...)
	return r
}

// AddRoute registers the route with gin and hydrates the supplied OpenAPI spec.
func AddRoute[D any, R any](spec *docs.OpenApi, g gin.IRoutes, route Route) {
	_, resp := reflect.TypeOf(new(D)), reflect.TypeOf(new(R))
	parsingKey := "json"

	responseContent := defaultResponseContent()
	if shouldRegisterSchema(resp) {
		if _, err := spec.Components.AddSchema(resp, parsingKey, "binding"); err == nil {
			responseContent = map[string]docs.MediaTypeObject{
				"application/json": {
					Schema: &docs.SchemaOrReference{
						ReferenceObject: &docs.ReferenceObject{
							Ref: fmt.Sprintf("%s_%s", utils.GetSchemaPath(resp), parsingKey),
						},
					},
				},
			}
		}
	}

	var basePath string
	switch v := g.(type) {
	case *gin.RouterGroup:
		basePath = v.BasePath()
	case *gin.Engine:
		basePath = v.BasePath()
	default:
		panic("Invalid gin.IRoutes type, expected *gin.RouterGroup or *gin.Engine")
	}

	baseOp := utils.MergeStructs(route.opts...)
	path := fmt.Sprintf("%s%s", basePath, route.Path)

	op := utils.MergeStructs(baseOp, docs.OperationObject{
		OperationId: fmt.Sprintf("%s_%s", route.Method, path),
		Responses: map[string]docs.ResponseOrReference{
			"200": {
				ResponseObject: &docs.ResponseObject{
					Content: responseContent,
				},
			},
		},
	})

	if spec.Paths == nil {
		spec.Paths = make(docs.Paths)
	}

	pathItem := docs.PathItemObject{}
	switch route.Method {
	case http.MethodGet:
		g.GET(route.Path, route.Handlers...)
		pathItem.Get = &op
	case http.MethodPost:
		g.POST(route.Path, route.Handlers...)
		pathItem.Post = &op
	case http.MethodPut:
		g.PUT(route.Path, route.Handlers...)
		pathItem.Put = &op
	case http.MethodPatch:
		g.PATCH(route.Path, route.Handlers...)
		pathItem.Patch = &op
	case http.MethodDelete:
		g.DELETE(route.Path, route.Handlers...)
		pathItem.Delete = &op
	case http.MethodOptions:
		g.OPTIONS(route.Path, route.Handlers...)
		pathItem.Options = &op
	case http.MethodHead:
		g.HEAD(route.Path, route.Handlers...)
		pathItem.Head = &op
	default:
		fmt.Printf("Unsupported HTTP method: %s\n", route.Method)
	}

	spec.Paths.Update(path, pathItem)
}

func defaultResponseContent() map[string]docs.MediaTypeObject {
	return map[string]docs.MediaTypeObject{
		"application/json": {
			Schema: &docs.SchemaOrReference{
				SchemaObject: &docs.SchemaObject{Type: "object"},
			},
		},
	}
}

func shouldRegisterSchema(t reflect.Type) bool {
	if t == nil {
		return false
	}

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t == nil {
		return false
	}

	return t.Kind() != reflect.Interface && t.Kind() != reflect.Invalid
}
