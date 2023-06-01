package endpoint

import (
	"bc-opp-api/internal/lib"
	"bc-opp-api/internal/model"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type creditDatail struct {
	OperatorID string `json:"operatorID"`
	AppSecret  string `json:"appSecret"`
	PlayerID   string `json:"playerID"`
	GameID     string `json:"gameID"`
	BetID      string `json:"betID"`
	Amount     string `json:"amount"`
	GameStatus string `json:"gameStatus"`
	GameResult string `json:"gameResult"`
	Currency   string `json:"currency"`
	Type       string `json:"type"`
	Time       string `json:"time"`
}

func Validate(c *gin.Context) {
	token := c.PostForm("token")
	appSecret := c.PostForm("appSecret")
	operatorID := c.PostForm("operatorID")

	// Step 1: Check the required parameters
	if k := lib.HasPostFormEmpty(c, "token", "appSecret", "operatorID"); k != "" {
		lib.ErrorResponse(c, 400, "Missing parameter:"+k, nil)
		return
	}

	var player_data model.Player
	var err error

	// Step 2: Check AppSecret
	if lib.CheckWEAppSecret(operatorID, appSecret) {
		lib.ErrorResponse(c, 401, "Incorrect appSecret", nil)
		return
	}

	// Step 3: Check Token
	player_info := model.GetPlayerInfo(token)
	player_data, err = model.GetPlayer(player_info.PlayerID)
	if err != nil || player_data.PlayerID == "" {
		lib.ErrorResponse(c, 404, "Token has expired", err)
		return
	}

	c.JSON(200, gin.H{"playerID": player_data.PlayerID, "nickname": player_data.Nickname, "currency": player_data.Currency, "test": lib.IntToBool(player_data.Test), "time": player_data.Created})
}

func GetBalance(c *gin.Context) {
	token := c.PostForm("token")
	appSecret := c.PostForm("appSecret")
	operatorID := c.PostForm("operatorID")
	playerID := c.PostForm("playerID")

	// Step 1: Check the required parameters
	if k := lib.HasPostFormEmpty(c, "token", "appSecret", "operatorID", "playerID"); k != "" {
		lib.ErrorResponse(c, 400, "Missing parameter:"+k, nil)
		return
	}

	// Step 2: Check AppSecret
	if lib.CheckWEAppSecret(operatorID, appSecret) {
		lib.ErrorResponse(c, 401, "Incorrect appSecret", nil)
		return
	}

	// Step 3: Check Token
	player_info := model.GetPlayerInfo(token)
	player_data, err := model.GetPlayer(player_info.PlayerID)
	if err != nil || player_data.PlayerID == "" {
		lib.ErrorResponse(c, 404, "Token has expired", err)
		return
	} else if player_data.PlayerID != playerID {
		lib.ErrorResponse(c, 400, "Incorrect playerID", nil)
		return
	}

	c.JSON(200, gin.H{"balance": player_data.Balance, "currency": player_data.Currency, "time": player_data.Created})
}

func Debit(c *gin.Context) {
	token := c.PostForm("token")
	operatorID := c.PostForm("operatorID")
	appSecret := c.PostForm("appSecret")
	playerID := c.PostForm("playerID")
	gameID := c.PostForm("gameID")
	betID := c.PostForm("betID")
	amount := c.PostForm("amount")
	currency := c.PostForm("currency")
	we_tran_type := c.PostForm("type")
	we_time := c.PostForm("time")

	// Step 1: Check the required parameters
	if k := lib.HasPostFormEmpty(c, "token", "amount", "appSecret", "betID", "operatorID", "playerID", "gameID", "time", "type"); k != "" {
		lib.ErrorResponse(c, 400, "Missing parameter:"+k, nil)
		return
	}

	// Step 2: Check AppSecret
	if lib.CheckWEAppSecret(operatorID, appSecret) {
		lib.ErrorResponse(c, 401, "Incorrect appSecret", nil)
		return
	}

	// Step 3: Check Token
	player_info := model.GetPlayerInfo(token)
	player_data, err := model.GetPlayer(player_info.PlayerID)
	if err != nil || player_data.PlayerID == "" {
		lib.ErrorResponse(c, 404, "Token has expired", err)
		return
	} else if player_data.PlayerID != playerID {
		lib.ErrorResponse(c, 400, "Incorrect playerID", nil)
		return
	}

	// Step 4: Check Amount
	amount_int, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		lib.ErrorResponse(c, 400, "Incorrect amount format:"+amount, err)
		return
	} else if amount_int <= 0 {
		lib.ErrorResponse(c, 400, "Incorrect amount format:"+amount, nil)
		return
	} else if amount_int > player_data.Balance {
		lib.ErrorResponse(c, 402, "Insufficient balance", nil)
		return
	}

	// Step 5: Check Duplicate transaction
	isExist, err := model.CheckIfBetTransferExist(betID, "Debit")
	if isExist {
		lib.ErrorResponse(c, 409, "Duplicate transaction: ", err)
		return
	}

	// Step 6: Add Transfer
	refID := time.Now().Format("20060102") + xid.New().String()
	transfer := model.BetTransfer{
		TransferID: refID,
		PlayerID:   playerID,
		Type:       "Debit",
		BetID:      betID,
		GameID:     gameID,
		WeType:     we_tran_type,
		WeTime:     we_time,
		Amount:     amount_int,
		Success:    true,
		Created:    time.Now().Unix(),
		Updated:    time.Now().Unix(),
	}
	AddErr := model.AddBetTransfer(transfer)
	if AddErr != nil {
		lib.ErrorResponse(c, 500, "Internal Server Error: ", AddErr)
		return
	}

	// Step 7: Update Balance
	balance, _ := model.UpdateBalance(playerID, -amount_int)

	c.JSON(200, gin.H{"balance": balance, "currency": currency, "time": time.Now().Unix(), "refID": refID})
}

func Credit(c *gin.Context) {
	token := c.PostForm("token")
	refID := ""
	// Step 1: Unmarshal The JSON Data
	data := c.PostForm("data")
	balance := int64(0)
	creditDatail := []creditDatail{}

	err := json.Unmarshal([]byte(data), &creditDatail)
	if err != nil {
		lib.ErrorResponse(c, 400, "JSON format failed: "+err.Error(), nil)
	}

	// Step 2: Read Credit Detail Data
	for _, v := range creditDatail {

		// Step 2.1: Check AppSecret
		if lib.CheckWEAppSecret(v.OperatorID, v.AppSecret) {
			lib.ErrorResponse(c, 401, "Incorrect appSecret: operatorID:"+v.OperatorID, nil)
			return
		}

		// Step 2.2: Check Token
		if token != "" {
			player_info := model.GetPlayerInfo(token)
			player_data, err := model.GetPlayer(player_info.PlayerID)
			if err != nil || player_data.PlayerID == "" {
				lib.ErrorResponse(c, 404, "Token has expired", err)
				return
			} else if player_data.PlayerID != v.PlayerID {
				lib.ErrorResponse(c, 400, "Incorrect playerID: playerID:"+v.PlayerID, nil)
				return
			}
		} else {
			player_data, err := model.GetPlayer(v.PlayerID)
			if err != nil || player_data.PlayerID == "" {
				lib.ErrorResponse(c, 400, "Incorrect playerID: playerID:"+v.PlayerID, nil)
				return
			}
		}

		// Step 2.3: Check Amount
		amount_int, err := strconv.ParseInt(v.Amount, 10, 64)
		if err != nil {
			lib.ErrorResponse(c, 400, "Incorrect amount format: "+v.Amount, err)
			return
		}

		// Step 2.4: Check Debit Record
		debitRecord, err := model.GetBetTransferByBetID(v.BetID, "Debit")
		if err != nil || debitRecord.PlayerID == "" {
			lib.ErrorResponse(c, 410, "Failed to read bet transaction: ", err)
			return
		}

		// Step 2.5: Check GameID
		if strings.Trim(debitRecord.GameID, " ") != strings.Trim(v.GameID, " ") {
			lib.ErrorResponse(c, 400, "Incorrect gameID: ", err)
			return
		}

		// Step 2.6: Check Transaction Already Paid
		isExist, err := model.CheckIfBetTransferExist(v.BetID, "Credit")
		if isExist {
			lib.ErrorResponse(c, 409, "Duplicate transaction: betID already paid ", err)
			return
		}

		// Step 2.7: Check Duplicate Transaction
		isRollbacked, err := model.CheckIfBetTransferExist(v.BetID, "Rollback")
		if isRollbacked {
			lib.ErrorResponse(c, 409, "Duplicate transaction: betID already rollbacked ", err)
			return
		}

		// Step 9: Add Transfer
		refID = time.Now().Format("20060102") + xid.New().String()
		transfer := model.BetTransfer{TransferID: refID,
			PlayerID: v.PlayerID,
			Type:     "Credit",
			BetID:    v.BetID,
			GameID:   v.GameID,
			WeType:   v.Type,
			WeTime:   v.Time,
			Amount:   amount_int,
			Success:  true,
			Created:  time.Now().Unix(),
			Updated:  time.Now().Unix(),
		}
		err = model.AddBetTransfer(transfer)
		if err != nil {
			lib.ErrorResponse(c, 500, "Internal Server Error: ", err)
			return
		}

		// Step 10: Update Balance
		balance, _ = model.UpdateBalance(v.PlayerID, amount_int)

	}

	c.JSON(200, gin.H{"balance": balance, "currency": creditDatail[0].Currency, "time": time.Now().Unix(), "refID": refID})
}

func Rollback(c *gin.Context) {
	operatorID := c.PostForm("operatorID")
	appSecret := c.PostForm("appSecret")
	playerID := c.PostForm("playerID")
	gameID := c.PostForm("gameID")
	betID := c.PostForm("betID")
	amount := c.PostForm("amount")
	currency := c.PostForm("currency")
	we_tran_type := c.PostForm("type")
	we_time := c.PostForm("time")

	// Step 1: Check the required parameters
	if k := lib.HasPostFormEmpty(c, "amount", "appSecret", "betID", "operatorID", "playerID", "gameID", "time", "type"); k != "" {
		lib.ErrorResponse(c, 400, "Missing parameter:"+k, nil)
		return
	}

	// Step 2: Check AppSecret
	if lib.CheckWEAppSecret(operatorID, appSecret) {
		lib.ErrorResponse(c, 401, "Incorrect appSecret", nil)
		return
	}

	// Step 3: Check Amount
	amount_int, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		lib.ErrorResponse(c, 400, "Incorrect amount format:"+amount, err)
		return
	}

	// Step 4: Check Debit Record
	debitRecord, err := model.GetBetTransferByBetID(betID, "Debit")
	if err != nil || debitRecord.PlayerID == "" {
		lib.ErrorResponse(c, 410, "Failed to read bet transaction: ", err)
		return
	}

	// Step 5: Check GameID
	if strings.Trim(debitRecord.GameID, " ") != strings.Trim(gameID, " ") {
		lib.ErrorResponse(c, 400, "Incorrect gameID: ", err)
		return
	}

	// Step 6: Check Duplicate transaction
	isExist, err := model.CheckIfBetTransferExist(betID, "Rollback")
	if isExist {
		lib.ErrorResponse(c, 409, "Duplicate transaction: betID already rollbacked", err)
		return
	}

	// Step 7: Get Credit Amount
	creditRecord, _ := model.GetBetTransferByBetID(betID, "Credit")

	// Step 7: Check Rollback Amount
	if amount_int > 0 && amount_int != (debitRecord.Amount-creditRecord.Amount) {
		lib.ErrorResponse(c, 409, "Incorrect Rollback Amount: ", err)
		return
	}

	// Step 8: Add Transfer
	refID := time.Now().Format("20060102") + xid.New().String()
	transfer := model.BetTransfer{TransferID: refID,
		PlayerID: playerID,
		Type:     "Rollback",
		BetID:    betID,
		GameID:   gameID,
		WeType:   we_tran_type,
		WeTime:   we_time,
		Amount:   amount_int,
		Success:  true,
		Created:  time.Now().Unix(),
		Updated:  time.Now().Unix(),
	}
	err = model.AddBetTransfer(transfer)
	if err != nil {
		lib.ErrorResponse(c, 500, "Internal Server Error: ", err)
		return
	}

	// Step 9: Update Balance
	balance, _ := model.UpdateBalance(playerID, amount_int)

	c.JSON(200, gin.H{"balance": balance, "currency": currency, "time": time.Now().Unix(), "refID": refID})
}
