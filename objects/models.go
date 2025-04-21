package objects

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID
	Email    string
	Password string
	Role     string
}

type Pvz struct {
	Id               uuid.UUID
	RegistrationDate time.Time
	City             string
	Receptions		 []*Reception
}

type Reception struct {
	Id       uuid.UUID
	DateTime time.Time
	PvzId    uuid.UUID
	Status   string
	Products []*Product	
}

type Product struct {
	Id          uuid.UUID
	DateTime    time.Time
	Type        string
	ReceptionId uuid.UUID
}
