package facades

import (
	"log"

	"github.com/TickLabVN/tonic/core/docs"
	"github.com/goravel/framework/facades"
)

func OpenApi() *docs.OpenApi {
	instance, err := facades.App().Make("goravel.gin.openapi")
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return instance.(*docs.OpenApi)
}
