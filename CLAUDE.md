# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Generic MongoDB repository library for Go, built on top of the official `go.mongodb.org/mongo-driver/v2`. Package name: `mongoclient`. Import path: `github.com/inc4/gomongo-client`.

## Commands

```bash
go build ./...    # Build
go vet ./...      # Lint
go test ./...     # Run tests (none exist yet)
```

## Architecture

This is a single-package library (`mongoclient`) with a generic `Repository[T]` pattern:

- **`database.go`** — `IRepository[T]` interface and `Document` interface; `Repository[T]` struct wrapping `*mongo.Collection`. `T` must be a pointer to a struct. The `Document` interface (`SetID`, `BeforeInsert`, `BeforeUpdate`) enables lifecycle hooks.
- **`field.go`** — `BaseField` struct (ID, CreatedAt, UpdatedAt) that implements `Document`. Embed this in model structs to get auto-generated IDs and timestamps.
- **`client.go`** — `Connect()` function. Automatically sets ServerAPI v1 for SRV URIs.
- **`crud.go`** — All CRUD operations on `Repository[T]`. Update methods auto-detect whether the update arg is a mongo operator (`$set`, etc.) or a struct (wraps in `$set`). `InsertOne` re-fetches the inserted doc to return it with server-side defaults.
- **`collection.go`** — Non-CRUD collection operations: aggregation, distinct, bulk write, change streams, count.
- **`indexes.go`** — Index management. `EnsureIndexesAssertType` uses an `IIndex` interface on the model type to self-declare indexes.
- **`transaction.go`** — Transaction helper (marked as untested via TODO).
- **`utils.go`** — Internal helpers: `isStructOrPtrToStruct`, `isMongoOperator`.

## Key Patterns

- **Lifecycle hooks**: If `T` implements `Document`, `BeforeInsert()`/`BeforeUpdate()` are called automatically before writes. Embed `BaseField` to get this for free.
- **Index self-declaration**: Models can implement `IIndex` to declare their own indexes, then call `EnsureIndexesAssertType`.
- **Update auto-wrapping**: Passing a struct to update methods auto-wraps it in `bson.M{"$set": ...}`. Passing a `bson.M` with `$`-prefixed keys passes through as-is.
