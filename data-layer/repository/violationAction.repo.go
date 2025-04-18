package repository

import (
	"context"
	// "github.com/mci-its/backend-service/dto"
	"BE-DroneMEQ-Violation-Management/data-layer/model"

	"gorm.io/gorm"
)

type (
	ViolationActionRepository interface {
		CreateViolationAction(ctx context.Context, tx *gorm.DB, violationAction model.ViolationAction) (model.ViolationAction, error)
	}

	violationActionRepository struct {
		db *gorm.DB
	}
)

func NewViolationActionRepository(db *gorm.DB) ViolationActionRepository {
	return &violationActionRepository{
		db: db,
	}
}

func (r *violationActionRepository) CreateViolationAction(ctx context.Context, tx *gorm.DB, violationAction model.ViolationAction) (model.ViolationAction, error) {
	if tx == nil {
		tx = r.db
	}

	// Create the violation action in the database
	if err := tx.WithContext(ctx).Create(&violationAction).Error; err != nil {
		return model.ViolationAction{}, err
	}
	return violationAction, nil
}
