package rest

import (
	"fmt"
	"net/http"
	"order/internal/ports/rest/swagger/docs"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title hr-report backend
// @version 1.0
// @description hr-report backend
// @license.name kanya384

// @contact.name API Support
// @contact.email kanya384@mail.ru

// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (d *Rest) initRouter(port int) *http.Server {

	var router = gin.Default()

	d.routerDocs(router.Group("/docs"))
	router.POST("/api/order", d.CreateOrder)
	router.GET("/api/order/:id", d.ReadOrderById)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	return srv
}

func (d *Rest) routerDocs(router *gin.RouterGroup) {
	docs.SwaggerInfo.BasePath = "/"

	router.Any("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
