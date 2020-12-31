package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
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
	goRoot := os.Getenv("GOROOT")
	if goRoot == "" {
		goRoot = `C:\dev\go1.15.2`
	}

	eg.Go(func() error {
		r := mux.NewRouter()
		addRoot := func(name string) {
			r.PathPrefix(fmt.Sprintf("/%s/", name)).Handler(http.StripPrefix(fmt.Sprintf("/%s", name), bootstrap.New(bootstrap.Config{
				MainPackage: fmt.Sprintf("github.com/PieterD/warp/app/%s", name),
				StaticPath:  fmt.Sprintf("app/%s/static", name),
				GoRoot:      goRoot,
			})))
		}
		addRoot("gltest")
		addRoot("particle")
		srv := &http.Server{
			Handler:      r,
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
