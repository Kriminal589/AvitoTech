package databases

import (
	"github.com/jackc/pgtype"
)

type DBInt interface {
	GetUserBanner(tagId uint64, featureId uint64) (pgtype.JSONB, error)
	GetBanner(featureId uint64) (pgtype.JSONB, error)
}
