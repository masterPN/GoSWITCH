package wholesaledatabase

import (
	"mssql-service/internal/helpers"
	"os"
)

var (
	dbname     = os.Getenv("WHOLESALE_DB_DATABASE")
	password   = os.Getenv("WHOLESALE_DB_PASSWORD")
	username   = os.Getenv("WHOLESALE_DB_USERNAME")
	port       = os.Getenv("WHOLESALE_DB_PORT")
	host       = os.Getenv("WHOLESALE_DB_HOST")
	dbInstance *helpers.ServiceModel
)

func New() helpers.Service {
	return helpers.New(dbname, password, username, port, host, dbInstance)
}
