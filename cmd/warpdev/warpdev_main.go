package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/PieterD/warp/pkg/bootstrap"
	"golang.org/x/sync/errgroup"
)

const (
	addr = "localhost:8080"
)

func main() {
	//TODO:
	// Command to create docker directory from package name.
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	goRoot := os.Getenv("GOROOT")
	if goRoot == "" {
		goRoot = `C:\dev\go1.15.2`
	}

	cfg := bootstrap.Config{
		MainPackage: "github.com/PieterD/warp/cmd/gltest",
		StaticPath:  "cmd/gltest/static",
		GoRoot:      goRoot,
	}
	bootstrapHandler := bootstrap.New(cfg)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		srv := &http.Server{
			Handler:      bootstrapHandler,
			Addr:         addr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		if err := srv.ListenAndServe(); err != nil {
			return fmt.Errorf("running http server: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("waiting for group: %w", err)
	}
	return nil
}
