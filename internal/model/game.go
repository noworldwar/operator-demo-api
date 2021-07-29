package model

import (
	"time"

	"github.com/pkg/errors"
)

type Game struct {
	GameID         string    `json:"gameID" xorm:"varchar(50)   pk"`
	Name           string    `json:"name" xorm:"varchar(100)  notnull"`
	Sort           int       `json:"sort" xorm:"int           notnull"`
	Image          string    `json:"image" xorm:"varchar(255)  notnull"`
	GameType       string    `json:"gameType" xorm:"varchar(20)   notnull"`
	Description    string    `json:"description" xorm:"varchar(255)  notnull"`
	ProviderID     string    `json:"providerID" xorm:"varchar(20)   notnull"`
	ProviderGameID string    `json:"providerGameID" xorm:"varchar(20)   notnull"`
	Disabled       bool      `json:"disabled"`
	Created        time.Time `json:"created"  xorm:"created"`
	Updated        time.Time `json:"updated"  xorm:"updated"`
}

func GetGameList(GameType string) (m []Game, err error) {
	datalist := []Game{}
	session := WGDB.NewSession()
	defer session.Close()

	if GameType != "" && len(GameType) > 0 {
		session = session.In("game_type", GameType)
	}

	err = session.Find(&datalist)
	if err != nil {
		return nil, errors.Errorf("GetGame fail: %v", err)
	}
	return datalist, nil
}

func AddGame(m Game) (err error) {
	session := WGDB.NewSession()
	defer session.Close()
	_, err = session.Insert(&m)
	return
}

func UpdateGame(m Game) (err error) {
	session := WGDB.NewSession()
	defer session.Close()
	_, err = session.ID(m.GameID).Update(&m)
	return
}
