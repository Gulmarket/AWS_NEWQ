package model

import (
    "context"
    "mime/multipart"
    "time"
)

type PlantationService interface {
    SignUp(ctx context.Context, p *Plantation) error
    SignIn(ctx context.Context, p *Plantation) error
    UpdateLocation(ctx context.Context, p *Plantation, id int) error
    UpdateInfo(ctx context.Context, p *Plantation, id int) error
    Delivery(ctx context.Context) (*[]Delivery, error)
    UploadPlantationDeliveryService(ctx context.Context, google_spreadsheet_id string) error
    ShowFilteredProducts(ctx context.Context, params Filter, id int) (*FlowersFilter, error)
    ShowPlantationFilteredProducts(ctx context.Context, params Filter, id int) (*PlantationFlowersFilter, error)
    Orders(ctx context.Context, id int, status string) (*[]PlantationOrder, error)
    DeclineOrder(ctx context.Context, orderId int) error
    ApproveOrder(ctx context.Context, orderId int) error
}

type PlantationRepository interface {
    FindByID(ctx context.Context, id int) (*Plantation, error)
    FindByEmail(ctx context.Context, email string) (*Plantation, error)
    Create(ctx context.Context, p *Plantation) error
    UpdateLocation(ctx context.Context, p *Plantation, id int) error
    UpdateInfo(ctx context.Context, p *Plantation, id int) error
    GetAllDelivery(ctx context.Context) (*[]Delivery, error)
    GetFilteredProducts(ctx context.Context, params Filter, id int) (*FlowersFilter, error)
    GetPlantationFilteredProducts(ctx context.Context, params Filter, id int) (*PlantationFlowersFilter, error)
    GetOrders(ctx context.Context, id int, status string) (*[]PlantationOrder, error)
    DeclineOrder(ctx context.Context, orderId int) error
    ApproveOrder(ctx context.Context, orderId int) error
}

type ProUserService interface {
    SignUp(ctx context.Context, p *ProUser) error
    SignIn(ctx context.Context, p *ProUser) error
    Onboarding(ctx context.Context, id int) (string, error)
    AddShop(ctx context.Context, shops []Shop, id int) error
    UpdateInitials(ctx context.Context, p *ProUser) error
    UpdateProUserCity(ctx context.Context, city string, id int) error
    CreateOrder(ctx context.Context, orders []Order, id int) error
    ShowCard(ctx context.Context, deliveryId int) (*Card, error)
    Orders(ctx context.Context, id int, status string) (*[]ProUserOrder, error)
    Wallet(ctx context.Context, id int) (float64, error)
    ProUserInfo(ctx context.Context, id int) (*PrivateOffice, error)
    UpdateShopsInfo(ctx context.Context, shops []Shop, id int) error
    AddFavoriteProduct(ctx context.Context, proUserId, deliveryId int) error
    RemoveFavoriteProduct(ctx context.Context, proUserId, deliveryId int) error
    GetFavoriteProducts(ctx context.Context, id int) ([]Deliveries, error)
    UpdateProUserRole(ctx context.Context, id int, role string) error
    GetProUserRole(ctx context.Context, id int) (*string, error)
    AddToCart(ctx context.Context, cart Cart) error
    GetCart(ctx context.Context, id int) ([]CartItem, error)
    RemoveFromCart(ctx context.Context, proUserId, deliveryId int) error
}

type ProUserRepository interface {
    FindByID(ctx context.Context, id int) (*ProUser, error)
    FindByEmail(ctx context.Context, email string) (*ProUser, error)
    Create(ctx context.Context, p *ProUser) error
    GetOnboardingPath(ctx context.Context, id int) (string, error)
    AddShop(ctx context.Context, shops []Shop, id int) error
    UpdateInitials(ctx context.Context, p *ProUser) error
    UpdateProUserRow(ctx context.Context, value interface{}, row string, id int) error
    CreateOrder(ctx context.Context, orders []Order, id int) error
    GetCardInfo(ctx context.Context, deliveryId int) (*Delivery, *int, *[]Address, error)
    GetOrders(ctx context.Context, id int, status string) (*[]ProUserOrder, error)
    GetFinancialCharges(ctx context.Context) (*FinancialCharges, error)
    GetWallet(ctx context.Context, id int) (float64, error)
    GetProUserInfo(ctx context.Context, id int) (*PrivateOffice, error)
    UpdateShopsInfo(ctx context.Context, shops []Shop, id int) error
    AddFavoriteProduct(ctx context.Context, proUserId, deliveryId int) error
    RemoveFavoriteProduct(ctx context.Context, proUserId, deliveryId int) error
    GetFavoriteProducts(ctx context.Context, id int) ([]Deliveries, error)
    UpdateProUserRole(ctx context.Context, id int, role string) error
    GetProUserRole(ctx context.Context, id int) (*string, error)
    AddToCart(ctx context.Context, cart Cart) error
    GetCart(ctx context.Context, id int) ([]CartItem, error)
    RemoveFromCart(ctx context.Context, proUserId, deliveryId int) error
}

type TokenService interface {
    NewPairFromPlantation(ctx context.Context, p *Plantation, prevTokenID string) (*TokenPair, error)
    NewPairFromProUser(ctx context.Context, p *ProUser, prevTokenID string) (*TokenPair, error)
    Signout(ctx context.Context, id int) error
    ValidateIDToken(tokenString string) (*Plantation, error)
    ValidateProUserIDToken(tokenString string) (*ProUser, error)
    ValidateRefreshToken(refreshTokenString string) (*RefreshToken, error)
}

type TokenRepository interface {
    SetRefreshToken(ctx context.Context, id string, tokenID string, expiresIn time.Duration) error
    DeleteRefreshToken(ctx context.Context, id string, prevTokenID string) error
    DeletePlantationRefreshTokens(ctx context.Context, id string) error
}

type AWSService interface {
    UploadImage(ctx context.Context, imageFileHeader *multipart.FileHeader) (string, error)
}

type AWSRepository interface {
}
