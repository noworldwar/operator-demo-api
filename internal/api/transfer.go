package api

import (
	"bc-opp-api/internal/endpoint"
	"bc-opp-api/internal/lib"
	"bc-opp-api/internal/model"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// 取得最後10筆轉帳紀錄
func GetTransfer(c *gin.Context) {
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

	// 取得轉帳紀錄
	data, err := model.GetTransferBy(info.PlayerID)
	if err != nil {
		lib.ErrorResponse(c, 500, "Failed to read transfer", nil)
		return
	}

	c.JSON(200, gin.H{"data": data})
}

// 新增轉帳交易
func CreateTransfer(c *gin.Context) {
	token := c.PostForm("token")
	fromBank := c.PostForm("fromBank")
	toBank := c.PostForm("toBank")
	amount := c.PostForm("amount")

	// 檢查請求參數
	if k := lib.HasPostFormEmpty(c, "token", "fromBank", "toBank", "amount"); k != "" {
		lib.ErrorResponse(c, 400, "Missing parameter:"+k, nil)
		return
	}

	iamount, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		lib.ErrorResponse(c, 400, "Incorrect amount format:"+amount, nil)
		return
	} else if iamount <= 0 {
		lib.ErrorResponse(c, 400, "Incorrect amount format:"+amount, nil)
		return
	}

	if !lib.CheckBank(fromBank) {
		lib.ErrorResponse(c, 400, "Incorrect fromBank:"+fromBank, nil)
		return
	}

	if !lib.CheckBank(toBank) {
		lib.ErrorResponse(c, 400, "Incorrect toBank:"+toBank, nil)
		return
	}

	info := model.GetPlayerInfo(token)
	if info.PlayerID == "" {
		lib.ErrorResponse(c, 403, "Token has expired", nil)
		return
	}

	// 取得錢包餘額
	checkWallet := GetBalanceByBank(fromBank, info)
	if !checkWallet.Success {
		lib.ErrorResponse(c, 500, "Failed to get balance", nil)
		return
	} else if float64(iamount) > checkWallet.Balance {
		lib.ErrorResponse(c, 402, "Insufficient balance", nil)
		return
	}

	// 取款
	guid := time.Now().Format("20060102") + xid.New().String()
	fromWallet := Withdraw(fromBank, info.PlayerID, guid, iamount)
	if !fromWallet.Success {
		lib.ErrorResponse(c, 500, "Failed to update balance (Withdraw)", nil)
		return
	}

	// 存款
	toWallet := Deposit(toBank, info.PlayerID, guid, iamount)

	// 寫入交易紀錄到DB
	tran := model.Transfer{
		TransferID:  guid,
		PlayerID:    info.PlayerID,
		FromBank:    fromBank,
		FromBalance: fromWallet.Balance,
		ToBank:      toBank,
		ToBalance:   toWallet.Balance,
		Amount:      iamount,
		Success:     toWallet.Success,
		Created:     time.Now().Unix(),
		Updated:     time.Now().Unix(),
	}
	err = model.AddTransfer(tran)
	if err != nil {
		fmt.Println("Add Transfer data Error:", err, tran)
	}

	c.JSON(200, gin.H{
		"transferID": guid,
		"fromBank":   fromBank, "fromBalance": fromWallet.Balance,
		"toBank": toBank, "toBalance": toWallet.Balance,
		"success": toWallet.Success,
	})
}

// 取款
func Withdraw(bank, playerID, uid string, amount int64) model.Wallet {
	wallet := model.Wallet{Bank: bank}

	switch strings.ToLower(bank) {
	case "main":
		bal, err := model.UpdateBalance(playerID, -amount*100)
		if err != nil {
			fmt.Println("Update Balance Error:" + err.Error())
		} else {
			wallet.Balance = float64(bal) / 100
			wallet.Success = true
		}
	case "we":
		success, bal := endpoint.WEWithdraw(playerID, uid, amount*100)
		wallet.Balance = bal / 100
		wallet.Success = success
	case "tpg":
		success, bal := endpoint.TPGWithdraw(playerID, uid, amount*100)
		wallet.Balance = bal / 100
		wallet.Success = success
	default:
		fmt.Println("Incorrect bank:" + bank)
	}

	return wallet
}

// 存款
func Deposit(bank, playerID, uid string, amount int64) model.Wallet {
	wallet := model.Wallet{Bank: bank}

	switch strings.ToLower(bank) {
	case "main":
		bal, err := model.UpdateBalance(playerID, amount*100)
		if err != nil {
			fmt.Println("Update Balance Error:" + err.Error())
		} else {
			wallet.Balance = float64(bal) / 100
			wallet.Success = true
		}
	case "we":
		success, bal := endpoint.WEDeposit(playerID, uid, amount*100)
		wallet.Balance = bal / 100
		wallet.Success = success
	case "tpg":
		success, bal := endpoint.TPGDeposit(playerID, uid, amount*100)
		wallet.Balance = bal / 100
		wallet.Success = success
	default:
		fmt.Println("Incorrect bank:" + bank)
	}

	return wallet
}
