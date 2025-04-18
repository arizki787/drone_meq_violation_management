package repository

import (
	"context"
	"log"
	"math"
	"time"

	"BE-DroneMEQ-Violation-Management/data-layer/model"
	"BE-DroneMEQ-Violation-Management/helpers/dto"

	"gorm.io/gorm"
)

type (
	// ViolationRepository interface for data access layer
	ViolationRepository interface {
		CreateViolation(ctx context.Context, tx *gorm.DB, violation model.Violation) (model.Violation, error)
		GetAllViolationWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest, OfficePkid int64) (dto.GetAllViolationRepositoryResponse, error)
		GetViolationById(ctx context.Context, tx *gorm.DB, pkid int64) (model.Violation, []model.ViolationMedia, error)
		GetViolationByFilter(ctx context.Context, tx *gorm.DB, status []string, violationType []string, severity []string, req dto.PaginationRequest, OfficePkid int64) (dto.GetAllViolationRepositoryResponse, error)
		UpdateViolation(ctx context.Context, tx *gorm.DB, pkid int64, updatedData model.Violation) (model.Violation, error)
		RestoreViolation(ctx context.Context, tx *gorm.DB, pkid int64, violation model.Violation) (model.Violation, error)
		HardDeleteViolation(ctx context.Context, pkid int64) error
	}

	violationRepository struct {
		db *gorm.DB
	}
)

func applyPaginationDefaults(page, perPage *int) {
	if *perPage == 0 {
		*perPage = 10
	}
	if *page == 0 {
		*page = 1
	}
}

func NewViolationRepository(db *gorm.DB) ViolationRepository {
	return &violationRepository{
		db: db,
	}
}

func (r *violationRepository) CreateViolation(ctx context.Context, tx *gorm.DB, violation model.Violation) (model.Violation, error) {
	// Check if transaction is nil
	if tx == nil {
		// If nil, use the default db
		tx = r.db
	}

	// Create the violation
	if err := tx.WithContext(ctx).Create(&violation).Error; err != nil {
		return model.Violation{}, err
	}

	return violation, nil
}

func (r *violationRepository) GetAllViolationWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest, OfficePkid int64) (dto.GetAllViolationRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var violations []model.Violation
	var count int64

	// Apply default values for pagination
	applyPaginationDefaults(&req.Page, &req.PerPage)

	query := tx.WithContext(ctx).Model(&model.Violation{}).Debug()

	log.Printf("filtering by OfficePkid: %v", OfficePkid)

	if OfficePkid != 0 {
		query = query.Where("office_pkid = ?", OfficePkid)
	}

	// Count the total of violations
	if err := query.Count(&count).Error; err != nil {
		return dto.GetAllViolationRepositoryResponse{}, err
	}

	// Get all violations with pagination
	if err := query.Scopes(Paginate(req.Page, req.PerPage)).Find(&violations).Error; err != nil {
		return dto.GetAllViolationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.GetAllViolationRepositoryResponse{
		Violations: violations,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, nil
}

func (r *violationRepository) GetViolationById(ctx context.Context, tx *gorm.DB, Pkid int64) (model.Violation, []model.ViolationMedia, error) {
	if tx == nil {
		tx = r.db
	}

	var violation model.Violation
	var violationMedia []model.ViolationMedia
	if err := tx.WithContext(ctx).Unscoped().Where("Pkid= ?", Pkid).Take(&violation).Error; err != nil {
		return model.Violation{}, []model.ViolationMedia{}, err
	}
	if err := tx.WithContext(ctx).Unscoped().Where("violation_pkid= ?", Pkid).Find(&violationMedia).Error; err != nil {
		return model.Violation{}, []model.ViolationMedia{}, err
	}
	if len(violationMedia) == 0 {
		return violation, []model.ViolationMedia{}, nil
	}
	return violation, violationMedia, nil
}

func (r *violationRepository) GetViolationByFilter(ctx context.Context, tx *gorm.DB, status []string, violationType []string, severity []string, req dto.PaginationRequest, OfficePkid int64) (dto.GetAllViolationRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var violation []model.Violation
	var count int64

	applyPaginationDefaults(&req.Page, &req.PerPage)

	query := tx.WithContext(ctx).Model(&model.Violation{}).Debug()

	log.Printf("filtering by status: %v, violationType: %v, severity: %v, OfficePkid: %v", status, violationType, severity, OfficePkid)

	if len(violationType) > 0 {
		log.Printf("filtering by violationType: %v", violationType)
		query = query.Where("violation_type IN (?)", violationType)
	}
	if len(severity) > 0 {
		log.Printf("filtering by severity: %v", severity)
		query = query.Where("severity IN (?)", severity)
	}
	if len(status) > 0 {
		log.Printf("filtering by status: %v", status)
		query = query.Where("status IN (?)", status)
	}

	if OfficePkid != 0 {
		log.Printf("filtering by OfficePkid: %v", OfficePkid)
		query = query.Where("office_pkid = ?", OfficePkid)
	}
	// Count the total of violations
	if err := query.Count(&count).Error; err != nil {
		return dto.GetAllViolationRepositoryResponse{}, err
	}
	// Get all violations with pagination
	if err := query.Limit(req.PerPage).
		Offset((req.Page - 1) * req.PerPage).
		Find(&violation).Error; err != nil {
		return dto.GetAllViolationRepositoryResponse{}, err
	}

	// calculate the total number of pages
	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.GetAllViolationRepositoryResponse{
		Violations: violation,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, nil
}

func (r *violationRepository) UpdateViolation(ctx context.Context, tx *gorm.DB, pkid int64, updatedData model.Violation) (model.Violation, error) {
	// Use r.db if no transaction is provided
	if tx == nil {
		tx = r.db
	}

	// Update violation where Pkid matches
	if err := tx.WithContext(ctx).
		Model(&model.Violation{}).
		Where("Pkid = ?", pkid).
		Updates(updatedData).Error; err != nil {
		return model.Violation{}, err
	}

	return updatedData, nil
}

func (r *violationRepository) RestoreViolation(ctx context.Context, tx *gorm.DB, pkid int64, violation model.Violation) (model.Violation, error) {
	if tx == nil {
		tx = r.db
	}

	now := time.Now()
	updateData := map[string]interface{}{
		"deleted_at":   nil,
		"deleted_by":   nil,
		"deleted_host": nil,
		"is_deleted":   false,
		"updated_by":   violation.UpdatedBy,
		"updated_host": violation.UpdatedHost,
		"updated_at":   now,
		"office_pkid":  violation.OfficePkid,
	}

	err := tx.WithContext(ctx).Model(&model.Violation{}).Where("pkid = ?", pkid).Updates(updateData).Error
	if err != nil {
		return model.Violation{}, err
	}

	// Fetch the updated violation record
	err = tx.WithContext(ctx).Unscoped().Where("pkid = ?", pkid).Take(&violation).Error
	if err != nil {
		return model.Violation{}, err
	}

	return violation, nil
}

func (r *violationRepository) HardDeleteViolation(ctx context.Context, pkid int64) error {
	var violation model.Violation
	if err := r.db.WithContext(ctx).Unscoped().Where("pkid = ?", pkid).Take(&violation).Error; err != nil {
		log.Printf("Error when trying to get violation by pkid: %v", err)
		return err
	}
	if err := r.db.WithContext(ctx).Unscoped().Delete(&violation).Error; err != nil {
		log.Printf("Error when trying to delete violation by pkid: %v", err)
		return err
	}
	return nil
}
