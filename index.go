package gomongo_client

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Indexer is implemented by models that define MongoDB indexes.
// If a model implements both [Model] and [Indexer], its indexes can
// be created automatically at startup via [AutoEnsureIndexes].
type Indexer interface {
	Indexes() []mongo.IndexModel
}

// IndexedModel is a model that implements both [Model] (for collection
// name) and [Indexer] (for index definitions). Use this as the
// constraint for [AutoEnsureIndexes].
type IndexedModel interface {
	Model
	Indexer
}

// EnsureIndexes creates the given indexes on the collection if they do
// not already exist. It returns the names of the created indexes.
func EnsureIndexes(
	ctx context.Context,
	collection *mongo.Collection,
	indexes []mongo.IndexModel,
	opts ...options.Lister[options.CreateIndexesOptions],
) ([]string, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}

	names, err := collection.Indexes().CreateMany(ctx, indexes, opts...)
	if err != nil {
		return nil, err
	}

	return names, nil
}

// AutoEnsureIndexes creates indexes for every provided model at
// application startup. Each model must implement both [Model] (to
// resolve the collection name) and [Indexer] (to provide index
// definitions).
//
// Call this once during application initialisation:
//
//	err := mongorm.AutoEnsureIndexes(ctx, client,
//	    UserModel{},
//	    ChatModel{},
//	    MessageModel{},
//	)
func AutoEnsureIndexes(ctx context.Context, client *Client, models ...IndexedModel) error {
	for _, model := range models {
		collection := client.Database().Collection(model.Collection())
		indexes := model.Indexes()

		if len(indexes) == 0 {
			continue
		}

		if _, err := EnsureIndexes(ctx, collection, indexes); err != nil {
			return fmt.Errorf("mongorm: failed to ensure indexes for %q: %w", model.Collection(), err)
		}
	}

	return nil
}

// ListIndexes returns the current indexes for the given collection.
func ListIndexes(
	ctx context.Context,
	collection *mongo.Collection,
	opts ...options.Lister[options.ListIndexesOptions],
) ([]bson.M, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}

	cursor, err := collection.Indexes().List(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("mongorm: failed to list indexes: %w", err)
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var indexes []bson.M
	if err = cursor.All(ctx, &indexes); err != nil {
		return nil, err
	}

	return indexes, nil
}
