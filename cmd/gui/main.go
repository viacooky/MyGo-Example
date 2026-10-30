package main

import (
	"log/slog"
	"my-app/internal/app"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug, // 设置最低输出级别
	}))
	slog.SetDefault(logger)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
