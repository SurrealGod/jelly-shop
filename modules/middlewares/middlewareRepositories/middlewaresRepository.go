package middlewaresRepositories

import "github.com/jmoiron/sqlx"

type IMiddlewareRepository interface {
}

type middlewareRepository struct {
	db *sqlx.DB
}

func MiddlewaresRepository(db *sqlx.DB) IMiddlewareRepository {
	
	return &middlewareRepository{
		db: db,
	}
}
