package cache

import (
	"crypto/sha256"
	"strconv"
	"time"

	cache "github.com/go-pkgz/expirable-cache/v3"
)

const (
	cacheCapacity = 10000
	ttl           = time.Minute * 5
)

type BannerCache struct {
	db    BannerGetter
	cache cache.Cache[[32]byte, []byte]
}

func New(db BannerGetter) *BannerCache {
	bannerCache := cache.NewCache[[32]byte, []byte]().WithMaxKeys(cacheCapacity).WithTTL(ttl)

	return &BannerCache{
		db:    db,
		cache: bannerCache,
	}
}

func (c BannerCache) GetUserBanner(tagID uint64, featureID uint64, lastRevision bool, isAdmin bool) ([]byte, error) {
	if lastRevision {
		return c.db.GetUserBanner(tagID, featureID, isAdmin)
	}

	key := keyCash(tagID, featureID, isAdmin)
	if cachedBanner, ok := c.cache.Get(key); ok {
		return cachedBanner, nil
	}

	res, err := c.db.GetUserBanner(tagID, featureID, isAdmin)

	if err != nil {
		return nil, err
	}

	c.cache.Set(key, res, 0)

	return res, nil
}

func keyCash(tagID uint64, featureID uint64, isAdmin bool) [32]byte {
	tag := strconv.FormatUint(tagID, 10)
	feature := strconv.FormatUint(featureID, 10)
	admin := strconv.FormatBool(isAdmin)
	key := []byte(tag + feature + admin)

	return sha256.Sum256(key)
}
