package databases

type DBInt interface {
	GetUserBanner(tagId uint64, featureId uint64) (uint64, error)
}
