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

func (p PgxDB) GetUserBanner(tagID uint64, featureID uint64, isAdmin bool) ([]byte, error) {
	var data pgtype.JSONB

	sqlQuery, _, err := p.psql.Select("content").From("banner b").
		LeftJoin("banner_tag_link t ON b.banner_id = t.banner_id").
		Where("($1 = true and t.tag_id = $2 and b.feature_id = $3) or " +
			"($1 = false and t.tag_id = $2 and b.feature_id = $3 and b.is_active = true)").ToSql()

	if err != nil {
		log.Errorf("Unable to build SELECT query: %v", err)
		return nil, err
	}

	err = p.QueryRow(context.Background(), sqlQuery, isAdmin, tagID, featureID).Scan(&data)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("db: reserve: no such banner with tag_id %d and feature_id %d, %w", tagID, featureID, err)
	} else if err != nil {
		return nil, err
	}

	result := data.Bytes

	return result, err
}

func (p PgxDB) GetBannersByFeatureID(id uint64, limit int, offset int) ([]models.Banner, error) {
	subquery := p.psql.Select("b.banner_id, content, is_active, created_at, updated_at, btl.tag_id").
		From("banner b").LeftJoin("banner_tag_link btl ON b.banner_id = btl.banner_id").Where(squirrel.Eq{
		"b.feature_id": id,
	})

	sqlQuery, args, err := p.psql.Select("tab.banner_id, content, is_active, created_at, updated_at, "+
		"array_agg(tag_id) tag_ids").
		FromSelect(subquery, "tab").GroupBy("content, is_active, created_at, updated_at, " +
		"tab.banner_id").ToSql()

	if err != nil {
		log.Errorf("Unable to build SELECT query: %v", err)
		return nil, err
	}

	rows, err := p.Query(context.Background(), sqlQuery, args...)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("db: reserve: no such banner with feature_id %d, %w", id, err)
	}
	if err != nil {
		return nil, err
	}

	var (
		bannerID             uint64
		result               []models.Banner
		content              pgtype.JSONB
		isActive             bool
		createdAt, updatedAt time.Time
		tagIDs               []uint64
	)

	countLine := 0
	lenArray := 0
	_, _ = pgxV5.ForEachRow(rows, []any{&bannerID, &content, &isActive, &createdAt, &updatedAt, &tagIDs}, func() error {
		if countLine >= offset {
			if limit != -1 {
				if lenArray < limit {
					lenArray++

					result = append(result, models.Banner{
						BannerID:  bannerID,
						TagIDs:    tagIDs,
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
					BannerID:  bannerID,
					TagIDs:    tagIDs,
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

	if len(result) == 0 {
		return nil, pgx.ErrNoRows
	}

	return result, err
}

func (p PgxDB) GetBannersByTagID(id uint64, limit int, offset int) ([]models.Banner, error) {
	subquery := p.psql.Select("b.banner_id, feature_id, content, is_active, created_at, updated_at, btl.tag_id").
		From("banner b").LeftJoin("banner_tag_link btl ON b.banner_id = btl.banner_id").
		Where("b.banner_id in (select banner_id from banner_tag_link where tag_id = $1)")

	sqlQuery, _, err := p.psql.Select("tab.banner_id, feature_id, content, is_active, created_at, "+
		"updated_at, array_agg(tag_id) tag_ids").FromSelect(subquery, "tab").
		GroupBy("feature_id, content, is_active, created_at, updated_at, tab.banner_id").ToSql()

	if err != nil {
		log.Errorf("Unable to build SELECT query: %v", err)
		return nil, err
	}

	rows, err := p.Query(context.Background(), sqlQuery, id)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("db: reserve: no such banner with tag_id %d, %w", id, err)
	}
	if err != nil {
		return nil, err
	}

	var (
		isActive             bool
		bannerID, featureID  uint64
		tagIDs               []uint64
		content              pgtype.JSONB
		createdAt, updatedAt time.Time
		result               []models.Banner
	)

	countLine := 0
	lenArray := 0
	_, err = pgxV5.ForEachRow(rows, []any{&bannerID, &featureID, &content,
		&isActive, &createdAt, &updatedAt, &tagIDs}, func() error {
		if countLine >= offset {
			if limit != -1 {
				if lenArray < limit {
					lenArray++

					result = append(result, models.Banner{
						IsActive:  isActive,
						BannerID:  bannerID,
						TagIDs:    tagIDs,
						FeatureID: featureID,
						Content:   content,
						CreatedAt: createdAt,
						UpdatedAt: updatedAt,
					})
				}
			} else {
				lenArray++

				result = append(result, models.Banner{
					IsActive:  isActive,
					BannerID:  bannerID,
					TagIDs:    tagIDs,
					FeatureID: featureID,
					Content:   content,
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				})
			}
		}

		countLine++

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p PgxDB) GetBanners(featureID uint64, tagID uint64, limit int, offset int) ([]models.Banner, error) {
	sub := p.psql.Select("b.banner_id, feature_id, content, is_active, created_at, updated_at, btl.tag_id").
		From("banner b").LeftJoin("banner_tag_link btl ON b.banner_id = btl.banner_id").Where(squirrel.Eq{
		"b.feature_id": featureID,
	})

	subquery := p.psql.Select("tab.banner_id, feature_id, content, is_active, created_at, updated_at, "+
		"array_agg(tag_id) tag_ids").
		FromSelect(sub, "tab").GroupBy("feature_id, content, is_active, created_at, updated_at, " +
		"tab.banner_id")

	sqlQuery, _, err := p.psql.Select("*").FromSelect(subquery, "ban").Where("$2 = ANY (tag_ids)").
		ToSql()

	if err != nil {
		log.Errorf("Unable to build SELECT query: %v", err)
		return nil, err
	}

	rows, err := p.Query(context.Background(), sqlQuery, featureID, tagID)

	if err != nil {
		return nil, err
	}

	var (
		bannerID             uint64
		createdAt, updatedAt time.Time
		content              pgtype.JSONB
		isActive             bool
		tagIDs               []uint64
		result               []models.Banner
	)

	countLine := 0
	lenArray := 0
	_, _ = pgxV5.ForEachRow(rows, []any{&bannerID, &featureID, &content, &isActive,
		&createdAt, &updatedAt, &tagIDs,
	}, func() error {
		if countLine >= offset {
			if limit != -1 {
				if lenArray < limit {
					lenArray++

					result = append(result, models.Banner{
						BannerID:  bannerID,
						TagIDs:    tagIDs,
						FeatureID: featureID,
						Content:   content,
						IsActive:  isActive,
						CreatedAt: createdAt,
						UpdatedAt: updatedAt,
					})
				}
			} else {
				lenArray++

				result = append(result, models.Banner{
					BannerID:  bannerID,
					TagIDs:    tagIDs,
					FeatureID: featureID,
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

	return result, nil
}

func (p PgxDB) DeleteBanner(id uint64) error {
	sqlQuery, args, err := p.psql.Delete("banner").Where(squirrel.Eq{
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

func (p PgxDB) PostBanner(banner models.Banner) (uint64, error) {
	sqlQuery, args, err := p.psql.Insert("banner").Columns("feature_id", "content", "is_active",
		"created_at", "updated_at").Values(
		banner.FeatureID,
		banner.Content,
		banner.IsActive,
		time.Now().UTC(),
		time.Now().UTC(),
	).Suffix("RETURNING \"banner_id\"").ToSql()

	if err != nil {
		log.Errorf("Unable to build INSERT query: %v", err)
		return 0, err
	}

	var bannerID uint64

	err = p.QueryRow(context.Background(), sqlQuery, args...).Scan(&bannerID)

	if err != nil {
		log.Errorf("failed to execute query: %v", err)
	}

	builder := p.psql.Insert("banner_tag_link").Columns("banner_id", "tag_id")

	for i := range banner.TagIDs {
		builder = builder.Values(bannerID, banner.TagIDs[i])
	}

	sqlQuery, args, err = builder.ToSql()
	if err != nil {
		log.Errorf("Unable to build INSERT query: %v", err)
		return 0, err
	}

	_ = p.QueryRow(context.Background(), sqlQuery, args...)

	return bannerID, err
}

func (p PgxDB) PatchBanner(banner models.Banner, id int) error {
	sqlQuery, args, err := p.psql.Update("banner").Set("feature_id", banner.FeatureID).Set("is_active", banner.IsActive).
		Set("content", banner.Content).Set("updated_at", time.Now().UTC()).Where(squirrel.Eq{"banner_id": id}).ToSql()

	if err != nil {
		log.Errorf("Unable to build UPDATE query: %v", err)
		return err
	}

	_, err = p.Exec(context.Background(), sqlQuery, args...)

	if err != nil {
		return err
	}

	sqlQuery, args, err = p.psql.Delete("banner_tag_link").Where(squirrel.Eq{"banner_id": id}).ToSql()
	if err != nil {
		log.Errorf("Unable to build DELETE query: %v", err)
		return err
	}
	_, err = p.Exec(context.Background(), sqlQuery, args...)

	if err != nil {
		return err
	}

	builder := p.psql.Insert("banner_tag_link").Columns("banner_id", "tag_id")

	for i := range banner.TagIDs {
		builder = builder.Values(id, banner.TagIDs[i])
	}

	sqlQuery, args, err = builder.ToSql()
	if err != nil {
		log.Errorf("Unable to build INSERT query: %v", err)
		return err
	}

	_, err = p.Exec(context.Background(), sqlQuery, args...)

	return err
}
