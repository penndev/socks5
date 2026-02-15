package main

import (
	"context"
)

type Server struct {
	Host     string
	Port     int
	Username string
	Password string
	appCtx   context.Context
}

func (t *Server) Startup(ctx context.Context) {
	t.appCtx = ctx
}
