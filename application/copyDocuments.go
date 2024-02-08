package application

import (
	"mongo-exporter-cli/domain/model"
)

func CopyDocuments(dbManager model.DbManager, mongoFrom string, mongoTo string, databaseName string, collectionName string, query string) bool {
	res, _ := dbManager.Find(mongoFrom, databaseName, collectionName, query)

	dbManager.Insert(mongoTo, databaseName, collectionName, res, true)

	return true
}
