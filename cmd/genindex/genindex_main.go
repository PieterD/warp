package main

import (
	"flag"
	"fmt"
	"github.com/PieterD/warp/pkg/bootstrap"
	"os"
)

func main() {
	cfg := bootstrap.IndexConfig{}
	fs := flag.NewFlagSet("genindex", flag.ContinueOnError)
	fs.BoolVar(&cfg.ManifestJson, "manifest", false, "Include manifest.json link")
	err := fs.Parse(os.Args[1:])
	switch err {
	case nil:
	case flag.ErrHelp:
		fs.PrintDefaults()
		os.Exit(1)
		return
	default:
		_, _ = fmt.Fprintf(os.Stderr, "error parsing command line arguments: %v", err)
		fs.PrintDefaults()
		os.Exit(1)
		return
	}

	fmt.Println(bootstrap.IndexHtml(cfg))
}
