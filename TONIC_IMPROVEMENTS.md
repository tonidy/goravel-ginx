# Tonic Integration Improvements

## Summary

Enhanced the Goravel Gin Tonic integration to provide a cleaner, more intuitive API that extends existing Goravel routes rather than requiring separate method calls.

## Key Changes

### 1. New `WithTonic()` Method

Extends existing Goravel route groups with Tonic documentation support:

```go
// Before: Required special Tonic methods
userGroup.TonicGet(openapi, "/:id", handlers...)

// After: Clean extension of existing routes
tonic := userGroup.WithTonic(openapi)
tonic.Get("/:id", handlers...)  // Familiar Goravel API!
```

### 2. TonicGroup Type

New `TonicGroup` type wraps `Group` and provides:
- All standard HTTP methods (`Get`, `Post`, `Put`, `Delete`, `Patch`, `Options`)
- Automatic OpenAPI documentation generation
- Full compatibility with Goravel middleware chains
- Embedded `Group` for seamless integration

### 3. Typed Route Support

Added typed route functions for explicit request/response documentation:

```go
ginx.GetTyped[GetUserRequest, UserResponse](tonic, "/:id", handlers...)
ginx.PostTyped[CreateUserRequest, UserResponse](tonic, "/", handlers...)
```

### 4. Fixed Handler Signatures

Updated all handlers and middleware to match Goravel's `HandlerFunc` signature:

```go
// HandlerFunc must return Response
func Handler(ctx contractshttp.Context) contractshttp.Response {
    return ctx.Response().Json(200, data)
}

// Middleware returns nil to continue, or Response to abort
func Middleware() contractshttp.HandlerFunc {
    return func(ctx contractshttp.Context) contractshttp.Response {
        if !authorized {
            return ctx.Response().Json(401, error)
        }
        ctx.Request().Next()
        return nil
    }
}
```

## Migration Guide

### Option 1: Use New Clean API (Recommended)

```go
// Old approach
userGroup := facades.Route("").Prefix("/users").(*ginx.Group)
userGroup.TonicGet(openapi, "/:id", handlers...)

// New clean approach
userGroup := facades.Route("").Prefix("/users").(*ginx.Group)
tonic := userGroup.WithTonic(openapi)
tonic.Get("/:id", handlers...)
```

### Option 2: Keep Existing Code

The old `TonicGet`, `TonicPost`, etc. methods still work for backward compatibility.

## Benefits

1. **Familiar API**: Uses standard Goravel route methods you already know
2. **Less Boilerplate**: No need to pass `openapi` to every route call
3. **Middleware Support**: Works seamlessly with existing Goravel middleware
4. **Type Safety**: Optional typed routes for better documentation
5. **Clean Separation**: Clear distinction between regular routes and documented routes

## Examples

See:
- `/examples/tonic_clean_example.go` - New clean API examples
- `/examples/tonic_usage_example.go` - Updated with correct signatures
- `TONIC_INTEGRATION.md` - Updated documentation

## Breaking Changes

### Handler Signatures

All handlers must now return `contractshttp.Response`:

```go
// Before
func Handler(ctx contractshttp.Context) {
    ctx.Response().Json(200, data)
}

// After
func Handler(ctx contractshttp.Context) contractshttp.Response {
    return ctx.Response().Json(200, data)
}
```

### Middleware Signatures

Middleware must return `Response` (nil to continue):

```go
// Before
func Middleware() contractshttp.HandlerFunc {
    return func(ctx contractshttp.Context) {
        ctx.Request().Next()
    }
}

// After
func Middleware() contractshttp.HandlerFunc {
    return func(ctx contractshttp.Context) contractshttp.Response {
        ctx.Request().Next()
        return nil
    }
}
```
