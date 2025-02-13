package postgre

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v4/pgxpool"
    "go.uber.org/zap"
    "tables/internal/model"
    "time"
)

type GulmarketRepo struct {
    logger zap.SugaredLogger
    pgdb   *pgxpool.Pool
}

func NewGulmarket(logger zap.SugaredLogger, pgdb *pgxpool.Pool) *GulmarketRepo {
    return &GulmarketRepo{logger: logger, pgdb: pgdb}
}

func (r *GulmarketRepo) GetAllDelivery(ctx context.Context) (*[]model.Delivery, error) {
    delivery := []model.Delivery{}
    query := `SELECT id, farm_box, box_type, box_size, mixed, species, product, color, length, price, boxes, packing, plantation_id FROM delivery;`

    rows, err := r.pgdb.Query(context.Background(), query)
    if err != nil {
        r.logger.Error(err)
    }
    defer rows.Close()
    for rows.Next() {
        var rep model.Delivery
        err = rows.Scan(
            &rep.Id,
            &rep.FarmBox,
            &rep.BoxType,
            &rep.BoxSize,
            &rep.Mixed,
            &rep.Species,
            &rep.Product,
            &rep.Color,
            &rep.Length,
            &rep.Price,
            &rep.Boxes,
            &rep.Packing,
            &rep.PlantationId,
        )
        delivery = append(delivery, rep)
        if err != nil {
            r.logger.Info(fmt.Sprintf("Delivery error repository: %s", err))
        }
    }

    return &delivery, err
}

func (r *GulmarketRepo) InsertDelivery(ctx context.Context, request model.Delivery) (id int, err error) {
    ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
    defer cancel()

    queryI := `INSERT INTO delivery (farm_box,  box_type,  box_size,  mixed, species,
    product, color,  length, price, boxes,packing, plantation_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id;`

    row := r.pgdb.QueryRow(ctx, queryI,
        request.FarmBox,
        request.BoxType,
        request.BoxSize,
        request.Mixed,
        request.Species,
        request.Product,
        request.Color,
        request.Length,
        request.Price,
        request.Boxes,
        request.Packing,
        request.PlantationId)
    err = row.Scan(&id)
    if err != nil {
        r.logger.Info(fmt.Sprintf("Insert user profile repositiry %s", err))
    }

    return id, err
}
