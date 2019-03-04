package route

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/yaoice/gocele/pkg/controller"
	"github.com/yaoice/gocele/pkg/log"
	"github.com/yaoice/gocele/pkg/middleware"
	"os"
)

// @title Swagger gocele
// @version 1.0
// @description This is a gocele server.
// @BasePath /apis/v1
func InstallRoutes(r *gin.Engine) {
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// a ping api test
	r.GET("/ping", controller.Ping)

	// config reload
	r.Any("/-/reload", func(c *gin.Context) {
		log.Info("===== Server Stop! Cause: Config Reload. =====")
		os.Exit(1)
	})

	authMiddleware := middleware.GetAuthMiddleware()
	// Login
	r.POST("/login", authMiddleware.LoginHandler)

	// Unauthenticated
//	r.GET("/", accessible)

	rootGroup := r.Group("/api/v1")
	rootGroup.Use(authMiddleware.MiddlewareFunc())
	{
		// for test
		rootGroup.GET("/hello", middleware.HelloHandler)
		rootGroup.GET("/refresh_token", authMiddleware.RefreshHandler)

		calC := controller.NewCalController()
		rootGroup.POST("add", calC.Add)
		rootGroup.POST("mul", calC.Mul)
		rootGroup.POST("tasks", calC.GetTask)
	}
}
