package gomongo_client

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// CountDocuments returns the number of documents matching the filter.
// A nil filter counts all documents in the collection.
func CountDocuments(
	ctx context.Context,
	collection *mongo.Collection,
	filters any,
	opts ...options.Lister[options.CountOptions],
) (int64, error) {
	if collection == nil {
		return 0, ErrNilCollection
	}
	if filters == nil {
		filters = bson.D{}
	}

	return collection.CountDocuments(ctx, filters, opts...)
}

// EstimatedDocumentCount returns an estimated count of documents in the
// collection using collection metadata rather than scanning.
func EstimatedDocumentCount(
	ctx context.Context,
	collection *mongo.Collection,
	opts ...options.Lister[options.EstimatedDocumentCountOptions],
) (int64, error) {
	if collection == nil {
		return 0, ErrNilCollection
	}

	return collection.EstimatedDocumentCount(ctx, opts...)
}

// DistinctRaw returns the raw distinct result for the given field name.
func DistinctRaw(
	ctx context.Context,
	collection *mongo.Collection,
	filters any,
	fieldName string,
	opts ...options.Lister[options.DistinctOptions],
) (*mongo.DistinctResult, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}
	if filters == nil {
		filters = bson.D{}
	}

	result := collection.Distinct(ctx, fieldName, filters, opts...)
	if err := result.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// Distinct decodes the distinct values for the given field into the
// provided documents pointer (typically *[]string or *[]any).
func Distinct(
	ctx context.Context,
	collection *mongo.Collection,
	documents any,
	filters any,
	fieldName string,
	opts ...options.Lister[options.DistinctOptions],
) error {
	result, err := DistinctRaw(ctx, collection, filters, fieldName, opts...)
	if err != nil {
		return err
	}

	return result.Decode(documents)
}

// FindOneRaw returns the raw [mongo.SingleResult] for the first
// matching document.
func FindOneRaw(
	ctx context.Context,
	collection *mongo.Collection,
	filters any,
	opts ...options.Lister[options.FindOneOptions],
) (*mongo.SingleResult, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}
	if filters == nil {
		filters = bson.D{}
	}

	result := collection.FindOne(ctx, filters, opts...)
	if err := result.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// FindOne finds a single matching document and decodes it into the
// provided document pointer.
func FindOne(
	ctx context.Context,
	collection *mongo.Collection,
	document any,
	filters any,
	opts ...options.Lister[options.FindOneOptions],
) error {
	if document == nil {
		return mongo.ErrNilDocument
	}

	result, err := FindOneRaw(ctx, collection, filters, opts...)
	if err != nil {
		return err
	}

	return result.Decode(document)
}

// FindRaw returns a [mongo.Cursor] over documents matching the filter.
// The caller is responsible for closing the cursor.
func FindRaw(
	ctx context.Context,
	collection *mongo.Collection,
	filters any,
	opts ...options.Lister[options.FindOptions],
) (*mongo.Cursor, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}
	if filters == nil {
		filters = bson.D{}
	}

	cursor, err := collection.Find(ctx, filters, opts...)
	if err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return cursor, nil
}

// Find decodes all matching documents into the provided slice pointer.
func Find(
	ctx context.Context,
	collection *mongo.Collection,
	documents any,
	filters any,
	opts ...options.Lister[options.FindOptions],
) error {
	if documents == nil {
		return mongo.ErrNilDocument
	}

	cursor, err := FindRaw(ctx, collection, filters, opts...)
	if err != nil {
		return err
	}
	defer func() { _ = cursor.Close(ctx) }()

	return cursor.All(ctx, documents)
}

// InsertOneRaw inserts a single document and returns the raw result.
func InsertOneRaw(
	ctx context.Context,
	collection *mongo.Collection,
	insert any,
	opts ...options.Lister[options.InsertOneOptions],
) (*mongo.InsertOneResult, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}
	if insert == nil {
		return nil, mongo.ErrNilDocument
	}

	return collection.InsertOne(ctx, insert, opts...)
}

// InsertOne inserts a single document and returns both the raw result
// and the generated [bson.ObjectID].
func InsertOne(
	ctx context.Context,
	collection *mongo.Collection,
	insert any,
	opts ...options.Lister[options.InsertOneOptions],
) (*mongo.InsertOneResult, bson.ObjectID, error) {
	result, err := InsertOneRaw(ctx, collection, insert, opts...)
	if err != nil {
		return nil, bson.NilObjectID, err
	}

	documentID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil, bson.NilObjectID, ErrInvalidObjectID
	}

	return result, documentID, nil
}

// InsertMany inserts multiple documents into the collection.
func InsertMany(
	ctx context.Context,
	collection *mongo.Collection,
	insert any,
	opts ...options.Lister[options.InsertManyOptions],
) (*mongo.InsertManyResult, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}
	if insert == nil {
		return nil, mongo.ErrNilDocument
	}

	return collection.InsertMany(ctx, insert, opts...)
}

// UpdateByID updates a single document identified by its _id.
func UpdateByID(
	ctx context.Context,
	collection *mongo.Collection,
	id any,
	update any,
	opts ...options.Lister[options.UpdateOneOptions],
) (*mongo.UpdateResult, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}
	if update == nil {
		return nil, mongo.ErrNilDocument
	}

	return collection.UpdateByID(ctx, id, update, opts...)
}

// UpdateOne updates the first document matching the filter.
func UpdateOne(
	ctx context.Context,
	collection *mongo.Collection,
	filters any,
	update any,
	opts ...options.Lister[options.UpdateOneOptions],
) (*mongo.UpdateResult, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}
	if update == nil {
		return nil, mongo.ErrNilDocument
	}
	if filters == nil {
		filters = bson.D{}
	}

	return collection.UpdateOne(ctx, filters, update, opts...)
}

// UpdateMany updates all documents matching the filter.
func UpdateMany(
	ctx context.Context,
	collection *mongo.Collection,
	filters any,
	update any,
	opts ...options.Lister[options.UpdateManyOptions],
) (*mongo.UpdateResult, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}
	if update == nil {
		return nil, mongo.ErrNilDocument
	}
	if filters == nil {
		filters = bson.D{}
	}

	return collection.UpdateMany(ctx, filters, update, opts...)
}

// DeleteOne deletes the first document matching the filter.
func DeleteOne(
	ctx context.Context,
	collection *mongo.Collection,
	filters any,
	opts ...options.Lister[options.DeleteOneOptions],
) (*mongo.DeleteResult, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}
	if filters == nil {
		filters = bson.D{}
	}

	return collection.DeleteOne(ctx, filters, opts...)
}

// DeleteMany deletes all documents matching the filter.
func DeleteMany(
	ctx context.Context,
	collection *mongo.Collection,
	filters any,
	opts ...options.Lister[options.DeleteManyOptions],
) (*mongo.DeleteResult, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}
	if filters == nil {
		filters = bson.D{}
	}

	return collection.DeleteMany(ctx, filters, opts...)
}

// FindOneAndUpdateRaw atomically finds and updates a document,
// returning the raw [mongo.SingleResult].
func FindOneAndUpdateRaw(
	ctx context.Context,
	collection *mongo.Collection,
	filters any,
	update any,
	opts ...options.Lister[options.FindOneAndUpdateOptions],
) (*mongo.SingleResult, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}
	if update == nil {
		return nil, mongo.ErrNilDocument
	}
	if filters == nil {
		filters = bson.D{}
	}

	result := collection.FindOneAndUpdate(ctx, filters, update, opts...)
	if err := result.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// FindOneAndUpdate atomically finds, updates, and decodes a document
// into the provided pointer.
func FindOneAndUpdate(
	ctx context.Context,
	collection *mongo.Collection,
	document any,
	filters any,
	update any,
	opts ...options.Lister[options.FindOneAndUpdateOptions],
) error {
	if document == nil {
		return mongo.ErrNilDocument
	}

	result, err := FindOneAndUpdateRaw(ctx, collection, filters, update, opts...)
	if err != nil {
		return err
	}

	return result.Decode(document)
}

// FindOneAndDeleteRaw atomically finds and deletes a document,
// returning the raw [mongo.SingleResult].
func FindOneAndDeleteRaw(
	ctx context.Context,
	collection *mongo.Collection,
	filters any,
	opts ...options.Lister[options.FindOneAndDeleteOptions],
) (*mongo.SingleResult, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}
	if filters == nil {
		filters = bson.D{}
	}

	result := collection.FindOneAndDelete(ctx, filters, opts...)
	if err := result.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// FindOneAndDelete atomically finds, deletes, and decodes a document
// into the provided pointer.
func FindOneAndDelete(
	ctx context.Context,
	collection *mongo.Collection,
	document any,
	filters any,
	opts ...options.Lister[options.FindOneAndDeleteOptions],
) error {
	if document == nil {
		return mongo.ErrNilDocument
	}

	result, err := FindOneAndDeleteRaw(ctx, collection, filters, opts...)
	if err != nil {
		return err
	}

	return result.Decode(document)
}

// FindOneAndReplaceRaw atomically finds and replaces a document,
// returning the raw [mongo.SingleResult].
func FindOneAndReplaceRaw(
	ctx context.Context,
	collection *mongo.Collection,
	filters any,
	replacement any,
	opts ...options.Lister[options.FindOneAndReplaceOptions],
) (*mongo.SingleResult, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}
	if replacement == nil {
		return nil, mongo.ErrNilDocument
	}
	if filters == nil {
		filters = bson.D{}
	}

	result := collection.FindOneAndReplace(ctx, filters, replacement, opts...)
	if err := result.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// FindOneAndReplace atomically finds, replaces, and decodes a document
// into the provided pointer.
func FindOneAndReplace(
	ctx context.Context,
	collection *mongo.Collection,
	document any,
	filters any,
	replacement any,
	opts ...options.Lister[options.FindOneAndReplaceOptions],
) error {
	if document == nil {
		return mongo.ErrNilDocument
	}

	result, err := FindOneAndReplaceRaw(ctx, collection, filters, replacement, opts...)
	if err != nil {
		return err
	}

	return result.Decode(document)
}

// AggregateRaw runs an aggregation pipeline and returns a [mongo.Cursor].
// The caller is responsible for closing the cursor.
func AggregateRaw(
	ctx context.Context,
	collection *mongo.Collection,
	pipeline any,
	opts ...options.Lister[options.AggregateOptions],
) (*mongo.Cursor, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}
	if pipeline == nil {
		return nil, ErrNilPipeline
	}

	cursor, err := collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return cursor, nil
}

// Aggregate runs an aggregation pipeline and decodes all results into
// the provided slice pointer.
func Aggregate(
	ctx context.Context,
	collection *mongo.Collection,
	documents any,
	pipeline any,
	opts ...options.Lister[options.AggregateOptions],
) error {
	if documents == nil {
		return mongo.ErrNilDocument
	}

	cursor, err := AggregateRaw(ctx, collection, pipeline, opts...)
	if err != nil {
		return err
	}
	defer func() { _ = cursor.Close(ctx) }()

	return cursor.All(ctx, documents)
}
