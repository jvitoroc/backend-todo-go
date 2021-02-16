package repository

import (
	"errors"
	"fmt"

	"github.com/jvitoroc/todo-go/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (u *Repository) UpsertVerificationRequest(vr *model.VerificationRequest) (*model.VerificationRequest, *model.AppError) {
	if err := u.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"code", "expires_at"}),
	}).Create(vr).Error; err != nil {
		return nil, model.NewGenericInternalError(err)
	}

	return vr, nil
}

func (u *Repository) GetVerificationRequest(userId int) (*model.VerificationRequest, *model.AppError) {
	vr := model.VerificationRequest{}
	if err := u.DB.First(&vr, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.NewNotFoundError(fmt.Sprintf(model.MSG_VERIFICATION_NOT_FOUND, userId))
		} else {
			return nil, model.NewGenericInternalError(err)
		}
	}

	return &vr, nil
}
