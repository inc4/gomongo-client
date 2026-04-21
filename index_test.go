package gomongo_client

import (
	"context"
	"errors"
	"testing"
)

func TestEnsureIndexes_NilCollection(t *testing.T) {
	_, err := EnsureIndexes(context.Background(), nil, nil)
	if !errors.Is(err, ErrNilCollection) {
		t.Errorf("EnsureIndexes(nil) error = %v, want %v", err, ErrNilCollection)
	}
}

func TestListIndexes_NilCollection(t *testing.T) {
	_, err := ListIndexes(context.Background(), nil)
	if !errors.Is(err, ErrNilCollection) {
		t.Errorf("ListIndexes(nil) error = %v, want %v", err, ErrNilCollection)
	}
}
