package onevoisdatabase

import (
	"mssql-service/internal/helpers"
	"os"
)

var (
	dbname     = os.Getenv("ONEVOIS_DB_DATABASE")
	password   = os.Getenv("ONEVOIS_DB_PASSWORD")
	username   = os.Getenv("ONEVOIS_DB_USERNAME")
	port       = os.Getenv("ONEVOIS_DB_PORT")
	host       = os.Getenv("ONEVOIS_DB_HOST")
	dbInstance *helpers.ServiceModel
)

func New() helpers.Service {
	return helpers.New(dbname, password, username, port, host, dbInstance)
}
