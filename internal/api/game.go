package api

import (
	"bc-opp-api/internal/endpoint"
	"bc-opp-api/internal/lib"
	"bc-opp-api/internal/model"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const tpg_game_url = "https://stagingweblobby.triple-pg.com/we/direct2Game/202"
const we_game_url = "https://uat-web-game-fe.bpweg.com/"
const we_s_id = "hwidradminjci70"

func GetGameLink(c *gin.Context) {
	token := c.Query("token")
	gp := c.Query("gp")
	gameID := c.Query("gameID")
	link := ""
	success := false

	// 檢查請求參數
	if k := lib.HasQueryEmpty(c, "token", "gp", "gameID"); k != "" {
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
	case "tpg":
		// success, tpgToken := endpoint.TPGToken(info)
		// if !success {
		// 	lib.ErrorResponse(c, 500, "Failed to TPG login", nil)
		// 	return
		// }
		// link = fmt.Sprintf("%s?gameid=%s&lang=zh-tw&playmode=real&token=%s", tpg_game_url, gameID, tpgToken)
		link = fmt.Sprintf("%s?gameid=%s&lang=zh-tw&playmode=demo", viper.GetString("tpg_game_url"), gameID)
	case "we_t":
		success, link = endpoint.WELogin(info)
		if !success {
			lib.ErrorResponse(c, 500, "Failed to WE login", nil)
			return
		}
	case "we_s":
		link = fmt.Sprintf("%s?operator=%s&token=%s", we_game_url, we_s_id, token)
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
