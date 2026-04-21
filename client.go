package gomongo_client

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Client wraps a [mongo.Client] with a default database name.
type Client struct {
	client   *mongo.Client
	database string
}

// Connect creates a new [Client] by connecting to MongoDB with the
// given database name and client options. It pings the server to verify
// connectivity before returning.
func Connect(ctx context.Context, database string, opts ...*options.ClientOptions) (*Client, error) {
	client, err := mongo.Connect(opts...)

	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Client{
		client:   client,
		database: database,
	}, nil
}

// ConnectWithURI creates a new [Client] by parsing the connection URI
// and connecting to MongoDB. For Atlas Serverless or any deployment
// that requires the Stable API, pass [WithServerAPI].
//
// Example:
//
//	client, err := mongorm.ConnectWithURI(ctx, "mydb",
//	    "mongodb+srv://user:pass@cluster.mongodb.net",
//	    mongorm.WithServerAPI(),
//	)
func ConnectWithURI(ctx context.Context, database, uri string, opts ...ClientOption) (*Client, error) {
	cfg := &clientConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	clientOptions := options.Client().ApplyURI(uri)

	if cfg.serverAPI {
		clientOptions.SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))
	}

	return Connect(ctx, database, clientOptions)
}

// ClientOption configures optional parameters for [ConnectWithURI].
type ClientOption func(*clientConfig)

type clientConfig struct {
	serverAPI bool
}

// WithServerAPI enables the MongoDB Stable API (ServerAPIVersion1).
// This is required for Atlas Serverless and recommended for Atlas
// Dedicated deployments.
func WithServerAPI() ClientOption {
	return func(c *clientConfig) {
		c.serverAPI = true
	}
}

// Disconnect closes the MongoDB connection.
func (c *Client) Disconnect(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

// Client returns the underlying [mongo.Client].
func (c *Client) Client() *mongo.Client {
	return c.client
}

// DatabaseName returns the configured database name.
func (c *Client) DatabaseName() string {
	return c.database
}

// Database returns a handle for the configured database.
func (c *Client) Database(opts ...options.Lister[options.DatabaseOptions]) *mongo.Database {
	return c.client.Database(c.database, opts...)
}
