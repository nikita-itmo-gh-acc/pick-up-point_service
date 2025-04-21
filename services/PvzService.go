package services

import (
	"context"

	"pvz_service/mappers"
	"pvz_service/objects"
	"pvz_service/repos"

	"time"

	"github.com/google/uuid"
)

type PvzService interface {
	CreatePvz(ctx context.Context, pvzData objects.PvzDto) (uuid.UUID, error)
	FilterPvz(ctx context.Context, query objects.PvzQuery) ([]*objects.PvzDto, error)
}

type PVZservice struct {
	PvzRepo repos.PvzRepository
}

func NewPvzService(repo repos.PvzRepository) *PVZservice{
	return &PVZservice{
		PvzRepo: repo,
	}
}

func (s *PVZservice) CreatePvz(ctx context.Context, pvzData objects.PvzDto) (uuid.UUID, error) {
	pvz, err := mappers.DtoToPvz(pvzData)
	if err != nil {
		return uuid.UUID{}, err
	}
	pvz.Id = uuid.New()
	if err := s.PvzRepo.Create(ctx, pvz); err != nil {
		return uuid.UUID{}, err
	}

	return pvz.Id, nil
}

func (s *PVZservice) FilterPvz(ctx context.Context, query objects.PvzQuery) ([]*objects.PvzDto, error) {
	offset := query.Limit * (query.Page - 1)
	start, _ := time.Parse("2006-01-02", query.StartDate)
	end, _ := time.Parse("2006-01-02", query.EndDate)
	foundPvz, err := s.PvzRepo.GetList(ctx, query.Limit, offset, start, end)
	if err != nil {
		return []*objects.PvzDto{}, err
	}
	dtoList := make([]*objects.PvzDto, 0, len(foundPvz))
	for _, pvz := range foundPvz {
		dto, _ := mappers.PvzToDto(*pvz)
		dtoList = append(dtoList, dto)
	}

	return dtoList, nil
}

