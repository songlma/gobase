package app

import (
	"context"

	"github.com/songlma/gobase/errorz"
)

type App interface {
	Name() string
	Once(ctx context.Context, params string) errorz.Error
	Start(ctx context.Context) errorz.Error
	Stop(ctx context.Context) errorz.Error
	Ready(ctx context.Context) bool
}
