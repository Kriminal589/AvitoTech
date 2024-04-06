package databases

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
)

func (p PgxDB) GetUserBanner(tagId uint64, featureId uint64) (uint64, error) {
	ctx := context.Background()

	var err error
	var bannerId uint64
	err = p.QueryRow(ctx, "SELECT id FROM banner WHERE tag_id = $1 AND feature_id = $2", tagId, featureId).Scan(&bannerId)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("db: reserve: no such banner with tag_id %d and feature_id %d", tagId, featureId)
	} else if err != nil {
		return 0, err
	}
	return bannerId, err
}
