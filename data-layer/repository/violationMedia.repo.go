package repository

import (
	"context"

	"BE-DroneMEQ-Violation-Management/data-layer/model"

	"gorm.io/gorm"
)

type (
	ViolationMediaRepository interface {
		CreateViolationMedia(ctx context.Context, tx *gorm.DB, violationMedia model.ViolationMedia) (model.ViolationMedia, error)
	}

	violationMediaRepository struct {
		db *gorm.DB
	}
)

func NewViolationMediaRepository(db *gorm.DB) ViolationMediaRepository {
	return &violationMediaRepository{
		db: db,
	}
}

func (r *violationMediaRepository) CreateViolationMedia(ctx context.Context, tx *gorm.DB, violationMedia model.ViolationMedia) (model.ViolationMedia, error) {
	if tx == nil {
		tx = r.db
	}

	// Create the violation media in the database
	if err := tx.WithContext(ctx).Create(&violationMedia).Error; err != nil {
		return model.ViolationMedia{}, err
	}
	// create the violation_violationMedia in the database
	if err := tx.WithContext(ctx).Create(&model.ViolationViolationMedia{
		ViolationPkid:      violationMedia.ViolationPkid,
		ViolationMediaPkid: violationMedia.Pkid,
	}).Error; err != nil {
		return model.ViolationMedia{}, err
	}

	return violationMedia, nil
}
