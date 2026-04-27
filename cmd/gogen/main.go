package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/vladimirkorzhenevskiy/gogen/internal/adapter/fs"
	"github.com/vladimirkorzhenevskiy/gogen/internal/controller/cli"
	"github.com/vladimirkorzhenevskiy/gogen/internal/usecase/generate"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	defer cancel()

	if err := cli.New(generate.New(fs.New())).Run(ctx); err != nil {
		log.Fatal(err)
	}
}
