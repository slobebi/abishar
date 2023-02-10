package app

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"abishar/internal/pkg/redigo"
	"abishar/internal/config"

	"github.com/jmoiron/sqlx"
)

func connectDatabase(config config.Database) *sqlx.DB {
	log.Print("Connecting Database")
	ctx := context.Background()

	connectString := fmt.Sprintf("user=%s dbname=%s host=%s password=%s port=%s sslmode=%s",
		config.User, config.DBName, config.Host, config.Password, config.Port, config.SSLMode,
	)

	var (
		db  *sqlx.DB
		err error
	)

	// connect with retry
	for t := 0; t <= config.Retry; t++ {
		db, err = sqlx.ConnectContext(ctx, "postgres", connectString)
		if err != nil {
			time.Sleep(time.Second * 3)
		} else {
			break
		}
	}

	if err != nil {
		log.Fatal("failed to connect to DB", err)
	}

	log.Print("Database Connected")

	return db
}

func connectRedis(config config.Redis) *redigo.Redis {
	log.Print("Connecting Redis")
	return redigo.New(config)
}
