package service

import (
	"BE-DroneMEQ-Violation-Management/data-layer/model"
	"BE-DroneMEQ-Violation-Management/data-layer/repository"
	"BE-DroneMEQ-Violation-Management/helpers/constants"
	"BE-DroneMEQ-Violation-Management/helpers/dto"
	"BE-DroneMEQ-Violation-Management/helpers/dto/external"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type (
	ViolationNoteService interface {
		CreateViolationNote(ctx context.Context, req dto.ViolationNoteCreateRequest, violationPkid int64, metadata external.MetadataDTO) (dto.ViolationNoteCreateResponse, error)
		GetViolationNoteByViolationId(ctx context.Context, req dto.PaginationRequest, violationPkid int64, metadata external.MetadataDTO) (dto.ViolationNotePaginationResponse, error)
		GetViolationNoteByFilter(ctx context.Context, noteType []string, req dto.PaginationRequest, violationPkid int64, metadata external.MetadataDTO) (dto.ViolationNotePaginationResponse, error)
		UpdateViolationNote(ctx context.Context, req dto.ViolationNoteUpdateRequest, violationPkid int64, notePkid int64, metadata external.MetadataDTO) (dto.ViolationNoteUpdateResponse, error)
		SoftDeleteViolationNote(ctx context.Context, violationPkid int64, notePkid int64, metadata external.MetadataDTO) (dto.ViolationNoteDeleteResponse, error)
		RestoreViolationNote(ctx context.Context, violationPkid int64, notePkid int64, metadata external.MetadataDTO) (dto.ViolationNoteDeleteResponse, error)
		HardDeleteViolationNote(ctx context.Context, violationPkid int64, notePkid int64) error
	}

	violationNoteService struct {
		violationNoteRepo repository.ViolationNoteRepository
	}
)

func NewViolationNoteService(violationNoteRepo repository.ViolationNoteRepository) ViolationNoteService {
	return &violationNoteService{
		violationNoteRepo: violationNoteRepo,
	}
}

func (s *violationNoteService) CreateViolationNote(ctx context.Context, req dto.ViolationNoteCreateRequest, ViolationID int64, metadata external.MetadataDTO) (dto.ViolationNoteCreateResponse, error) {
	violationNote := model.ViolationNote{
		ViolationPkid: ViolationID,
		NoteType:      constants.ViolationNoteType(req.NoteType),
		Note:          req.Note,
	}
	now := time.Now()
	violationNote.OfficePkid = &metadata.OfficePKID
	violationNote.CreatedBy = metadata.CreatedBy
	violationNote.CreatedHost = metadata.CreatedHost
	violationNote.CreatedAt = now

	violationNoteCreate, err := s.violationNoteRepo.CreateViolationNotes(ctx, nil, violationNote)
	if err != nil {
		return dto.ViolationNoteCreateResponse{}, err
	}
	return dto.ViolationNoteCreateResponse{
		Pkid:          violationNoteCreate.Pkid,
		ViolationPkid: violationNoteCreate.ViolationPkid,
		NoteType:      string(violationNoteCreate.NoteType),
		Note:          violationNoteCreate.Note,
		OfficePkid:    *violationNoteCreate.OfficePkid,
	}, nil
}

func (s *violationNoteService) GetViolationNoteByViolationId(ctx context.Context, req dto.PaginationRequest, violationPkid int64, metadata external.MetadataDTO) (dto.ViolationNotePaginationResponse, error) {
	dataWithPaginate, err := s.violationNoteRepo.GetAllViolationNoteWithPagination(ctx, nil, req, violationPkid, metadata)
	if err != nil {
		return dto.ViolationNotePaginationResponse{}, err
	}

	var datas []dto.ViolationNoteResponse
	for _, violationNote := range dataWithPaginate.ViolationNotes {
		data := dto.ViolationNoteResponse{
			Pkid:          violationNote.Pkid,
			ViolationPkid: violationNote.ViolationPkid,
			NoteType:      string(violationNote.NoteType),
			Note:          violationNote.Note,
			OfficePkid:    *violationNote.OfficePkid,
		}
		datas = append(datas, data)
	}

	return dto.ViolationNotePaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (s *violationNoteService) GetViolationNoteByFilter(ctx context.Context, noteType []string, req dto.PaginationRequest, violationPkid int64, metadata external.MetadataDTO) (dto.ViolationNotePaginationResponse, error) {
	dataWithPaginate, err := s.violationNoteRepo.GetViolationNoteByFilter(ctx, nil, noteType, req, violationPkid, metadata)
	if err != nil {
		return dto.ViolationNotePaginationResponse{}, err
	}

	var datas []dto.ViolationNoteResponse
	for _, violationNote := range dataWithPaginate.ViolationNotes {
		data := dto.ViolationNoteResponse{
			Pkid:          violationNote.Pkid,
			ViolationPkid: violationNote.ViolationPkid,
			NoteType:      string(violationNote.NoteType),
			Note:          violationNote.Note,
			OfficePkid:    *violationNote.OfficePkid,
		}
		datas = append(datas, data)
	}

	return dto.ViolationNotePaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (s *violationNoteService) UpdateViolationNote(ctx context.Context, req dto.ViolationNoteUpdateRequest, violationPkid int64, notePkid int64, metadata external.MetadataDTO) (dto.ViolationNoteUpdateResponse, error) {
	violationNote := model.ViolationNote{
		NoteType:    constants.ViolationNoteType(req.NoteType),
		Note:        req.Note,
		UpdatedBy:   &metadata.CreatedBy,
		UpdatedHost: &metadata.CreatedHost,
		UpdatedAt:   time.Now(),
	}
	violationNoteUpdate, err := s.violationNoteRepo.UpdateViolationNotes(ctx, nil, notePkid, violationPkid, violationNote)
	if err != nil {
		return dto.ViolationNoteUpdateResponse{}, err
	}

	return dto.ViolationNoteUpdateResponse{
		Pkid:          violationNoteUpdate.Pkid,
		ViolationPkid: violationNoteUpdate.ViolationPkid,
		NoteType:      string(violationNoteUpdate.NoteType),
		Note:          violationNoteUpdate.Note,
		OfficePkid:    *violationNoteUpdate.OfficePkid,
	}, nil
}

func (s *violationNoteService) SoftDeleteViolationNote(ctx context.Context, violationPkid int64, notePkid int64, metadata external.MetadataDTO) (dto.ViolationNoteDeleteResponse, error) {
	violationNote, err := s.violationNoteRepo.GetViolationNoteById(ctx, nil, violationPkid, notePkid)
	if err != nil {
		return dto.ViolationNoteDeleteResponse{}, err
	}
	now := time.Now()
	violationNote.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	violationNote.DeletedBy = &metadata.CreatedBy
	violationNote.DeletedHost = &metadata.CreatedBy
	violationNote.IsDeleted = true

	violationNote.UpdatedBy = &metadata.CreatedBy
	violationNote.UpdatedHost = &metadata.CreatedHost
	violationNote.UpdatedAt = now
	violationNote.OfficePkid = &metadata.OfficePKID

	_, err = s.violationNoteRepo.UpdateViolationNotes(ctx, nil, notePkid, violationPkid, violationNote)
	if err != nil {
		return dto.ViolationNoteDeleteResponse{}, err
	}

	response := dto.ViolationNoteDeleteResponse{
		Pkid:          violationNote.Pkid,
		ViolationPkid: violationNote.ViolationPkid,
		Note:          violationNote.Note,
		NoteType:      string(violationNote.NoteType),
		IsDeleted:     violationNote.IsDeleted,
		OfficePkid:    *violationNote.OfficePkid,
	}

	return response, nil
}

func (s *violationNoteService) RestoreViolationNote(ctx context.Context, violationPkid int64, notePkid int64, metadata external.MetadataDTO) (dto.ViolationNoteDeleteResponse, error) {
	violationNote, err := s.violationNoteRepo.GetViolationNoteById(ctx, nil, violationPkid, notePkid)
	if err != nil {
		return dto.ViolationNoteDeleteResponse{}, errors.New("violation note not found")
	}
	now := time.Now()
	violationNote.DeletedAt = gorm.DeletedAt{Time: now, Valid: false}
	violationNote.DeletedBy = nil
	violationNote.DeletedHost = nil
	violationNote.IsDeleted = false

	violationNote.UpdatedBy = &metadata.CreatedBy
	violationNote.UpdatedHost = &metadata.CreatedHost
	violationNote.UpdatedAt = now
	violationNote.OfficePkid = &metadata.OfficePKID

	_, err = s.violationNoteRepo.RestoreViolationNotes(ctx, nil, notePkid, violationPkid, violationNote)
	if err != nil {
		return dto.ViolationNoteDeleteResponse{}, err
	}

	response := dto.ViolationNoteDeleteResponse{
		Pkid:          violationNote.Pkid,
		ViolationPkid: violationNote.ViolationPkid,
		Note:          violationNote.Note,
		NoteType:      string(violationNote.NoteType),
		IsDeleted:     violationNote.IsDeleted,
		OfficePkid:    *violationNote.OfficePkid,
	}

	return response, nil
}

func (s *violationNoteService) HardDeleteViolationNote(ctx context.Context, violationPkid int64, notePkid int64) error {
	err := s.violationNoteRepo.HardDeleteViolationNotes(ctx, notePkid, violationPkid)
	if err != nil {
		return err
	}
	return nil
}
