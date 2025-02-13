package service

import (
    "context"
    "encoding/xml"
    "errors"
    "fmt"
    "github.com/gulmarket/model"
    "github.com/gulmarket/model/apperrors"
    "github.com/gulmarket/types"
    "io"
    "log"
    "net/http"
    "strconv"
    "time"
)

type proUserService struct {
    ProUserRepository model.ProUserRepository
}

type PUSConfig struct {
    ProUserRepository model.ProUserRepository
}

func NewProUserService(c *PUSConfig) model.ProUserService {
    return &proUserService{
        ProUserRepository: c.ProUserRepository,
    }
}

func (s *proUserService) SignUp(ctx context.Context, p *model.ProUser) error {
    pw, err := hashPassword(p.Password)

    if err != nil {
        log.Printf("Unable to signup Pro User for email: %v\n", p.Email)
        return err
    }

    p.Password = pw

    if err := s.ProUserRepository.Create(ctx, p); err != nil {
        return err
    }

    return nil
}

func (s *proUserService) SignIn(ctx context.Context, p *model.ProUser) error {
    pFetched, err := s.ProUserRepository.FindByEmail(ctx, p.Email)

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

func (s *proUserService) AddShop(ctx context.Context, shops []model.Shop, id int) error {
    err := s.ProUserRepository.AddShop(ctx, shops, id)

    if err != nil {
        return err
    }

    return nil
}

func (s *proUserService) Onboarding(ctx context.Context, id int) (string, error) {
    path, err := s.ProUserRepository.GetOnboardingPath(ctx, id)

    if err != nil {
        return "", err
    }

    return path, nil
}

func (s *proUserService) UpdateInitials(ctx context.Context, p *model.ProUser) error {
    err := s.ProUserRepository.UpdateInitials(ctx, p)

    if err != nil {
        return err
    }

    return nil
}

func (s *proUserService) UpdateProUserCity(ctx context.Context, city string, id int) error {
    err := s.ProUserRepository.UpdateProUserRow(ctx, city, "city", id)

    if err != nil {
        return err
    }

    return nil
}

func (s *proUserService) CreateOrder(ctx context.Context, orders []model.Order, id int) error {
    err := s.ProUserRepository.CreateOrder(ctx, orders, id)

    if err != nil {
        return err
    }

    return nil
}

func (s *proUserService) ShowCard(ctx context.Context, deliveryId int) (*model.Card, error) {
    delivery, residue, addresses, err := s.ProUserRepository.GetCardInfo(ctx, deliveryId)

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

    formattedPrice, err := getTengePrice(delivery.BoxSize, delivery.Packing, delivery.Price, rate, charges)

    if err != nil {
        return nil, err
    }

    card := &model.Card{
        Delivery:    *delivery,
        FlowersLeft: *residue,
        TengePrice:  *formattedPrice,
        Addresses:   *addresses,
    }

    return card, nil
}

func (s *proUserService) Orders(ctx context.Context, id int, status string) (*[]model.ProUserOrder, error) {
    orders, err := s.ProUserRepository.GetOrders(ctx, id, status)

    if err != nil {
        return nil, err
    }

    return orders, nil
}

func (s *proUserService) Wallet(ctx context.Context, id int) (float64, error) {
    wallet, err := s.ProUserRepository.GetWallet(ctx, id)

    if err != nil {
        return 0, err
    }

    return wallet, nil
}

func (s *proUserService) ProUserInfo(ctx context.Context, id int) (*model.PrivateOffice, error) {
    info, err := s.ProUserRepository.GetProUserInfo(ctx, id)

    if err != nil {
        return nil, err
    }

    return info, nil
}

func (s *proUserService) UpdateShopsInfo(ctx context.Context, shops []model.Shop, id int) error {
    err := s.ProUserRepository.UpdateShopsInfo(ctx, shops, id)

    if err != nil {
        return err
    }

    return nil
}

func (s *proUserService) AddFavoriteProduct(ctx context.Context, proUserId, deliveryId int) error {
    err := s.ProUserRepository.AddFavoriteProduct(ctx, proUserId, deliveryId)

    if err != nil {
        return err
    }

    return nil
}

func (s *proUserService) RemoveFavoriteProduct(ctx context.Context, proUserId, deliveryId int) error {
    err := s.ProUserRepository.RemoveFavoriteProduct(ctx, proUserId, deliveryId)

    if err != nil {
        return err
    }

    return nil
}

func (s *proUserService) GetFavoriteProducts(ctx context.Context, id int) ([]model.Deliveries, error) {
    deliveries, err := s.ProUserRepository.GetFavoriteProducts(ctx, id)

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

    for i := range deliveries {
        tengePrice, err := getTengePrice(deliveries[i].BoxSize, deliveries[i].Packing, deliveries[i].Price, rate, charges)
        if err != nil {
            return nil, err
        }
        deliveries[i].TengePrice = *tengePrice
    }

    return deliveries, nil
}

func (s *proUserService) UpdateProUserRole(ctx context.Context, id int, role string) error {
    err := s.ProUserRepository.UpdateProUserRole(ctx, id, role)

    if err != nil {
        return err
    }

    return nil
}

func (s *proUserService) GetProUserRole(ctx context.Context, id int) (*string, error) {
    role, err := s.ProUserRepository.GetProUserRole(ctx, id)

    if err != nil {
        return nil, err
    }

    return role, nil
}

func (s *proUserService) AddToCart(ctx context.Context, cart model.Cart) error {
    err := s.ProUserRepository.AddToCart(ctx, cart)

    if err != nil {
        return err
    }

    return nil
}

func (s *proUserService) GetCart(ctx context.Context, id int) ([]model.CartItem, error) {
    cartItems, err := s.ProUserRepository.GetCart(ctx, id)

    rate, err := getUsdExchangeRate()

    if err != nil {
        return nil, err
    }

    charges, err := s.ProUserRepository.GetFinancialCharges(ctx)

    if err != nil {
        return nil, err
    }

    for i := range cartItems {
        tengePrice, err := getTengePrice(cartItems[i].BoxSize, cartItems[i].Packing, cartItems[i].Price, rate, charges)
        if err != nil {
            return nil, err
        }
        cartItems[i].TengePriceForOne = *tengePrice
        cartItems[i].TotalTengePrice = *tengePrice * float64(cartItems[i].Quantity)
    }

    return cartItems, err
}

func (s *proUserService) RemoveFromCart(ctx context.Context, proUserId, deliveryId int) error {
    err := s.ProUserRepository.RemoveFromCart(ctx, proUserId, deliveryId)

    if err != nil {
        return err
    }

    return nil
}

func getUsdExchangeRate() (float64, error) {
    now := time.Now()
    formattedDate := now.Format("02.01.2006")
    url := fmt.Sprintf("https://nationalbank.kz/rss/get_rates.cfm?fdate=%s", formattedDate)

    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Error fetching the URL: %v", err)
        return 0, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading the response body: %v", err)
        return 0, err
    }

    var rss types.RSS
    err = xml.Unmarshal(body, &rss)
    if err != nil {
        log.Printf("Error decoding XML: %v", err)
        return 0, err
    }

    for _, item := range rss.Items {
        if item.Title == "USD" {
            rate, err := strconv.ParseFloat(item.Description, 64)
            if err != nil {
                log.Printf("Couldn't parse string to float: %v", err)
                return 0, err
            }
            return rate, nil
        }
    }

    err = errors.New("couldn't get USD exchange rate")
    return 0, err
}

func getTengePrice(boxSize string, packing int, price float64, rate float64, charges *model.FinancialCharges) (*float64, error) {
    var boxWeight int

    if boxSize == "HB" {
        boxWeight = 27
    } else if boxSize == "QB" {
        boxWeight = 10
    }

    tengePrice := (charges.Tariff*float64(boxWeight)/float64(packing) + price) * charges.WarehouseCommission * charges.BankInstallment * (rate + 2)
    formattedPriceStr := fmt.Sprintf("%.0f", tengePrice)
    formattedPrice, err := strconv.ParseFloat(formattedPriceStr, 64)
    if err != nil {
        fmt.Println("Error converting formatted number:", err)
        return nil, err
    }

    return &formattedPrice, nil
}
