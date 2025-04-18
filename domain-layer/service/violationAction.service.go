package service

import (
	"BE-DroneMEQ-Violation-Management/data-layer/model"
	"BE-DroneMEQ-Violation-Management/data-layer/repository"
	"BE-DroneMEQ-Violation-Management/helpers/constants"
	"BE-DroneMEQ-Violation-Management/helpers/dto"
	"BE-DroneMEQ-Violation-Management/helpers/dto/external"
	"context"
)

type (
	ViolationActionService interface {
		CreateViolationAction(ctx context.Context, req dto.ViolationActionCreateRequest, Pkid int64, metadata external.MetadataDTO) (dto.ViolationActionCreateResponse, error)
	}

	violationActionService struct {
		violationActionRepo repository.ViolationActionRepository
	}
)

func NewViolationActionService(violationActionRepo repository.ViolationActionRepository) ViolationActionService {
	return &violationActionService{
		violationActionRepo: violationActionRepo,
	}
}

func (s *violationActionService) CreateViolationAction(ctx context.Context, req dto.ViolationActionCreateRequest, ViolationID int64, metadata external.MetadataDTO) (dto.ViolationActionCreateResponse, error) {
	violationAction := model.ViolationAction{
		ViolationPkid: ViolationID,
		ActionTaken:   req.ActionTaken,
		Outcome:       constants.ViolationOutcome(req.Outcome),
		TakenByPkid:   req.TakenByPkid,
	}
	violationAction.OfficePkid = &metadata.OfficePKID
	violationAction.CreatedBy = metadata.CreatedBy
	violationAction.CreatedHost = metadata.CreatedHost
	violationAction.CreatedAt = metadata.CreatedAt

	violationActionCreate, err := s.violationActionRepo.CreateViolationAction(ctx, nil, violationAction)
	if err != nil {
		return dto.ViolationActionCreateResponse{}, err
	}
	return dto.ViolationActionCreateResponse{
		Pkid:          violationActionCreate.Pkid,
		ViolationPkid: violationActionCreate.ViolationPkid,
		ActionTaken:   violationActionCreate.ActionTaken,
		Outcome:       string(violationActionCreate.Outcome),
		TakenByPkid:   violationActionCreate.TakenByPkid,
	}, nil
}
