package main

import (
	"context"
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/nats-io/nats.go"
	"good_api_test/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	natsURL := os.Getenv("NATS_URL")
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	clickhouseHost := os.Getenv("CLICKHOUSE_HOST")
	clickhouseUser := os.Getenv("CLICKHOUSE_USER")
	clickhousePassword := os.Getenv("CLICKHOUSE_PASSWORD")
	clickhouseDB := os.Getenv("CLICKHOUSE_DB")
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{clickhouseHost + ":9000"},
		Auth: clickhouse.Auth{
			Database: clickhouseDB,
			Username: clickhouseUser,
			Password: clickhousePassword,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = nc.Subscribe("goods", func(m *nats.Msg) {
		var good models.Good
		if err := json.Unmarshal(m.Data, &good); err != nil {
			log.Println(err)
			return
		}

		batch, err := conn.PrepareBatch(context.Background(), "INSERT INTO testdb.goods_log")
		if err != nil {
			log.Println(err)
			return
		}
		err = batch.Append(good.Id, good.ProjectId, good.Name, good.Description, good.Priority, good.Removed, good.CreatedAt)
		if err != nil {
			log.Println(err)
			return
		}
		if err := batch.Send(); err != nil {
			log.Println(err)
			return
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Subscriber started")
	select {}
}
