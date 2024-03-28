package domain

// User 领域对象
type User struct {
	Email    string
	Password string
}

type Address struct {
	Id     int64
	UserId int64
}
