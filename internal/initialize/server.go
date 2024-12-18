package initialize

import (
	"borrow_book/internal/config"
	"borrow_book/internal/infra/server/http"
	"borrow_book/internal/router"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitServer configures and returns the HTTP server instance.
func InitServer(config *config.Config, appRouter *router.AppRouter) *http.Server {
	router := gin.Default()

	// Apply CORS middleware with configured settings
	corsConfig := config.CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     corsConfig.AllowedOrigins,
		AllowMethods:     corsConfig.AllowedMethods,
		AllowHeaders:     corsConfig.AllowedHeaders,
		ExposeHeaders:    corsConfig.ExposeHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
		MaxAge:           corsConfig.MaxAge,
	}))

	apiGroup := router.Group("/api")
	registerAPIRoutes(apiGroup, appRouter)
	registerSwaggerRoutes(router, appRouter)

	addr := fmt.Sprintf(":%s", config.Server.Port)

	// Create and return the HTTP server
	server := http.NewServer(router, addr)
	return server
}

func registerAPIRoutes(group *gin.RouterGroup, appRouter *router.AppRouter) {
	appRouter.RegisterBookRoutes(group)
	appRouter.RegisterAuthorRoutes(group)
	appRouter.RegisterExtraBookRoutes(group)
}

func registerSwaggerRoutes(router *gin.Engine, appRouter *router.AppRouter) {
	swaggerGroup := router.Group("/")
	appRouter.RegisterSwaggerRoutes(swaggerGroup)
}
