package repository

import (
	"context"
	"log"
	"math"
	"time"

	"BE-DroneMEQ-Violation-Management/data-layer/model"
	"BE-DroneMEQ-Violation-Management/helpers/dto"
	"BE-DroneMEQ-Violation-Management/helpers/dto/external"

	"gorm.io/gorm"
)

type (
	ViolationNoteRepository interface {
		CreateViolationNotes(ctx context.Context, tx *gorm.DB, violationNotes model.ViolationNote) (model.ViolationNote, error)
		GetAllViolationNoteWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest, violationId int64, metadata external.MetadataDTO) (dto.GetAllViolationNoteRepositoryResponse, error)
		GetViolationNoteByFilter(ctx context.Context, tx *gorm.DB, noteType []string, req dto.PaginationRequest, violationId int64, metadata external.MetadataDTO) (dto.GetAllViolationNoteRepositoryResponse, error)
		GetViolationNoteById(ctx context.Context, tx *gorm.DB, violationPkid int64, notePkid int64) (model.ViolationNote, error)
		UpdateViolationNotes(ctx context.Context, tx *gorm.DB, notePkid int64, violationId int64, updatedData model.ViolationNote) (model.ViolationNote, error)
		RestoreViolationNotes(ctx context.Context, tx *gorm.DB, notePkid int64, violationId int64, violationNote model.ViolationNote) (model.ViolationNote, error)
		HardDeleteViolationNotes(ctx context.Context, notePkid int64, violationId int64) error
	}

	violationNoteRepository struct {
		db *gorm.DB
	}
)

func NewViolationNoteRepository(db *gorm.DB) ViolationNoteRepository {
	return &violationNoteRepository{
		db: db,
	}
}

func (r *violationNoteRepository) CreateViolationNotes(ctx context.Context, tx *gorm.DB, violationNote model.ViolationNote) (model.ViolationNote, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&violationNote).Error; err != nil {
		return model.ViolationNote{}, err
	}
	return violationNote, nil
}

func (r *violationNoteRepository) GetAllViolationNoteWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest, violationId int64, metadata external.MetadataDTO) (dto.GetAllViolationNoteRepositoryResponse, error) {

	if tx == nil {
		tx = r.db
	}

	var total int64
	var violationNotes []model.ViolationNote

	applyPaginationDefaults(&req.Page, &req.PerPage)

	query := tx.WithContext(ctx).Model(&model.ViolationNote{}).Where("violation_pkid = ? AND office_pkid = ?", violationId, metadata.OfficePKID)

	if err := query.Count(&total).Error; err != nil {
		return dto.GetAllViolationNoteRepositoryResponse{}, err
	}

	if err := query.Scopes(Paginate(req.Page, req.PerPage)).Find(&violationNotes).Error; err != nil {
		return dto.GetAllViolationNoteRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(total) / float64(req.PerPage)))

	return dto.GetAllViolationNoteRepositoryResponse{
		ViolationNotes: violationNotes,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   total,
			MaxPage: totalPage,
		},
	}, nil
}

func (r *violationNoteRepository) GetViolationNoteByFilter(ctx context.Context, tx *gorm.DB, noteType []string, req dto.PaginationRequest, violationId int64, metadata external.MetadataDTO) (dto.GetAllViolationNoteRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var violationNote []model.ViolationNote
	var total int64

	// Log the inputs for debugging
	log.Printf("Filter parameters - noteType: %v, violationId: %d, officePKID: %d, page: %d, perPage: %d", noteType, violationId, metadata.OfficePKID, req.Page, req.PerPage)

	applyPaginationDefaults(&req.Page, &req.PerPage)

	// Log after pagination default applied
	log.Printf("After applying pagination defaults - page: %d, perPage: %d", req.Page, req.PerPage)

	// Build the query with filters
	query := tx.WithContext(ctx).Model(&model.ViolationNote{}).
		Where("violation_pkid = ? AND note_type IN (?) AND office_pkid = ?", violationId, noteType, metadata.OfficePKID)

	// Log query for debugging
	log.Println("Built query:", query.Statement.SQL.String())

	// Count total records with the applied filters
	if err := query.Count(&total).Error; err != nil {
		log.Printf("Error counting total records: %v", err)
		return dto.GetAllViolationNoteRepositoryResponse{}, err
	}

	// Log total count for debugging
	log.Printf("Total records after filtering: %d", total)

	// Fetch the records with pagination
	if err := query.Scopes(Paginate(req.Page, req.PerPage)).Find(&violationNote).Error; err != nil {
		log.Printf("Error fetching paginated records: %v", err)
		return dto.GetAllViolationNoteRepositoryResponse{}, err
	}

	// Log the results after query execution
	log.Printf("Fetched violation notes: %v", violationNote)

	// Calculate total pages
	totalPage := int64(math.Ceil(float64(total) / float64(req.PerPage)))

	// Log pagination response for debugging
	log.Printf("Pagination - Page: %d, PerPage: %d, Total: %d, TotalPage: %d", req.Page, req.PerPage, total, totalPage)

	// Return the response
	return dto.GetAllViolationNoteRepositoryResponse{
		ViolationNotes: violationNote,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   total,
			MaxPage: totalPage,
		},
	}, nil
}

func (r *violationNoteRepository) GetViolationNoteById(ctx context.Context, tx *gorm.DB, violationPkid int64, notePkid int64) (model.ViolationNote, error) {
	if tx == nil {
		tx = r.db
	}

	var violationNote model.ViolationNote
	if err := tx.WithContext(ctx).Unscoped().Where("pkid = ? AND violation_pkid = ?", notePkid, violationPkid).Take(&violationNote).Error; err != nil {
		return model.ViolationNote{}, err
	}

	return violationNote, nil
}

func (r *violationNoteRepository) UpdateViolationNotes(ctx context.Context, tx *gorm.DB, notePkid int64, violationPkid int64, updatedData model.ViolationNote) (model.ViolationNote, error) {
	if tx == nil {
		tx = r.db
	}

	var violationNote model.ViolationNote
	if err := tx.WithContext(ctx).Unscoped().Where("pkid = ? AND violation_pkid = ?", notePkid, violationPkid).Take(&violationNote).Error; err != nil {
		return model.ViolationNote{}, err
	}

	if err := tx.WithContext(ctx).Model(&violationNote).Updates(updatedData).Error; err != nil {
		return model.ViolationNote{}, err
	}

	return violationNote, nil
}

func (r *violationNoteRepository) RestoreViolationNotes(ctx context.Context, tx *gorm.DB, notePkid int64, violationPkid int64, violationNote model.ViolationNote) (model.ViolationNote, error) {
	if tx == nil {
		tx = r.db
	}

	now := time.Now()
	updateData := map[string]interface{}{
		"deleted_at":   nil,
		"deleted_by":   nil,
		"deleted_host": nil,
		"is_deleted":   false,
		"updated_by":   violationNote.UpdatedBy,
		"updated_at":   now,
		"office_pkid":  violationNote.OfficePkid,
	}

	err := tx.WithContext(ctx).Model(&model.ViolationNote{}).Unscoped().Where("pkid = ? AND violation_pkid = ?", notePkid, violationPkid).Updates(updateData).Error
	if err != nil {
		return model.ViolationNote{}, err
	}

	return violationNote, nil
}

func (r *violationNoteRepository) HardDeleteViolationNotes(ctx context.Context, notePkid int64, violationPkid int64) error {
	var violationNote model.ViolationNote
	if err := r.db.WithContext(ctx).Unscoped().Where("pkid = ? AND violation_pkid = ?", notePkid, violationPkid).Delete(&violationNote).Error; err != nil {
		return err
	}

	return nil
}
