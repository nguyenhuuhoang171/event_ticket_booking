package server

import (
	"event_ticket_booking/config"
	authHandler "event_ticket_booking/internal/domain/auth/handler"
	"event_ticket_booking/internal/domain/ping"
	"event_ticket_booking/middleware"
	commonModel "event_ticket_booking/model"
	"os"
	"strings"

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
	configureTrustedProxies(router)

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
	router.POST("/refresh-token", initAuthHandler.RefreshToken)
}

func configureTrustedProxies(router *gin.Engine) {
	trustedProxies := os.Getenv("TRUSTED_PROXIES")

	if trustedProxies != "" {
		proxyList := strings.Split(trustedProxies, ",")
		for i, proxy := range proxyList {
			proxyList[i] = strings.TrimSpace(proxy)
		}
		router.SetTrustedProxies(proxyList)
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
