package model

type Player struct {
	PlayerID   string `json:"playerID" xorm:"varchar(30) pk"`
	OpPlayerID string `json:"opPlayerID" xorm:"varchar(25) notnull"`
	Nickname   string `json:"nickname" xorm:"varchar(30)"`
	Currency   string `json:"currency" xorm:"varchar(5) notnull"`
	Password   string `json:"password" xorm:"varchar(255)"`
	Balance    int64  `json:"balance"`
	Test       int    `xorm:"int notnull"`
	Disabled   bool   `json:"disabled"`
	Created    int64  `xorm:"bigint"`
	Updated    int64  `xorm:"bigint"`
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
