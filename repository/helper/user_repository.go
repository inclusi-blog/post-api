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
	CreateUser = "insert into users(id, email, role_id) VALUES (uuid_generate_v4(), $1, (select id from roles where name = $2)) returning id"
)

type userRepository struct {
	db *sqlx.DB
}

func (repo userRepository) CreateUser(ctx context.Context, request CreateUserRequest) (uuid.UUID, error) {
	var userUUID uuid.UUID
	err := repo.db.GetContext(ctx, &userUUID, CreateUser, request.Email, request.Role)

	if err != nil{
		return userUUID, err
	}

	return userUUID, nil
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return userRepository{db: db}
}
