package repository

import (
	"testing"
)

func TestNonExistentAccount(t *testing.T) {
	_ = NewInMemoryRepository()

}
