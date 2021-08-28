package helper

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(ctx context.Context, request CreateUserRequest) (uuid.UUID, error)
}

const (
	CreateUser = "insert into users(id, email, role_id, password, username,is_active, created_at) VALUES (uuid_generate_v4(), $1, (select id from roles where name = $2), $3, $4, true, current_timestamp) returning id"
)

type userRepository struct {
	db *sqlx.DB
}

func (repo userRepository) CreateUser(ctx context.Context, request CreateUserRequest) (uuid.UUID, error) {
	var userUUID uuid.UUID
	err := repo.db.GetContext(ctx, &userUUID, CreateUser, request.Email, request.Role, request.Password, request.Username)

	if err != nil {
		return userUUID, err
	}

	return userUUID, nil
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return userRepository{db: db}
}
