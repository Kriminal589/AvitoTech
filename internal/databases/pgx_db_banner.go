package databases

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func (p PgxDB) GetUserBanner(tagId uint64, featureId uint64) (pgtype.JSONB, error) {
	var message pgtype.JSONB

	sqlTest, args, err := psql.Select("message").From("banner b").LeftJoin("banner_tmp t ON b.id = t.banner_id").Where(squirrel.Eq{
		"t.tag_id": tagId, "b.feature_id": featureId}).ToSql()

	if err != nil {
		log.Errorf("Unable to build SELECT query: %v", err)
		return pgtype.JSONB{}, err
	}

	err = p.QueryRow(context.Background(), sqlTest, args[0], args[1]).Scan(&message)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return pgtype.JSONB{}, fmt.Errorf("db: reserve: no such banner with tag_id %d and feature_id %d", tagId, featureId)
	} else if err != nil {
		return pgtype.JSONB{}, err
	}
	return message, err
}

func (p PgxDB) GetBanner(featureId uint64) (pgtype.JSONB, error) {
	sqlTest, args, err := psql.Select("message").From("banner").Where(squirrel.Eq{
		"feature_id": featureId}).ToSql()

	if err != nil {
		log.Errorf("Unable to build SELECT query: %v", err)
		return pgtype.JSONB{}, err
	}

	var messages pgtype.JSONB

	err = p.QueryRow(context.Background(), sqlTest, args[0]).Scan(&messages)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return pgtype.JSONB{}, fmt.Errorf("db: reserve: no such banner with feature_id %d", featureId)
	}
	if err != nil {
		return pgtype.JSONB{}, err
	}

	return messages, err
}
