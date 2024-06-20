package main

import (
	"log/slog"

	"github.com/karamaru-alpha/mocka/internal"
)

func main() {
	if err := internal.Run(); err != nil {
		slog.Error(err.Error())
	}
}
