package databases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	pgxV5 "github.com/jackc/pgx/v5"

	"AvitoTech/internal/models"
)

func (p PgxDB) GetUserBanner(tagId uint64, featureId uint64, isAdmin bool) ([]byte, error) {
	var data pgtype.JSONB

	sqlQuery, _, err := psql.Select("content").From("banner b").LeftJoin("banner_tag_link t ON b.banner_id = t.banner_id").
		Where("($1 = true and t.tag_id = $2 and b.feature_id = $3) or " +
			"($1 = false and t.tag_id = $2 and b.feature_id = $3 and b.is_active = true)").ToSql()

	if err != nil {
		log.Errorf("Unable to build SELECT query: %v", err)
		return nil, err
	}

	err = p.QueryRow(context.Background(), sqlQuery, isAdmin, tagId, featureId).Scan(&data)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("db: reserve: no such banner with tag_id %d and feature_id %d", tagId, featureId)
	} else if err != nil {
		return nil, err
	}

	result := data.Bytes

	return result, err
}

func (p PgxDB) GetBannersByFeatureId(id uint64, limit int, offset int) ([]models.Banner, error) {
	subquery := psql.Select("b.banner_id, content, is_active, created_at, updated_at, btl.tag_id").
		From("banner b").LeftJoin("banner_tag_link btl ON b.banner_id = btl.banner_id").Where(squirrel.Eq{
		"b.feature_id": id,
	})

	sqlQuery, args, err := psql.Select("tab.banner_id, content, is_active, created_at, updated_at, array_agg(tag_id) tag_ids").
		FromSelect(subquery, "tab").GroupBy("content, is_active, created_at, updated_at, tab.banner_id").ToSql()

	if err != nil {
		log.Errorf("Unable to build SELECT query: %v", err)
		return nil, err
	}

	rows, err := p.Query(context.Background(), sqlQuery, args...)

	var (
		bannerId             uint64
		result               []models.Banner
		content              pgtype.JSONB
		isActive             bool
		createdAt, updatedAt time.Time
		tagIds               []uint64
	)

	countLine := 0
	lenArray := 0
	_, _ = pgxV5.ForEachRow(rows, []any{&bannerId, &content, &isActive, &createdAt, &updatedAt, &tagIds}, func() error {
		if countLine >= offset {
			if limit != -1 {
				if lenArray < limit {
					lenArray++

					result = append(result, models.Banner{
						BannerID:  bannerId,
						TagIDS:    tagIds,
						FeatureID: id,
						Content:   content,
						IsActive:  isActive,
						CreatedAt: createdAt,
						UpdatedAt: updatedAt,
					})
				}
			} else {
				lenArray++

				result = append(result, models.Banner{
					BannerID:  bannerId,
					TagIDS:    tagIds,
					FeatureID: id,
					Content:   content,
					IsActive:  isActive,
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				})
			}
		}

		countLine++

		return nil
	})

	return result, err
}

func (p PgxDB) GetBannersByTagId(id uint64, limit int, offset int) ([]models.Banner, error) {
	var bannerId uint64

	sqlQuery, args, err := psql.Select("banner_id").From("banner_tag_link").Where("tag_id = ?", id).ToSql()
	if err != nil {
		log.Errorf("Unable to build SELECT query: %v", err)
		return nil, err
	}

	rows, err := p.Query(context.Background(), sqlQuery, args...)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("db: reserve: no such banner with tag %d", id)
	}

	subquery := psql.Select("feature_id, content, is_active, created_at, updated_at, btl.tag_id").
		From("banner b").LeftJoin("banner_tag_link btl ON b.banner_id = btl.banner_id").Where(squirrel.Eq{
		"b.banner_id": bannerId,
	})

	sqlQuery, args, err = psql.Select("feature_id, content, is_active, created_at, updated_at, array_agg(tag_id) tag_ids").
		FromSelect(subquery, "tab").GroupBy("feature_id, content, is_active, created_at, updated_at").ToSql()

	if err != nil {
		log.Errorf("Unable to build SELECT query: %v", err)
		return nil, err
	}

	var (
		featureId            uint64
		createdAt, updatedAt time.Time
		content              pgtype.JSONB
		isActive             bool
		tagIds               []uint64
		result               []models.Banner
	)

	_, _ = pgxV5.ForEachRow(rows, []any{&bannerId}, func() error {

		err = p.QueryRow(context.Background(), sqlQuery, bannerId).Scan(&featureId, &content, &isActive, &createdAt, &updatedAt, &tagIds)

		if err != nil && errors.Is(err, pgx.ErrNoRows) {
			log.Errorf("Unable to build SELECT query: %v", err)
			return err
		}

		result = append(result, models.Banner{
			BannerID:  bannerId,
			TagIDS:    tagIds,
			FeatureID: featureId,
			Content:   content,
			IsActive:  isActive,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
		return nil
	})

	return result, nil
}

func (p PgxDB) GetBanners(featureId uint64, tagId uint64, limit int, offset int) ([]models.Banner, error) {
	sub := psql.Select("b.banner_id, feature_id, content, is_active, created_at, updated_at, btl.tag_id").
		From("banner b").LeftJoin("banner_tag_link btl ON b.banner_id = btl.banner_id").Where(squirrel.Eq{
		"b.feature_id": featureId,
	})

	subquery := psql.Select("tab.banner_id, feature_id, content, is_active, created_at, updated_at, array_agg(tag_id) tag_ids").
		FromSelect(sub, "tab").GroupBy("feature_id, content, is_active, created_at, updated_at, tab.banner_id")

	sqlQuery, _, err := psql.Select("*").FromSelect(subquery, "ban").Where("$2 = ANY (tag_ids)").ToSql()

	if err != nil {
		log.Errorf("Unable to build SELECT query: %v", err)
		return nil, err
	}

	rows, err := p.Query(context.Background(), sqlQuery, featureId, tagId)

	var (
		bannerId             uint64
		createdAt, updatedAt time.Time
		content              pgtype.JSONB
		isActive             bool
		tagIds               []uint64
		result               []models.Banner
	)

	_, _ = pgxV5.ForEachRow(rows, []any{&bannerId, &featureId, &content, &isActive, &createdAt, &updatedAt, &tagIds}, func() error {
		result = append(result, models.Banner{
			BannerID:  bannerId,
			TagIDS:    tagIds,
			FeatureID: featureId,
			Content:   content,
			IsActive:  isActive,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})

		return nil
	})

	return result, nil
}

func (p PgxDB) DeleteBanner(id uint64) error {
	sqlQuery, args, err := psql.Delete("banner").Where(squirrel.Eq{
		"banner_id": id}).ToSql()

	if err != nil {
		log.Errorf("Unable to build DELETE query: %v", err)
		return err
	}

	_, err = p.Exec(context.Background(), sqlQuery, args[0])

	if err != nil {
		return err
	}
	return err
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

func (p PgxDB) PatchBanner(banner models.RequestBanner, id int) error {
	sqlQuery, args, err := psql.Update("banner").Set("feature_id", banner.FeatureID).Set("is_active", banner.IsActive).
		Set("content", banner.Content).Set("updated_at", time.Now().UTC()).Where(squirrel.Eq{"banner_id": id}).ToSql()

	if err != nil {
		log.Errorf("Unable to build UPDATE query: %v", err)
		return err
	}

	_, err = p.Exec(context.Background(), sqlQuery, args...)

	if err != nil {
		return err
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
		return err
	}

	_, err = p.Exec(context.Background(), sqlQuery, args...)

	return err
}
