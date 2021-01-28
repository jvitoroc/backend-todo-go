package repo

import "gorm.io/gorm"

func Initialize(db *gorm.DB) {
	db.AutoMigrate(&Todo{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&UserActivationRequest{})
}
