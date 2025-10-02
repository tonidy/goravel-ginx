package gin

import (
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/contracts/validation"
)

const BindingRoute = "goravel.gin.route"
const BindingOpenApi = "goravel.gin.openapi"

var (
	App              foundation.Application
	ConfigFacade     config.Config
	LogFacade        log.Log
	ValidationFacade validation.Validation
	ViewFacade       http.View
	OpenApiInstance  *docs.OpenApi
)

type ServiceProvider struct{}

func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings: []string{
			BindingRoute,
			BindingOpenApi,
		},
		Dependencies: []string{
			binding.Config,
			binding.Log,
			binding.Validation,
			binding.View,
		},
		ProvideFor: []string{
			binding.Route,
		},
	}
}

func (r *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.BindWith(BindingOpenApi, func(app foundation.Application, parameters map[string]any) (any, error) {
		cfg := app.MakeConfig()

		openapi := &docs.OpenApi{
			OpenAPI: cfg.GetString("swagger.openapi_version", "3.0.1"),
			Info: docs.InfoObject{
				Version:     cfg.GetString("swagger.info.version", "1.0.0"),
				Title:       cfg.GetString("swagger.info.title", "API Documentation"),
				Description: cfg.GetString("swagger.info.description", ""),
			},
		}

		// Set contact info if configured
		if cfg.GetString("swagger.info.contact", "") != "" {
			openapi.Info.Contact = &docs.ContactObject{
				Name:  cfg.GetString("swagger.info.contact.name", ""),
				URL:   cfg.GetString("swagger.info.contact.url", ""),
				Email: cfg.GetString("swagger.info.contact.email", ""),
			}
		}

		// Set license info if configured
		if cfg.GetString("swagger.info.license", "") != "" {
			openapi.Info.License = &docs.LicenseObject{
				Name: cfg.GetString("swagger.info.license.name", ""),
				URL:  cfg.GetString("swagger.info.license.url", ""),
			}
		}

		OpenApiInstance = openapi
		return openapi, nil
	})

	app.BindWith(BindingRoute, func(app foundation.Application, parameters map[string]any) (any, error) {
		return NewRoute(app.MakeConfig(), parameters)
	})
}

func (r *ServiceProvider) Boot(app foundation.Application) {
	ConfigFacade = app.MakeConfig()
	LogFacade = app.MakeLog()
	ValidationFacade = app.MakeValidation()
	ViewFacade = app.MakeView()

	// Initialize OpenAPI instance if not already done
	if OpenApiInstance == nil {
		if instance, err := app.Make(BindingOpenApi); err == nil {
			if openapi, ok := instance.(*docs.OpenApi); ok {
				OpenApiInstance = openapi
			}
		}
	}
}
