package gomongo_client

import (
	"context"
	"slices"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Model is implemented by types that map to a MongoDB collection. The
// returned string is the collection name used for all CRUD operations.
//
// Example:
//
//	type UserModel struct {
//	    mongorm.BaseField `bson:",inline"`
//	    Name string       `bson:"name"`
//	}
//
//	func (UserModel) Collection() string { return "users" }
type Model interface {
	Collection() string
}

// Repository provides type-safe CRUD, aggregation, and pagination
// methods for a single MongoDB collection. Create one with
// [NewRepository].
type Repository[T Model] struct {
	client *Client
	model  T
}

// NewRepository creates a [Repository] for the given model. The model
// determines the collection name.
func NewRepository[T Model](client *Client, model T) *Repository[T] {
	return &Repository[T]{
		client: client,
		model:  model,
	}
}

// CollectionName returns the MongoDB collection name for this
// repository's model.
func (r *Repository[T]) CollectionName() string {
	return r.model.Collection()
}

// Collection returns a [mongo.Collection] handle for this repository.
func (r *Repository[T]) Collection(opts ...options.Lister[options.CollectionOptions]) *mongo.Collection {
	return r.client.Database().Collection(r.model.Collection(), opts...)
}

// EnsureIndexes creates indexes defined by the model's [Indexer]
// implementation. Returns [ErrNotIndexer] if the model does not
// implement [Indexer].
func (r *Repository[T]) EnsureIndexes(
	ctx context.Context,
	opts ...options.Lister[options.CreateIndexesOptions],
) ([]string, error) {
	var model T
	indexer, ok := any(model).(Indexer)
	if !ok {
		return nil, ErrNotIndexer
	}

	return EnsureIndexes(ctx, r.Collection(), indexer.Indexes(), opts...)
}

// ListIndexes returns the current indexes on this repository's
// collection.
func (r *Repository[T]) ListIndexes(
	ctx context.Context,
	opts ...options.Lister[options.ListIndexesOptions],
) ([]bson.M, error) {
	return ListIndexes(ctx, r.Collection(), opts...)
}

// CountDocuments returns the number of documents matching the filter.
func (r *Repository[T]) CountDocuments(
	ctx context.Context,
	filters any,
	opts ...options.Lister[options.CountOptions],
) (int64, error) {
	return CountDocuments(ctx, r.Collection(), filters, opts...)
}

// EstimatedDocumentCount returns an estimated document count using
// collection metadata.
func (r *Repository[T]) EstimatedDocumentCount(
	ctx context.Context,
	opts ...options.Lister[options.EstimatedDocumentCountOptions],
) (int64, error) {
	return EstimatedDocumentCount(ctx, r.Collection(), opts...)
}

// Distinct decodes distinct values for the given field.
func (r *Repository[T]) Distinct(
	ctx context.Context,
	documents any,
	filters any,
	fieldName string,
	opts ...options.Lister[options.DistinctOptions],
) error {
	return Distinct(ctx, r.Collection(), documents, filters, fieldName, opts...)
}

// FindByID finds a single document by its _id field.
func (r *Repository[T]) FindByID(
	ctx context.Context,
	id any,
	opts ...options.Lister[options.FindOneOptions],
) (T, error) {
	var document T
	err := FindOne(ctx, r.Collection(), &document, bson.M{"_id": id}, opts...)
	return document, err
}

// FindOne finds the first document matching the filter.
func (r *Repository[T]) FindOne(
	ctx context.Context,
	filters any,
	opts ...options.Lister[options.FindOneOptions],
) (T, error) {
	var document T
	err := FindOne(ctx, r.Collection(), &document, filters, opts...)
	return document, err
}

// Find returns all documents matching the filter.
func (r *Repository[T]) Find(
	ctx context.Context,
	filters any,
	opts ...options.Lister[options.FindOptions],
) ([]T, error) {
	documents := make([]T, 0)
	err := Find(ctx, r.Collection(), &documents, filters, opts...)
	return documents, err
}

// InsertOne inserts a single document and returns the result with the
// generated ObjectID.
func (r *Repository[T]) InsertOne(
	ctx context.Context,
	insert any,
	opts ...options.Lister[options.InsertOneOptions],
) (*mongo.InsertOneResult, bson.ObjectID, error) {
	return InsertOne(ctx, r.Collection(), insert, opts...)
}

// InsertMany inserts multiple documents.
func (r *Repository[T]) InsertMany(
	ctx context.Context,
	insert any,
	opts ...options.Lister[options.InsertManyOptions],
) (*mongo.InsertManyResult, error) {
	return InsertMany(ctx, r.Collection(), insert, opts...)
}

// UpdateByID updates a single document by its _id field.
func (r *Repository[T]) UpdateByID(
	ctx context.Context,
	id any,
	update any,
	opts ...options.Lister[options.UpdateOneOptions],
) (*mongo.UpdateResult, error) {
	return UpdateByID(ctx, r.Collection(), id, update, opts...)
}

// UpdateOne updates the first document matching the filter.
func (r *Repository[T]) UpdateOne(
	ctx context.Context,
	filters any,
	update any,
	opts ...options.Lister[options.UpdateOneOptions],
) (*mongo.UpdateResult, error) {
	return UpdateOne(ctx, r.Collection(), filters, update, opts...)
}

// UpdateMany updates all documents matching the filter.
func (r *Repository[T]) UpdateMany(
	ctx context.Context,
	filters any,
	update any,
	opts ...options.Lister[options.UpdateManyOptions],
) (*mongo.UpdateResult, error) {
	return UpdateMany(ctx, r.Collection(), filters, update, opts...)
}

// DeleteOne deletes the first document matching the filter.
func (r *Repository[T]) DeleteOne(
	ctx context.Context,
	filters any,
	opts ...options.Lister[options.DeleteOneOptions],
) (*mongo.DeleteResult, error) {
	return DeleteOne(ctx, r.Collection(), filters, opts...)
}

// DeleteMany deletes all documents matching the filter.
func (r *Repository[T]) DeleteMany(
	ctx context.Context,
	filters any,
	opts ...options.Lister[options.DeleteManyOptions],
) (*mongo.DeleteResult, error) {
	return DeleteMany(ctx, r.Collection(), filters, opts...)
}

// FindAndUpdateByID atomically finds, updates, and returns a document
// by its _id.
func (r *Repository[T]) FindAndUpdateByID(
	ctx context.Context,
	id any,
	update any,
	opts ...options.Lister[options.FindOneAndUpdateOptions],
) (T, error) {
	var document T
	err := FindOneAndUpdate(ctx, r.Collection(), &document, bson.M{"_id": id}, update, opts...)
	return document, err
}

// FindOneAndUpdate atomically finds, updates, and returns the first
// matching document.
func (r *Repository[T]) FindOneAndUpdate(
	ctx context.Context,
	filters any,
	update any,
	opts ...options.Lister[options.FindOneAndUpdateOptions],
) (T, error) {
	var document T
	err := FindOneAndUpdate(ctx, r.Collection(), &document, filters, update, opts...)
	return document, err
}

// FindAndDeleteByID atomically finds, deletes, and returns a document
// by its _id.
func (r *Repository[T]) FindAndDeleteByID(
	ctx context.Context,
	id any,
	opts ...options.Lister[options.FindOneAndDeleteOptions],
) (T, error) {
	var document T
	err := FindOneAndDelete(ctx, r.Collection(), &document, bson.M{"_id": id}, opts...)
	return document, err
}

// FindOneAndDelete atomically finds, deletes, and returns the first
// matching document.
func (r *Repository[T]) FindOneAndDelete(
	ctx context.Context,
	filters any,
	opts ...options.Lister[options.FindOneAndDeleteOptions],
) (T, error) {
	var document T
	err := FindOneAndDelete(ctx, r.Collection(), &document, filters, opts...)
	return document, err
}

// FindAndReplaceByID atomically finds, replaces, and returns a document
// by its _id.
func (r *Repository[T]) FindAndReplaceByID(
	ctx context.Context,
	id any,
	replacement any,
	opts ...options.Lister[options.FindOneAndReplaceOptions],
) (T, error) {
	var document T
	err := FindOneAndReplace(ctx, r.Collection(), &document, bson.M{"_id": id}, replacement, opts...)
	return document, err
}

// FindOneAndReplace atomically finds, replaces, and returns the first
// matching document.
func (r *Repository[T]) FindOneAndReplace(
	ctx context.Context,
	filters any,
	replacement any,
	opts ...options.Lister[options.FindOneAndReplaceOptions],
) (T, error) {
	var document T
	err := FindOneAndReplace(ctx, r.Collection(), &document, filters, replacement, opts...)
	return document, err
}

// Aggregate runs an aggregation pipeline and decodes the results into
// the provided slice pointer.
func (r *Repository[T]) Aggregate(
	ctx context.Context,
	documents any,
	pipeline any,
	opts ...options.Lister[options.AggregateOptions],
) error {
	return Aggregate(ctx, r.Collection(), documents, pipeline, opts...)
}

// FindWithOffset returns documents using skip/limit pagination.
func (r *Repository[T]) FindWithOffset(
	ctx context.Context,
	filters any,
	pagination *OffsetPagination,
	sort any,
) ([]T, error) {
	if pagination == nil {
		pagination = DefaultOffsetPagination()
	}
	if pagination.Limit <= 0 {
		pagination.Limit = DefaultPaginationLimit
	}
	if pagination.Skip < 0 {
		pagination.Skip = 0
	}

	opts := options.Find().
		SetSkip(pagination.Skip).
		SetLimit(pagination.Limit)

	if sort != nil {
		opts.SetSort(sort)
	}

	return r.Find(ctx, filters, opts)
}

// FindWithOffsetAndTotal returns paginated documents together with the
// total matching count, using a $facet aggregation pipeline. This
// executes the count and the paginated fetch in a single round-trip.
func (r *Repository[T]) FindWithOffsetAndTotal(
	ctx context.Context,
	filters any,
	pagination *OffsetPagination,
	sort any,
) (*OffsetPaginationResult[T], error) {
	if pagination == nil {
		pagination = DefaultOffsetPagination()
	}
	if pagination.Limit <= 0 {
		pagination.Limit = DefaultPaginationLimit
	}
	if pagination.Skip < 0 {
		pagination.Skip = 0
	}

	// Build the items sub-pipeline (optional sort → skip → limit).
	itemsPipeline := bson.A{}
	if sort != nil {
		itemsPipeline = append(itemsPipeline, bson.M{"$sort": sort})
	}
	itemsPipeline = append(itemsPipeline,
		bson.M{"$skip": pagination.Skip},
		bson.M{"$limit": pagination.Limit},
	)

	// Build the top-level pipeline.
	pipeline := bson.A{}
	if filters != nil {
		pipeline = append(pipeline, bson.M{"$match": filters})
	}

	pipeline = append(pipeline,
		bson.M{"$facet": bson.D{
			{Key: "items", Value: itemsPipeline},
			{Key: "pagination", Value: bson.A{
				bson.M{"$count": "totalCount"},
			}},
		}},
		bson.M{"$unwind": bson.M{
			"path":                       "$pagination",
			"preserveNullAndEmptyArrays": true,
		}},
		bson.M{"$addFields": bson.M{
			"totalCount": bson.M{
				"$ifNull": bson.A{"$pagination.totalCount", 0},
			},
		}},
	)

	// $facet always returns exactly one document, but Aggregate decodes
	// via cursor.All which expects a slice.
	var results []OffsetPaginationResult[T]
	if err := r.Aggregate(ctx, &results, pipeline); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &OffsetPaginationResult[T]{
			Items:      make([]T, 0),
			TotalCount: 0,
		}, nil
	}

	result := &results[0]
	if result.Items == nil {
		result.Items = make([]T, 0)
	}

	return result, nil
}

// FindWithCursor returns paginated documents using cursor-based
// pagination with the "limit + 1" technique to detect additional pages.
//
// IMPORTANT: The caller is responsible for decoding the cursor string
// (via [DecodeCursor]) and incorporating it into the filters parameter
// (e.g. bson.M{"_id": bson.M{"$gt": cursorID}}). This function does
// NOT apply cursor-based filtering automatically — it only handles
// the limit+1 fetch, direction reversal, and result envelope.
//
// The cursorEncoder function converts an item into an opaque cursor
// string (typically via [EncodeCursor]).
//
// Example:
//
//	filters := bson.M{}
//	if req.Cursor != "" {
//	    cur, _ := mongorm.DecodeCursor[MyCursor](req.Cursor)
//	    filters["_id"] = bson.M{"$gt": cur.ID}
//	}
//	result, err := repo.FindWithCursor(ctx, filters, pagination, sort,
//	    func(item MyModel) (string, error) {
//	        return mongorm.EncodeCursor(MyCursor{ID: item.ID})
//	    },
//	)
func (r *Repository[T]) FindWithCursor(
	ctx context.Context,
	filters any,
	pagination *CursorPagination,
	sort any,
	cursorEncoder func(T) (string, error),
) (*CursorPaginationResult[T], error) {
	if pagination == nil {
		pagination = DefaultCursorPagination()
	}
	if pagination.Limit <= 0 {
		pagination.Limit = DefaultPaginationLimit
	}

	// Fetch limit+1 to determine if there are more results.
	opts := options.Find().SetLimit(pagination.Limit + 1)
	if sort != nil {
		opts.SetSort(sort)
	}

	items, err := r.Find(ctx, filters, opts)
	if err != nil {
		return nil, err
	}

	hasMore := len(items) > int(pagination.Limit)
	if hasMore {
		items = items[:pagination.Limit]
	}

	// When paging backward, reverse to maintain consistent order.
	if pagination.Direction == CursorPrev {
		slices.Reverse(items)
	}

	result := &CursorPaginationResult[T]{
		Items: items,
	}

	// Ensure Items is never nil for JSON serialisation.
	if result.Items == nil {
		result.Items = make([]T, 0)
	}

	if len(items) == 0 {
		return result, nil
	}

	firstItem := items[0]
	lastItem := items[len(items)-1]

	switch pagination.Direction {
	case CursorNext:
		result.HasNext = hasMore
		result.HasPrev = pagination.Cursor != ""

		if result.HasNext {
			if result.NextCursor, err = cursorEncoder(lastItem); err != nil {
				return nil, err
			}
		}
		if result.HasPrev {
			if result.PrevCursor, err = cursorEncoder(firstItem); err != nil {
				return nil, err
			}
		}

	case CursorPrev:
		result.HasPrev = hasMore
		result.HasNext = pagination.Cursor != ""

		if result.HasPrev {
			if result.PrevCursor, err = cursorEncoder(firstItem); err != nil {
				return nil, err
			}
		}
		if result.HasNext {
			if result.NextCursor, err = cursorEncoder(lastItem); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}
