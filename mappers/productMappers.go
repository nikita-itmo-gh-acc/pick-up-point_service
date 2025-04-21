package mappers

import (
	"pvz_service/objects"
	"time"
)

func ProductToDto(product objects.Product) (*objects.ProductDto, error) {
	timeStr := product.DateTime.Format(Layout)
	return &objects.ProductDto{
		Id: product.Id,
		DateTime: timeStr,
		Type: product.Type,
		ReceptionId: product.ReceptionId,
	}, nil
}

func DtoToProduct(dto objects.ProductDto) (*objects.Product, error) {
	dateTime, err := time.Parse(Layout, dto.DateTime)
	if err != nil {
		dateTime = time.Now()
	}
	return &objects.Product{
		Id: dto.Id,
		DateTime: dateTime,
		Type: dto.Type,
		ReceptionId: dto.ReceptionId,
	}, nil
}
