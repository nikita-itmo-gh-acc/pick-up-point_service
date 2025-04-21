package repos

import (
	"context"
	"database/sql"
	"fmt"
	"pvz_service/objects"

	"github.com/google/uuid"
)

type ProductRepository interface{
	Create(ctx context.Context, product *objects.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindLastByTime(ctx context.Context, pvzId uuid.UUID) (*objects.Product, error)
}

type productRepo struct {
	db *sql.DB
}

func NewProductRepo (db *sql.DB) *productRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) Create(ctx context.Context, product *objects.Product) error {
	query := `INSERT INTO "product" ("id", "dateTime", "type", "receptionId") VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, product.Id, product.DateTime, product.Type, product.ReceptionId)
	if err != nil {
		fmt.Println("Can't create product -", err)
		return err
	}

	return nil
}

func (r *productRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM "product" WHERE "id" = $1`
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		fmt.Println("can't delete from product table -", err)
		return err
	}
	return nil
}

func (r *productRepo) FindLastByTime(ctx context.Context, pvzId uuid.UUID) (*objects.Product, error) {
	query := `
		SELECT p."id", p."dateTime", p."type", p."receptionId" FROM "product" p
		JOIN "reception" r ON p."receptionId" = r."id"
		WHERE r."pvzId" = $1
		ORDER BY "dateTime" DESC LIMIT 1
	`
	product := objects.Product{}
	err := r.db.QueryRowContext(ctx, query, pvzId).Scan(&product.Id, &product.DateTime, &product.Type, &product.ReceptionId)
	if err != nil {
		fmt.Println("Can't find product -", err)
		return nil, err
	}
	return &product, nil
}
