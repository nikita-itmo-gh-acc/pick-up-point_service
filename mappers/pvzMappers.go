package mappers

import (
	"fmt"
	"pvz_service/objects"
	"time"
)

var Layout string = "2006-01-02T15:04:05.000Z"

func PvzToDto(pvz objects.Pvz) (*objects.PvzDto, error) {
	regDateStr := pvz.RegistrationDate.Format(Layout)
	receptionsDto := make([]*objects.ReceptionDto, 0, len(pvz.Receptions))
	for _, recp := range pvz.Receptions {
		dto, _ := ReceptionToDto(*recp)
		receptionsDto = append(receptionsDto, dto)
	}
	return &objects.PvzDto{
		Id: pvz.Id,
		RegistrationDate: regDateStr,
		City: pvz.City,
		Receptions: receptionsDto,
	}, nil
}

func DtoToPvz(dto objects.PvzDto) (*objects.Pvz, error) {
	regDateTime, err := time.Parse(Layout, dto.RegistrationDate)
	if err != nil {
		regDateTime = time.Now()
	}
	receptions := make([]*objects.Reception, 0, len(dto.Receptions))
	for _, recp := range dto.Receptions {
		model, err := DtoToReception(*recp)
		if err != nil {
			continue
		}
		receptions = append(receptions, model)
	}
	fmt.Println("In mapper:", dto, regDateTime, receptions)
	return &objects.Pvz{
		Id: dto.Id,
		RegistrationDate: regDateTime,
		City: dto.City,
		Receptions: receptions,
	}, nil
}
