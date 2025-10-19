package database

// Repository provides an interface to interact with the database
type Repository[E any] interface {
	Create(entity E) RepositoryOperationResult[E]
	FindBy(key any) RepositoryOperationResult[E]
	FindManyBy(key any) RepositoryOperationResult[E]
	Update(entity E) RepositoryOperationResult[E]
	Delete(entity E) RepositoryOperationResult[E]
}

type RepositoryOperationResult[E any] struct {
	Data  *E
	Error error
}
