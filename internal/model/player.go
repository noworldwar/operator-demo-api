package model

import (
	"time"
)

type Player struct {
	PlayerID string    `json:"playerID" xorm:"varchar(30) pk"`
	Nickname string    `json:"nickname" xorm:"varchar(30)"`
	Password string    `json:"password" xorm:"varchar(255)"`
	Balance  int64     `json:"balance"`
	Disabled bool      `json:"disabled"`
	Created  time.Time `json:"created"  xorm:"created"`
	Updated  time.Time `json:"updated"  xorm:"updated"`
}

func GetPlayer(playerID string) (m Player, err error) {
	session := WGDB.NewSession()
	defer session.Close()
	_, err = session.ID(playerID).Get(&m)
	return
}

func AddPlayer(m Player) (err error) {
	session := WGDB.NewSession()
	defer session.Close()
	_, err = session.Insert(&m)
	return
}

func UpdatePlayer(m Player) (err error) {
	session := WGDB.NewSession()
	defer session.Close()
	_, err = session.ID(m.PlayerID).Update(&m)
	return
}
