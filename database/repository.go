package database

// Repository defines a generic interface for CRUD operations on entities.
// E represents the entity type, and C represents the criteria type used for queries.
type Repository[E any, C any] interface {
	// Create persists a new entity and returns the result of the operation.
	Create(entity E) RepositoryOperationResult[E]

	// FindBy retrieves a single entity matching the given criteria.
	FindBy(criteria C) RepositoryOperationResult[E]

	// FindManyBy retrieves multiple entities matching the given criteria.
	FindManyBy(criteria C) RepositoryOperationResult[[]E]

	// Update modifies an existing entity and returns the result of the operation.
	Update(entity E) RepositoryOperationResult[E]

	// Delete removes an entity and returns the result of the operation.
	Delete(entity E) RepositoryOperationResult[E]
}

// RepositoryOperationResult represents the outcome of a repository operation.
// It contains either the resulting data or an error if the operation failed.
type RepositoryOperationResult[E any] struct {
	Data  *E    // The resulting data, if the operation was successful
	Error error // Any error that occurred during the operation
}
