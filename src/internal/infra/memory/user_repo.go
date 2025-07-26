package memory

import (
	"context"
	"t03/internal/domain"
	"time"
)

func (repo *GameRepositoryImpl) SaveUser(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	entity := userToEntity(user)

	_, err := repo.storage.pool.Exec(ctx, `
		INSERT INTO users (id,user_login, user_password)
		VALUES ($1, $2, $3)
	`, entity.ID, entity.Login, entity.Password)

	return err
}

func (repo *GameRepositoryImpl) GetUser(login string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	var entity UserEntity
	err := repo.storage.pool.QueryRow(ctx, `
		SELECT id,user_login,user_password
		FROM users
		WHERE user_login = $1
	`, login).Scan(&entity.ID, &entity.Login, &entity.Password)

	if err != nil {
		return nil, err
	}

	return userToDomain(&entity), nil
}
