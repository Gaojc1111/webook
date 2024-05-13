package domain

// User 领域对象
type User struct {
	ID       int64
	Email    string
	Password string
	Phone    string
}

type Address struct {
	Id     int64
	UserId int64
}
