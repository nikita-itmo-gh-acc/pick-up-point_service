package mappers

import (
	"pvz_service/objects"
	"time"
)

func ReceptionToDto(reception objects.Reception) (*objects.ReceptionDto, error) {
	timeStr := reception.DateTime.Format(Layout)
	productsDto := make([]*objects.ProductDto, 0, len(reception.Products))
	for _, prod := range reception.Products {
		dto, _ := ProductToDto(*prod)
		productsDto = append(productsDto, dto)
	}
	return &objects.ReceptionDto{
		Id: reception.Id,
		DateTime: timeStr,
		PvzId: reception.PvzId,
		Status: reception.Status,
		Products: productsDto,
	}, nil
}

func DtoToReception(dto objects.ReceptionDto) (*objects.Reception, error) {
	dateTime, err := time.Parse(Layout, dto.DateTime)
	if err != nil {
		dateTime = time.Now()
	}
	products := make([]*objects.Product, 0, len(dto.Products))
	for _, prod := range dto.Products {
		model, err := DtoToProduct(*prod)
		if err != nil {
			continue
		}
		products = append(products, model)
	}
	return &objects.Reception{
		Id: dto.Id,
		DateTime: dateTime,
		PvzId: dto.PvzId,
		Status: dto.Status,
		Products: products,
	}, nil
}
