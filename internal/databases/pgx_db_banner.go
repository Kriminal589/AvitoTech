package databases

import (
	"AvitoTech/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	pgxV5 "github.com/jackc/pgx/v5"
	"time"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func (p PgxDB) GetUserBanner(tagId uint64, featureId uint64) (pgtype.JSONB, error) {
	var data pgtype.JSONB

	sqlQuery, args, err := psql.Select("content").From("banner b").LeftJoin("banner_tag_link t ON b.banner_id = t.id").Where(squirrel.Eq{
		"t.tag_id": tagId, "b.feature_id": featureId}).ToSql()

	if err != nil {
		log.Errorf("Unable to build SELECT query: %v", err)
		return pgtype.JSONB{}, err
	}

	err = p.QueryRow(context.Background(), sqlQuery, args[0], args[1]).Scan(&data)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return pgtype.JSONB{}, fmt.Errorf("db: reserve: no such banner with tag_id %d and feature_id %d", tagId, featureId)
	} else if err != nil {
		return pgtype.JSONB{}, err
	}
	return data, err
}

func (p PgxDB) GetBannersWithFeatureId(id uint64) ([]models.Banner, error) {
	sqlQuery, args, err := psql.Select("banner_id, feature_id, content, is_active, created_at, updated_at").From("banner").Where(squirrel.Eq{
		"feature_id": id}).ToSql()

	if err != nil {
		log.Errorf("Unable to build SELECT query: %v", err)
		return nil, err
	}

	var (
		result               []models.Banner
		featureId, bannerId  uint64
		createdAt, updatedAt time.Time
		content              pgtype.JSONB
		isActive             bool
		tagId                uint64
	)

	rows, err := p.Query(context.Background(), sqlQuery, args[0])
	_, err = pgxV5.ForEachRow(rows, []any{&bannerId, &featureId, &content, &isActive, &createdAt, &updatedAt}, func() error {
		var tagsId []uint64

		sqlQuery, args, err := psql.Select("tag_id").From("banner_tag_link").Where(squirrel.Eq{
			"banner_id": bannerId}).ToSql()

		if err != nil {
			log.Errorf("Unable to build SELECT query: %v", err)
			return err
		}

		r, err := p.Query(context.Background(), sqlQuery, args[0])
		_, err = pgxV5.ForEachRow(r, []any{&tagId}, func() error {
			tagsId = append(tagsId, tagId)
			return nil
		})

		result = append(result, models.Banner{BannerID: bannerId, FeatureID: featureId, Content: content,
			IsActive: isActive, TagIDS: tagsId, CreatedAt: createdAt, UpdatedAt: updatedAt})

		return nil
	})

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("db: reserve: no such banner with featureId %d", id)
	}
	if err != nil {
		return nil, err
	}

	return result, err
}

func (p PgxDB) DeleteBanner(id uint64) (string, error) {
	sqlQuery, args, err := psql.Delete("banner").Where(squirrel.Eq{
		"banner_id": id}).ToSql()

	if err != nil {
		log.Errorf("Unable to build DELETE query: %v", err)
		return "", err
	}

	_, err = p.Exec(context.Background(), sqlQuery, args[0])

	if err != nil {
		return "", err
	}
	return "204", err
}

func (p PgxDB) PostBanner(banner models.RequestBanner) (uint64, error) {
	sqlQuery, args, err := psql.Insert("banner").Columns("feature_id", "content", "is_active",
		"created_at", "updated_at").Values(banner.FeatureID, banner.Content, banner.IsActive, time.Now().UTC(), time.Now().UTC()).Suffix("RETURNING \"banner_id\"").ToSql()

	if err != nil {
		log.Errorf("Unable to build INSERT query: %v", err)
		return 0, err
	}

	var bannerId uint64

	err = p.QueryRow(context.Background(), sqlQuery, args...).Scan(&bannerId)

	builder := psql.Insert("banner_tag_link").Columns("banner_id", "tag_id")

	for i := range banner.TagIDS {
		builder = builder.Values(bannerId, banner.TagIDS[i])
	}

	sqlQuery, args, err = builder.ToSql()
	if err != nil {
		log.Errorf("Unable to build INSERT query: %v", err)
		return 0, err
	}

	_ = p.QueryRow(context.Background(), sqlQuery, args...)

	return bannerId, err
}

func (p PgxDB) PatchBanner(banner models.RequestBanner, id int) (string, error) {
	sqlQuery, args, err := psql.Update("banner").Set("feature_id", banner.FeatureID).Set("is_active", banner.IsActive).
		Set("content", banner.Content).Set("updated_at", time.Now().UTC()).Where(squirrel.Eq{"banner_id": id}).ToSql()

	if err != nil {
		log.Errorf("Unable to build UPDATE query: %v", err)
		return "nil", err
	}

	_, err = p.Exec(context.Background(), sqlQuery, args...)

	if err != nil {
		return "", err
	}

	sqlQuery, args, err = psql.Delete("banner_tag_link").Where(squirrel.Eq{"banner_id": id}).ToSql()
	if err != nil {
		log.Errorf("Unable to build DELETE query: %v", err)
	}
	_, err = p.Exec(context.Background(), sqlQuery, args...)

	builder := psql.Insert("banner_tag_link").Columns("banner_id", "tag_id")

	for i := range banner.TagIDS {
		builder = builder.Values(id, banner.TagIDS[i])
	}

	sqlQuery, args, err = builder.ToSql()
	if err != nil {
		log.Errorf("Unable to build INSERT query: %v", err)
		return "nil", err
	}

	_ = p.QueryRow(context.Background(), sqlQuery, args...)

	return "200", nil
}
