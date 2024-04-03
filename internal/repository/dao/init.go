package dao

import "gorm.io/gorm"

// InitTable 建表
func InitTable(db *gorm.DB) error {
	// Gorm会默认给表名添加复数 user -> users
	return db.AutoMigrate(&User{})
}
