package main

import (
	"context"
	"fmt"
	"pvz_service/objects"
	"pvz_service/repos"
	"pvz_service/repos/mocks"
	"testing"

	"github.com/google/uuid"

	"os"
	"pvz_service/database"
)

var (
	dbUrl = os.Getenv("DB_URL")
	baseUrl = os.Getenv("POSTGRES_URL")
	dbName = os.Getenv("DB_NAME")
	port = os.Getenv("SERVE_PORT")
	migrationsDir = os.Getenv("MIGRATION_DIR")
)

var conn = &database.DBConnection{
	DbName: dbName,
	URL: dbUrl,
	BaseURL: baseUrl,
}


func TestUserRepo_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		user objects.User
	}

	testCases := []struct{
		name string
		args args
		wantErr	 bool
	}{
		{
			name: "basic user creation",
			args: args{
				ctx: context.Background(),
				user: objects.User{
					Id: uuid.New(), Email: "dafa@mail.ru", Password: "12345", Role: "employee",
				},
			}, 
			wantErr: false,
		},
		{
			name: "again with same password",
			args: args{
				ctx: context.Background(),
				user: objects.User{
					Id: uuid.New(), Email: "dafa@mail.ru", Password: "23423", Role: "moderator",
				},
			},
			wantErr: true,
		},
	}

	mockUserRepo := new(mocks.UserRepository)
	userRepo := repos.NewUserRepo(conn.DB)
	conn.InitPostgresConn()

	for _, ts := range testCases {
		t.Run(ts.name, func(t *testing.T) {
			
			var err error = nil
			if ts.wantErr {
				err = fmt.Errorf("db error")
			}
			mockUserRepo.On("Create", ts.args.ctx, ts.args.user).Return(err)

			realErr := userRepo.Create(context.Background(), testUser) 
		})
	}
}