package types

import "github.com/gulmarket/model"

type ProUserSignupReq struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,gte=6,lte=30"`
    Phone    string `json:"phone" binding:"required"`
    Role     string `json:"role"`
}

type Initials struct {
    Name     string `json:"name"`
    Surname  string `json:"surname"`
    Patronym string `json:"patronym"`
}

type ProUserCity struct {
    City string `json:"city"`
}

type ShopsReq struct {
    Shops []model.Shop `json:"shops"`
}

type OrderReq struct {
    Orders []model.Order `json:"orders"`
}

type CartReq struct {
    DeliveryId      int     `json:"delivery_id"`
    DeliveryAddress string  `json:"delivery_address"`
    Quantity        int     `json:"quantity"`
    TotalPrice      float64 `json:"total_price"`
}

type AddFavoriteProductReq struct {
    DeliveryId int `json:"delivery_id"`
}

type RemoveFavoriteProductReq struct {
    DeliveryId int `json:"delivery_id"`
}

type RemoveFromCartReq struct {
    DeliveryId int `json:"delivery_id"`
}

type UpdateProUserRoleReq struct {
    Role string `json:"role"`
}
