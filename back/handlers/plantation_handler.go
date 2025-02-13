package handlers

import (
    "fmt"
    "github.com/gulmarket/model"
    "github.com/gulmarket/model/apperrors"
    "github.com/gulmarket/types"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

func (h *Handler) PlantationSignUp(c *gin.Context) {

    var req types.SignupReq

    if ok := bindData(c, &req); !ok {
        return
    }

    p := &model.Plantation{
        Email:    req.Email,
        Password: req.Password,
        Phone:    req.Phone,
    }

    ctx := c.Request.Context()
    err := h.PlantationService.SignUp(ctx, p)

    if err != nil {
        log.Printf("Failed to sign up plantation")
        c.JSON(apperrors.Status(err), gin.H{
            "error": fmt.Sprintf("%v", err),
        })
        return
    }

    tokens, err := h.TokenService.NewPairFromPlantation(ctx, p, "")

    if err != nil {
        log.Printf("Failed to create tokens for plantation: %v\n", err.Error())

        c.JSON(apperrors.Status(err), gin.H{
            "error": err,
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "tokens": tokens,
    })
}

func (h *Handler) PlantationSignIn(c *gin.Context) {
    var req types.SigninReq

    if ok := bindData(c, &req); !ok {
        return
    }

    p := &model.Plantation{
        Password: req.Password,
        Email:    req.Email,
    }

    ctx := c.Request.Context()
    err := h.PlantationService.SignIn(ctx, p)

    if err != nil {
        log.Printf("Failed to sign in plantation: %v\n", err.Error())
        c.JSON(apperrors.Status(err), gin.H{
            "error": err,
        })
        return
    }

    tokens, err := h.TokenService.NewPairFromPlantation(ctx, p, "")

    if err != nil {
        log.Printf("Failed to create tokens for plantations: %v\n", err.Error())

        c.JSON(apperrors.Status(err), gin.H{
            "error": err,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "tokens": tokens,
    })
}

func (h *Handler) UpdatePlantationLocation(c *gin.Context) {
    var req types.UpdateLocationReq

    if ok := bindData(c, &req); !ok {
        return
    }

    p := &model.Plantation{
        Country: &req.Country,
        City:    &req.City,
    }

    plantation, exists := c.Get("plantation")

    if !exists {
        log.Printf("Unable to extract plantation from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := plantation.(*model.Plantation).Id

    err := h.PlantationService.UpdateLocation(c, p, id)

    if err != nil {
        c.JSON(apperrors.Status(err), gin.H{
            "error": err,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
    })
}

func (h *Handler) UpdatePlantationInfo(c *gin.Context) {
    var req model.PlantationInfo

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    p := &model.Plantation{
        Name:         &req.Name,
        Description:  &req.Description,
        LogoUrl:      &req.LogoUrl,
        Country:      &req.Country,
        City:         &req.City,
        WorkSchedule: &req.WorkSchedule,
    }

    plantation, exists := c.Get("plantation")

    if !exists {
        log.Printf("Unable to extract plantation from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := plantation.(*model.Plantation).Id

    err := h.PlantationService.UpdateInfo(c, p, id)

    if err != nil {
        c.JSON(apperrors.Status(err), gin.H{
            "error": err,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
    })
}

func (h *Handler) PlantationDelivery(c *gin.Context) {
    res, err := h.PlantationService.Delivery(c.Request.Context())
    if err != nil {
        c.JSON(apperrors.Status(err), gin.H{
            "error": err,
        })
        return
    }
    return
    c.JSON(http.StatusOK, res)
}

func (h *Handler) UploadPlantationDelivery(c *gin.Context) {

    id := c.Param("google_spreadsheet_id")
    err := h.PlantationService.UploadPlantationDeliveryService(c.Request.Context(), id)
    if err != nil {
        c.JSON(apperrors.Status(err), gin.H{
            "error": err,
        })
        return
    }
    return
    c.JSON(http.StatusOK, gin.H{
        "spread sheed downloading": "done",
    })
}

func (h *Handler) ShowFilteredProducts(c *gin.Context) {
    var req model.Filter

    if ok := bindData(c, &req); !ok {
        return
    }

    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract plantation from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    res, err := h.PlantationService.ShowFilteredProducts(c, req, id)

    if err != nil {
        c.JSON(apperrors.Status(err), gin.H{
            "error": err,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "filtered_products": res,
    })
}

func (h *Handler) PlantationOrders(c *gin.Context) {
    status := c.Query("status")
    plantation, exists := c.Get("plantation")

    if !exists {
        log.Printf("Unable to extract plantation from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := plantation.(*model.Plantation).Id

    orders, err := h.PlantationService.Orders(c, id, status)

    if err != nil {
        c.JSON(apperrors.Status(err), gin.H{
            "error": fmt.Sprintf("%v", err),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "orders":  orders,
    })
}

func (h *Handler) DeclineOrder(c *gin.Context) {
    var req types.DeclineOrderReq

    if ok := bindData(c, &req); !ok {
        return
    }

    _, exists := c.Get("plantation")

    if !exists {
        log.Printf("Unable to extract plantation from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    err := h.PlantationService.DeclineOrder(c, req.OrderId)

    if err != nil {
        c.JSON(apperrors.Status(err), gin.H{
            "error": fmt.Sprintf("%v", err),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
    })
}

func (h *Handler) ApproveOrder(c *gin.Context) {
    var req types.ApproveOrderReq

    if ok := bindData(c, &req); !ok {
        return
    }

    _, exists := c.Get("plantation")

    if !exists {
        log.Printf("Unable to extract plantation from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    err := h.PlantationService.ApproveOrder(c, req.OrderId)

    if err != nil {
        c.JSON(apperrors.Status(err), gin.H{
            "error": fmt.Sprintf("%v", err),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
    })
}

func (h *Handler) PlantationFilteredProducts(c *gin.Context) {
    var req model.Filter

    if ok := bindData(c, &req); !ok {
        return
    }

    plantation, exists := c.Get("plantation")

    if !exists {
        log.Printf("Unable to extract plantation from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := plantation.(*model.Plantation).Id

    res, err := h.PlantationService.ShowPlantationFilteredProducts(c, req, id)

    if err != nil {
        c.JSON(apperrors.Status(err), gin.H{
            "error": err,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "filtered_products": res,
    })
}
