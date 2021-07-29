package app

import (
	"bc-opp-api/internal/api"
	"bc-opp-api/internal/lib"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var APIServer http.Server

func InitRouter() {
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(lib.APILogger())

	r.POST("/login", api.Login)
	r.POST("/player", api.CreatePlayer)
	r.PUT("/password", api.ChangePassword)
	r.GET("/player/wallet", api.GetWallet)

	r.GET("/transfer", api.GetTransfer)
	r.POST("/transfer", api.CreateTransfer)

	r.GET("/gamelink", api.GetGameLink)
	r.GET("/gamelist", api.GetGameList)

	r.GET("/info", api.GetSystemInfo)

	r.NoRoute(func(c *gin.Context) { c.AbortWithStatus(400) })
	APIServer = http.Server{Handler: r, Addr: ":7901"}
}

func RunServer() {
	if err := APIServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
