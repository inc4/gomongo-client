package gomongo_client

import (
	"context"
	"errors"
	"testing"
)

func TestOperations_NilCollection(t *testing.T) {
	ctx := context.Background()

	t.Run("CountDocuments", func(t *testing.T) {
		_, err := CountDocuments(ctx, nil, nil)
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})

	t.Run("EstimatedDocumentCount", func(t *testing.T) {
		_, err := EstimatedDocumentCount(ctx, nil)
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})

	t.Run("DistinctRaw", func(t *testing.T) {
		_, err := DistinctRaw(ctx, nil, nil, "field")
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})

	t.Run("FindOneRaw", func(t *testing.T) {
		_, err := FindOneRaw(ctx, nil, nil)
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})

	t.Run("FindRaw", func(t *testing.T) {
		_, err := FindRaw(ctx, nil, nil)
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})

	t.Run("InsertOneRaw", func(t *testing.T) {
		_, err := InsertOneRaw(ctx, nil, struct{}{})
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})

	t.Run("InsertMany", func(t *testing.T) {
		_, err := InsertMany(ctx, nil, []any{})
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})

	t.Run("UpdateByID", func(t *testing.T) {
		_, err := UpdateByID(ctx, nil, "id", struct{}{})
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})

	t.Run("UpdateOne", func(t *testing.T) {
		_, err := UpdateOne(ctx, nil, nil, struct{}{})
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})

	t.Run("UpdateMany", func(t *testing.T) {
		_, err := UpdateMany(ctx, nil, nil, struct{}{})
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})

	t.Run("DeleteOne", func(t *testing.T) {
		_, err := DeleteOne(ctx, nil, nil)
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})

	t.Run("DeleteMany", func(t *testing.T) {
		_, err := DeleteMany(ctx, nil, nil)
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})

	t.Run("AggregateRaw nil collection", func(t *testing.T) {
		_, err := AggregateRaw(ctx, nil, []any{})
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})

	t.Run("FindOneAndUpdateRaw", func(t *testing.T) {
		_, err := FindOneAndUpdateRaw(ctx, nil, nil, struct{}{})
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})

	t.Run("FindOneAndDeleteRaw", func(t *testing.T) {
		_, err := FindOneAndDeleteRaw(ctx, nil, nil)
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})

	t.Run("FindOneAndReplaceRaw", func(t *testing.T) {
		_, err := FindOneAndReplaceRaw(ctx, nil, nil, struct{}{})
		if !errors.Is(err, ErrNilCollection) {
			t.Errorf("error = %v, want %v", err, ErrNilCollection)
		}
	})
}

func TestOperations_NilDocument(t *testing.T) {
	ctx := context.Background()

	t.Run("FindOne nil document", func(t *testing.T) {
		err := FindOne(ctx, nil, nil, nil)
		if err == nil {
			t.Error("FindOne(nil document) should error")
		}
	})

	t.Run("Find nil documents", func(t *testing.T) {
		err := Find(ctx, nil, nil, nil)
		if err == nil {
			t.Error("Find(nil documents) should error")
		}
	})

	t.Run("FindOneAndUpdate nil document", func(t *testing.T) {
		err := FindOneAndUpdate(ctx, nil, nil, nil, struct{}{})
		if err == nil {
			t.Error("FindOneAndUpdate(nil document) should error")
		}
	})

	t.Run("FindOneAndDelete nil document", func(t *testing.T) {
		err := FindOneAndDelete(ctx, nil, nil, nil)
		if err == nil {
			t.Error("FindOneAndDelete(nil document) should error")
		}
	})

	t.Run("FindOneAndReplace nil document", func(t *testing.T) {
		err := FindOneAndReplace(ctx, nil, nil, nil, struct{}{})
		if err == nil {
			t.Error("FindOneAndReplace(nil document) should error")
		}
	})

	t.Run("Aggregate nil documents", func(t *testing.T) {
		err := Aggregate(ctx, nil, nil, []any{})
		if err == nil {
			t.Error("Aggregate(nil documents) should error")
		}
	})
}
