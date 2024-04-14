package databases

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxDB struct {
	*pgxpool.Pool
	logger pgx.Logger
	psql   squirrel.StatementBuilderType
}

func NewPgxDB(pool *pgxpool.Pool, logger pgx.Logger) *PgxDB {
	return &PgxDB{
		Pool:   pool,
		logger: logger,
		psql:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
