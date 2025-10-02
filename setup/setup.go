package main

import (
	"os"

	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/packages/match"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support/path"
)

var config = `map[string]any{
        // Optional, default is 4096 KB
        "body_limit": 4096,
        "header_limit": 4096,
        "route": func() (route.Route, error) {
            return ginxfacades.Route("ginx"), nil
        },
        // Optional, default is http/template
        "template": func() (render.HTMLRender, error) {
            return gin.DefaultTemplate()
        },
    }`

func main() {
	packages.Setup(os.Args).
		Install(
			modify.GoFile(path.Config("app.go")).
				Find(match.Imports()).Modify(modify.AddImport(packages.GetModulePath())).
				Find(match.Providers()).Modify(modify.Register("&ginx.ServiceProvider{}")),
			modify.GoFile(path.Config("http.go")).
				Find(match.Imports()).
				Modify(
					modify.AddImport("github.com/goravel/framework/contracts/route"), modify.AddImport(packages.GetModulePath()),
					modify.AddImport("github.com/tonidy/goravel-ginx/facades", "ginxfacades"), modify.AddImport("github.com/gin-gonic/gin/render"),
				).
				Find(match.Config("http.drivers")).Modify(modify.AddConfig("ginx", config)).
				Find(match.Config("http")).Modify(modify.AddConfig("default", `"ginx"`)),
		).
		Uninstall(
			modify.GoFile(path.Config("app.go")).
				Find(match.Providers()).Modify(modify.Unregister("&ginx.ServiceProvider{}")).
				Find(match.Imports()).Modify(modify.RemoveImport(packages.GetModulePath())),
			modify.GoFile(path.Config("http.go")).
				Find(match.Config("http.drivers")).Modify(modify.RemoveConfig("ginx")).
				Find(match.Config("http")).Modify(modify.AddConfig("default", `""`)).
				Find(match.Imports()).
				Modify(
					modify.RemoveImport("github.com/goravel/framework/contracts/route"), modify.RemoveImport(packages.GetModulePath()),
					modify.RemoveImport("github.com/tonidy/goravel-ginx/facades", "ginxfacades"), modify.RemoveImport("github.com/gin-gonic/gin/render"),
				),
		).
		Execute()
}
