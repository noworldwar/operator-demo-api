package api

import (
	"bc-opp-api/internal/endpoint"
	"bc-opp-api/internal/lib"
	"bc-opp-api/internal/model"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// 取得玩家所有錢包
func GetWallet(c *gin.Context) {
	token := c.Query("token")

	// 檢查請求參數
	if k := lib.HasQueryEmpty(c, "token"); k != "" {
		lib.ErrorResponse(c, 400, "Missing parameter:"+k, nil)
		return
	}

	info := model.GetPlayerInfo(token)
	if info.PlayerID == "" {
		lib.ErrorResponse(c, 403, "Token has expired", nil)
		return
	}

	// 取得錢包餘額
	var data []model.Wallet

	data = append(data, GetBalanceByBank("Main", info))
	data = append(data, GetBalanceByBank("WE", info))
	data = append(data, GetBalanceByBank("TPG", info))

	c.JSON(200, gin.H{"data": data})
}

func GetBalanceByBank(bank string, info model.PlayerInfo) model.Wallet {
	wallet := model.Wallet{Bank: bank}

	switch strings.ToLower(bank) {
	case "main":
		player, err := model.GetPlayer(info.PlayerID)
		if err != nil {
			fmt.Println("Get Player Error:" + err.Error())
		} else if player.PlayerID != "" {
			wallet.Balance = float64(player.Balance) / 100
			wallet.Success = true
		}
	case "we":
		success, bal := endpoint.WEGetBalance(info)
		wallet.Balance = bal / 100
		wallet.Success = success
	case "tpg":
		success, bal := endpoint.TPGGetBalance(info)
		wallet.Balance = bal / 100
		wallet.Success = success
	default:
		fmt.Println("Incorrect bank:" + bank)
	}

	return wallet
}
