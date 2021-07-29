package model

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var RDB *redis.Client

type PlayerInfo struct {
	PlayerID string
	Nickname string
}

func SetPlayerInfo(player Player) string {
	token := uuid.New().String()
	token = strings.Replace(token, "-", "", -1)
	b, _ := json.Marshal(PlayerInfo{PlayerID: player.PlayerID, Nickname: player.Nickname})
	_ = RDB.Set(context.Background(), token, string(b), time.Hour*1).Err()
	return token
}

func GetPlayerInfo(token string) (info PlayerInfo) {
	info = PlayerInfo{}

	res, err := RDB.Get(context.Background(), token).Result()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(res), &info)
	if err != nil {
		fmt.Println("Convert Error:", err)
		return
	}

	_ = RDB.Expire(context.Background(), token, time.Hour*1).Err()
	return
}
