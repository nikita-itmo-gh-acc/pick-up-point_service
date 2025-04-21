package repos

import (
	"context"
	"database/sql"
	"fmt"
	"pvz_service/objects"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *objects.User) error
	GetById(ctx context.Context, id uuid.UUID) (*objects.User, error)
	GetByEmail(ctx context.Context, email string) (*objects.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo{
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, user *objects.User) error {
	query := `INSERT INTO "user" ("id", "email", "password", "role") VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, user.Id, user.Email, user.Password, user.Role)
	if err != nil {
		fmt.Println("Can't create user -", err)
		return err
	}

	return nil
}

func (r *userRepo) GetById(ctx context.Context, id uuid.UUID) (*objects.User, error) {
	query := `SELECT * FROM "user" WHERE id = $1`
	user := objects.User{}
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Email, &user.Password, &user.Role); err != nil {
		fmt.Println("Can't find user -", err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*objects.User, error) {
	query := `SELECT * FROM "user" WHERE email = $1 LIMIT 1`
	user := objects.User{}
	if err := r.db.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Email, &user.Password, &user.Role); err != nil {
		fmt.Println("Can't find user -", err)
		return nil, err
	}
	return &user, nil
}
