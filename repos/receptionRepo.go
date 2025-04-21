package repos

import (
	"context"
	"database/sql"
	"pvz_service/objects"

	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

type ReceptionRepository interface {
	Create(ctx context.Context, reception *objects.Reception) error
	FastUpdate(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, reception *objects.Reception) error
	FindLastByTime(ctx context.Context, pvzId uuid.UUID) (*objects.Reception, error)
}

type receptionRepo struct {
	db *sql.DB
}

func NewReceptionRepo (db *sql.DB) *receptionRepo {
	return &receptionRepo{
		db: db,
	}
}

func (r *receptionRepo) Create(ctx context.Context, reception *objects.Reception) error {
	query := `INSERT INTO "reception" ("id", "dateTime", "pvzId", "status") VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, reception.Id, reception.DateTime, reception.PvzId, reception.Status)
	if err != nil {
		fmt.Println("Can't create reception -", err)
		return err
	}

	return nil
}

func (r *receptionRepo) FastUpdate(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE "reception" SET status = 'close' WHERE "id" = $1`
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		fmt.Println("Can't fast update reception -", err)
		return err
	}

	return nil
}

func (r *receptionRepo) Update(ctx context.Context, reception *objects.Reception) error {
	buf := new(bytes.Buffer)
	buf.WriteString(`UPDATE "reception" SET `)
	v := reflect.ValueOf(*reception)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		if field.IsZero() {
			continue
		}
		fieldName := fieldType.Name
		fmt.Fprintf(buf, `"%s%s" = %v, `, strings.ToLower(fieldName[:1]), fieldName[1:], field.Interface())
	}
	fmt.Fprintf(buf, `WHERE "id" = %v`, reception.Id)
	if _, err := r.db.ExecContext(ctx, buf.String()); err != nil {
		fmt.Println("Can't update reception -", err)
		return err
	}
	return nil
}

func (r *receptionRepo) FindLastByTime(ctx context.Context, pvzId uuid.UUID) (*objects.Reception, error) {
	query := `SELECT * FROM "reception" WHERE "pvzId" = $1 ORDER BY "dateTime" DESC LIMIT 1`
	reception := objects.Reception{}
	err := r.db.QueryRowContext(ctx, query, pvzId).Scan(&reception.Id, &reception.DateTime, &reception.PvzId, &reception.Status)
	if err != nil {
		fmt.Println("Can't find reception -", err)
		return nil, err
	}
	return &reception, nil
}
