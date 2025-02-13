package postgre

import (
    "context"
    "github.com/jackc/pgx/v4/pgxpool"
    "time"
)

func Dial(url string) (*pgxpool.Pool, error) {
    conf, cfgErr := pgxpool.ParseConfig(url)
    if cfgErr != nil {
        return nil, cfgErr
    }
    conf.MaxConns = 20
    conf.MinConns = 10
    conf.MaxConnIdleTime = 10 * time.Second

    conn, connErr := pgxpool.ConnectConfig(context.Background(), conf)

    if connErr != nil {
        return nil, connErr
    }
    return conn, nil
}
