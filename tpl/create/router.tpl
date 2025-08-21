package router

import (
	"github.com/gin-gonic/gin"
    gloableRouter "{{ .ProjectName }}/common/base/router"
    "{{ .ProjectName }}/pkg/middleware"
)

func Init{{ .StructName }}Router(
	r *gin.RouterGroup,
	allRouter gloableRouter.Routers,
) {
	// No route group has permission
	// noAuthRouter := r.Group("/")
	{
		// noAuthRouter.GET("/{{ .StructNameSnakeCase }}", allRouter.{{ .StructName }}Handler.Get{{ .StructName }})
	}

	// Non-strict permission routing group
    noStrictAuthRouter := r.Group("/").Use(middleware.NoStrictAuth(allRouter.JWT, allRouter.Logger))
    {
		noStrictAuthRouter.GET("/{{ .StructNameSnakeCase }}", allRouter.{{ .StructName }}Handler.Get{{ .StructName }})
    }

    // Strict permission routing group
    // strictAuthRouter := r.Group("/").Use(middleware.StrictAuth(allRouter.JWT, allRouter.Logger))
    {
		// strictAuthRouter.GET("/{{ .StructNameSnakeCase }}", allRouter.{{ .StructName }}Handler.Get{{ .StructName }})
    }
}