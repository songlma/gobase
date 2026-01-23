package app

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/songlma/gobase/errorz"
	"github.com/songlma/gobase/logger"
)

type DefaultApp struct {
	addr   string
	prefix string

	otherApps []App
}

func NewDefaultApp(ctx context.Context, addr, prefix string, app ...App) *DefaultApp {
	return &DefaultApp{
		prefix:    prefix,
		addr:      addr,
		otherApps: app,
	}
}

func GetPromHttpHandler() http.Handler {
	return promhttp.Handler()
}

func GetReadinessHandler(ready func(ctx context.Context) bool) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		if ready(ctx) {
			writer.Write([]byte(time.Now().Format("2006-01-02 15:04:05")))
			return
		}
		writer.WriteHeader(http.StatusInternalServerError)
		return

	})
}

func (app *DefaultApp) Name() string {
	return "gov2-defaultApp"
}

func (app *DefaultApp) Once(ctx context.Context, params string) errorz.Error {
	req := httptest.NewRequest(http.MethodGet, app.prefix+"/metrics", nil)
	resp := httptest.NewRecorder()
	GetPromHttpHandler().ServeHTTP(resp, req)
	resp.Body.String()
	return nil
}
func (app *DefaultApp) Start(ctx context.Context) errorz.Error {
	if app.addr != "" {
		logger.Info(ctx, "start ListenAndServe addr:", app.addr)
		http.Handle(app.prefix+"/metrics", GetPromHttpHandler())
		http.Handle(app.prefix+"/k8s_readiness", GetReadinessHandler(app.Ready))
		err := http.ListenAndServe(app.addr, nil)
		if err != nil || err != http.ErrServerClosed {
			logger.Error(ctx, "web start or accept error:", err)
			log.Fatalf("web start fail:%v", err)
		}
	} else {
		logger.Error(ctx, "default app addr is empty")
		log.Fatal("default app addr is empty")
	}
	return nil

}
func (app *DefaultApp) Stop(ctx context.Context) errorz.Error {
	return nil

}
func (app *DefaultApp) Ready(ctx context.Context) bool {
	for _, myapp := range app.otherApps {
		if myapp.Ready(ctx) == false {
			return false
		}
	}
	return true
}
