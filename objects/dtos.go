package objects

import (
	"github.com/google/uuid"
)

type UserDto struct {
	Id       uuid.UUID `json:"id,omitempty"`
	Email    string `json:"email,omitempty" valid:"email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty" valid:"in(employee|moderator)"`
}

type TokenDto struct {
	Token	string `json:"token"`
}

type PvzDto struct {
	Id 					uuid.UUID `json:"id,omitempty"`
	RegistrationDate	string `json:"registrationDate"`
	City				string `json:"city" valid:"in(Москва|Санкт-Петербург|Казань)"`
	Receptions 			[]*ReceptionDto `json:"receptions,omitempty"`
}

type ReceptionDto struct {
	Id			uuid.UUID `json:"id,omitempty"`
	DateTime 	string `json:"dateTime,omitempty"`
	PvzId		uuid.UUID `json:"pvzId,omitempty"`
	Status		string `json:"status,omitempty" valid:"in(in_progress|close)"`
	Products 	[]*ProductDto `json:"products,omitempty"`
}

type ProductDto struct {
	Id			uuid.UUID `json:"id,omitempty"`
	DateTime 	string `json:"dateTime,omitempty"`
	Type		string `json:"type,omitempty" valid:"in(электроника|одежда|обувь)"`
	ReceptionId	uuid.UUID `json:"ReceptionId,omitempty"`
}

type AddProductDto struct {
	Type		string `json:"type,omitempty" valid:"in(электроника|одежда|обувь)"`
	PvzId		uuid.UUID `json:"pvzId,omitempty"`
}
