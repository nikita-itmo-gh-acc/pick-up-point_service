package services

import (
	"context"
	"fmt"
	"pvz_service/mappers"
	"pvz_service/objects"
	"pvz_service/repos"

	"time"

	"github.com/google/uuid"
)

type ProductService interface {
	DeleteLast(ctx context.Context, pvzId uuid.UUID) error
	Add(ctx context.Context, productData objects.AddProductDto) (*objects.ProductDto, error)
}

type ProdService struct {
	ProductRepo repos.ProductRepository
	ReceptionRepo repos.ReceptionRepository
}

func NewProductService(repo repos.ProductRepository, recpRepo repos.ReceptionRepository) *ProdService{
	return &ProdService{
		ProductRepo: repo,
		ReceptionRepo: recpRepo,
	}
}

func (s *ProdService) DeleteLast(ctx context.Context, pvzId uuid.UUID) error {
	product, err := s.ProductRepo.FindLastByTime(ctx, pvzId)
	if err != nil {
		return err
	}
	
	if err := s.ProductRepo.Delete(ctx, product.Id); err != nil {
		return err
	}

	return nil
}

func (s *ProdService) Add(ctx context.Context, productData objects.AddProductDto) (*objects.ProductDto, error) {
	product := &objects.Product{ Type: productData.Type }

	reception, _ := s.ReceptionRepo.FindLastByTime(ctx, productData.PvzId)
	if reception == nil || reception.Status == "closed" {
		return nil, fmt.Errorf("can't add product - no receptions available")
	}

	product.ReceptionId = reception.Id
	product.Id = uuid.New()
	product.DateTime = time.Now()

	if err := s.ProductRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	responseDto, _ := mappers.ProductToDto(*product)
	return responseDto, nil
}
