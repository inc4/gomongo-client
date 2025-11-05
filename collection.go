package mongoclient

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *Repository[T]) Collection() *mongo.Collection {
	return r.collection
}

// EstimatedCount returns the estimated number of documents in the collection
func (r *Repository[T]) EstimatedCount(ctx context.Context, opts ...options.Lister[options.EstimatedDocumentCountOptions]) (int64, error) {
	count, err := r.collection.EstimatedDocumentCount(ctx, opts...)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	return count, nil
}

// CountDocuments returns the exact number of documents matching the filter
func (r *Repository[T]) CountDocuments(ctx context.Context, filter any, opts ...options.Lister[options.CountOptions]) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter, opts...)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	return count, nil
}

// Aggregate performs an aggregation pipeline and returns raw bson.M results
func (r *Repository[T]) Aggregate(ctx context.Context, pipeline any, opts ...options.Lister[options.AggregateOptions]) ([]bson.M, error) {
	cursor, err := r.collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode aggregate results: %w", err)
	}
	return results, nil
}

// AggregateTyped performs an aggregation pipeline and returns typed results
func (r *Repository[T]) AggregateTyped(ctx context.Context, pipeline any, opts ...options.Lister[options.AggregateOptions]) ([]T, error) {
	cursor, err := r.collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode aggregate results: %w", err)
	}
	return results, nil
}

// Distinct finds the distinct values for a specified field
func (r *Repository[T]) Distinct(ctx context.Context, fieldName string, filter any, opts ...options.Lister[options.DistinctOptions]) ([]any, error) {
	var arr []any
	err := r.collection.Distinct(ctx, fieldName, filter, opts...).Decode(&arr)
	if err != nil {
		return nil, fmt.Errorf("failed to find distinct values: %w", err)
	}
	return arr, nil
}

// BulkWrite performs multiple write operations in bulk
func (r *Repository[T]) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...options.Lister[options.BulkWriteOptions]) (*mongo.BulkWriteResult, error) {
	result, err := r.collection.BulkWrite(ctx, models, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to perform bulk write: %w", err)
	}
	return result, nil
}

// Watch creates a change stream for the collection
func (r *Repository[T]) Watch(ctx context.Context, pipeline any, opts ...options.Lister[options.ChangeStreamOptions]) (*mongo.ChangeStream, error) {
	stream, err := r.collection.Watch(ctx, pipeline, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create change stream: %w", err)
	}
	return stream, nil
}
