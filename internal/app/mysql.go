package app

import (
	"bc-opp-api/internal/model"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/core"
	"xorm.io/xorm"
)

func InitMySQL() {
	model.WGDB, _ = xorm.NewEngine("mysql", "root:password@tcp(127.0.0.1:3306)/wgdb?charset=utf8")
	model.WGDB.SetMapper(core.GonicMapper{})

	err := model.WGDB.Ping()
	if err != nil {
		log.Fatalln("Init MySQL Error:", err)
	}

	err = model.WGDB.Sync2(new(model.Player), new(model.Transfer), new(model.Wallet), new(model.Game))
	if err != nil {
		log.Fatalln("AutoMigrate Error:", err)
	}
}
