package cache

type BannerGetter interface {
	GetUserBanner(tagID uint64, featureID uint64, isAdmin bool) ([]byte, error)
}
