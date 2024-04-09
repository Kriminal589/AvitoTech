package databases

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	pgxV5 "github.com/jackc/pgx/v5"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func (p PgxDB) GetUserBanner(tagId uint64, featureId uint64) (pgtype.JSONB, error) {
	var data pgtype.JSONB

	sqlRequest, args, err := psql.Select("content").From("banner b").LeftJoin("banner_tag_link t ON b.id = t.banner_id").Where(squirrel.Eq{
		"t.tag_id": tagId, "b.feature_id": featureId}).ToSql()

	if err != nil {
		log.Errorf("Unable to build SELECT query: %v", err)
		return pgtype.JSONB{}, err
	}

	err = p.QueryRow(context.Background(), sqlRequest, args[0], args[1]).Scan(&data)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return pgtype.JSONB{}, fmt.Errorf("db: reserve: no such banner with tag_id %d and feature_id %d", tagId, featureId)
	} else if err != nil {
		return pgtype.JSONB{}, err
	}
	return data, err
}

func (p PgxDB) GetBanner(featureId uint64) ([]pgtype.JSONB, error) {
	sqlRequest, args, err := psql.Select("*").From("banner").Where(squirrel.Eq{
		"feature_id": featureId}).ToSql()

	if err != nil {
		log.Errorf("Unable to build SELECT query: %v", err)
		return nil, err
	}

	var result []pgtype.JSONB
	var data []pgtype.JSONB

	rows, err := p.Query(context.Background(), sqlRequest, args[0])

	_, err = pgxV5.ForEachRow(rows, []any{&data}, func() error {
		fmt.Println(data)
		result = append(result, data...)
		return nil
	})

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("db: reserve: no such banner with feature_id %d", featureId)
	}
	if err != nil {
		return nil, err
	}

	return result, err
}
