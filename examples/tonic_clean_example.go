package examples

import (
	"net/http"

	contractshttp "github.com/goravel/framework/contracts/http"
	ginx "github.com/goravel/gin"
	"github.com/goravel/gin/facades"
)

// Clean example request/response types
type CleanGetUserRequest struct {
	ID string `uri:"id" binding:"required" json:"id" description:"User ID"`
}

type CleanCreateUserRequest struct {
	Name  string `json:"name" binding:"required" description:"User's full name"`
	Email string `json:"email" binding:"required,email" description:"User's email address"`
}

type CleanUserResponse struct {
	ID    string `json:"id" description:"User ID"`
	Name  string `json:"name" description:"User's full name"`
	Email string `json:"email" description:"User's email address"`
}

// Example 1: Clean API - extends existing Goravel routes with Tonic
func SetupCleanTonicRoutes() {
	openapi := facades.OpenApi()

	// Get the base route group
	userGroup := facades.Route("").Prefix("/api/users").(*ginx.Group)

	// Enable Tonic documentation - returns a TonicGroup
	tonic := userGroup.WithTonic(openapi)

	// Now use familiar Goravel route methods that automatically generate OpenAPI docs
	tonic.Get("/:id",
		ginx.BindMiddleware[CleanGetUserRequest](),
		CleanGetUserHandler,
	)

	tonic.Post("/",
		ginx.BindMiddleware[CleanCreateUserRequest](),
		CleanCreateUserHandler,
	)

	tonic.Put("/:id",
		ginx.BindMiddleware[CleanCreateUserRequest](),
		CleanUpdateUserHandler,
	)

	tonic.Delete("/:id",
		ginx.BindMiddleware[CleanGetUserRequest](),
		CleanDeleteUserHandler,
	)
}

// Example 2: Clean API with typed routes
func SetupCleanTypedRoutes() {
	openapi := facades.OpenApi()

	userGroup := facades.Route("").Prefix("/api/users").(*ginx.Group)
	tonic := userGroup.WithTonic(openapi)

	// Typed routes for better OpenAPI documentation
	ginx.GetTyped[CleanGetUserRequest, CleanUserResponse](tonic, "/:id",
		ginx.BindMiddleware[CleanGetUserRequest](),
		CleanGetUserHandler,
	)

	ginx.PostTyped[CleanCreateUserRequest, CleanUserResponse](tonic, "/",
		ginx.BindMiddleware[CleanCreateUserRequest](),
		CleanCreateUserHandler,
	)

	ginx.PutTyped[CleanCreateUserRequest, CleanUserResponse](tonic, "/:id",
		ginx.BindMiddleware[CleanCreateUserRequest](),
		CleanUpdateUserHandler,
	)

	ginx.DeleteTyped[CleanGetUserRequest, CleanUserResponse](tonic, "/:id",
		ginx.BindMiddleware[CleanGetUserRequest](),
		CleanDeleteUserHandler,
	)
}

// Example 3: Works with existing Goravel middleware
func SetupCleanWithMiddleware() {
	openapi := facades.OpenApi()

	// Use standard Goravel middleware chain
	authGroup := facades.Route("").
		Prefix("/api/v1/users").
		Middleware(CleanAuthMiddleware()).(*ginx.Group)

	// Enable Tonic and continue using familiar route methods
	tonic := authGroup.WithTonic(openapi)

	tonic.Get("/:id",
		ginx.BindMiddleware[CleanGetUserRequest](),
		CleanGetUserHandler,
	)

	tonic.Post("/",
		ginx.BindMiddleware[CleanCreateUserRequest](),
		CleanCreateUserHandler,
	)
}

// Handlers
func CleanGetUserHandler(ctx contractshttp.Context) contractshttp.Response {
	req, ok := ginx.GetRequest[CleanGetUserRequest](ctx)
	if !ok {
		return ctx.Response().Json(http.StatusBadRequest, contractshttp.Json{
			"error": "invalid request",
		})
	}

	user := CleanUserResponse{
		ID:    req.ID,
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	return ctx.Response().Json(http.StatusOK, user)
}

func CleanCreateUserHandler(ctx contractshttp.Context) contractshttp.Response {
	req, ok := ginx.GetRequest[CleanCreateUserRequest](ctx)
	if !ok {
		return ctx.Response().Json(http.StatusBadRequest, contractshttp.Json{
			"error": "invalid request",
		})
	}

	user := CleanUserResponse{
		ID:    "generated-id",
		Name:  req.Name,
		Email: req.Email,
	}

	return ctx.Response().Json(http.StatusCreated, user)
}

func CleanUpdateUserHandler(ctx contractshttp.Context) contractshttp.Response {
	idReq, _ := ginx.GetRequest[CleanGetUserRequest](ctx)

	user := CleanUserResponse{
		ID:    idReq.ID,
		Name:  "Updated Name",
		Email: "updated@example.com",
	}

	return ctx.Response().Json(http.StatusOK, user)
}

func CleanDeleteUserHandler(ctx contractshttp.Context) contractshttp.Response {
	req, ok := ginx.GetRequest[CleanGetUserRequest](ctx)
	if !ok {
		return ctx.Response().Json(http.StatusBadRequest, contractshttp.Json{
			"error": "invalid request",
		})
	}

	return ctx.Response().Json(http.StatusOK, contractshttp.Json{
		"message": "User " + req.ID + " deleted successfully",
	})
}

func CleanAuthMiddleware() contractshttp.Middleware {
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
