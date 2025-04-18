package service

import (
	"BE-DroneMEQ-Violation-Management/data-layer/config"
	"BE-DroneMEQ-Violation-Management/data-layer/model"
	"BE-DroneMEQ-Violation-Management/data-layer/repository"
	"mime/multipart"
	"strconv"
	"time"

	"BE-DroneMEQ-Violation-Management/helpers/dto"
	"BE-DroneMEQ-Violation-Management/helpers/dto/external"
	"context"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type (
	ViolationMediaService interface {
		CreateViolationMedia(ctx context.Context, medias []*multipart.FileHeader, ViolationPkid int64, metadata external.MetadataDTO) ([]dto.ViolationMediaCreateResponse, error)
	}

	violationMediaService struct {
		violationMediaRepo repository.ViolationMediaRepository
		violationRepo      repository.ViolationRepository
	}
)

func NewViolationMediaService(violationMediaRepo repository.ViolationMediaRepository, violationRepo repository.ViolationRepository) ViolationMediaService {
	return &violationMediaService{
		violationMediaRepo: violationMediaRepo,
		violationRepo:      violationRepo,
	}
}

func (s *violationMediaService) CreateViolationMedia(ctx context.Context, medias []*multipart.FileHeader, violationPkid int64, metadata external.MetadataDTO) ([]dto.ViolationMediaCreateResponse, error) {
	//get violation by pkid
	violation, _, err := s.violationRepo.GetViolationById(ctx, nil, violationPkid)
	if err != nil {
		return nil, err
	}

	var responses []dto.ViolationMediaCreateResponse
	cldService, ctxBg := config.Credentials()

	for _, media := range medias {
		file, err := media.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		res, err := cldService.Upload.Upload(ctxBg, file, uploader.UploadParams{
			PublicID:       "violation_media/" + strconv.FormatInt(violation.Pkid, 10) + "/" + media.Filename,
			UniqueFilename: api.Bool(false),
			Overwrite:      api.Bool(true),
		})

		if err != nil {
			return nil, err
		}

		violationMedia := model.ViolationMedia{
			ViolationPkid:  violation.Pkid,
			UploadedByPkid: violation.ReportedByPkid,
			MediaURL:       res.SecureURL,
			MediaType:      media.Header.Get("Content-Type"),
			OfficePkid:     &metadata.OfficePKID,
			CreatedBy:      metadata.CreatedBy,
			CreatedHost:    metadata.CreatedHost,
			UploadedAt:     time.Now(),
			CreatedAt:      time.Now(),
		}

		violationMediaCreate, err := s.violationMediaRepo.CreateViolationMedia(ctx, nil, violationMedia)
		if err != nil {
			return nil, err
		}


		responses = append(responses, dto.ViolationMediaCreateResponse{
			Pkid:          violationMediaCreate.Pkid,
			ViolationPkid: violationMediaCreate.ViolationPkid,
			MediaURL:      violationMediaCreate.MediaURL,
			MediaType:     violationMediaCreate.MediaType,
		})
	}

	return responses, nil
}
