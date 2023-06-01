package lib

import "github.com/spf13/viper"

func CheckBank(val string) bool {
	banks := []string{"Main", "WE"}

	for _, bank := range banks {
		if bank == val {
			return true
		}
	}

	return false
}

func CheckWEAppSecret(operatorID, appSecret string) bool {
	if operatorID == viper.GetString("gameprovider.we.s_id") && appSecret == viper.GetString("gameprovider.we.s_appsecret") {
		return false
	} else {
		return true
	}
}
