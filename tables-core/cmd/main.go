package main

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "log"
    "tables/internal/config"
    "tables/internal/server"
    "time"
)

// @title Gulmarket
// @version 0.1
// @description Сервис загрузки таблиц плантаций
// @termsOfService http://swagger.io/terms/

// @contact.name Baurzhan Alzhanov
// @host Пример = api.gulmarket.com
// @BasePath /api/v1
func main() {
    var err error
    time.Local, err = time.LoadLocation("Asia/Almaty")
    if err != nil {
        log.Printf("error loading '%s': %v\n", time.Local, err)
    }

    encoderCfg := zap.NewProductionConfig()
    encoderCfg.EncoderConfig.TimeKey = "timestamp"
    encoderCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    encoderCfg.EncoderConfig.StacktraceKey = ""

    l, err := encoderCfg.Build()
    if err != nil {
        log.Fatal(err)
    }

    logger := l.Sugar()

    defer logger.Sync()
    conf, err := config.NewAppConfig()
    if err != nil {
        log.Fatal("[app] Ошибка при инициализации конфигурации приложения: ", err)
    }
    httpServer, err := server.NewServer(conf, logger)
    if err != nil {
        log.Fatal("Ошибка при инициализации http сервера: ", err)
    }
    err = httpServer.RunBlocking()
    if err != nil {
        log.Fatal("Ошибка при запуске http сервера: ", err)
    }
}
