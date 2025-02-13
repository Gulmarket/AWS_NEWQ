package server

import (
    "go.uber.org/zap"
    "tables/internal/config"
    "tables/internal/service"
    "tables/internal/storage"
)

// Все сервисы
type ServerServices struct {
    Srv service.IService
}

// Создать все сервисы
func newServices(conf *config.Config, logger *zap.SugaredLogger, gulmarketStorage *storage.GulmarketStorage) (*ServerServices, error) {
    // Создаем сервисы
    srv := service.NewGulmarketService(conf, logger, gulmarketStorage)

    return &ServerServices{
        Srv: srv,
    }, nil
}
