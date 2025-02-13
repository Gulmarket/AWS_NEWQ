package server

import (
    "go.uber.org/zap"
    "tables/internal/handler"
)

type ServerHandlers struct {
    gulmarket handler.Handlers
}

func newHandlers(svc *ServerServices, logger *zap.SugaredLogger) (*ServerHandlers, error) {

    h := &ServerHandlers{}
    h.gulmarket = handler.NewHandler(svc.Srv, logger)

    return h, nil
}
