package mongoclient

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// TODO: need to test.

// Transaction executes operations within a transaction
func (r *Repository[T]) Transaction(ctx context.Context, fn func(sessCtx context.Context) error, opts ...options.Lister[options.SessionOptions]) error {
	session, err := r.collection.Database().Client().StartSession(opts...)
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (any, error) {
		return nil, fn(sessCtx)
	})

	return err
}
