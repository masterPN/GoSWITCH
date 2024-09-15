package server

import (
	"fmt"
	"mssql-service/internal/data/onevoisdata"
	"mssql-service/internal/data/wholesaledata"
	"mssql-service/internal/database/onevoisdatabase"
	"mssql-service/internal/database/wholesaledatabase"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port            int
	onevoisModels   onevoisdata.Models
	wholesaleModels wholesaledata.Models
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:            port,
		onevoisModels:   onevoisdata.NewModels(onevoisdatabase.New().GetDbInstance()),
		wholesaleModels: wholesaledata.NewModels(wholesaledatabase.New().GetDbInstance()),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
