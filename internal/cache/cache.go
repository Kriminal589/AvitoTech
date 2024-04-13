package cache

import (
	"crypto/sha256"
	"github.com/go-pkgz/expirable-cache/v3"
	"strconv"
	"time"
)

const (
	cacheCapacity = 10000
	ttl           = time.Minute * 5
)

type cacheDB struct {
	db    BannerGetter
	cache cache.Cache[[32]byte, []byte]
}

func New(db BannerGetter) *cacheDB {
	bannerCache := cache.NewCache[[32]byte, []byte]().WithMaxKeys(cacheCapacity).WithTTL(ttl)

	return &cacheDB{
		db:    db,
		cache: bannerCache,
	}
}

func (c cacheDB) GetUserBanner(tagId uint64, featureId uint64, lastRevision bool, isAdmin bool) ([]byte, error) {
	if lastRevision {
		return c.db.GetUserBanner(tagId, featureId, isAdmin)
	}

	key := keyCash(tagId, featureId, isAdmin)
	if cachedBanner, ok := c.cache.Get(key); ok {
		return cachedBanner, nil
	}

	res, err := c.db.GetUserBanner(tagId, featureId, isAdmin)

	if err != nil {
		return nil, err
	}

	c.cache.Set(key, res, 0)

	return res, nil
}

func keyCash(tagId uint64, featureId uint64, isAdmin bool) [32]byte {
	tag := strconv.FormatUint(tagId, 10)
	feature := strconv.FormatUint(featureId, 10)
	admin := strconv.FormatBool(isAdmin)
	key := []byte(tag + feature + admin)

	return sha256.Sum256(key)
}
