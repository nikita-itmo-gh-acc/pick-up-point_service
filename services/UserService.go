package services

import (
	"context"
	"pvz_service/mappers"
	"pvz_service/objects"
	"pvz_service/repos"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"fmt"
)

type UserService interface {
	Register(ctx context.Context, regData objects.UserDto) (uuid.UUID, error)
	Login(ctx context.Context, loginData objects.UserDto) error
}

	type UService struct {
		UserRepo repos.UserRepository
	}

func NewUserService(repo repos.UserRepository) *UService{
	return &UService{
		UserRepo: repo,
	}
}

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func verifyPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func (s *UService) Register(ctx context.Context, regData objects.UserDto) (uuid.UUID, error) {
	user := mappers.DtoToUser(regData)
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		fmt.Println("Can't hash user password -", err)
		return uuid.UUID{}, err
	}

	user.Password = hashedPassword
	user.Id = uuid.New()
	if err := s.UserRepo.Create(ctx, user); err != nil {
		fmt.Println("Can't register new user -", err)
		return uuid.UUID{}, err
	}

	return user.Id, nil
}

func (s *UService) Login(ctx context.Context, loginData objects.UserDto) error {
	user := mappers.DtoToUser(loginData)
	found, err := s.UserRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		fmt.Println("Can't find user by email -", err)
		return err
	}

	if !verifyPassword(user.Password, found.Password) {
		return fmt.Errorf("wrong password")
	}

	return nil
}
