package gomongo_client

import (
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Sentinel errors for the mongorm package.
var (
	// ErrNilCollection is returned when a nil *mongo.Collection is
	// passed to an operation function.
	ErrNilCollection = errors.New("nil collection")

	// ErrNilPipeline is returned when a nil aggregation pipeline is
	// passed to an aggregate function.
	ErrNilPipeline = errors.New("nil pipeline")

	// ErrInvalidObjectID is returned when an InsertOne result cannot
	// be converted to a bson.ObjectID.
	ErrInvalidObjectID = errors.New("invalid object id")

	// ErrNotIndexer is returned when a model does not implement the
	// [Indexer] interface but index operations are requested.
	ErrNotIndexer = errors.New("model does not implement Indexer")

	// ErrNoResults is returned when an aggregation expected at least
	// one result document but received none.
	ErrNoResults = errors.New("no results")
)

// MapNotFoundErr translates [mongo.ErrNoDocuments] into the given domain error.
// All other errors are returned unchanged.
func MapNotFoundErr(err, notFound error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return notFound
	}
	return err
}
