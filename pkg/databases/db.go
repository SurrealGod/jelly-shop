package databases

import (
	"log"

	"github.com/SurrealGod/jelly-shop-code/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	
)

func DbConnect(cfg config.IDbConfig) *sqlx.DB{
	//connect db
	db, err:= sqlx.Connect("pgx", cfg.Url())

	if err != nil{
		log.Fatalf("connect to Database failed: %v\n", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxOpenConns())
	return db
}