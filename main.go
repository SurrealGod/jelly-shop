package main

import (
	"os"

	"github.com/SurrealGod/jelly-shop-code/config"
	"github.com/SurrealGod/jelly-shop-code/modules/servers"
	"github.com/SurrealGod/jelly-shop-code/pkg/databases"
)

func envPath() string {

	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())
	db := databases.DbConnect(cfg.Db())
	defer db.Close()
	servers.NewServer(cfg, db).Start()
}