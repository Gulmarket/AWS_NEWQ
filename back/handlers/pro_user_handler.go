package handlers

import (
    "errors"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/gulmarket/model"
    "github.com/gulmarket/model/apperrors"
    "github.com/gulmarket/types"
    "log"
    "net/http"
    "strconv"
)

func (h *Handler) ProUserSignUp(c *gin.Context) {

    var req types.ProUserSignupReq

    if ok := bindData(c, &req); !ok {
        return
    }

    p := &model.ProUser{
        Email:    req.Email,
        Password: req.Password,
        Phone:    req.Phone,
    }

    ctx := c.Request.Context()
    err := h.ProUserService.SignUp(ctx, p)

    if err != nil {
        log.Printf("Failed to sign up pro user")

        c.JSON(apperrors.Status(err), gin.H{
            "error": err.Error(),
        })

        return
    }

    tokens, err := h.TokenService.NewPairFromProUser(ctx, p, "")

    if err != nil {
        log.Printf("Failed to create tokens for pro user: %v\n", err.Error())

        c.JSON(apperrors.Status(err), gin.H{
            "error": err,
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "tokens": tokens,
    })
}

func (h *Handler) ProUserSignIn(c *gin.Context) {
    var req types.SigninReq

    if ok := bindData(c, &req); !ok {
        return
    }

    p := &model.ProUser{
        Password: req.Password,
        Email:    req.Email,
    }

    ctx := c.Request.Context()
    err := h.ProUserService.SignIn(ctx, p)

    if err != nil {
        log.Printf("Failed to sign in Pro User: %v\n", err.Error())
        c.JSON(apperrors.Status(err), gin.H{
            "error": err,
        })
        return
    }

    tokens, err := h.TokenService.NewPairFromProUser(ctx, p, "")

    if err != nil {
        log.Printf("Failed to create tokens for pro user: %v\n", err.Error())

        c.JSON(apperrors.Status(err), gin.H{
            "error": err,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "tokens": tokens,
    })
}

func (h *Handler) AddShop(c *gin.Context) {
    var req types.ShopsReq

    if ok := bindData(c, &req); !ok {
        return
    }

    id, err := getIdFromToken(c)

    if err != nil {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": err,
        })
        return
    }

    err = h.ProUserService.AddShop(c, req.Shops, *id)

    if err != nil {
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": err,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
    })

}

func (h *Handler) Onboarding(c *gin.Context) {
    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": err,
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    path, err := h.ProUserService.Onboarding(c, id)

    if err != nil {
        status := apperrors.NewInternal()
        c.JSON(status.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "path":    path,
    })
}

func (h *Handler) UpdateProUserCity(c *gin.Context) {
    var req types.ProUserCity

    if ok := bindData(c, &req); !ok {
        return
    }

    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": err,
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    err := h.ProUserService.UpdateProUserCity(c, req.City, id)

    if err != nil {
        log.Printf("Failed to update city: %v\n", err.Error())

        c.JSON(apperrors.Status(err), gin.H{
            "error": err,
        })

        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
    })

}

func (h *Handler) UpdateInitials(c *gin.Context) {
    var req types.Initials

    if ok := bindData(c, &req); !ok {
        return
    }

    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": err,
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    p := model.ProUser{
        Id:       id,
        Name:     &req.Name,
        Surname:  &req.Surname,
        Patronym: &req.Patronym,
    }

    err := h.ProUserService.UpdateInitials(c, &p)

    if err != nil {
        log.Printf("Failed to update initials: %v\n", err.Error())

        c.JSON(apperrors.Status(err), gin.H{
            "error": err,
        })
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
    })
}

func (h *Handler) CreateOrder(c *gin.Context) {
    var req types.OrderReq

    if ok := bindData(c, &req); !ok {
        return
    }

    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    err := h.ProUserService.CreateOrder(c, req.Orders, id)

    if err != nil {
        log.Printf("Failed to create orders: %v\n", err.Error())

        c.JSON(apperrors.Status(err), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
    })
}

func (h *Handler) ShowProUserCard(c *gin.Context) {
    paramId := c.Param("delivery_id")
    deliveryId, err := strconv.Atoi(paramId)

    if err != nil {
        log.Printf("Failed to convert string to int: %v", err)

        c.JSON(apperrors.Status(err), gin.H{
            "error": fmt.Sprintf("%v", err),
        })
        return
    }

    card, err := h.ProUserService.ShowCard(c, deliveryId)

    if err != nil {
        log.Printf("Failed to get card info: %v", err)

        c.JSON(apperrors.Status(err), gin.H{
            "error": fmt.Sprintf("%v", err),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "card":    card,
    })
}

func (h *Handler) ProUserOrders(c *gin.Context) {
    status := c.Query("status")
    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    orders, err := h.ProUserService.Orders(c, id, status)

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

func (h *Handler) Wallet(c *gin.Context) {
    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    wallet, err := h.ProUserService.Wallet(c, id)

    if err != nil {
        c.JSON(apperrors.Status(err), gin.H{
            "error": fmt.Sprintf("%v", err),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "wallet":  wallet,
    })
}

func (h *Handler) GetProUserInfo(c *gin.Context) {
    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    proUserInfo, err := h.ProUserService.ProUserInfo(c, id)

    if err != nil {
        c.JSON(apperrors.Status(err), gin.H{
            "error": fmt.Sprintf("%v", err),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success":       true,
        "pro_user_info": proUserInfo,
    })
}

func (h *Handler) UpdateShopsInfo(c *gin.Context) {
    var req types.ShopsReq

    if ok := bindData(c, &req); !ok {
        return
    }

    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    err := h.ProUserService.UpdateShopsInfo(c, req.Shops, id)

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

func (h *Handler) AddFavoriteProduct(c *gin.Context) {
    var req types.AddFavoriteProductReq

    if ok := bindData(c, &req); !ok {
        return
    }

    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    err := h.ProUserService.AddFavoriteProduct(c, id, req.DeliveryId)

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

func (h *Handler) GetFavoriteProducts(c *gin.Context) {
    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    deliveries, err := h.ProUserService.GetFavoriteProducts(c, id)

    if err != nil {
        c.JSON(apperrors.Status(err), gin.H{
            "error": fmt.Sprintf("%v", err),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success":    true,
        "deliveries": deliveries,
    })
}

func (h *Handler) UpdateProUserRole(c *gin.Context) {
    var req types.UpdateProUserRoleReq

    if ok := bindData(c, &req); !ok {
        return
    }

    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := proUser.(*model.ProUser).Id
    email := proUser.(*model.ProUser).Email
    phone := proUser.(*model.ProUser).Phone
    password := proUser.(*model.ProUser).Password

    err := h.ProUserService.UpdateProUserRole(c, id, req.Role)

    if err != nil {
        c.JSON(apperrors.Status(err), gin.H{
            "error": fmt.Sprintf("%v", err),
        })
        return
    }

    p := &model.ProUser{
        Id:       id,
        Email:    email,
        Phone:    phone,
        Password: password,
        Role:     req.Role,
    }

    tokens, err := h.TokenService.NewPairFromProUser(c, p, "")

    if err != nil {
        log.Printf("Failed to create tokens for pro user: %v\n", err.Error())

        c.JSON(apperrors.Status(err), gin.H{
            "error": err,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success":        true,
        "updated_tokens": tokens,
    })
}

func (h *Handler) GetProUserRole(c *gin.Context) {
    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    role, err := h.ProUserService.GetProUserRole(c, id)

    if err != nil {
        c.JSON(apperrors.Status(err), gin.H{
            "error": fmt.Sprintf("%v", err),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "role":    role,
    })
}

func (h *Handler) AddToCart(c *gin.Context) {
    var req types.CartReq

    if ok := bindData(c, &req); !ok {
        return
    }

    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    cart := model.Cart{
        ProUserId:       id,
        DeliveryId:      req.DeliveryId,
        DeliveryAddress: req.DeliveryAddress,
        Quantity:        req.Quantity,
        TotalPrice:      req.TotalPrice,
    }

    err := h.ProUserService.AddToCart(c, cart)

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

func (h *Handler) RemoveFavoriteProduct(c *gin.Context) {
    var req types.RemoveFavoriteProductReq

    if ok := bindData(c, &req); !ok {
        return
    }

    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    err := h.ProUserService.RemoveFavoriteProduct(c, id, req.DeliveryId)

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

func (h *Handler) Cart(c *gin.Context) {
    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    cartItems, err := h.ProUserService.GetCart(c, id)

    if err != nil {
        c.JSON(apperrors.Status(err), gin.H{
            "error": fmt.Sprintf("%v", err),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success":    true,
        "cart_items": cartItems,
    })
}

func (h *Handler) RemoveFromCart(c *gin.Context) {
    var req types.RemoveFromCartReq

    if ok := bindData(c, &req); !ok {
        return
    }

    proUser, exists := c.Get("pro_user")

    if !exists {
        log.Printf("Unable to extract pro user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": fmt.Sprintf("%v", err),
        })

        return
    }

    id := proUser.(*model.ProUser).Id

    err := h.ProUserService.RemoveFromCart(c, id, req.DeliveryId)

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

func (h *Handler) UploadImage(c *gin.Context) {

    c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.MaxBodyBytes)

    log.Println("Starting file upload process...")

    imageFileHeader, err := c.FormFile("imageFile")
    if err != nil {
        log.Printf("Unable to parse multipart/form-data: %+v", err)
        if err.Error() == "multipart: NextPart: http: request body too large" {
            c.JSON(http.StatusRequestEntityTooLarge, gin.H{
                "error": fmt.Sprintf("Max request body size is %v bytes\n", h.MaxBodyBytes),
            })
            return
        }
        e := apperrors.NewBadRequest("Unable to parse multipart/form-data")
        c.JSON(e.Status(), gin.H{"error": e})
        return
    }

    if imageFileHeader == nil {
        log.Println("No file was uploaded")
        e := apperrors.NewBadRequest("Must include an imageFile")
        c.JSON(e.Status(), gin.H{"error": e})
        return
    }

    mimeType := imageFileHeader.Header.Get("Content-Type")
    if valid := isAllowedImageType(mimeType); !valid {
        log.Println("Image is not an allowable mime-type")
        e := apperrors.NewBadRequest("imageFile must be 'image/jpeg' or 'image/png'")
        c.JSON(e.Status(), gin.H{"error": e})
        return
    }

    // Check if AWSService is initialized
    if h.AWSService == nil {
        log.Println("AWSService is nil")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    log.Printf("Uploading file: %s with size %d and MIME type %s\n", imageFileHeader.Filename, imageFileHeader.Size, mimeType)

    imageUrl, err := h.AWSService.UploadImage(c, imageFileHeader)
    if err != nil {
        log.Printf("Failed to upload image: %v\n", err)
        c.JSON(apperrors.Status(err), gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "image_url": imageUrl,
        "success":   true,
    })
}

func getIdFromToken(c *gin.Context) (*int, error) {
    proUser, exists := c.Get("pro_user")

    if !exists {
        err := errors.New("unable to extract pro user from request context")

        return nil, err
    }

    id := proUser.(*model.ProUser).Id

    return &id, nil
}
