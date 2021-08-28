package util

import (
	"fmt"
	"github.com/google/uuid"
)

type UUIDGenerator interface {
	Generate() string
}

type uuidGenerator struct {
}

func (u uuidGenerator) Generate() string {
	return fmt.Sprintf("%s", uuid.New())
}

func NewUUIDGenerator() UUIDGenerator {
	return uuidGenerator{}
}
