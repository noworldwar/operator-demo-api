package endpoint

import (
	"bc-opp-api/internal/lib"
	"bc-opp-api/internal/model"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const we_id = "16e5shi0"
const we_url = "https://op-api.aeweg.com"
const we_secret = "4XMdW16LWRMkKn_BnSbSC-q9BMsLSCQISFN06ANsK18="

func WELogin(info model.PlayerInfo) (bool, string) {
	data := url.Values{
		"operatorID": {we_id},
		"appSecret":  {we_secret},
		"playerID":   {info.PlayerID},
	}

	status, result := CallWEAPI("login", data)

	if status != 200 {
		return false, ""
	}

	var m map[string]interface{}
	err := json.Unmarshal(result, &m)
	if err != nil {
		fmt.Println("WE response format error")
		return false, ""
	}

	gameURL, ok := m["url"].(string)
	if !ok {
		fmt.Println("WE response balance format error")
		return false, ""
	}

	return true, gameURL
}

func WEGetBalance(info model.PlayerInfo) (bool, float64) {
	data := url.Values{
		"operatorID": {we_id},
		"appSecret":  {we_secret},
		"playerID":   {info.PlayerID},
	}

	status, result := CallWEAPI("balance", data)

	if status == 404 {
		return WECreatePlayer(info), 0
	}

	if status != 200 {
		return false, 0
	}

	var m map[string]interface{}
	err := json.Unmarshal(result, &m)
	if err != nil {
		fmt.Println("WE response format error")
		return false, 0
	}

	balance, ok := m["balance"].(float64)
	if !ok {
		fmt.Println("WE response balance format error")
		return false, 0
	}

	return true, balance
}

func WECreatePlayer(info model.PlayerInfo) bool {
	data := url.Values{
		"operatorID": {we_id},
		"appSecret":  {we_secret},
		"playerID":   {info.PlayerID},
		"nickname":   {info.Nickname},
	}

	status, _ := CallWEAPI("create", data)
	return status == 200
}

func WEDeposit(playerID, uid string, amount int64) (bool, float64) {
	data := url.Values{
		"operatorID": {we_id},
		"appSecret":  {we_secret},
		"playerID":   {playerID},
		"uid":        {uid},
		"amount":     {fmt.Sprintf("%d", amount)},
	}

	status, result := CallWEAPI("deposit", data)
	if status != 200 {
		return false, 0
	}

	var m map[string]interface{}
	err := json.Unmarshal(result, &m)
	if err != nil {
		fmt.Println("WE response format error")
		return false, 0
	}

	balance, ok := m["balance"].(float64)
	if !ok {
		fmt.Println("WE response balance format error")
		return false, 0
	}

	return true, balance
}

func WEWithdraw(playerID, uid string, amount int64) (bool, float64) {
	data := url.Values{
		"operatorID": {we_id},
		"appSecret":  {we_secret},
		"playerID":   {playerID},
		"uid":        {uid},
		"amount":     {fmt.Sprintf("%d", amount)},
	}

	status, result := CallWEAPI("withdraw", data)
	if status != 200 {
		return false, 0
	}

	var m map[string]interface{}
	err := json.Unmarshal(result, &m)
	if err != nil {
		fmt.Println("WE response format error")
		return false, 0
	}

	balance, ok := m["balance"].(float64)
	if !ok {
		fmt.Println("WE response balance format error")
		return false, 0
	}

	return true, balance
}

func CallWEAPI(funcName string, data url.Values) (int, []byte) {
	url := we_url
	status := 999
	errmsg := ""
	var result []byte

	switch funcName {
	case "login":
		url += "/player/login"
	case "create":
		url += "/player/create"
	case "balance":
		url += "/player/balance"
	case "deposit":
		url += "/player/deposit"
	case "withdraw":
		url += "/player/withdraw"
	}

	singB64 := base64.StdEncoding.EncodeToString([]byte(we_secret + we_id + data.Get("playerID")))

	req, _ := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("signature", singB64)

	clt := http.Client{}
	rsp, err := clt.Do(req)
	if err != nil {
		errmsg = err.Error()
	} else {
		status = rsp.StatusCode
		rspBody, _ := ioutil.ReadAll(rsp.Body)
		errmsg = string(rspBody)
		if rsp.StatusCode == 200 {
			result = rspBody
		}
	}

	msg := "------------------------------------------------------------\r\n"
	msg += fmt.Sprintf("[%s] \r\n\r\n", time.Now().Format("2006/01/02 15:04:05"))
	msg += fmt.Sprintf("[Request] \r\nPOST %s\r\n\r\n", url)
	msg += fmt.Sprintf("[Signature] \r\n%s\r\n\r\n", singB64)
	msg += fmt.Sprintf("[Body] \r\n%v\r\n\r\n", data)
	msg += fmt.Sprintf("[Status] \r\n%v\r\n\r\n", status)
	msg += fmt.Sprintf("[Response Data] \r\n%s\r\n\r\n", strings.TrimRight(errmsg, "\n"))
	msg += "------------------------------------------------------------\r\n"

	go lib.WriteLog("we_", msg)

	return status, result
}
