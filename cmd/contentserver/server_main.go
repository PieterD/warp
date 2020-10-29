package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

func run() error {
	var (
		addr      string
		staticDir string
		certFile  string
		keyFile   string
	)
	flag.StringVar(&addr, "addr", "127.0.0.1:12345", "Addr to start webserver on")
	flag.StringVar(&staticDir, "static", "./static", "Path to static files directory.")
	flag.StringVar(&certFile, "cert", "", "SSL Certificate.")
	flag.StringVar(&keyFile, "key", "", "TLS private key.")
	flag.Parse()

	if (certFile == "") != (keyFile == "") {
		_, _ = fmt.Fprintf(os.Stderr, "requires either both -cert and -key, or neither.")
		flag.PrintDefaults()
	}
	ssl := false
	if certFile != "" && keyFile != "" {
		ssl = true
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		r := mux.NewRouter()
		r.PathPrefix("/").Handler(http.StripPrefix("", http.FileServer(http.Dir(staticDir))))
		srv := &http.Server{
			Handler:      r,
			Addr:         addr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		if !ssl {
			if err := srv.ListenAndServe(); err != nil {
				return fmt.Errorf("running http server: %w", err)
			}
			return nil
		}
		if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil {
			return fmt.Errorf("running https server: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("waiting for group: %w", err)
	}
	return nil
}
