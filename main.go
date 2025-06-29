package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"good_api_test/api"
	"good_api_test/broker"
	"good_api_test/cache"
	"good_api_test/service"
	"good_api_test/storage"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbName)

	var db *storage.PostgresDB
	var err error
	for i := 0; i < 5; i++ {
		db, err = storage.NewPostgresDB(connString)
		if err == nil {
			break
		}
		log.Printf("Could not initialize database (attempt %d): %v. Retrying in 5 seconds...", i+1, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not initialize database after multiple attempts: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	cache := cache.NewRedisCache(redisAddr)

	natsURL := os.Getenv("NATS_URL")
	broker, err := broker.NewNatsBroker(natsURL)
	if err != nil {
		log.Fatalf("Could not initialize broker: %v", err)
	}

	service := service.New(db, cache, broker)

	server := api.NewServer(service)

	log.Println("Server starting on port 8080...")
	if err := server.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
