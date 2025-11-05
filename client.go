package mongoclient

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/connstring"
)

// Connect returns a new instance of the Mongo database client.
func Connect(ctx context.Context, uri string, databaseName string) (*mongo.Database, error) {
	opts := options.Client().ApplyURI(uri)

	// Configure the MongoDB API version.
	if strings.Contains(uri, connstring.SchemeMongoDBSRV) {
		opts.SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))
	}

	// Create and initialize the MongoDB client.
	client, err := mongo.Connect(opts)

	if err != nil {
		return nil, err
	}

	// Check the MongoDB client connection.
	err = client.Ping(ctx, nil)

	if err != nil {
		return nil, err
	}

	return client.Database(databaseName), nil
}
