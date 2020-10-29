package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "please provide one file name.\n")
		os.Exit(1)
	}
	fileName := os.Args[1]
	if err := run(fileName); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "running hashfile command: %v\n", err)
		os.Exit(1)
	}
}

func run(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer func() { _ = f.Close() }()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("hashing file: %w", err)
	}
	sum := h.Sum(nil)
	fmt.Printf("%02x\n", sum)
	return nil
}
