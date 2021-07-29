package api

import (
	"bc-opp-api/internal/lib"
	"bc-opp-api/internal/model"
	"strings"

	"github.com/gin-gonic/gin"
)

// 玩家登入
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 檢查請求參數
	if k := lib.HasPostFormEmpty(c, "username", "password"); k != "" {
		lib.ErrorResponse(c, 400, "Missing parameter:"+k, nil)
		return
	}

	// 取得DB玩家資料
	player, err := model.GetPlayer(username)
	if err != nil {
		lib.ErrorResponse(c, 500, "Failed to load player", err)
		return
	} else if player.PlayerID == "" {
		lib.ErrorResponse(c, 404, "Incorrect username", nil)
		return
	} else if player.Disabled {
		lib.ErrorResponse(c, 403, "Player is disabled", nil)
		return
	} else if err := lib.ComparePassword(player.Password, password); err != nil {
		lib.ErrorResponse(c, 401, "Incorrect password", err)
		return
	}

	// 寫入Redis
	token := model.SetPlayerInfo(player)

	c.JSON(200, gin.H{
		"token":    token,
		"nickname": player.Nickname,
		"balance":  float64(player.Balance) / 100,
	})
}

// 創建玩家
func CreatePlayer(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	nickname := c.PostForm("nickname")

	// 檢查請求參數
	if k := lib.HasPostFormEmpty(c, "username", "password", "nickname"); k != "" {
		lib.ErrorResponse(c, 400, "Missing parameter:"+k, nil)
		return
	}

	if len(username) > 30 {
		lib.ErrorResponse(c, 400, "username maximum is 30 characters", nil)
		return
	}

	if len(password) > 20 {
		lib.ErrorResponse(c, 400, "password maximum is 20 characters", nil)
		return
	}

	if len(nickname) > 30 {
		lib.ErrorResponse(c, 400, "nickname maximum is 30 characters", nil)
		return
	}

	passwordHash, err := lib.GeneratePasswordHash(password)
	if err != nil {
		lib.ErrorResponse(c, 400, "Incorrect password", err)
		return
	}

	// 新增玩家到DB
	player := model.Player{PlayerID: username, Password: passwordHash, Nickname: nickname, Balance: 500000, Disabled: false}
	err = model.AddPlayer(player)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			lib.ErrorResponse(c, 409, "Username already exists", err)
			return
		}

		lib.ErrorResponse(c, 500, "Failed to create player", err)
		return
	}

	// 寫入Redis
	token := model.SetPlayerInfo(player)

	c.JSON(200, gin.H{
		"token":    token,
		"nickname": player.Nickname,
		"balance":  float64(player.Balance) / 100,
	})
}

// 變更密碼
func ChangePassword(c *gin.Context) {
	token := c.PostForm("token")
	password := c.PostForm("password")

	// 檢查請求參數
	if k := lib.HasPostFormEmpty(c, "token", "password"); k != "" {
		lib.ErrorResponse(c, 400, "Missing parameter:"+k, nil)
		return
	}

	passwordHash, err := lib.GeneratePasswordHash(password)
	if err != nil {
		lib.ErrorResponse(c, 400, "Incorrect password", err)
		return
	}

	info := model.GetPlayerInfo(token)
	if info.PlayerID == "" {
		lib.ErrorResponse(c, 403, "Token has expired", nil)
		return
	}

	player, err := model.GetPlayer(info.PlayerID)
	if err != nil {
		lib.ErrorResponse(c, 500, "Failed to load player", err)
		return
	}

	player.Password = passwordHash

	err = model.UpdatePlayer(player)
	if err != nil {
		lib.ErrorResponse(c, 500, "Failed to update player", err)
		return
	}

	c.JSON(200, gin.H{
		"token":    token,
		"nickname": player.Nickname,
	})
}
