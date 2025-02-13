package storage

import (
    "context"
    "fmt"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "github.com/redis/go-redis/v9"
    "log"
    "os"
)

type DataSources struct {
    DB          *sqlx.DB
    RedisClient *redis.Client
}

func InitDS() (*DataSources, error) {
    log.Printf("Initializing data sources\n")
    pgPort := os.Getenv("PG_PORT")
    pgUser := os.Getenv("PG_USER")
    pgPassword := os.Getenv("PG_PASSWORD")
    pgDB := os.Getenv("PG_DB")
    pgSSL := os.Getenv("PG_SSL")

    pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", "postgres_gulmarket", pgPort, pgUser, pgPassword, pgDB, pgSSL)

    log.Printf("Connecting to Postgresql\n")
    db, err := sqlx.Open("postgres", pgConnString)

    if err != nil {
        return nil, fmt.Errorf("error opening db: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error connecting to db: %w", err)
    }

    redisPort := os.Getenv("REDIS_PORT")

    log.Printf("Connecting to Redis\n")
    rdb := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", "redis-back", redisPort),
        Password: "",
        DB:       0,
    })

    _, err = rdb.Ping(context.Background()).Result()

    if err != nil {
        return nil, fmt.Errorf("error connecting to redis: %w", err)
    }

    return &DataSources{
        DB:          db,
        RedisClient: rdb,
    }, nil
}

func (d *DataSources) Close() error {
    if err := d.DB.Close(); err != nil {
        return fmt.Errorf("error closing Postgresql: %w", err)
    }

    if err := d.RedisClient.Close(); err != nil {
        return fmt.Errorf("error closing Redis Client: %w", err)
    }

    return nil
}
