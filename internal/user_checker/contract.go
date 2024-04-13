package user_checker

type UserDB interface {
	GetUserRole(id uint64) (bool, error)
}
