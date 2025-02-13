package storage

import (
    "context"
    "github.com/jackc/pgx/v4/pgxpool"
    "go.uber.org/zap"
    "tables/internal/config"
    "tables/internal/model"
    "tables/internal/storage/postgre"
)

type GulmarketStorage struct {
    Pg            *pgxpool.Pool
    GulmarketList IGulmarketRepository
}

type IGulmarketRepository interface {
    GetAllDelivery(ctx context.Context) (*[]model.Delivery, error)
    InsertDelivery(ctx context.Context, request model.Delivery) (id int, err error)
}

func NewGulmarketRepo(conf *config.Config, logger *zap.SugaredLogger) (*GulmarketStorage, error) {
    pgDB, err := postgre.Dial(conf.Gulmarketdb.Url)
    if err != nil {
        return nil, err
    }

    var storage GulmarketStorage

    if pgDB != nil {
        storage.Pg = pgDB
        storage.GulmarketList = postgre.NewGulmarket(*logger, storage.Pg)
    }

    return &storage, nil
}
