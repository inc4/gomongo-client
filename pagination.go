package gomongo_client

import (
	"encoding/base64"
	"encoding/json"
)

// OffsetPagination holds skip/limit parameters for offset-based
// pagination queries.
type OffsetPagination struct {
	Skip  int64 `json:"skip" bson:"skip"`
	Limit int64 `json:"limit" bson:"limit"`
}

// DefaultOffsetPagination returns an [OffsetPagination] with sensible
// defaults (skip 0, limit [DefaultPaginationLimit]).
func DefaultOffsetPagination() *OffsetPagination {
	return &OffsetPagination{
		Limit: DefaultPaginationLimit,
	}
}

// OffsetPaginationResult holds paginated items and the total count.
// Used as the decoded result of a $facet aggregation.
type OffsetPaginationResult[T any] struct {
	Items      []T   `json:"items" bson:"items"`
	TotalCount int64 `json:"totalCount" bson:"totalCount"`
}

// CursorDirection represents the traversal direction for cursor-based
// pagination.
type CursorDirection string

const (
	// CursorNext pages forward through the result set.
	CursorNext CursorDirection = "next"

	// CursorPrev pages backward through the result set.
	CursorPrev CursorDirection = "prev"
)

// AllCursorDirectionValues returns all valid cursor direction values.
func AllCursorDirectionValues() []CursorDirection {
	return []CursorDirection{CursorNext, CursorPrev}
}

// String returns the string representation.
func (d CursorDirection) String() string {
	return string(d)
}

// IsValid reports whether d is a recognized direction.
func (d CursorDirection) IsValid() bool {
	switch d {
	case CursorNext, CursorPrev:
		return true
	}
	return false
}

// CursorPagination holds parameters for cursor-based pagination queries.
type CursorPagination struct {
	Cursor    string          `json:"cursor" bson:"cursor"`
	Direction CursorDirection `json:"direction" bson:"direction"`
	Limit     int64           `json:"limit" bson:"limit"`
}

// DefaultCursorPagination returns a [CursorPagination] with sensible
// defaults (no cursor, limit [DefaultPaginationLimit], direction next).
func DefaultCursorPagination() *CursorPagination {
	return &CursorPagination{
		Direction: CursorNext,
		Limit:     DefaultPaginationLimit,
	}
}

// CursorPaginationResult holds paginated items with cursor metadata.
type CursorPaginationResult[T any] struct {
	Items      []T    `json:"items" bson:"items"`
	NextCursor string `json:"nextCursor,omitempty" bson:"nextCursor,omitempty"`
	PrevCursor string `json:"prevCursor,omitempty" bson:"prevCursor,omitempty"`
	HasNext    bool   `json:"hasNext" bson:"hasNext"`
	HasPrev    bool   `json:"hasPrev" bson:"hasPrev"`
}

// DecodeCursor decodes a base64url-encoded JSON cursor into the
// target type T. This is the inverse of [EncodeCursor].
//
// Example:
//
//	type MyCursor struct {
//	    ID bson.ObjectID `json:"id"`
//	}
//	cursor, err := mongorm.DecodeCursor[MyCursor](rawCursor)
func DecodeCursor[T any](cursor string) (T, error) {
	var result T

	data, err := base64.URLEncoding.DecodeString(cursor)

	if err != nil {
		return result, err
	}

	if err = json.Unmarshal(data, &result); err != nil {
		return result, err
	}

	return result, nil
}

// EncodeCursor encodes a cursor value as base64url-encoded JSON.
// The cursor can be any JSON-serialisable type (typically a struct
// containing the sort fields).
//
// Example:
//
//	encoded, err := mongorm.EncodeCursor(MyCursor{ID: lastDoc.ID})
func EncodeCursor[T any](cursor T) (string, error) {
	data, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(data), nil
}
