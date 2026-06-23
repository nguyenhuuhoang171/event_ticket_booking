package server

import (
	"event_ticket_booking/config"
	authHandler "event_ticket_booking/internal/domain/auth/handler"
	bookingHandler "event_ticket_booking/internal/domain/booking/handler"
	eventHandler "event_ticket_booking/internal/domain/event/handler"
	"event_ticket_booking/internal/domain/ping"
	"event_ticket_booking/middleware"
	commonModel "event_ticket_booking/model"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config, lib commonModel.Lib) *gin.Engine {
	router := gin.New()

	// Use gin's default middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.SecurityHeaders())
	router.Use(cors.Default())

	// Configure trusted proxies
	configureTrustedProxies(router, cfg)

	setDomainRoute(router, cfg, lib)

	return router
}

func setDomainRoute(router *gin.Engine, cfg config.Config, lib commonModel.Lib) {
	// ping
	router.GET("/ping", ping.Ping)

	// auth
	initAuthHandler := authHandler.NewHandler(cfg, lib)
	router.POST("/signup", initAuthHandler.Signup)
	router.POST("/login", initAuthHandler.Login)
	router.POST("/logout", initAuthHandler.Logout)

	// event (requires authentication)
	initEventHandler := eventHandler.NewHandler(cfg, lib)
	eventGroup := router.Group("/events", middleware.Authorize(cfg.Authentication.AccessSecret, lib.Redis))
	{
		eventGroup.POST("", initEventHandler.Create)
		eventGroup.GET("", initEventHandler.List)
		eventGroup.GET("/:id", initEventHandler.GetByID)
		eventGroup.GET("/:id/stats", initEventHandler.Stats)
		eventGroup.PUT("/:id", initEventHandler.Update)
		eventGroup.DELETE("/:id", initEventHandler.Delete)
	}

	// booking (requires authentication)
	initBookingHandler := bookingHandler.NewHandler(cfg, lib)
	bookingGroup := router.Group("/bookings", middleware.Authorize(cfg.Authentication.AccessSecret, lib.Redis))
	{
		bookingGroup.POST("", initBookingHandler.Create)
		bookingGroup.GET("", initBookingHandler.List)
		bookingGroup.POST("/:id/cancel", initBookingHandler.Cancel)
	}
}

func configureTrustedProxies(router *gin.Engine, cfg config.Config) {
	if len(cfg.Server.TrustedProxies) > 0 {
		router.SetTrustedProxies(cfg.Server.TrustedProxies)
	} else {
		router.SetTrustedProxies([]string{
			"127.0.0.1",      // localhost IPv4
			"::1",            // localhost IPv6
			"10.0.0.0/8",     // private network
			"172.16.0.0/12",  // private network
			"192.168.0.0/16", // private network
		})
	}
}
