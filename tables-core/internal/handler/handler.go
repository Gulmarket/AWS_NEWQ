package handler

import (
    "fmt"
    "github.com/labstack/echo/v4"
    "go.uber.org/zap"
    "net/http"
    "tables/internal/service"
)

type Handlers interface {
    UploadPlantationDelivery(ctx echo.Context) error
    PlantationDelivery(ctx echo.Context) error
    //UploadExcel(ctx echo.Context) error
}

type Handler struct {
    svc service.IService
    log *zap.SugaredLogger
}

func NewHandler(svc service.IService, log *zap.SugaredLogger) *Handler {
    return &Handler{
        svc: svc,
        log: log,
    }
}

func (h *Handler) PlantationDelivery(ctx echo.Context) error {
    res, err := h.svc.Delivery(ctx.Request().Context())
    if err != nil {
        return ctx.JSON(http.StatusBadRequest, err)
    }
    return ctx.JSON(http.StatusOK, res)
}

func (h *Handler) UploadPlantationDelivery(ctx echo.Context) error {
    spreadsheetId := ctx.Param("spreadsheetId")
    err := h.svc.UploadPlantationDeliveryService(ctx.Request().Context(), spreadsheetId)
    if err != nil {
        return ctx.JSON(http.StatusBadRequest, err)
    }
    message := fmt.Sprintf("google sheet with  id=%s succefully wroten in database", spreadsheetId)
    return ctx.JSON(http.StatusOK, message)
}
