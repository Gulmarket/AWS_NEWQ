package types

import "github.com/gulmarket/model"

type SignupReq struct {
    Email        string          `json:"email" binding:"required,email"`
    Password     string          `json:"password" binding:"required,gte=6,lte=30"`
    Name         *string         `json:"name"`
    Phone        string          `json:"phone" binding:"required"`
    Description  *string         `json:"description"`
    Country      *string         `json:"country"`
    City         *string         `json:"city"`
    LogoUrl      *string         `json:"logo_url"`
    WorkSchedule *model.Schedule `json:"work_schedule"`
}

type SigninReq struct {
    Password string `json:"password" binding:"required,gte=6,lte=30"`
    Email    string `json:"email"`
}

type UpdateLocationReq struct {
    Country string `json:"country"`
    City    string `json:"city"`
}

type DeclineOrderReq struct {
    OrderId int `json:"order_id"`
}

type ApproveOrderReq struct {
    OrderId int `json:"order_id"`
}
