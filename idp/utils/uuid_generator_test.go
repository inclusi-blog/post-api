package util

import "testing"

func TestUuidGenerator_Generate(t *testing.T) {
	uuid := NewUUIDGenerator().Generate()
	if len(uuid) != 36 {
		t.Errorf("UUID %s is invalid with length %d", uuid, len(uuid))
	}
}
