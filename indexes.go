package mongoclient

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type IIndex interface {
	Indexes() []mongo.IndexModel
}

// EnsureIndexesAssertType creates indexes for the collection if they don't exist
func (r *Repository[T]) EnsureIndexesAssertType(ctx context.Context, opts ...options.Lister[options.CreateIndexesOptions]) error {
	var document T
	documentWithIndexes, ok := any(document).(IIndex)
	if !ok {
		return fmt.Errorf("repository type %T does not implement IIndex interface", document)
	}
	_, err := r.collection.Indexes().CreateMany(ctx, documentWithIndexes.Indexes(), opts...)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}
	return nil
}

// EnsureIndexes creates indexes for the collection if they don't exist
func (r *Repository[T]) EnsureIndexes(ctx context.Context, indexes []mongo.IndexModel, opts ...options.Lister[options.CreateIndexesOptions]) error {
	_, err := r.collection.Indexes().CreateMany(ctx, indexes, opts...)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}
	return nil
}

// GetIndexes returns all indexes for the collection
func (r *Repository[T]) GetIndexes(ctx context.Context, opts ...options.Lister[options.ListIndexesOptions]) ([]bson.M, error) {
	cursor, err := r.collection.Indexes().List(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to list indexes: %w", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode indexes: %w", err)
	}

	return results, nil
}
