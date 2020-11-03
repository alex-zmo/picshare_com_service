package main

import (
	"github.com/gmo-personal/picshare_com_service/database"
	"github.com/gmo-personal/picshare_com_service/server"
)

func main() {
	database.InitDatabase()
	server.InitServer()
}
