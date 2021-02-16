package repository

import (
	"log"

	"github.com/jvitoroc/todo-go/config"
	"github.com/jvitoroc/todo-go/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(cfg *config.Config) *Repository {
	db := getDb(cfg)

	db.AutoMigrate(model.User{})
	db.AutoMigrate(model.VerificationRequest{})
	db.AutoMigrate(model.Todo{})

	return &Repository{
		DB: db,
	}
}

func (u *Repository) BeginTran(fn func(*Repository) *model.AppError) *model.AppError {
	tx := u.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return model.NewGenericInternalError(err)
	}

	tran := &Repository{DB: tx}
	if err := fn(tran); err != nil {
		tran.DB.Rollback()
		return err
	}

	if err := tran.DB.Commit().Error; err != nil {
		return model.NewGenericInternalError(err)
	}

	return nil
}

func getDb(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:test.db?&cache=shared&_fk=1"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %s", err.Error())
	}
	return db
}
