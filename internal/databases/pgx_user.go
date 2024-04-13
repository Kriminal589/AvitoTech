package databases

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
)

func (p PgxDB) GetUserRole(id uint64) (bool, error) {
	sqlQuery, args, err := psql.Select("admin").From("roles").Where(squirrel.Eq{"user_id": id}).ToSql()

	if err != nil {
		log.Errorf("Error while building SQL query: %v", err)
		return false, err
	}

	var admin bool

	err = p.QueryRow(context.Background(), sqlQuery, args...).Scan(&admin)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("db: reserve: no such user with id %d", id)
		}
	}

	return admin, err
}
