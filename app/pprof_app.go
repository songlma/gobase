package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/songlma/gobase/errorz"
	"github.com/songlma/gobase/logger"
)

type PprofApp struct {
	addr   string
	server *http.Server
}

func NewPprofApp(ctx context.Context, addr string) *PprofApp {
	return &PprofApp{
		addr: addr,
	}
}

func (app *PprofApp) Name() string {
	return "gov2-pprofApp"
}

func (app *PprofApp) Once(ctx context.Context, params string) errorz.Error {
	return nil
}

func (app *PprofApp) Start(ctx context.Context) errorz.Error {
	if app.addr != "" {
		logger.Info(ctx, "PprofApp start ListenAndServe addr:", app.addr)
		app.server = &http.Server{
			Addr: app.addr,
		}
		err := app.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln(fmt.Sprintf("PprofApp StartWebApp %s ListenAndServe err %+v", app.addr, err))
		}
	} else {
		logger.Error(ctx, "PprofApp app addr is empty")
		log.Fatal("PprofApp app addr is empty")
	}
	return nil

}
func (app *PprofApp) Stop(ctx context.Context) errorz.Error {
	if app.server == nil {
		return nil
	}
	err := app.server.Shutdown(ctx)
	if err != nil {
		log.Fatalln(fmt.Sprintf("ShutdownErr:%v", err))
	}
	return nil

}
func (app *PprofApp) Ready(ctx context.Context) bool {
	return true
}
