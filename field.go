package gomongo_client

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// BaseField provides common fields for MongoDB documents: a unique ID
// and creation/update timestamps. Embed this struct in your document
// models to get automatic ID generation and timestamp management.
//
// Example:
//
//	type User struct {
//	    mongorm.BaseField `bson:",inline"`
//	    Name  string      `json:"name" bson:"name"`
//	    Email string      `json:"email" bson:"email"`
//	}
type BaseField struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// SetID sets the document's ID to the given value.
func (b *BaseField) SetID(id bson.ObjectID) {
	b.ID = id
}

// DefaultID generates a new ObjectID if the current ID is the zero value.
func (b *BaseField) DefaultID() {
	if b.ID.IsZero() {
		b.ID = bson.NewObjectID()
	}
}

// DefaultCreatedAt sets CreatedAt to the current time.
func (b *BaseField) DefaultCreatedAt() {
	b.CreatedAt = time.Now()
}

// DefaultUpdatedAt sets UpdatedAt to the current time.
func (b *BaseField) DefaultUpdatedAt() {
	b.UpdatedAt = time.Now()
}

// BeforeInsert populates default fields before inserting a new document.
// It generates an ID (if zero) and sets both timestamps to now.
func (b *BaseField) BeforeInsert() {
	b.DefaultID()
	b.DefaultCreatedAt()
	b.DefaultUpdatedAt()
}

// BeforeUpdate refreshes the UpdatedAt timestamp before updating an
// existing document.
func (b *BaseField) BeforeUpdate() {
	b.DefaultUpdatedAt()
}
