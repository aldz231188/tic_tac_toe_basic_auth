package memory

import (
	"t03/internal/domain"
)

func userToEntity(user *domain.User) *UserEntity {
	return &UserEntity{
		ID:       user.ID,
		Login:    user.Login,
		Password: user.Password,
	}
}

func userToDomain(entity *UserEntity) *domain.User {
	return &domain.User{
		ID:       entity.ID,
		Login:    entity.Login,
		Password: entity.Password,
	}
}
