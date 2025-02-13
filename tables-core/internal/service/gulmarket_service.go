package service

import (
    "context"
    "fmt"
    "github.com/xuri/excelize/v2"
    "go.uber.org/zap"
    "net/http"
    "strconv"
    "strings"
    "tables/internal/config"
    "tables/internal/model"
    "tables/internal/storage"
)

const (
    farm_column     = 0
    box_type_column = 1
    box_size_column = 2
    mixed_column    = 3
    species_column  = 4
    product_column  = 5
    color_column    = 6
    length_column   = 7
    price_column    = 8
    boxes_column    = 9
    packing_column  = 10
)

type IService interface {
    Delivery(ctx context.Context) (*[]model.Delivery, error)
    UploadPlantationDeliveryService(ctx context.Context, spreadsheetId string) error
}

type Service struct {
    conf   *config.Config
    logger *zap.SugaredLogger
    sprepo *storage.GulmarketStorage
}

func NewGulmarketService(conf *config.Config, logger *zap.SugaredLogger, sprepo *storage.GulmarketStorage) IService {
    return &Service{
        conf:   conf,
        logger: logger,
        sprepo: sprepo,
    }
}

func (s *Service) Delivery(ctx context.Context) (*[]model.Delivery, error) {
    var res *[]model.Delivery
    var err error
    res, err = s.sprepo.GulmarketList.GetAllDelivery(ctx)

    return res, err
}

func (s *Service) UploadPlantationDeliveryService(_ context.Context, spreadsheetId string) error {

    resp, err := http.Get(fmt.Sprintf("%s%s/export?format=xlsx", "https://docs.google.com/spreadsheets/d/", spreadsheetId))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    f, err := excelize.OpenReader(resp.Body)
    if err != nil {
        return err
    }

    sheet_name := f.GetSheetName(0)
    rows, err := f.GetRows(sheet_name)

    if err != nil {
        return err
    }
    var deliveries []model.Delivery

    for i, row := range rows {
        if i == 0 {
            continue
        }

        var (
            delivery model.Delivery
        )
        for j, cellStrVal := range row {
            cellStrVal = strings.TrimSpace(cellStrVal)

            switch j {
            case farm_column:
                delivery.FarmBox = cellStrVal
            case box_type_column:
                delivery.BoxType = cellStrVal
            case box_size_column:
                delivery.BoxSize = cellStrVal
            case mixed_column:
                delivery.Mixed = cellStrVal
            case species_column:
                delivery.Species = cellStrVal
            case product_column:
                delivery.Product = cellStrVal
            case color_column:
                delivery.Color = cellStrVal
            case length_column:
                delivery.Length = cellStrVal
            case price_column:
                delivery.Price, err = strconv.ParseFloat(strings.ReplaceAll(cellStrVal, ",", "."), 64)
                if err != nil {
                    return err
                }
            case boxes_column:
                delivery.Boxes, err = strconv.Atoi(cellStrVal)
                if err != nil {
                    return err
                }
            case packing_column:
                delivery.Packing, err = strconv.Atoi(cellStrVal)
                if err != nil {
                    return err
                }
            }
            delivery.PlantationId = spreadsheetId
        }
        deliveries = append(deliveries, delivery)
    }

    for _, delivery := range deliveries {

        id, err := s.sprepo.GulmarketList.InsertDelivery(nil, delivery)
        s.logger.Info(id)
        if err != nil {
            return err
        }
    }

    return nil
}
