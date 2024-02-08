package application

import "mongo-exporter-cli/domain/model"

func GetCollectionsList(dbManager model.DbManager, mongoUri string, databaseName string) []string {
	res, _ := dbManager.GetCollections(mongoUri, databaseName)
	return res
}
