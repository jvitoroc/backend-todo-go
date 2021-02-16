package repository

import (
	"errors"
	"fmt"

	"github.com/jvitoroc/todo-go/model"
	"gorm.io/gorm"
)

func (u *Repository) CreateUser(user *model.User) (*model.User, *model.AppError) {
	if err := u.DB.Create(user).Error; err != nil {
		return nil, model.NewGenericInternalError(err)
	}

	return user, nil
}

func (u *Repository) GetUser(userId int) (*model.User, *model.AppError) {
	user := model.User{}
	if err := u.DB.First(&user, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.NewNotFoundError(fmt.Sprintf(model.MSG_USER_NOT_FOUND, userId))
		} else {
			return nil, model.NewGenericInternalError(err)
		}
	}

	return &user, nil
}

func (u *Repository) GetUserByUsername(username string) (*model.User, *model.AppError) {
	user := model.User{}
	if err := u.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.NewNotFoundError(fmt.Sprintf(model.MSG_USERNAME_NOT_FOUND, username))
		} else {
			return nil, model.NewGenericInternalError(err)
		}
	}

	return &user, nil
}

func (u *Repository) GetUserByGoogleSub(googleSub string) (*model.User, *model.AppError) {
	user := model.User{}
	if err := u.DB.Where("google_sub = ?", googleSub).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.NewNotFoundError(fmt.Sprintf(model.MSG_USER_GOOGLE_ID_NOT_FOUND, googleSub))
		} else {
			return nil, model.NewGenericInternalError(err)
		}
	}

	return &user, nil
}

func (u *Repository) UpdateUser(user *model.User) *model.AppError {
	if err := u.DB.Save(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.NewNotFoundError(fmt.Sprintf(model.MSG_USER_NOT_FOUND, user.ID))
		} else {
			return model.NewGenericInternalError(err)
		}
	}

	return nil
}

func (u *Repository) CheckIfUsernameExists(username string) (bool, *model.AppError) {
	var count int64
	if err := u.DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, model.NewGenericInternalError(err)
	}

	return count > 0, nil
}

func (u *Repository) CheckIfEmailExists(email string) (bool, *model.AppError) {
	var count int64
	if err := u.DB.Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, model.NewGenericInternalError(err)
	}

	return count > 0, nil
}
