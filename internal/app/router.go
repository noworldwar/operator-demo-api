package app

import (
	"bc-opp-api/internal/api"
	"bc-opp-api/internal/endpoint"
	"bc-opp-api/internal/lib"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

	group := r.Group("/we/api")

	group.POST("/validate", endpoint.Validate)
	group.POST("/balance", endpoint.GetBalance)
	group.POST("/debit", endpoint.Debit)
	group.POST("/credit", endpoint.Credit)
	group.POST("/rollback", endpoint.Rollback)

	r.NoRoute(func(c *gin.Context) { c.AbortWithStatus(400) })
	APIServer = http.Server{Handler: r, Addr: viper.GetString("api.port")}
}

func RunServer() {
	if err := APIServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
