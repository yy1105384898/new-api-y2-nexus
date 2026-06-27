package router

import (
	"github.com/QuantumNous/new-api/controller"
	"github.com/QuantumNous/new-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetImageRouter(router *gin.Engine) {
	imageProxyRouter := router.Group("/v1")
	imageProxyRouter.Use(middleware.CORS())
	imageProxyRouter.Use(middleware.RouteTag("relay"))
	imageProxyRouter.Use(middleware.TokenOrUserAuth())
	{
		imageProxyRouter.GET("/images/:task_id/content", controller.ImageProxy)
	}

	imageV1Router := router.Group("/v1")
	imageV1Router.Use(middleware.CORS())
	imageV1Router.Use(middleware.RouteTag("relay"))
	imageV1Router.Use(middleware.TokenAuth(), middleware.PublicModelName(), middleware.Distribute())
	{
		imageV1Router.GET("/images/generations/:task_id", controller.RelayImageTaskFetch)
		imageV1Router.GET("/images/edits/:task_id", controller.RelayImageTaskFetch)
	}
}
