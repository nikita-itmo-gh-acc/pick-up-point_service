package services

import (
	"context"
	"fmt"
	"pvz_service/mappers"
	"pvz_service/objects"
	"pvz_service/repos"

	"github.com/google/uuid"
)

type ReceptionService interface {
	StartReception(ctx context.Context, receptionData objects.ReceptionDto) (*objects.ReceptionDto, error)
	CloseReception(ctx context.Context, pvzId uuid.UUID) (*objects.ReceptionDto, error)
}

type RecpService struct {
	ReceptionRepo repos.ReceptionRepository
}

func NewReceptionService(repo repos.ReceptionRepository) *RecpService{
	return &RecpService{
		ReceptionRepo: repo,
	}
}

func (s *RecpService) StartReception(ctx context.Context, receptionData objects.ReceptionDto) (*objects.ReceptionDto, error) {
	reception, err := mappers.DtoToReception(receptionData)
	if err != nil {
		fmt.Println("Can't map reception -", err)
		return nil, err 
	}
	lastReception, _ := s.ReceptionRepo.FindLastByTime(ctx, reception.PvzId)
	if lastReception != nil && lastReception.Status == "in_progress" {
		return nil, fmt.Errorf("can't start new reception, other in progress")
	}
	reception.Id = uuid.New()
	reception.Status = "in_progress"
	if err := s.ReceptionRepo.Create(ctx, reception); err != nil {
		fmt.Println("Can't start reception -", err)
		return nil, err
	}
	responseDto, _ := mappers.ReceptionToDto(*reception)
	return responseDto, nil
}

func (s *RecpService) CloseReception(ctx context.Context, pvzId uuid.UUID) (*objects.ReceptionDto, error) {
	reception, err := s.ReceptionRepo.FindLastByTime(ctx, pvzId)
	if err != nil {
		return nil, err
	}
	if reception.Status == "close" {
		return nil, fmt.Errorf("no open receptions found")
	}

	if err := s.ReceptionRepo.FastUpdate(ctx, reception.Id); err != nil {
		return nil, err
	}

	responseDto, _ := mappers.ReceptionToDto(*reception)
	responseDto.Status = "close"
	return responseDto, nil
}
