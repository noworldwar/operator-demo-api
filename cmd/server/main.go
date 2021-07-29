package main

import (
	"bc-opp-api/internal/app"
)

func main() {
	app.InitConfig()
	app.InitMySQL()
	app.InitRedis()
	app.InitRouter()
	app.RunServer()
}
