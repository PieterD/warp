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
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		h := bootstrap.New(bootstrap.Config{
			MainPackage: "github.com/PieterD/warp/cmd/gltest",
			StaticPath:  "cmd/gltest/static",
			GoRoot:      `C:\dev\go1.15.2`, //TODO: take from env if empty
		})
		srv := &http.Server{
			Handler:      h,
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
