package util

//go:generate mockgen -source=uuid_generator.go -destination=./../mocks/mock_uuid_generator.go -package=mocks

import (
	"github.com/google/uuid"
)

type UUIDGenerator interface {
	Generate() uuid.UUID
}

type uuidGenerator struct {
}

func (u uuidGenerator) Generate() uuid.UUID {
	return uuid.New()
}

func NewUUIDGenerator() UUIDGenerator {
	return uuidGenerator{}
}
