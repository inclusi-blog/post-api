package util

import "testing"

func TestUuidGenerator_Generate(t *testing.T) {
	uuid := NewUUIDGenerator().Generate()
	if len(uuid.String()) != 36 {
		t.Errorf("Id %s is invalid with length %d", uuid, len(uuid))
	}
}
