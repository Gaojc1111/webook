package domain

// User 领域对象
type User struct {
	WechatInfo
	ID        int64
	Email     string
	Password  string
	Phone     string
	CreatedAt int64
	UpdatedAt int64
}

type Address struct {
	Id     int64
	UserId int64
}
