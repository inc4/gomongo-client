// Package mongorm provides a lightweight ORM-like layer on top of the
// official MongoDB Go driver (v2). It offers:
//
//   - A thin [Client] wrapper for connection management.
//   - A generic [Repository] with type-safe CRUD, aggregation, and
//     pagination methods.
//   - [BaseField] for common document fields (_id, createdAt, updatedAt).
//   - Offset-based and cursor-based pagination helpers.
//   - Index management with auto-ensure on startup via [AutoEnsureIndexes].
//
// # Quick Start
//
//	client, err := mongorm.Connect(ctx, "mydb",
//	    options.Client().ApplyURI("mongodb://localhost:27017"))
//	defer client.Disconnect(ctx)
//
//	repo := mongorm.NewRepository[UserModel](client, UserModel{})
//	user, err := repo.FindByID(ctx, id)
//
// # Auto-Index
//
// Implement the [Indexer] interface on your model and call
// [AutoEnsureIndexes] at application startup:
//
//	err := mongorm.AutoEnsureIndexes(ctx, client, UserModel{}, ChatModel{})
//
// # Cursor Pagination
//
// The caller is responsible for decoding the cursor and building
// appropriate query filters (e.g. _id > lastSeenID). [Repository.FindWithCursor]
// handles the limit+1 fetch trick, direction reversal, and result
// envelope construction.
package gomongo_client
