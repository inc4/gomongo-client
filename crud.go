package mongoclient

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// InsertOne inserts a new document
func (r *Repository[T]) InsertOne(ctx context.Context, document T, opts ...options.Lister[options.InsertOneOptions]) (T, error) {
	var zero T

	// call BeforeInsert hook if it exists
	if hook, ok := any(document).(Document); ok {
		hook.BeforeInsert()
	}

	// Insert the document into the collection
	result, err := r.collection.InsertOne(ctx, document, opts...)
	if err != nil {
		return zero, fmt.Errorf("failed to insert document: %w", err)
	}

	// retrieve the inserted document
	var inserted T
	err = r.collection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&inserted)
	if err != nil {
		return zero, fmt.Errorf("failed to fetch inserted document: %w", err)
	}

	return inserted, nil
}

// InsertMany inserts multiple documents
func (r *Repository[T]) InsertMany(ctx context.Context, documents []T, opts ...options.Lister[options.InsertManyOptions]) ([]any, error) {
	if len(documents) == 0 {
		return nil, fmt.Errorf("no documents to insert")
	}

	interfaces := make([]any, len(documents))
	for i, doc := range documents {
		if hook, ok := any(doc).(Document); ok {
			hook.BeforeInsert()
		}
		interfaces[i] = doc
	}

	result, err := r.collection.InsertMany(ctx, interfaces, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert documents: %w", err)
	}

	return result.InsertedIDs, nil
}

// FindOne retrieves a single document
func (r *Repository[T]) FindOne(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) (T, error) {
	var result T
	if filter == nil {
		filter = bson.M{}
	}
	return result, r.collection.FindOne(ctx, filter, opts...).Decode(&result)
}

// Find retrieves multiple documents
func (r *Repository[T]) Find(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) ([]T, error) {
	if filter == nil {
		filter = bson.M{}
	}
	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute find: %w", err)
	}
	defer cursor.Close(ctx)

	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode results: %w", err)
	}

	return results, nil
}

func (r *Repository[T]) FindPaginated(ctx context.Context, filter any, page, pageSize int64) ([]T, error) {
	if page < 1 {
		return nil, fmt.Errorf("invalid page: must be >= 1")
	}
	if pageSize < 1 {
		return nil, fmt.Errorf("invalid pageSize: must be >= 1")
	}

	opts := options.Find()
	opts.SetSkip((page - 1) * pageSize)
	opts.SetLimit(pageSize)
	return r.FindDecoded(ctx, filter, opts)

}

func (r *Repository[T]) FindPaginatedWithTotal(ctx context.Context, filter any, page, pageSize int64) ([]T, int64, error) {
	// Validation of page/pageSize happens in FindPaginated
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	results, err := r.FindPaginated(ctx, filter, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *Repository[T]) FindDecoded(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) ([]T, error) {
	if filter == nil {
		filter = bson.M{}
	}
	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute find: %w", err)
	}
	defer cursor.Close(ctx)

	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode results: %w", err)
	}
	return results, nil
}

func (r *Repository[T]) FindDecodedWithTotal(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) ([]T, int64, error) {
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	results, err := r.FindDecoded(ctx, filter, opts...)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// FindByID finds a document by its ID
func (r *Repository[T]) FindByID(ctx context.Context, id any, opts ...options.Lister[options.FindOneOptions]) (T, error) {
	return r.FindOne(ctx, bson.M{"_id": id}, opts...)
}

// FindOneAndUpdate finds a document and updates it
func (r *Repository[T]) FindOneAndUpdate(ctx context.Context, filter, update any, opts ...options.Lister[options.FindOneAndUpdateOptions]) (T, error) {
	var result T

	switch {
	case isMongoOperator(update):
		// Already formatted update (e.g., $set, $inc) -> do nothing
	case isStructOrPtrToStruct(update):
		if hook, ok := update.(Document); ok {
			hook.BeforeUpdate()
		}
		update = bson.M{"$set": update}
	default:
		return result, fmt.Errorf("unsupported update type: %T", update)
	}

	opts = append(opts, options.FindOneAndUpdate().SetReturnDocument(options.After))
	return result, r.collection.FindOneAndUpdate(ctx, filter, update, opts...).Decode(&result)
}

func (r *Repository[T]) FindOneAndUpdateByID(ctx context.Context, id, update any, opts ...options.Lister[options.FindOneAndUpdateOptions]) (T, error) {
	return r.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts...)
}

// FindOneAndDelete finds a document and deletes it
func (r *Repository[T]) FindOneAndDelete(ctx context.Context, filter any, opts ...options.Lister[options.FindOneAndDeleteOptions]) (T, error) {
	var result T
	return result, r.collection.FindOneAndDelete(ctx, filter, opts...).Decode(&result)
}

// UpdateOne updates a single document
func (r *Repository[T]) UpdateOne(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
	switch {
	case isMongoOperator(update):
		// Already formatted update (e.g., $set, $inc) -> do nothing
	case isStructOrPtrToStruct(update):
		if hook, ok := update.(Document); ok {
			hook.BeforeUpdate()
		}
		update = bson.M{"$set": update}
	default:
		return nil, fmt.Errorf("unsupported update type: %T", update)
	}
	return r.collection.UpdateOne(ctx, filter, update, opts...)
}

// UpdateByID finds a document by its ID
func (r *Repository[T]) UpdateByID(ctx context.Context, id any, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
	return r.UpdateOne(ctx, bson.M{"_id": id}, update, opts...)
}

// UpdateMany updates multiple documents
func (r *Repository[T]) UpdateMany(ctx context.Context, filter any, update any, opts ...options.Lister[options.UpdateManyOptions]) (*mongo.UpdateResult, error) {
	switch {
	case isMongoOperator(update):
		// already formatted (e.g., $inc, $set) → pass through
	case isStructOrPtrToStruct(update):
		if hook, ok := update.(Document); ok {
			hook.BeforeUpdate()
		}
		update = bson.M{"$set": update}
	default:
		return nil, fmt.Errorf("unsupported update type: %T", update)
	}

	result, err := r.collection.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to update documents: %w", err)
	}
	return result, nil
}

// DeleteOne removes a single document
func (r *Repository[T]) DeleteOne(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) error {
	result, err := r.collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// DeleteByID removes a document by its ID
func (r *Repository[T]) DeleteByID(ctx context.Context, id any, opts ...options.Lister[options.DeleteOneOptions]) error {
	return r.DeleteOne(ctx, bson.M{"_id": id}, opts...)
}

// DeleteMany removes multiple documents
func (r *Repository[T]) DeleteMany(ctx context.Context, filter any, opts ...options.Lister[options.DeleteManyOptions]) (int64, error) {
	result, err := r.collection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return 0, fmt.Errorf("failed to delete documents: %w", err)
	}

	return result.DeletedCount, nil
}
