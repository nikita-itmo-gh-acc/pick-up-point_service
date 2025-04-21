package mappers

import (
	"pvz_service/objects"
)

func UserToDto(user objects.User) (*objects.UserDto) {
	return &objects.UserDto{
		Id: user.Id,
		Email: user.Email,
		Password: user.Password,
		Role: user.Role,
	}
}

func DtoToUser(dto objects.UserDto) (*objects.User) {
	return &objects.User{
		Id: dto.Id,
		Email: dto.Email,
		Password: dto.Password,
		Role: dto.Role,
	}
}
