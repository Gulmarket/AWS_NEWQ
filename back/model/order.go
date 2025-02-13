package model

import "time"

type Order struct {
    Id               int     `db:"id" json:"id"`
    ProUserId        int     `db:"pro_user_id" json:"pro_user_id"`
    PlantationId     int     `db:"plantation_id" json:"plantation_id"`
    DeliveryId       int     `db:"delivery_id" json:"delivery_id"`
    DeliveryAddress  string  `db:"delivery_address" json:"delivery_address"`
    Quantity         int     `db:"quantity" json:"quantity"`
    PriceForOne      float64 `db:"price_for_one" json:"price_for_one"`
    TengePriceForOne float64 `db:"tenge_price_for_one" json:"tenge_price_for_one"`
    TotalPrice       float64 `db:"total_price" json:"total_price"`
    TotalTengePrice  float64 `db:"total_tenge_price" json:"total_tenge_price"`
    Status           string  `db:"status" json:"status"`
}

type ProUserOrder struct {
    Species          string  `db:"species" json:"species"`
    BoxSize          string  `db:"box_size" json:"box_size"`
    OrderId          int     `db:"id" json:"order_id"`
    PlantationName   string  `db:"name" json:"plantation_name"`
    PlantationId     int     `db:"plantation_id" json:"plantation_id"`
    PriceForOne      float64 `db:"price_for_one" json:"price_for_one"`
    TotalPrice       float64 `db:"total_price" json:"total_price"`
    TengePriceForOne float64 `db:"tenge_price_for_one" json:"tenge_price_for_one"`
    TotalTengePrice  float64 `db:"total_tenge_price" json:"total_tenge_price"`
    DeliveryAddress  string  `db:"delivery_address" json:"delivery_address"`
    Status           string  `db:"status" json:"status"`
}

type PlantationOrder struct {
    Species         string  `db:"species" json:"species"`
    BoxSize         string  `db:"box_size" json:"box_size"`
    OrderId         int     `db:"id" json:"order_id"`
    Quantity        int     `db:"quantity" json:"quantity"`
    TotalPrice      float64 `db:"total_price" json:"total_price"`
    DeliveryAddress string  `db:"delivery_address" json:"delivery_address"`
    Status          string  `db:"status" json:"status"`
}

type Cart struct {
    Id              int       `db:"id" json:"id"`
    ProUserId       int       `db:"pro_user_id" json:"pro_user_id"`
    DeliveryId      int       `db:"delivery_id" json:"delivery_id"`
    DeliveryAddress string    `db:"delivery_address" json:"delivery_address"`
    Quantity        int       `db:"quantity" json:"quantity"`
    TotalPrice      float64   `db:"total_price" json:"total_price"`
    Status          string    `db:"status" json:"status"`
    CreatedAt       time.Time `db:"created_at" json:"created_at"`
    TimeToOverdue   time.Time `db:"time_to_overdue" json:"time_to_overdue"`
}

type CartItem struct {
    Delivery
    DeliveryAddress  string    `db:"delivery_address" json:"delivery_address"`
    Quantity         int       `db:"quantity" json:"quantity"`
    TengePriceForOne float64   `db:"tenge_price_for_one" json:"tenge_price_for_one"`
    TotalPrice       float64   `db:"total_price" json:"total_price"`
    TotalTengePrice  float64   `db:"total_tenge_price" json:"total_tenge_price"`
    Status           string    `db:"status" json:"status"`
    CreatedAt        time.Time `db:"created_at" json:"created_at"`
    TimeToOverdue    time.Time `db:"time_to_overdue" json:"time_to_overdue"`
}

type CartItems struct {
    CartItems []CartItem `json:"cart_items"`
}
