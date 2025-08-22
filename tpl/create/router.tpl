package router

import (
	"github.com/gin-gonic/gin"
    gloableRouter "{{ .ProjectName }}/common/base/router"
)

func Init{{ .StructName }}Router(
    publicRouter gin.IRoutes,
    privateRouter gin.IRoutes,
	allRouter gloableRouter.Routers,
) {
	{
		publicRouter.GET("/{{ .StructNameSnakeCase }}", allRouter.{{ .StructName }}Handler.Get{{ .StructName }})
	}

    {
    }
}