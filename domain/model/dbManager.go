package model

type DbManager interface {
	GetDatabases(mongoUri string) ([]string, error)
	GetCollections(mongoUri string, databaseName string) ([]string, error)
	Find(mongoUri string, databaseName string, collectionName string, query string) ([]any, error)
	Insert(mongoUri string, databaseName string, collectionName string, documents []any, clearCollection bool) error
}
