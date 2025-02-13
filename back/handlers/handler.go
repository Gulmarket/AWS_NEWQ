package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/gulmarket/handlers/middleware"
    "github.com/gulmarket/model"
    "net/http"
    "time"
)

type Handler struct {
    PlantationService model.PlantationService
    ProUserService    model.ProUserService
    TokenService      model.TokenService
    AWSService        model.AWSService
    MaxBodyBytes      int64
}

type Config struct {
    R                 *gin.Engine
    PlantationService model.PlantationService
    ProUserService    model.ProUserService
    TokenService      model.TokenService
    AWSService        model.AWSService
    TimeoutDuration   time.Duration
    BaseURL           string
    MaxBodyBytes      int64
}

func NewHandler(c *Config) {

    h := &Handler{
        PlantationService: c.PlantationService,
        ProUserService:    c.ProUserService,
        TokenService:      c.TokenService,
        AWSService:        c.AWSService,
        MaxBodyBytes:      c.MaxBodyBytes,
    }

    g := c.R.Group(c.BaseURL)

    g.POST("/plantation-signup", h.PlantationSignUp)
    g.GET("/onboarding", middleware.AuthProUser(h.TokenService), h.Onboarding)
    g.POST("/plantation-signin", h.PlantationSignIn)
    g.PUT("/update-plantation-location", middleware.AuthPlantation(h.TokenService), h.UpdatePlantationLocation)
    g.PUT("/update-plantation-info", middleware.AuthPlantation(h.TokenService), h.UpdatePlantationInfo)
    g.POST("/pro-user-signup", h.ProUserSignUp)
    g.POST("/pro-user-signin", h.ProUserSignIn)
    g.PUT("/update-initials", middleware.AuthProUser(h.TokenService), h.UpdateInitials)
    g.PUT("/update-pro-user-city", middleware.AuthProUser(h.TokenService), h.UpdateProUserCity)
    g.POST("/add-shop", middleware.AuthProUser(h.TokenService), h.AddShop)
    g.GET("/test", h.Test)
    g.GET("/plantation-delivery", h.PlantationDelivery)
    g.GET("/google-sheet/:google_spreadsheet_id", h.UploadPlantationDelivery)
    g.POST("/show-filtered-products", middleware.AuthProUser(h.TokenService), h.ShowFilteredProducts)
    g.GET("/show-pro-user-card/:delivery_id", h.ShowProUserCard)
    g.POST("/create-order", middleware.AuthProUser(h.TokenService), h.CreateOrder)
    g.GET("/pro-user-orders", middleware.AuthProUser(h.TokenService), h.ProUserOrders)
    g.GET("/plantation-orders", middleware.AuthPlantation(h.TokenService), h.PlantationOrders)
    g.POST("/decline-order", middleware.AuthPlantation(h.TokenService), h.DeclineOrder)
    g.POST("/approve-order", middleware.AuthPlantation(h.TokenService), h.ApproveOrder)
    g.GET("/wallet", middleware.AuthProUser(h.TokenService), h.Wallet)
    g.GET("/pro-user-info", middleware.AuthProUser(h.TokenService), h.GetProUserInfo)
    g.PUT("/update-shops-info", middleware.AuthProUser(h.TokenService), h.UpdateShopsInfo)
    g.POST("/add-favorite-product", middleware.AuthProUser(h.TokenService), h.AddFavoriteProduct)
    g.GET("/get-favorite-products", middleware.AuthProUser(h.TokenService), h.GetFavoriteProducts)
    g.PUT("/update-pro-user-role", middleware.AuthProUser(h.TokenService), h.UpdateProUserRole)
    g.GET("/pro-user-role", middleware.AuthProUser(h.TokenService), h.GetProUserRole)
    g.POST("/add-to-cart", middleware.AuthProUser(h.TokenService), h.AddToCart)
    g.DELETE("/remove-favorite-product", middleware.AuthProUser(h.TokenService), h.RemoveFavoriteProduct)
    g.GET("/cart", middleware.AuthProUser(h.TokenService), h.Cart)
    g.DELETE("/remove-from-cart", middleware.AuthProUser(h.TokenService), h.RemoveFromCart)
    g.POST("/upload-image", h.UploadImage)
    g.POST("/show-plantation-filtered-products", middleware.AuthPlantation(h.TokenService), h.PlantationFilteredProducts)

}

func (h *Handler) Test(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "test": "aufff",
    })
}
