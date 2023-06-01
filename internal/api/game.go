package api

import (
	"bc-opp-api/internal/lib"
	"bc-opp-api/internal/model"
	"bc-opp-api/internal/request"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GetGameLink(c *gin.Context) {
	token := c.Query("token")
	gp := c.Query("gp")
	// gameID := c.Query("gameID")
	link := ""
	success := false

	// 檢查請求參數
	if k := lib.HasQueryEmpty(c, "token", "gp"); k != "" {
		lib.ErrorResponse(c, 400, "Missing parameter:"+k, nil)
		return
	}

	info := model.GetPlayerInfo(token)
	if info.PlayerID == "" {
		lib.ErrorResponse(c, 403, "Token has expired", nil)
		return
	}

	// 組合遊戲url
	switch strings.ToLower(gp) {
	case "we_t":
		success, link = request.WELogin(info)
		if !success {
			lib.ErrorResponse(c, 500, "Failed to WE login", nil)
			return
		}
	case "we_s":
		link = fmt.Sprintf("%s?operator=%s&token=%s", viper.GetString("gameprovider.we.game_url"), viper.GetString("gameprovider.we.s_id"), token)
	default:
		lib.ErrorResponse(c, 400, "Incorrect gp:"+gp, nil)
		return
	}

	c.JSON(200, gin.H{"link": link})
}

func GetGameList(c *gin.Context) {
	gameType := c.Query("gameType")

	// 取得遊戲
	data, err := model.GetGameList(gameType)
	if err != nil {
		lib.ErrorResponse(c, 500, "Failed to get game list", nil)
		return
	}

	c.JSON(200, gin.H{"data": data})
}
