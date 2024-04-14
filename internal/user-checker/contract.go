package userchecker

type UserDB interface {
	GetUserRole(id uint64) (bool, error)
}
