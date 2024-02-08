package application

import "mongo-exporter-cli/domain/model"

func GetDatabasesList(dbManager model.DbManager, mongoUri string) []string {
	res, _ := dbManager.GetDatabases(mongoUri)
	return res
}
