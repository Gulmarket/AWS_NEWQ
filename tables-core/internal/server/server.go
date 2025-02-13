package server

import (
    "context"
    "fmt"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "go.uber.org/zap"
    "net/http"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "tables/internal/config"
    "tables/internal/storage"
    "time"
)

type Server struct {
    app      *echo.Echo
    conf     *config.Config
    quit     chan os.Signal
    logger   *zap.SugaredLogger
    created  bool
    running  bool
    services *ServerServices
    handlers *ServerHandlers
}

func NewServer(conf *config.Config, logger *zap.SugaredLogger) (*Server, error) {
    var err error
    serv := &Server{
        conf:     conf,
        quit:     make(chan os.Signal),
        logger:   logger,
        created:  false,
        running:  false,
        services: nil,
        handlers: nil,
    }
    //Создаем коннект к базе
    gulmarketRepo, err := storage.NewGulmarketRepo(conf, logger)
    if err != nil {
        return nil, err
    }

    // Создаем сервисы
    serv.services, err = newServices(conf, logger, gulmarketRepo)
    if err != nil {
        return nil, err
    }

    // Создаем контроллеры
    serv.handlers, err = newHandlers(serv.services, logger)
    if err != nil {
        return nil, err
    }

    // Настройка http сервера
    err = serv.Setup()
    if err != nil {
        return nil, err
    }
    return serv, nil
}

// Настройка HTTP сервера
func (s *Server) Setup() error {
    e := echo.New()
    e.Use(middleware.Recover())
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"*"},
        AllowHeaders: []string{"*"},
    }))

    excludeUrls := make(map[string]interface{})
    excludeUrls["/ready"] = nil
    excludeUrls["/live"] = nil
    excludeUrls["/metrics"] = nil
    excludeUrls["/swagger/index.html"] = nil
    excludeUrls["/swagger/swagger-ui.css"] = nil
    excludeUrls["/swagger/swagger-ui-bundle.js"] = nil
    excludeUrls["/swagger/swagger-ui-standalone-preset.js"] = nil
    excludeUrls["/swagger/favicon-32x32.png"] = nil
    excludeUrls["/swagger/doc.json"] = nil

    s.routeSwagger(e)
    s.routeHealth(e)
    s.routeApiV1(e)
    s.app = e

    return nil
}

func (s *Server) RunBlocking() error {
    wg := &sync.WaitGroup{}
    wg.Add(1)

    go s.start(wg)

    wg.Wait()

    return nil
}

// Запускает физический сервер
func (s *Server) start(wg *sync.WaitGroup) {
    defer wg.Done()

    s.logger.Debug("[http-server] запускаем задачи до старта http сервера...")

    go func() {
        s.logger.Debug("[http-server] запускаем http сервер...")

        timeout, _ := time.ParseDuration(s.conf.Server.Timeout)
        srv := &http.Server{
            Addr:         s.conf.Server.Addr,
            ReadTimeout:  timeout,
            WriteTimeout: timeout,
            IdleTimeout:  timeout,
        }

        err := s.app.StartServer(srv)
        if err != nil {
            s.logger.Debug(fmt.Sprintf("[http-server] останавливаем http сервер: ", err))
            time.Sleep(1 * time.Second)
            s.quit <- syscall.SIGTERM
        }
    }()

    s.running = true

    signal.Notify(s.quit, syscall.SIGINT, syscall.SIGTERM)
    <-s.quit // Блокируемся туту пока не получим сообщение об останоовке

    s.logger.Debug("[http-server] останавливаем http сервер...")

    if err := s.app.Shutdown(context.Background()); err != nil {
        s.logger.Debug(fmt.Sprintf("[http-server] ошибка при остановке http сервера: %s", err))
    }

    s.logger.Debug("[http-server] запускаем задачи до остановки http сервера...")

    s.logger.Info("[http-server] сервер остановлен")

    s.running = false
}
