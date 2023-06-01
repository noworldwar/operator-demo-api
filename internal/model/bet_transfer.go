package model

type BetTransfer struct {
	TransferID string `json:"transferID" xorm:"varchar(255) pk"`
	PlayerID   string `json:"playerID"   xorm:"varchar(30) notnull"`
	Type       string `json:"type"   xorm:"varchar(20) notnull"`
	BetID      string `json:"betID"   xorm:"varchar(64) notnull"`
	GameID     string `json:"gameID"   xorm:"varchar(64) notnull"`
	Amount     int64  `json:"amount" xorm:"float notnull" `
	WeType     string `json:"wetype"   xorm:"varchar(20) notnull"`
	WeTime     string `json:"wetime"   xorm:"bigint"`
	Success    bool   `json:"success" xorm:"tinyint(1) notnull"`
	Created    int64  `xorm:"bigint"`
	Updated    int64  `xorm:"bigint"`
}

func CheckIfBetTransferExist(betID, apitype string) (bool, error) {
	session := WGDB.NewSession()
	defer session.Close()
	return session.Exist(&BetTransfer{BetID: betID, Type: apitype})
}

func GetBetTransferByBetID(betID, apitype string) (m BetTransfer, err error) {
	session := WGDB.NewSession()
	defer session.Close()
	_, err = session.Where("bet_id = ?", betID).And("type = ?", apitype).Get(&m)
	return
}

func GetBetTransferBy(playerID string) (m []BetTransfer, err error) {
	session := WGDB.NewSession()
	defer session.Close()
	err = session.Where("PlayerID=?", playerID).Desc("Created").Limit(10, 0).Find(&m)
	return
}

func AddBetTransfer(m BetTransfer) (err error) {
	session := WGDB.NewSession()
	defer session.Close()
	_, err = session.Insert(&m)
	return
}
