package database

import (
	"context"
)

// Driver is an interface for database.
// It contains all methods that database should implement.
type Driver interface {
	// Database specific methods
	Ping(ctx context.Context) error
	IsInitialized(ctx context.Context) bool
	Close() error

	// Other methods related to database operation
	Student
}
