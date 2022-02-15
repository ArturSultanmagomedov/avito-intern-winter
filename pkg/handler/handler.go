package handler

import (
	_ "for_avito_tech_with_gin/docs"
	"for_avito_tech_with_gin/pkg"
	"for_avito_tech_with_gin/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h Handler) InitRouters() *gin.Engine {
	router := gin.New()

	api := router.Group("/api/v1", h.middleware)
	{
		api.POST("/add_funds", h.addFundsHandler)
		api.POST("/write_off_funds", h.writeOffFundsHandler)
		api.POST("/funds_transfer", h.fundsTransferHandler)
		api.GET("/get_balance", h.getBalanceHandler(&pkg.DefaultCurrencyCalculator{}))
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
