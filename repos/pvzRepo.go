package repos

import (
	"context"
	"fmt"
	"pvz_service/objects"

	"database/sql"

	"time"

	"github.com/google/uuid"
)

type PvzRepository interface{
	Create(ctx context.Context, pvz *objects.Pvz) error
	GetById(ctx context.Context, id uuid.UUID) (*objects.Pvz, error)
	GetList(ctx context.Context, limit, offset int, from, to time.Time) ([]*objects.Pvz, error)
}

type pvzRepo struct {
	db *sql.DB
}

func NewPvzRepo(db *sql.DB) *pvzRepo{
	return &pvzRepo{
		db: db,
	}
}

func (r *pvzRepo) Create(ctx context.Context, pvz *objects.Pvz) error {
	query := `INSERT INTO "pvz" ("id", "registrationDate", "city") VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, pvz.Id, pvz.RegistrationDate, pvz.City)
	if err != nil {
		fmt.Println("Can't create pvz -", err)
		return err
	}

	return nil
}

func (r *pvzRepo) GetById(ctx context.Context, id uuid.UUID) (*objects.Pvz, error) {
	query := `SELECT * FROM "pvz" WHERE id = $1`
	pvz := objects.Pvz{}
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&pvz.Id, &pvz.RegistrationDate, &pvz.City); err != nil {
		fmt.Println("Can't find pvz -", err)
		return nil, err
	}

	return &pvz, nil
}

func (r *pvzRepo) GetList(ctx context.Context, limit, offset int, from, to time.Time) ([]*objects.Pvz, error) {
	queryStr := ` 
				SELECT 	p."id", p."registrationDate", p."city", 
						r."id", r."dateTime", r."status", 
						pr."id", pr."dateTime", pr."type"
				FROM "pvz" p
				JOIN "reception" r ON p."id" = r."pvzId"
				JOIN "product" pr ON r."id" = pr."receptionId"
				WHERE r."dateTime" BETWEEN $1 AND $2
				LIMIT $3 OFFSET $4
	`
	rows, err := r.db.QueryContext(ctx, queryStr, from, to, limit, offset)

	if err != nil {
		fmt.Println("error during pvz search - ", err)
		return nil, err
	}

	defer rows.Close()

	reserved_capacity := 100
	noRowsFound := true
	foundPvz := make([]*objects.Pvz, 0, reserved_capacity)
	pvzMap := map[uuid.UUID]*objects.Pvz{}
    receptionMap := map[uuid.UUID]*objects.Reception{}

	for rows.Next() {
		noRowsFound = false
		var (
			pvzId, receptionId, productId uuid.UUID
			pvzRegDate, receptionTime, productTime time.Time
			city, status, productType string
		)

		if err := rows.Scan(&pvzId, &pvzRegDate, &city,
							&receptionId, &receptionTime, &status,
							&productId, &productTime, &productType); err != nil {
			fmt.Println("can't scan joined table row:", err)
            continue
		}
		
		if _, ok := pvzMap[pvzId]; !ok {
			pvz := &objects.Pvz{Id: pvzId, RegistrationDate: pvzRegDate, City: city}
			pvzMap[pvzId] = pvz
			foundPvz = append(foundPvz, pvz)
		}

		if _, ok := receptionMap[receptionId]; !ok {
			r := &objects.Reception{Id: receptionId, DateTime: receptionTime, Status: status, PvzId: pvzId}
			receptionMap[receptionId] = r
			pvzMap[pvzId].Receptions = append(pvzMap[pvzId].Receptions, r)
		}

		product := objects.Product{Id: productId, DateTime: productTime, Type: productType, ReceptionId: receptionId}
		receptionMap[receptionId].Products = append(receptionMap[receptionId].Products, &product)
	}

	if noRowsFound {
		return nil, sql.ErrNoRows
	}

	return foundPvz, nil
}
