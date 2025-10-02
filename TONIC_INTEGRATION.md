# Tonic Swagger/OpenAPI Integration

This package now includes built-in support for [Tonic](https://github.com/TickLabVN/tonic), enabling automatic OpenAPI/Swagger documentation generation for your Goravel Gin routes.

## Features

- 🚀 Automatic OpenAPI 3.0.1 documentation generation
- 📝 Type-safe request/response definitions
- 🔄 Seamless integration with Goravel's routing system
- 🎯 Built-in request binding and validation
- 🏷️ Support for tags, descriptions, and metadata
- 📦 Configuration via Goravel's config system

## Installation

1. Install Tonic package:
```bash
go get github.com/TickLabVN/tonic
```

2. The service provider is already registered in goravel-ginx

## Configuration

Create or update your configuration file (e.g., `config/swagger.go`):

```go
package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()

	config.Add("swagger", map[string]any{
		"openapi_version": "3.0.3",

		"info": map[string]any{
			"version":     "1.0.0",
			"title":       "My API",
			"description": "My API Description",

			// Optional: Contact information
			"contact": map[string]any{
				"name":  "Your Name",
				"url":   "https://github.com/yourusername",
				"email": "your.email@example.com",
			},

			// Optional: License information
			"license": map[string]any{
				"name": "MIT",
				"url":  "https://opensource.org/licenses/MIT",
			},
		},
	})
}
```

## Usage

### Clean API (Recommended)

The cleanest way to use Tonic is by extending your existing Goravel route groups with `.WithTonic()`:

```go
package routes

import (
	"net/http"

	ginx "github.com/goravel/gin"
	"github.com/goravel/gin/facades"
	contractshttp "github.com/goravel/framework/contracts/http"
)

// Define request/response types
type GetUserRequest struct {
	ID string `uri:"id" binding:"required" json:"id" description:"User ID"`
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required" description:"User's full name"`
	Email string `json:"email" binding:"required,email" description:"User's email address"`
}

type UserResponse struct {
	ID    string `json:"id" description:"User ID"`
	Name  string `json:"name" description:"User's name"`
	Email string `json:"email" description:"User's email"`
}

func SetupRoutes() {
	openapi := facades.OpenApi()

	// Create your normal Goravel route group
	userGroup := facades.Route("").Prefix("/api/users").(*ginx.Group)

	// Enable Tonic documentation - keeps the familiar Goravel API
	tonic := userGroup.WithTonic(openapi)

	// Use standard route methods - they now generate OpenAPI docs automatically!
	tonic.Get("/:id",
		ginx.BindMiddleware[GetUserRequest](),
		GetUserHandler,
	)

	tonic.Post("/",
		ginx.BindMiddleware[CreateUserRequest](),
		CreateUserHandler,
	)

	tonic.Put("/:id",
		ginx.BindMiddleware[CreateUserRequest](),
		UpdateUserHandler,
	)

	tonic.Delete("/:id",
		ginx.BindMiddleware[GetUserRequest](),
		DeleteUserHandler,
	)
}

// Handler example
func GetUserHandler(ctx contractshttp.Context) contractshttp.Response {
	req, ok := ginx.GetRequest[GetUserRequest](ctx)
	if !ok {
		return ctx.Response().Json(http.StatusBadRequest, contractshttp.Json{
			"error": "invalid request",
		})
	}

	user := UserResponse{
		ID:    req.ID,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	return ctx.Response().Json(http.StatusOK, user)
}
```

### Using Typed Routes

For better OpenAPI documentation with explicit request/response types:

```go
func SetupTypedRoutes() {
	openapi := facades.OpenApi()
	userGroup := facades.Route("").Prefix("/api/users").(*ginx.Group)
	tonic := userGroup.WithTonic(openapi)

	// Typed routes provide explicit type information for OpenAPI docs
	ginx.GetTyped[GetUserRequest, UserResponse](tonic, "/:id",
		ginx.BindMiddleware[GetUserRequest](),
		GetUserHandler,
	)

	ginx.PostTyped[CreateUserRequest, UserResponse](tonic, "/",
		ginx.BindMiddleware[CreateUserRequest](),
		CreateUserHandler,
	)

	ginx.PutTyped[CreateUserRequest, UserResponse](tonic, "/:id",
		ginx.BindMiddleware[CreateUserRequest](),
		UpdateUserHandler,
	)

	ginx.DeleteTyped[GetUserRequest, UserResponse](tonic, "/:id",
		ginx.BindMiddleware[GetUserRequest](),
		DeleteUserHandler,
	)
}
```

### Works with Goravel Middleware

The Tonic integration works seamlessly with existing Goravel middleware:

```go
func SetupWithMiddleware() {
	openapi := facades.OpenApi()

	// Standard Goravel middleware chain
	authGroup := facades.Route("").
		Prefix("/api/v1/users").
		Middleware(AuthMiddleware()).(*ginx.Group)

	// Enable Tonic - middleware is automatically applied
	tonic := authGroup.WithTonic(openapi)

	tonic.Get("/:id",
		ginx.BindMiddleware[GetUserRequest](),
		GetUserHandler,
	)

	tonic.Post("/",
		ginx.BindMiddleware[CreateUserRequest](),
		CreateUserHandler,
	)
}
```

## API Reference

### TonicGroup Methods

Once you call `.WithTonic()` on a route group, you get a `TonicGroup` with these methods:

```go
// Standard route methods - generate OpenAPI docs automatically
tonic.Get(path, handlers...)
tonic.Post(path, handlers...)
tonic.Put(path, handlers...)
tonic.Delete(path, handlers...)
tonic.Patch(path, handlers...)
tonic.Options(path, handlers...)

// Typed route functions - for explicit request/response types
ginx.GetTyped[Req, Res](tonic, path, handlers...)
ginx.PostTyped[Req, Res](tonic, path, handlers...)
ginx.PutTyped[Req, Res](tonic, path, handlers...)
ginx.DeleteTyped[Req, Res](tonic, path, handlers...)
ginx.PatchTyped[Req, Res](tonic, path, handlers...)
```

### Request Binding

Use `BindMiddleware` to automatically bind and validate requests:

```go
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

userGroup.TonicPost(openapi, "/",
	ginx.BindMiddleware[CreateUserRequest](),
	CreateUserHandler,
)

func CreateUserHandler(ctx contractshttp.Context) {
	// Get validated request
	req, ok := ginx.GetRequest[CreateUserRequest](ctx)
	if !ok {
		// Request validation failed
		return
	}

	// Use req.Name and req.Email
}
```

The `BindMiddleware` automatically tries multiple binding methods:
1. URI parameters (`binding:"required"` in struct tag)
2. JSON body
3. Query parameters
4. General binding

### With Middleware

Combine Tonic routes with Goravel middleware:

```go
authGroup := facades.Route().
	Prefix("/api/v1/users").
	Middleware(AuthMiddleware()).(*ginx.Group)

authGroup.TonicGet(openapi, "/:id",
	ginx.BindMiddleware[GetUserRequest](),
	GetUserHandler,
)
```

### Accessing OpenAPI Instance

Three ways to access the OpenAPI instance:

```go
// 1. Via facade (recommended)
openapi := facades.OpenApi()

// 2. Via global variable
openapi := ginx.OpenApiInstance

// 3. Via application make
instance, _ := facades.App().Make("goravel.gin.openapi")
openapi := instance.(*docs.OpenApi)
```

## Struct Tags for Documentation

Use struct tags to enhance API documentation:

```go
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required" description:"User's full name"`
	Email string `json:"email" binding:"required,email" description:"Valid email address"`
	Age   int    `json:"age" binding:"min=18,max=120" description:"User's age (18-120)"`
}
```

### Supported Tags

- `json`: JSON field name
- `binding`: Validation rules (required, email, min, max, etc.)
- `description`: Field description in OpenAPI spec
- `uri`: URI parameter binding
- `form`: Form parameter binding
- `query`: Query parameter binding

## Complete Example

```go
package routes

import (
	"net/http"

	ginx "github.com/goravel/gin"
	"github.com/goravel/gin/facades"
	contractshttp "github.com/goravel/framework/contracts/http"
)

// Request/Response types
type GetUserRequest struct {
	ID string `uri:"id" binding:"required" description:"User ID"`
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required" description:"Full name"`
	Email string `json:"email" binding:"required,email" description:"Email address"`
}

type UserResponse struct {
	ID        string `json:"id" description:"User ID"`
	Name      string `json:"name" description:"Full name"`
	Email     string `json:"email" description:"Email address"`
	CreatedAt string `json:"created_at" description:"Creation timestamp"`
}

func RegisterUserRoutes() {
	openapi := facades.OpenApi()

	// Public routes
	publicGroup := facades.Route().Prefix("/api/v1/users").(*ginx.Group)

	publicGroup.TonicGet(openapi, "/:id",
		ginx.BindMiddleware[GetUserRequest](),
		GetUser,
	)

	// Protected routes
	protectedGroup := facades.Route().
		Prefix("/api/v1/users").
		Middleware(AuthMiddleware()).(*ginx.Group)

	protectedGroup.TonicPost(openapi, "/",
		ginx.BindMiddleware[CreateUserRequest](),
		CreateUser,
	)

	protectedGroup.TonicPut(openapi, "/:id",
		ginx.BindMiddleware[CreateUserRequest](),
		UpdateUser,
	)

	protectedGroup.TonicDelete(openapi, "/:id",
		ginx.BindMiddleware[GetUserRequest](),
		DeleteUser,
	)
}

func GetUser(ctx contractshttp.Context) {
	req, ok := ginx.GetRequest[GetUserRequest](ctx)
	if !ok {
		ctx.Request().AbortWithStatusJson(http.StatusBadRequest, contractshttp.Json{
			"error": "invalid request",
		})
		return
	}

	user := UserResponse{
		ID:        req.ID,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: "2024-01-01T00:00:00Z",
	}

	ctx.Response().Json(http.StatusOK, user)
}

func CreateUser(ctx contractshttp.Context) {
	req, ok := ginx.GetRequest[CreateUserRequest](ctx)
	if !ok {
		ctx.Request().AbortWithStatusJson(http.StatusBadRequest, contractshttp.Json{
			"error": "invalid request",
		})
		return
	}

	user := UserResponse{
		ID:        "generated-id",
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: "2024-01-01T00:00:00Z",
	}

	ctx.Response().Json(http.StatusCreated, user)
}

func UpdateUser(ctx contractshttp.Context) {
	// Implementation here
	ctx.Response().Json(http.StatusOK, contractshttp.Json{
		"message": "user updated",
	})
}

func DeleteUser(ctx contractshttp.Context) {
	req, ok := ginx.GetRequest[GetUserRequest](ctx)
	if !ok {
		ctx.Request().AbortWithStatusJson(http.StatusBadRequest, contractshttp.Json{
			"error": "invalid request",
		})
		return
	}

	ctx.Response().Json(http.StatusOK, contractshttp.Json{
		"message": "user " + req.ID + " deleted",
	})
}

func AuthMiddleware() contractshttp.HandlerFunc {
	return func(ctx contractshttp.Context) {
		token := ctx.Request().Header("Authorization", "")
		if token == "" {
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, contractshttp.Json{
				"error": "unauthorized",
			})
			return
		}
		ctx.Request().Next()
	}
}
```

## Benefits

1. **Type Safety**: Compile-time type checking for requests and responses
2. **Auto Documentation**: OpenAPI spec generated automatically from types
3. **Validation**: Built-in request validation using Gin's binding tags
4. **Integration**: Seamless integration with existing Goravel middleware and routing
5. **Standards**: OpenAPI 3.0.1 compliant documentation

## Notes

- The OpenAPI instance is automatically initialized when the service provider boots
- Configuration is optional - sensible defaults are provided
- All existing Goravel routing features remain available
- Tonic routes can be mixed with standard routes in the same application

## Troubleshooting

### OpenAPI instance is nil
Make sure the service provider is registered and the application has been booted.

### Routes not documented
Ensure you're passing the OpenAPI instance to the Tonic route methods.

### Type casting error
When using `TonicGet`, `TonicPost`, etc., you need to cast the route group:
```go
group := facades.Route().Prefix("/api").(*ginx.Group)
```

## Learn More

- [Tonic Documentation](https://github.com/TickLabVN/tonic)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Goravel Documentation](https://www.goravel.dev/)
