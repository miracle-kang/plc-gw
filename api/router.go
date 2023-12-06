package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	docs "github.com/miracle-kang/plc-gw/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Startup(port int) error {
	router := gin.Default()

	// Config
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Router
	route(router.Group("/api/v1"))

	// Swagger
	swaggerDoc(router)

	// Run
	return router.Run(":" + strconv.Itoa(port))
}

func swaggerDoc(router *gin.Engine) {

	docs.SwaggerInfo.Title = "PLC Web Gateway Service API"
	docs.SwaggerInfo.Description = "PLC Web Gateway Service API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// Route API
func route(api *gin.RouterGroup) {

	// Report
	report := api.Group("/report")
	{
		report.POST("", ReportData)
	}

	// Control
	control := api.Group("/control")
	{
		control.POST("/readTagsValueByNameList", ReadTags)
		control.POST("/writeTagValueByName", WriteTag)
		control.POST("/clearGateway", ClearGateway)
		control.POST("/clearGatewayFile", ClearGatewayFile)
	}

	// Manager
	manager := api.Group("/manager")
	{
		manager.GET("/gateways", ListGateways)
		manager.GET("/gateways/:sn/plcs", ListGatewayPLCs)
	}
}
