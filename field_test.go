package gomongo_client

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestBaseField_DefaultID(t *testing.T) {
	b := &BaseField{}

	if !b.ID.IsZero() {
		t.Fatal("new BaseField should have zero ID")
	}

	b.DefaultID()
	if b.ID.IsZero() {
		t.Error("DefaultID() should set a non-zero ID")
	}

	// Calling again should not change the ID.
	original := b.ID
	b.DefaultID()
	if b.ID != original {
		t.Error("DefaultID() should not overwrite an existing ID")
	}
}

func TestBaseField_SetID(t *testing.T) {
	b := &BaseField{}
	id := bson.NewObjectID()

	b.SetID(id)
	if b.ID != id {
		t.Errorf("SetID() ID = %v, want %v", b.ID, id)
	}
}

func TestBaseField_DefaultCreatedAt(t *testing.T) {
	b := &BaseField{}
	before := time.Now()

	b.DefaultCreatedAt()

	if b.CreatedAt.Before(before) {
		t.Error("DefaultCreatedAt() set a time in the past")
	}
}

func TestBaseField_DefaultUpdatedAt(t *testing.T) {
	b := &BaseField{}
	before := time.Now()

	b.DefaultUpdatedAt()

	if b.UpdatedAt.Before(before) {
		t.Error("DefaultUpdatedAt() set a time in the past")
	}
}

func TestBaseField_BeforeInsert(t *testing.T) {
	b := &BaseField{}
	b.BeforeInsert()

	if b.ID.IsZero() {
		t.Error("BeforeInsert() should set ID")
	}
	if b.CreatedAt.IsZero() {
		t.Error("BeforeInsert() should set CreatedAt")
	}
	if b.UpdatedAt.IsZero() {
		t.Error("BeforeInsert() should set UpdatedAt")
	}
}

func TestBaseField_BeforeUpdate(t *testing.T) {
	b := &BaseField{}
	b.BeforeInsert()

	originalCreatedAt := b.CreatedAt
	time.Sleep(time.Millisecond) // Ensure different timestamp.

	b.BeforeUpdate()

	if b.CreatedAt != originalCreatedAt {
		t.Error("BeforeUpdate() should not change CreatedAt")
	}
	if b.UpdatedAt.Equal(originalCreatedAt) || b.UpdatedAt.Before(originalCreatedAt) {
		t.Error("BeforeUpdate() should refresh UpdatedAt")
	}
}
