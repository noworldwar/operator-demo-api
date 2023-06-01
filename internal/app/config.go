package app

import (
	"bc-opp-api/internal/model"
	"log"
	"time"

	"github.com/spf13/viper"
)

func InitConfig() {
	model.Creatd = time.Now()
	viper.SetConfigName("config") // for development / production environment
	// viper.SetConfigName("config_local") // for local environment

	viper.AddConfigPath(".")      // for production structure
	viper.AddConfigPath("../")    // for dev structure
	viper.AddConfigPath("../../") // for local environment

	// Find and read the config file
	err := viper.ReadInConfig()

	// Handle errors reading the config file
	if err != nil {
		log.Fatalln(err)
	}
}
