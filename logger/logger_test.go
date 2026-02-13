package logger

import (
	"context"
	"database/sql"
	"testing"

	"github.com/songlma/gobase/errorz"
)

var testCtx context.Context

func init() {
	testCtx = context.Background()
}

func TestName(t *testing.T) {
	ctx := context.Background()

	Info(ctx, "测试日志")
}

func TestName2(t *testing.T) {
	ctx := context.Background()
	InitLog(ctx,
		"gov2",
		func(ctx context.Context) string {
			return "1231"
		},
		Opt.Fmt(FmtJson),
	)
	Trace(ctx, "测试日志")
}

func TestName3(t *testing.T) {
	ctx := context.Background()
	InitLog(ctx,
		"gov2", func(ctx context.Context) string {
			return "1231"
		}, Opt.Fmt(FmtJson),
		Opt.File("/Users/songbaokang/Desktop/log.json"),
	)
	Info(ctx, "测试日志")
}

func TestName4(t *testing.T) {
	ctx := context.Background()
	InitLog(ctx,
		"gov2",
		func(ctx context.Context) string {
			return "1231"
		},
		Opt.Fmt(FmtText),
		Opt.Level(LevelDebug),
	)
	Info(ctx, "测试日志")
}

func TestName5(t *testing.T) {
	ctx := context.Background()
	InitLog(ctx,
		"gov2",
		func(ctx context.Context) string {
			return "1231"
		},
		Opt.Fmt(FmtText),
		Opt.Level(LevelDebug),
	)
	Trace(ctx, "测试日志")
}

func TestName6(t *testing.T) {
	ctx := context.Background()
	InitLog(ctx,
		"gov2",
		func(ctx context.Context) string {
			return "1231"
		},
		Opt.Fmt(FmtText),
		Opt.Level(LevelDebug),
	)
	Info(ctx, "测试日志")
}
func TestError(t *testing.T) {
	ctx := context.Background()
	InitLog(ctx,
		"gov2",
		func(ctx context.Context) string {
			return "1231"
		},
		//Opt.Fmt(FMT_TEXT),
		//Opt.Level(LevelDebug),
	)
	errz := errorz.FromStd(sql.ErrNoRows)
	Error(ctx, "error log", errz)
}
func TestErrorf(t *testing.T) {
	ctx := context.Background()
	InitLog(ctx,
		"gov2",
		func(ctx context.Context) string {
			return "1231"
		},
		Opt.Fmt(FmtJson),
		Opt.Level(LevelDebug),
		Opt.File("/Users/songbaokang/Desktop/log.json"),
		Opt.File("/Users/songbaokang/Desktop/err.log"),
	)
	errz := errorz.FromStd(sql.ErrNoRows)
	Errorf(ctx, "error log:%s", errz)
}
