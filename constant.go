package gomongo_client

const (
	// DefaultPaginationLimit is the default number of documents returned
	// when no explicit limit is set. A value of 1 is almost never
	// desired — 20 is a safe, commonly used default.
	DefaultPaginationLimit int64 = 20
)
