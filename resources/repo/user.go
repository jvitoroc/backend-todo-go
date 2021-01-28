package repo

import (
	"errors"
	"fmt"
	"time"

	"github.com/jvitoroc/todo-api/resources/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	ID           int       `json:"userId"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Email        string    `json:"email"`
	Active       bool      `json:"active"`
	Todos        []Todo    `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UserActivationRequest struct {
	UserID    int       `json:"userId" gorm:"primaryKey"`
	User      User      `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func InsertUser(db *gorm.DB, username string, email string, passwordHash string) (*User, *common.Error) {
	now := time.Now()
	tx := db.Begin()

	if err := CheckIfUsernameExists(tx, username); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := CheckIfEmailExists(tx, email); err != nil {
		tx.Rollback()
		return nil, err
	}

	user := User{Username: username, Email: email, PasswordHash: passwordHash, CreatedAt: now, UpdatedAt: now}
	if result := tx.Create(&user); result.Error != nil {
		tx.Rollback()
		return nil, common.CreateGenericInternalError(result.Error)
	}

	tx.Commit()
	return &user, nil
}

func UpsertUserActivationRequest(db *gorm.DB, userId int, code string, expiresAt time.Time) (*UserActivationRequest, *common.Error) {
	request := UserActivationRequest{UserID: userId, ExpiresAt: expiresAt, Code: code}

	if result := db.Clauses(clause.OnConflict{UpdateAll: true, Columns: []clause.Column{{Name: "userId"}}}).Create(&request); result.Error != nil {
		return nil, common.CreateGenericInternalError(result.Error)
	}

	return &request, nil
}

func UpdateUser(db *gorm.DB, user *User) *common.Error {
	var result *gorm.DB

	if result = db.Model(user).Updates(user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return common.CreateNotFoundError(fmt.Sprintf("User not found under given id (%d).", user.ID))
		} else {
			return common.CreateGenericInternalError(result.Error)
		}
	}

	return nil
}

func GetUserActivationRequest(db *gorm.DB, userId int) (*UserActivationRequest, *common.Error) {
	var result *gorm.DB
	request := UserActivationRequest{}

	if result = db.Where("user_id = ?", userId).First(&request); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, common.CreateNotFoundError(fmt.Sprintf("User's activation request not found under given user id (%d).", userId))
		} else {
			return nil, common.CreateGenericInternalError(result.Error)
		}
	}

	return &request, nil
}

func GetUser(db *gorm.DB, userId int) (*User, *common.Error) {
	var result *gorm.DB
	user := User{}

	if result = db.Where("id = ?", userId).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, common.CreateNotFoundError(fmt.Sprintf("User not found under given id (%d).", userId))
		} else {
			return nil, common.CreateGenericInternalError(result.Error)
		}
	}

	return &user, nil
}

func GetUserByUsername(db *gorm.DB, username string) (*User, *common.Error) {
	var result *gorm.DB
	user := User{}

	if result = db.Where("username = ?", username).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, common.CreateNotFoundError(fmt.Sprintf("User not found under given username (%s).", username))
		} else {
			return nil, common.CreateGenericInternalError(result.Error)
		}
	}

	return &user, nil
}

func CheckIfUsernameExists(db *gorm.DB, username string) *common.Error {
	var count int64

	if result := db.Model(&User{}).Select("id").Where("username = ?", username).Count(&count); result.Error != nil {
		return common.CreateGenericInternalError(result.Error)
	}

	if count > 0 {
		return common.CreateFormError(map[string]string{"username": "Username already exists."})
	}

	return nil
}

func CheckIfEmailExists(db *gorm.DB, email string) *common.Error {
	var count int64

	if result := db.Model(&User{}).Select("id").Where("email = ?", email).Count(&count); result.Error != nil {
		return common.CreateGenericInternalError(result.Error)
	}

	if count > 0 {
		return common.CreateFormError(map[string]string{"email": "Email already exists."})
	}

	return nil
}
