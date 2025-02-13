package server

import (
    "github.com/labstack/echo/v4"
    echoSwagger "github.com/swaggo/echo-swagger"
    "net/http"
)

func (s *Server) routeApiV1(r *echo.Echo) {
    apiv1 := r.Group("api/v1")
    apiv1.GET("/google-sheet/:spreadsheetId", s.handlers.gulmarket.UploadPlantationDelivery)
    apiv1.GET("/plantation-delivery", s.handlers.gulmarket.PlantationDelivery)
    //apiv1.POST("/excel-upload", s.handlers.gulmarket.UploadExcel)
}

func (s *Server) routeSwagger(r *echo.Echo) {
    r.GET("/swagger/*", echoSwagger.WrapHandler)
}

func (s *Server) routeHealth(r *echo.Echo) {
    r.GET("/ready", func(c echo.Context) error {
        return c.NoContent(http.StatusOK)
    })
    r.GET("/live", func(c echo.Context) error {
        return c.NoContent(http.StatusOK)
    })
}
