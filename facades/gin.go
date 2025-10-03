package facades

import (
	"log"

	"github.com/tonidy/goravel-ginx"

	"github.com/goravel/framework/contracts/route"
)

func Route(driver string) route.Route {
	if gin.App == nil {
		log.Fatalln("gin.App is not initialized")
		return nil
	}

	if gin.BindingRoute == "" {
		log.Fatalln("gin.BindingRoute is empty")
		return nil
	}

	instance, err := gin.App.MakeWith(gin.BindingRoute, map[string]any{
		"driver": driver,
	})

	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return instance.(route.Route)
}
