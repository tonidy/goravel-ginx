package examples

import (
	"net/http"

	contractshttp "github.com/goravel/framework/contracts/http"
	ginx "github.com/goravel/gin"
	"github.com/goravel/gin/facades"
)

// Example request/response types
type GetUserRequest struct {
	ID string `uri:"id" binding:"required" json:"id" description:"User ID"`
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required" description:"User's full name"`
	Email string `json:"email" binding:"required,email" description:"User's email address"`
}

type UserResponse struct {
	ID    string `json:"id" description:"User ID"`
	Name  string `json:"name" description:"User's full name"`
	Email string `json:"email" description:"User's email address"`
}

type ErrorResponse struct {
	Error string `json:"error" description:"Error message"`
}

// Example 1: Basic usage with Tonic routes
func SetupBasicRoutes() {
	// Get OpenAPI instance
	openapi := facades.OpenApi()

	// Create user group
	userGroup := facades.Route("").Prefix("/users")

	// Add typed routes with Tonic documentation
	ginx.AddTonicRouteTyped[GetUserRequest, UserResponse](
		userGroup,
		openapi,
		http.MethodGet,
		"/:id",
		ginx.BindMiddleware[GetUserRequest](),
		GetUserHandler,
	)

	ginx.AddTonicRouteTyped[CreateUserRequest, UserResponse](
		userGroup,
		openapi,
		http.MethodPost,
		"/",
		ginx.BindMiddleware[CreateUserRequest](),
		CreateUserHandler,
	)
}

// Example 2: Using Group methods for Tonic routes
func SetupGroupRoutes() {
	openapi := facades.OpenApi()

	// Cast to *ginx.Group to access Tonic methods
	userGroup := facades.Route("").Prefix("/api/users").(*ginx.Group)

	// Using convenience methods
	userGroup.TonicGet(openapi, "/:id",
		ginx.BindMiddleware[GetUserRequest](),
		GetUserHandler,
	)

	userGroup.TonicPost(openapi, "/",
		ginx.BindMiddleware[CreateUserRequest](),
		CreateUserHandler,
	)

	userGroup.TonicPut(openapi, "/:id",
		ginx.BindMiddleware[CreateUserRequest](),
		UpdateUserHandler,
	)

	userGroup.TonicDelete(openapi, "/:id",
		ginx.BindMiddleware[GetUserRequest](),
		DeleteUserHandler,
	)
}

// Example 3: With middleware
func SetupRoutesWithMiddleware() {
	openapi := facades.OpenApi()

	// Create authenticated group
	authGroup := facades.Route("").
		Prefix("/api/v1/users").(*ginx.Group)

	// Add routes with authentication
	authGroup.TonicGet(openapi, "/:id",
		ginx.BindMiddleware[GetUserRequest](),
		GetUserHandler,
	)

	authGroup.TonicPost(openapi, "/",
		ginx.BindMiddleware[CreateUserRequest](),
		CreateUserHandler,
	)
}

// Handlers
func GetUserHandler(ctx contractshttp.Context) contractshttp.Response {
	// Retrieve bound request
	req, ok := ginx.GetRequest[GetUserRequest](ctx)
	if !ok {
		return ctx.Response().Json(http.StatusBadRequest, contractshttp.Json{
			"error": "invalid request",
		})
	}

	// Business logic
	user := UserResponse{
		ID:    req.ID,
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	return ctx.Response().Json(http.StatusOK, user)
}

func CreateUserHandler(ctx contractshttp.Context) contractshttp.Response {
	req, ok := ginx.GetRequest[CreateUserRequest](ctx)
	if !ok {
		return ctx.Response().Json(http.StatusBadRequest, contractshttp.Json{
			"error": "invalid request",
		})
	}

	// Business logic
	user := UserResponse{
		ID:    "generated-id",
		Name:  req.Name,
		Email: req.Email,
	}

	return ctx.Response().Json(http.StatusCreated, user)
}

func UpdateUserHandler(ctx contractshttp.Context) contractshttp.Response {
	// Retrieve user ID from path
	idReq, _ := ginx.GetRequest[GetUserRequest](ctx)

	// In real implementation, you might need to bind body separately
	// or create a combined request type

	user := UserResponse{
		ID:    idReq.ID,
		Name:  "Updated Name",
		Email: "updated@example.com",
	}

	return ctx.Response().Json(http.StatusOK, user)
}

func DeleteUserHandler(ctx contractshttp.Context) contractshttp.Response {
	req, ok := ginx.GetRequest[GetUserRequest](ctx)
	if !ok {
		return ctx.Response().Json(http.StatusBadRequest, contractshttp.Json{
			"error": "invalid request",
		})
	}

	return ctx.Response().Json(http.StatusOK, contractshttp.Json{
		"message": "User " + req.ID + " deleted successfully",
	})
}

// Middleware example
func AuthMiddleware() contractshttp.HandlerFunc {
	return func(ctx contractshttp.Context) contractshttp.Response {
		// Your authentication logic here
		token := ctx.Request().Header("Authorization", "")
		if token == "" {
			return ctx.Response().Json(http.StatusUnauthorized, contractshttp.Json{
				"error": "unauthorized",
			})
		}
		ctx.Request().Next()
		return nil
	}
}
