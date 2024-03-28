package dao

import "gorm.io/gorm"

// InitTable 建表
func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
