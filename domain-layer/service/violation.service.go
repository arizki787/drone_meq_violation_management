package service

import (
	"context"
	"log"
	"sync"
	"time"

	"BE-DroneMEQ-Violation-Management/data-layer/model"
	"BE-DroneMEQ-Violation-Management/data-layer/repository"
	"BE-DroneMEQ-Violation-Management/helpers/constants"
	"BE-DroneMEQ-Violation-Management/helpers/dto"
	"BE-DroneMEQ-Violation-Management/helpers/dto/external"

	"gorm.io/gorm"
)

type (
	ViolationService interface {
		CreateViolation(ctx context.Context, req dto.ViolationCreateRequest, metadata external.MetadataDTO) (dto.ViolationResponse, error)
		GetAllViolationWithPagination(ctx context.Context, req dto.PaginationRequest, metadata external.MetadataDTO) (dto.ViolationPaginationResponse, error)
		GetViolationById(ctx context.Context, pkid int64) (dto.ViolationResponseGetId, error)
		GetViolationByFilter(ctx context.Context, status []string, violationType []string, severity []string, req dto.PaginationRequest, metadata external.MetadataDTO) (dto.ViolationPaginationResponse, error)
		UpdateViolation(ctx context.Context, req dto.ViolationUpdateRequest, pkid int64, metadata external.MetadataDTO) (dto.ViolationUpdateResponse, error)
		SoftDeleteViolation(ctx context.Context, pkid int64, metadata external.MetadataDTO) (dto.ViolationDeleteResponse, error)
		RestoreViolation(ctx context.Context, pkid int64, metadata external.MetadataDTO) (dto.ViolationDeleteResponse, error)
		HardDeleteViolation(ctx context.Context, pkid int64) error
	}

	violationService struct {
		violationRepo repository.ViolationRepository
	}
)

func NewViolationService(violationRepo repository.ViolationRepository) ViolationService {
	return &violationService{
		violationRepo: violationRepo,
	}
}

func (s *violationService) CreateViolation(ctx context.Context, req dto.ViolationCreateRequest, metadata external.MetadataDTO) (dto.ViolationResponse, error) {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	violation := model.Violation{
		ReportedByPkid: req.ReportedByPkid,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		ViolationType:  constants.ViolationType(req.ViolationType),
		ReportedAt:     req.ReportedAt,
		Status:         constants.ViolationStatus(req.Status),
		Severity:       constants.ViolationSeverity(req.Severity),
		Description:    req.Description,
	}
	now := time.Now()
	violation.OfficePkid = &metadata.OfficePKID
	violation.CreatedBy = metadata.CreatedBy
	violation.CreatedHost = metadata.CreatedHost
	violation.CreatedAt = now

	violationCreate, err := s.violationRepo.CreateViolation(ctx, nil, violation)
	if err != nil {
		return dto.ViolationResponse{}, err
	}

	return dto.ViolationResponse{
		Pkid:           violationCreate.Pkid,
		ReportedByPkid: violationCreate.ReportedByPkid,
		Latitude:       violationCreate.Latitude,
		Longitude:      violationCreate.Longitude,
		ViolationType:  string(violationCreate.ViolationType),
		ReportedAt:     violationCreate.ReportedAt,
		Status:         string(violationCreate.Status),
		Severity:       string(violationCreate.Severity),
		Description:    violationCreate.Description,
		OfficePkid:     *violationCreate.OfficePkid,
	}, nil
}

func (s *violationService) GetAllViolationWithPagination(ctx context.Context, req dto.PaginationRequest, metadata external.MetadataDTO) (dto.ViolationPaginationResponse, error) {
	dataWithPaginate, err := s.violationRepo.GetAllViolationWithPagination(ctx, nil, req, metadata.OfficePKID)
	if err != nil {
		return dto.ViolationPaginationResponse{}, err
	}

	var datas []dto.ViolationResponse
	for _, violation := range dataWithPaginate.Violations {
		data := dto.ViolationResponse{
			Pkid:           violation.Pkid,
			ReportedByPkid: violation.ReportedByPkid,
			Latitude:       violation.Latitude,
			Longitude:      violation.Longitude,
			ViolationType:  string(violation.ViolationType),
			ReportedAt:     violation.ReportedAt,
			Status:         string(violation.Status),
			Severity:       string(violation.Severity),
			Description:    violation.Description,
			OfficePkid:     *violation.OfficePkid,
		}
		datas = append(datas, data)
	}

	return dto.ViolationPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (s *violationService) GetViolationById(ctx context.Context, pkid int64) (dto.ViolationResponseGetId, error) {
	violation, violationMedia, err := s.violationRepo.GetViolationById(ctx, nil, pkid)
	if err != nil {
		return dto.ViolationResponseGetId{}, err
	}

	var mediaResponses []dto.ViolationMediaResponse
	for _, media := range violationMedia {
		mediaResponses = append(mediaResponses, dto.ViolationMediaResponse{
			MediaURL:  media.MediaURL,
			MediaType: media.MediaType,
		})
	}

	response := dto.ViolationResponseGetId{
		Violation:      violation,
		ViolationMedia: mediaResponses,
	}

	return response, nil
}

func (s *violationService) GetViolationByFilter(ctx context.Context, status []string, violationType []string, severity []string, req dto.PaginationRequest, metadata external.MetadataDTO) (dto.ViolationPaginationResponse, error) {
	dataWithPaginate, err := s.violationRepo.GetViolationByFilter(ctx, nil, status, violationType, severity, req, metadata.OfficePKID)
	if err != nil {
		return dto.ViolationPaginationResponse{}, err
	}

	var datas []dto.ViolationResponse
	for _, violation := range dataWithPaginate.Violations {
		data := dto.ViolationResponse{
			Pkid:           violation.Pkid,
			ReportedByPkid: violation.ReportedByPkid,
			Latitude:       violation.Latitude,
			Longitude:      violation.Longitude,
			ViolationType:  string(violation.ViolationType),
			Status:         string(violation.Status),
			Severity:       string(violation.Severity),
			Description:    violation.Description,
			OfficePkid:     *violation.OfficePkid,
		}
		datas = append(datas, data)
	}

	return dto.ViolationPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (s *violationService) UpdateViolation(ctx context.Context, req dto.ViolationUpdateRequest, pkid int64, metadata external.MetadataDTO) (dto.ViolationUpdateResponse, error) {
	// Update the violation
	data := model.Violation{
		ReportedByPkid: req.ReportedByPkid,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		ReportedAt:     req.ReportedAt,
		ViolationType:  constants.ViolationType(req.ViolationType),
		Status:         constants.ViolationStatus(req.Status),
		Severity:       constants.ViolationSeverity(req.Severity),
		Description:    req.Description,
		UpdatedBy:      &metadata.CreatedBy,
		UpdatedHost:    &metadata.CreatedHost,
		UpdatedAt:      time.Now(),
	}

	violationUpdate, err := s.violationRepo.UpdateViolation(ctx, nil, pkid, data)
	if err != nil {
		return dto.ViolationUpdateResponse{}, err
	}

	return dto.ViolationUpdateResponse{
		Pkid:           pkid,
		ReportedByPkid: violationUpdate.ReportedByPkid,
		Latitude:       violationUpdate.Latitude,
		Longitude:      violationUpdate.Longitude,
		ReportedAt:     violationUpdate.ReportedAt,
		ViolationType:  string(violationUpdate.ViolationType),
		Status:         string(violationUpdate.Status),
		Severity:       string(violationUpdate.Severity),
		Description:    violationUpdate.Description,
	}, nil
}

func (s *violationService) SoftDeleteViolation(ctx context.Context, pkid int64, metadata external.MetadataDTO) (dto.ViolationDeleteResponse, error) {
	violation, _, err := s.violationRepo.GetViolationById(ctx, nil, pkid)
	if err != nil {
		return dto.ViolationDeleteResponse{}, err
	}

	now := time.Now()
	violation.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	violation.DeletedBy = &metadata.CreatedBy
	violation.DeletedHost = &metadata.CreatedHost
	violation.IsDeleted = true

	violation.UpdatedBy = &metadata.CreatedBy
	violation.UpdatedHost = &metadata.CreatedHost
	violation.UpdatedAt = now
	violation.OfficePkid = &metadata.OfficePKID

	_, err = s.violationRepo.UpdateViolation(ctx, nil, pkid, violation)
	if err != nil {
		log.Printf("Failed to soft delete violation with pkid: %d, error: %v", pkid, err)
		return dto.ViolationDeleteResponse{}, err
	}

	response := dto.ViolationDeleteResponse{
		Pkid:      pkid,
		IsDeleted: true,
	}

	return response, nil
}

func (s *violationService) RestoreViolation(ctx context.Context, pkid int64, metadata external.MetadataDTO) (dto.ViolationDeleteResponse, error) {
	violation, _, err := s.violationRepo.GetViolationById(ctx, nil, pkid)
	if err != nil {
		log.Printf("Error get violation by id: %v", err)
		return dto.ViolationDeleteResponse{}, err
	}

	violation.UpdatedAt = time.Now()
	violation.UpdatedBy = &metadata.CreatedBy
	violation.UpdatedHost = &metadata.CreatedHost
	violation.DeletedAt = gorm.DeletedAt{Time: time.Time{}, Valid: false}
	violation.IsDeleted = false
	violation.DeletedHost = nil
	violation.DeletedBy = nil
	violation.OfficePkid = &metadata.OfficePKID

	_, err = s.violationRepo.RestoreViolation(ctx, nil, pkid, violation)
	if err != nil {
		log.Printf("Error restore violation: %v", err)
		return dto.ViolationDeleteResponse{}, err
	}

	response := dto.ViolationDeleteResponse{
		Pkid:      pkid,
		IsDeleted: false,
	}
	return response, nil
}

func (s *violationService) HardDeleteViolation(ctx context.Context, pkid int64) error {
	err := s.violationRepo.HardDeleteViolation(ctx, pkid)
	if err != nil {
		return err
	}

	return nil
}
