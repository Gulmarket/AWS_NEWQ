package service

import (
    "context"
    "fmt"
    "github.com/gulmarket/model"
    "github.com/gulmarket/model/apperrors"
    "github.com/xuri/excelize/v2"
    "log"
    "net/http"
    "strings"
)

type plantationService struct {
    PlantationRepository model.PlantationRepository
    ProUserRepository    model.ProUserRepository
}

type PSConfig struct {
    PlantationRepository model.PlantationRepository
    ProUserRepository    model.ProUserRepository
}

func NewPlantationService(c *PSConfig) model.PlantationService {
    return &plantationService{
        PlantationRepository: c.PlantationRepository,
        ProUserRepository:    c.ProUserRepository,
    }
}

func (s *plantationService) SignUp(ctx context.Context, p *model.Plantation) error {
    pw, err := hashPassword(p.Password)

    if err != nil {
        log.Printf("Unable to signup plantation for email: %v\n", p.Email)
        return err
    }

    p.Password = pw

    if err := s.PlantationRepository.Create(ctx, p); err != nil {
        return err
    }

    return nil
}

func (s *plantationService) SignIn(ctx context.Context, p *model.Plantation) error {
    pFetched, err := s.PlantationRepository.FindByEmail(ctx, p.Email)

    if err != nil {
        return apperrors.NewAuthorization("Invalid email")
    }

    match, err := comparePasswords(pFetched.Password, p.Password)

    if err != nil {
        return apperrors.NewInternal()
    }

    if !match {
        return apperrors.NewAuthorization("Invalid email and password combination")
    }

    *p = *pFetched
    return nil
}

func (s *plantationService) UpdateLocation(ctx context.Context, p *model.Plantation, id int) error {
    err := s.PlantationRepository.UpdateLocation(ctx, p, id)

    if err != nil {
        return err
    }

    return nil
}

func (s *plantationService) UpdateInfo(ctx context.Context, p *model.Plantation, id int) error {
    err := s.PlantationRepository.UpdateInfo(ctx, p, id)

    if err != nil {
        return err
    }

    return nil
}

func (s *plantationService) Test() {

}

func (s *plantationService) Delivery(ctx context.Context) (*[]model.Delivery, error) {
    var res *[]model.Delivery
    var err error
    res, err = s.PlantationRepository.GetAllDelivery(ctx)

    return res, err
}

func (s *plantationService) UploadPlantationDeliveryService(_ context.Context, google_spreadsheet_id string) error {

    var err error

    resp, err := http.Get(fmt.Sprintf("%s%s/export?format=xlsx", "https://docs.google.com/spreadsheets/d/", google_spreadsheet_id))
    if err != nil {
        return apperrors.NewInternal()
    }
    defer resp.Body.Close()

    f, err := excelize.OpenReader(resp.Body)
    if err != nil {
        return apperrors.NewInternal()
    }

    sheetName := f.GetSheetName(1)
    rows, err := f.GetRows(sheetName)
    if err != nil {
        return apperrors.NewInternal()
    }
    var css []model.Delivery

    for i, row := range rows {
        if i == 0 {
            continue
        }
        var (
            conditions model.Delivery
        )
        for j, cellStrVal := range row {
            cellStrVal = strings.TrimSpace(cellStrVal)
            switch j {

            }
        }
        css = append(css, conditions)
    }

    return nil
}

func (s *plantationService) ShowFilteredProducts(ctx context.Context, params model.Filter, id int) (*model.FlowersFilter, error) {
    res, err := s.PlantationRepository.GetFilteredProducts(ctx, params, id)
    if err != nil {
        return nil, err
    }

    rate, err := getUsdExchangeRate()
    if err != nil {
        return nil, err
    }

    charges, err := s.ProUserRepository.GetFinancialCharges(ctx)

    if err != nil {
        return nil, err
    }

    for i := range res.Deliveries {
        tengePrice, err := getTengePrice(res.Deliveries[i].BoxSize, res.Deliveries[i].Packing, res.Deliveries[i].Price, rate, charges)
        if err != nil {
            return nil, err
        }
        res.Deliveries[i].TengePrice = *tengePrice
    }

    return res, nil
}

func (s *plantationService) Orders(ctx context.Context, id int, status string) (*[]model.PlantationOrder, error) {
    orders, err := s.PlantationRepository.GetOrders(ctx, id, status)

    if err != nil {
        return nil, err
    }

    return orders, nil
}

func (s *plantationService) DeclineOrder(ctx context.Context, orderId int) error {
    err := s.PlantationRepository.DeclineOrder(ctx, orderId)

    if err != nil {
        return err
    }

    return nil
}

func (s *plantationService) ApproveOrder(ctx context.Context, orderId int) error {
    err := s.PlantationRepository.ApproveOrder(ctx, orderId)

    if err != nil {
        return err
    }

    return nil
}

func (s *plantationService) ShowPlantationFilteredProducts(ctx context.Context, params model.Filter, id int) (*model.PlantationFlowersFilter, error) {
    res, err := s.PlantationRepository.GetPlantationFilteredProducts(ctx, params, id)

    if err != nil {
        return nil, err
    }

    return res, nil
}
