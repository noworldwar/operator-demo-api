package request

import (
	"bc-opp-api/internal/lib"
	"bc-opp-api/internal/model"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func WELogin(info model.PlayerInfo) (bool, string) {
	requestTime := strconv.FormatInt(time.Now().Unix(), 10)
	data := url.Values{
		"operatorID":  {viper.GetString("gameprovider.we.t_id")},
		"appSecret":   {viper.GetString("gameprovider.we.appsecret")},
		"playerID":    {info.PlayerID},
		"requestTime": {requestTime},
	}
	singMD5 := MD5encode(viper.GetString("gameprovider.we.appsecret"), viper.GetString("gameprovider.we.t_id"), info.PlayerID, requestTime)
	status, result := CallWEAPI("login", singMD5, data)

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
	requestTime := strconv.FormatInt(time.Now().Unix(), 10)
	data := url.Values{
		"operatorID":  {viper.GetString("gameprovider.we.t_id")},
		"appSecret":   {viper.GetString("gameprovider.we.appsecret")},
		"playerID":    {info.PlayerID},
		"requestTime": {requestTime},
	}
	singMD5 := MD5encode(viper.GetString("gameprovider.we.appsecret"), viper.GetString("gameprovider.we.t_id"), info.PlayerID, requestTime)
	status, result := CallWEAPI("balance", singMD5, data)

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
	requestTime := strconv.FormatInt(time.Now().Unix(), 10)
	data := url.Values{
		"operatorID":  {viper.GetString("gameprovider.we.t_id")},
		"appSecret":   {viper.GetString("gameprovider.we.appsecret")},
		"playerID":    {info.PlayerID},
		"nickname":    {info.Nickname},
		"requestTime": {requestTime},
	}
	singMD5 := MD5encode(viper.GetString("gameprovider.we.appsecret"), info.Nickname, viper.GetString("gameprovider.we.t_id"), info.PlayerID, requestTime)
	status, _ := CallWEAPI("create", singMD5, data)
	return status == 200
}

func WEDeposit(playerID, uid string, amount int64) (bool, float64) {
	requestTime := strconv.FormatInt(time.Now().Unix(), 10)
	data := url.Values{
		"operatorID":  {viper.GetString("gameprovider.we.t_id")},
		"appSecret":   {viper.GetString("gameprovider.we.appsecret")},
		"playerID":    {playerID},
		"uid":         {uid},
		"amount":      {fmt.Sprintf("%d", amount)},
		"requestTime": {requestTime},
	}

	singMD5 := MD5encode(fmt.Sprintf("%d", amount), viper.GetString("gameprovider.we.appsecret"), viper.GetString("gameprovider.we.t_id"), playerID, requestTime, uid)
	status, result := CallWEAPI("deposit", singMD5, data)
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
	requestTime := strconv.FormatInt(time.Now().Unix(), 10)
	data := url.Values{
		"operatorID":  {viper.GetString("gameprovider.we.t_id")},
		"appSecret":   {viper.GetString("gameprovider.we.appsecret")},
		"playerID":    {playerID},
		"uid":         {uid},
		"amount":      {fmt.Sprintf("%d", amount)},
		"requestTime": {requestTime},
	}
	singMD5 := MD5encode(fmt.Sprintf("%d", amount), viper.GetString("gameprovider.we.appsecret"), viper.GetString("gameprovider.we.t_id"), playerID, requestTime, uid)
	status, result := CallWEAPI("withdraw", singMD5, data)
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

func CallWEAPI(funcName, singMD5 string, data url.Values) (int, []byte) {
	url := viper.GetString("gameprovider.we.api")
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

	// singMD5 := base64.StdEncoding.EncodeToString([]byte(we_secret + we_id + data.Get("playerID")))

	req, _ := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("signature", singMD5)

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
	msg += fmt.Sprintf("[Signature] \r\n%s\r\n\r\n", singMD5)
	msg += fmt.Sprintf("[Body] \r\n%v\r\n\r\n", data)
	msg += fmt.Sprintf("[Status] \r\n%v\r\n\r\n", status)
	msg += fmt.Sprintf("[Response Data] \r\n%s\r\n\r\n", strings.TrimRight(errmsg, "\n"))
	msg += "------------------------------------------------------------\r\n"

	go lib.WriteLog("we_", msg)

	return status, result
}

func MD5encode(inputs ...string) string {
	str := ""
	for _, v := range inputs {
		str += v
	}
	has := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", has)
}
