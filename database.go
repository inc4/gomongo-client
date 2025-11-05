package mongoclient

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// IRepository defines the interface for database operations
type IRepository[T any] interface {
	Collection() *mongo.Collection
	InsertOne(ctx context.Context, document T, opts ...options.Lister[options.InsertOneOptions]) (T, error)
	InsertMany(ctx context.Context, documents []T, opts ...options.Lister[options.InsertManyOptions]) ([]any, error)
	Find(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) ([]T, error)
	FindPaginated(ctx context.Context, filter any, page, pageSize int64) ([]T, error)
	FindPaginatedWithTotal(ctx context.Context, filter any, page, pageSize int64) ([]T, int64, error)
	FindOne(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) (T, error)
	FindByID(ctx context.Context, id any, opts ...options.Lister[options.FindOneOptions]) (T, error)
	FindOneAndUpdate(ctx context.Context, filter any, update any, opts ...options.Lister[options.FindOneAndUpdateOptions]) (T, error)
	FindOneAndUpdateByID(ctx context.Context, id, update any, opts ...options.Lister[options.FindOneAndUpdateOptions]) (T, error)
	FindOneAndDelete(ctx context.Context, filter any, opts ...options.Lister[options.FindOneAndDeleteOptions]) (T, error)
	UpdateOne(ctx context.Context, filter any, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
	UpdateByID(ctx context.Context, id any, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter any, update any, opts ...options.Lister[options.UpdateManyOptions]) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) error
	DeleteByID(ctx context.Context, id any, opts ...options.Lister[options.DeleteOneOptions]) error
	DeleteMany(ctx context.Context, filter any, opts ...options.Lister[options.DeleteManyOptions]) (int64, error)
	EstimatedCount(ctx context.Context, opts ...options.Lister[options.EstimatedDocumentCountOptions]) (int64, error)
	CountDocuments(ctx context.Context, filter any, opts ...options.Lister[options.CountOptions]) (int64, error)
	Aggregate(ctx context.Context, pipeline any, opts ...options.Lister[options.AggregateOptions]) ([]bson.M, error)
	AggregateTyped(ctx context.Context, pipeline any, opts ...options.Lister[options.AggregateOptions]) ([]T, error)
	Distinct(ctx context.Context, fieldName string, filter any, opts ...options.Lister[options.DistinctOptions]) ([]any, error)
	Transaction(ctx context.Context, fn func(sessCtx context.Context) error, opts ...options.Lister[options.SessionOptions]) error
	EnsureIndexes(ctx context.Context, indexes []mongo.IndexModel, opts ...options.Lister[options.CreateIndexesOptions]) error
	GetIndexes(ctx context.Context, opts ...options.Lister[options.ListIndexesOptions]) ([]bson.M, error)
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...options.Lister[options.BulkWriteOptions]) (*mongo.BulkWriteResult, error)
	Watch(ctx context.Context, pipeline any, opts ...options.Lister[options.ChangeStreamOptions]) (*mongo.ChangeStream, error)

	EnsureIndexesAssertType(ctx context.Context, opts ...options.Lister[options.CreateIndexesOptions]) error
}

type Document interface {
	SetID(id bson.ObjectID)
	BeforeInsert()
	BeforeUpdate()
}

// Repository implements IRepository interface for MongoDB
type Repository[T any] struct {
	collection *mongo.Collection
}

// NewRepository creates a new MongoDB repository
func NewRepository[T any](collectionName *mongo.Collection) *Repository[T] {
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		panic("repository must be a struct pointer to a struct")
	}
	return &Repository[T]{collection: collectionName}
}
