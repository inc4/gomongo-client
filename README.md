# gomongo-client

Generic MongoDB repository library for Go, built on top of the official [mongo-driver/v2](https://pkg.go.dev/go.mongodb.org/mongo-driver/v2).

## Installation

```bash
go get github.com/inc4/gomongo-client
```

## Quick Start

### Define a Model

Embed `BaseField` to get auto-generated IDs and timestamps (`createdAt`, `updatedAt`).

```go
package main

import (
    "context"
    "log"

    "github.com/inc4/gomongo-client"
    "go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
    mongoclient.BaseField `bson:",inline"`
    Name                  string `bson:"name" json:"name"`
    Email                 string `bson:"email" json:"email"`
    Age                   int    `bson:"age" json:"age"`
}

// Implement IIndex to self-declare indexes (optional).
func (u *User) Indexes() []mongo.IndexModel {
    return []mongo.IndexModel{
        {Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
        {Keys: bson.D{{Key: "name", Value: 1}}},
    }
}
```

### Connect and Create a Repository

```go
ctx := context.Background()

db, err := mongoclient.Connect(ctx, "mongodb://localhost:27017", "mydb")
if err != nil {
    log.Fatal(err)
}

userRepo := mongoclient.NewRepository[*User](db.Collection("users"))
```

### Ensure Indexes

```go
// Using self-declared indexes from the IIndex interface:
err = userRepo.EnsureIndexesAssertType(ctx)

// Or pass indexes explicitly:
err = userRepo.EnsureIndexes(ctx, []mongo.IndexModel{
    {Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
})
```

## CRUD Operations

### Insert

```go
// Insert one — returns the inserted document with server-side defaults.
user, err := userRepo.InsertOne(ctx, &User{
    Name:  "Alice",
    Email: "alice@example.com",
    Age:   30,
})
// user.ID, user.CreatedAt, user.UpdatedAt are auto-populated.

// Insert many
ids, err := userRepo.InsertMany(ctx, []*User{
    {Name: "Bob", Email: "bob@example.com", Age: 25},
    {Name: "Charlie", Email: "charlie@example.com", Age: 35},
})
```

### Find

```go
// Find all
users, err := userRepo.Find(ctx, bson.M{})

// Find with filter
users, err := userRepo.Find(ctx, bson.M{"age": bson.M{"$gte": 25}})

// Find one
user, err := userRepo.FindOne(ctx, bson.M{"email": "alice@example.com"})

// Find by ID
user, err := userRepo.FindByID(ctx, objectID)

// Paginated find (page 1, 10 items per page)
users, err := userRepo.FindPaginated(ctx, bson.M{}, 1, 10)

// Paginated find with total count
users, total, err := userRepo.FindPaginatedWithTotal(ctx, bson.M{}, 1, 10)
```

### Update

Update methods auto-detect the update argument: pass a struct and it gets wrapped in `$set`; pass a `bson.M` with `$`-prefixed operators and it passes through as-is.

```go
// Update by ID with a struct (auto-wrapped in $set, BeforeUpdate hook called)
result, err := userRepo.UpdateByID(ctx, user.ID, &User{Name: "Alice Updated"})

// Update by ID with mongo operators
result, err := userRepo.UpdateByID(ctx, user.ID, bson.M{
    "$set": bson.M{"name": "Alice Updated"},
    "$inc": bson.M{"age": 1},
})

// Find one and update — returns the updated document
updated, err := userRepo.FindOneAndUpdate(ctx,
    bson.M{"email": "alice@example.com"},
    bson.M{"$set": bson.M{"age": 31}},
)

// Find one and update by ID
updated, err := userRepo.FindOneAndUpdateByID(ctx, user.ID, bson.M{"$set": bson.M{"age": 31}})

// Update many
result, err := userRepo.UpdateMany(ctx,
    bson.M{"age": bson.M{"$lt": 30}},
    bson.M{"$set": bson.M{"name": "Young User"}},
)
```

### Delete

```go
// Delete one
err = userRepo.DeleteOne(ctx, bson.M{"email": "alice@example.com"})

// Delete by ID
err = userRepo.DeleteByID(ctx, user.ID)

// Delete many — returns deleted count
count, err := userRepo.DeleteMany(ctx, bson.M{"age": bson.M{"$lt": 18}})

// Find one and delete — returns the deleted document
deleted, err := userRepo.FindOneAndDelete(ctx, bson.M{"email": "bob@example.com"})
```

## Aggregation

```go
// Raw aggregation (returns []bson.M)
pipeline := bson.A{
    bson.M{"$match": bson.M{"age": bson.M{"$gte": 25}}},
    bson.M{"$group": bson.M{"_id": nil, "avgAge": bson.M{"$avg": "$age"}}},
}
results, err := userRepo.Aggregate(ctx, pipeline)

// Typed aggregation (returns []*User)
pipeline = bson.A{
    bson.M{"$match": bson.M{"age": bson.M{"$gte": 25}}},
    bson.M{"$sort": bson.M{"age": 1}},
}
users, err := userRepo.AggregateTyped(ctx, pipeline)
```

## Other Operations

```go
// Count
count, err := userRepo.CountDocuments(ctx, bson.M{"age": bson.M{"$gte": 18}})
estimated, err := userRepo.EstimatedCount(ctx)

// Distinct values
emails, err := userRepo.Distinct(ctx, "email", bson.M{})

// Bulk write
bulkResult, err := userRepo.BulkWrite(ctx, []mongo.WriteModel{
    mongo.NewInsertOneModel().SetDocument(&User{Name: "Dave", Email: "dave@example.com"}),
    mongo.NewUpdateOneModel().SetFilter(bson.M{"name": "Alice"}).SetUpdate(bson.M{"$set": bson.M{"age": 32}}),
})

// Change streams
stream, err := userRepo.Watch(ctx, mongo.Pipeline{})

// Access underlying collection
col := userRepo.Collection()
```

## Transactions

```go
err = userRepo.Transaction(ctx, func(sessCtx context.Context) error {
    _, err := userRepo.InsertOne(sessCtx, &User{Name: "Eve", Email: "eve@example.com", Age: 28})
    if err != nil {
        return err
    }
    _, err = userRepo.UpdateOne(sessCtx,
        bson.M{"name": "Alice"},
        bson.M{"$inc": bson.M{"age": 1}},
    )
    return err
})
```

## Lifecycle Hooks

Models embedding `BaseField` automatically get:

- **`BeforeInsert()`** — generates an `_id` (if zero), sets `createdAt` and `updatedAt`
- **`BeforeUpdate()`** — updates `updatedAt`

These hooks are called automatically by `InsertOne`, `InsertMany`, and all update methods (when passing a struct).

## License

MIT
