package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type Service interface {
	GetDbInstance() *redis.Client

	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type service struct {
	db *redis.Client
}

var (
	host       = os.Getenv("REDIS_HOST")
	password   = os.Getenv("REDIS_PASSWORD")
	dbInstance *service
)

func New() Service {
	//  Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	// Create a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,     // Redis server address
		Password: password, // Redis server password
		DB:       0,        // Default DB
	})

	dbInstance = &service{
		db: rdb,
	}
	return dbInstance
}

func (s *service) GetDbInstance() *redis.Client {
	return s.db
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the Redis server to check if the connection is working
	_, err := s.db.Ping(ctx).Result()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("Could not connect to Redis: %v", err)
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from host: %s", host)
	return s.db.Close()
}
