package gomongo_client

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

type testCursor struct {
	ID    string `json:"id"`
	Score int    `json:"score"`
}

func TestEncodeCursor(t *testing.T) {
	cursor := testCursor{ID: "abc123", Score: 42}

	encoded, err := EncodeCursor(cursor)
	if err != nil {
		t.Fatalf("EncodeCursor() error: %v", err)
	}
	if encoded == "" {
		t.Fatal("EncodeCursor() returned empty string")
	}

	// Verify it's valid base64url.
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("result is not valid base64url: %v", err)
	}

	// Verify it's valid JSON.
	var decoded testCursor
	if err = json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("decoded base64 is not valid JSON: %v", err)
	}
	if decoded.ID != cursor.ID || decoded.Score != cursor.Score {
		t.Errorf("decoded = %+v, want %+v", decoded, cursor)
	}
}

func TestDecodeCursor(t *testing.T) {
	original := testCursor{ID: "xyz789", Score: 99}

	encoded, err := EncodeCursor(original)
	if err != nil {
		t.Fatalf("EncodeCursor() error: %v", err)
	}

	decoded, err := DecodeCursor[testCursor](encoded)
	if err != nil {
		t.Fatalf("DecodeCursor() error: %v", err)
	}

	if decoded.ID != original.ID {
		t.Errorf("ID = %q, want %q", decoded.ID, original.ID)
	}
	if decoded.Score != original.Score {
		t.Errorf("Score = %d, want %d", decoded.Score, original.Score)
	}
}

func TestDecodeCursor_InvalidBase64(t *testing.T) {
	_, err := DecodeCursor[testCursor]("!!!not-base64!!!")
	if err == nil {
		t.Error("DecodeCursor() should fail for invalid base64")
	}
}

func TestDecodeCursor_InvalidJSON(t *testing.T) {
	encoded := base64.URLEncoding.EncodeToString([]byte("not-json"))

	_, err := DecodeCursor[testCursor](encoded)
	if err == nil {
		t.Error("DecodeCursor() should fail for invalid JSON")
	}
}

func TestEncodeDecode_RoundTrip(t *testing.T) {
	cursors := []testCursor{
		{ID: "", Score: 0},
		{ID: "simple", Score: 1},
		{ID: "unicode-日本語", Score: 999},
	}

	for _, original := range cursors {
		encoded, err := EncodeCursor(original)
		if err != nil {
			t.Fatalf("EncodeCursor(%+v) error: %v", original, err)
		}

		decoded, err := DecodeCursor[testCursor](encoded)
		if err != nil {
			t.Fatalf("DecodeCursor() error: %v", err)
		}

		if decoded != original {
			t.Errorf("round-trip: got %+v, want %+v", decoded, original)
		}
	}
}

func TestCursorDirection_IsValid(t *testing.T) {
	tests := []struct {
		d    CursorDirection
		want bool
	}{
		{CursorNext, true},
		{CursorPrev, true},
		{"", false},
		{"sideways", false},
		{"NEXT", false},
	}

	for _, tt := range tests {
		if got := tt.d.IsValid(); got != tt.want {
			t.Errorf("CursorDirection(%q).IsValid() = %v, want %v", tt.d, got, tt.want)
		}
	}
}

func TestCursorDirection_String(t *testing.T) {
	if s := CursorNext.String(); s != "next" {
		t.Errorf("CursorNext.String() = %q, want %q", s, "next")
	}
	if s := CursorPrev.String(); s != "prev" {
		t.Errorf("CursorPrev.String() = %q, want %q", s, "prev")
	}
}

func TestAllCursorDirectionValues(t *testing.T) {
	values := AllCursorDirectionValues()
	if len(values) != 2 {
		t.Errorf("AllCursorDirectionValues() len = %d, want 2", len(values))
	}
}

func TestDefaultOffsetPagination(t *testing.T) {
	p := DefaultOffsetPagination()
	if p.Skip != 0 {
		t.Errorf("Skip = %d, want 0", p.Skip)
	}
	if p.Limit != DefaultPaginationLimit {
		t.Errorf("Limit = %d, want %d", p.Limit, DefaultPaginationLimit)
	}
}

func TestDefaultCursorPagination(t *testing.T) {
	p := DefaultCursorPagination()
	if p.Cursor != "" {
		t.Errorf("Cursor = %q, want empty", p.Cursor)
	}
	if p.Limit != DefaultPaginationLimit {
		t.Errorf("Limit = %d, want %d", p.Limit, DefaultPaginationLimit)
	}
	if p.Direction != CursorNext {
		t.Errorf("Direction = %q, want %q", p.Direction, CursorNext)
	}
}
