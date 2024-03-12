package cli

import (
	"bufio"
	"flag"
	"fmt"
	"mongo-exporter-cli/application"
	"mongo-exporter-cli/domain/model"
	"mongo-exporter-cli/infrastructure/mongo"
	"os"
	"strconv"
	"strings"
)

var mongoManager model.DbManager = mongo.MongoManager{}

var (
	mongoTo   string
	mongoFrom string
)

func Init() {
	fmt.Println("\n                                                                    _                      _ _ \n _ __ ___   ___  _ __   __ _  ___         _____  ___ __   ___  _ __| |_ ___ _ __       ___| (_)\n| '_ ` _ \\ / _ \\| '_ \\ / _` |/ _ \\ _____ / _ \\ \\/ / '_ \\ / _ \\| '__| __/ _ \\ '__|____ / __| | |\n| | | | | | (_) | | | | (_| | (_) |_____|  __/>  <| |_) | (_) | |  | ||  __/ | |_____| (__| | |\n|_| |_| |_|\\___/|_| |_|\\__, |\\___/       \\___/_/\\_\\ .__/ \\___/|_|   \\__\\___|_|        \\___|_|_|\n                       |___/                      |_|                                          \n")
	parseVariables()

	databases := application.GetDatabasesList(mongoManager, mongoFrom)

	fmt.Println("Choose database:")
	for i, database := range databases {
		fmt.Printf("%d) %s \n", i, database)
	}
	var database int
	fmt.Scan(&database)
	databaseName := databases[database]

	collections := application.GetCollectionsList(mongoManager, mongoFrom, databaseName)
	fmt.Println("Choose collection:")
	for i, collection := range collections {
		fmt.Printf("%d) %s \n", i, collection)
	}

	var collection int
	fmt.Scan(&collection)
	collectionName := collections[collection]

	fmt.Println("Query(leave it empty to get last " + strconv.Itoa(int(mongo.SearchLimit)) + ")")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	query := scanner.Text()

	if application.CopyDocuments(mongoManager, mongoFrom, mongoTo, databaseName, collectionName, query) {
		fmt.Println("Documents copied successfully")
	}
}

func parseVariables() {
	mongoToPointer := flag.String("mongo-to", "", "Mongo that you are going to copy documents to (default os environment variable MONGO_TO)")
	mongoFromPointer := flag.String("mongo-from", "", "Mongo that you are going to copy documents from (default os environment variable MONGO_FROM)")
	searchLimitPointer := flag.Int64("limit", 500, "Limit of copied documents")

	flag.Parse()

	mongoTo, mongoFrom, mongo.SearchLimit = derefString(mongoToPointer), derefString(mongoFromPointer), derefInt64(searchLimitPointer)

	mongoTo = parseVariable("mongo-to", mongoTo)
	mongoFrom = parseVariable("mongo-from", mongoFrom)
}

func parseVariable(variableName string, variable string) string {
	if variable == "" {
		envName := strings.ToUpper(variableName)
		envName = strings.Replace(envName, "-", "_", -1)

		variable = os.Getenv(strings.ToUpper(envName))

		if variable == "" {
			panic("You need to set up all variables")
		}
	}

	return variable
}

func derefString(s *string) string {
	if s != nil {
		return *s
	}

	panic("Invalid arguments")
}

func derefInt64(s *int64) int64 {
	if s != nil {
		return *s
	}

	panic("Invalid arguments")
}
