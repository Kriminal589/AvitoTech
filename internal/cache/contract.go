package cache

type BannerGetter interface {
	GetUserBanner(tagId uint64, featureId uint64, isAdmin bool) ([]byte, error)
}
