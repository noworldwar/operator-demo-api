package app

import (
	"bc-opp-api/internal/model"
	"time"

	"github.com/spf13/viper"
)

func InitConfig() {
	model.Creatd = time.Now()
	viper.SetDefault("tpg_game_url", "https://stagingweblobby.triple-pg.com/we/direct2Game/202")
}
