package bootstrap

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

type Config struct {
	MainPackage string
	StaticPath  string
	GoRoot      string
}

func New(cfg Config) http.Handler {
	r := mux.NewRouter()
	bh := &binaryHandler{
		mainPackage: cfg.MainPackage,
		goRoot:      cfg.GoRoot,
	}
	r.Path("/").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		if _, err := io.Copy(writer, strings.NewReader(indexHtml)); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error copying index.html: %v", err)
		}
	})
	r.Path("/_binary").Handler(bh)
	r.Path("/_wasm_exec.js").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/javascript")
		wasmExecPath := filepath.Join(bh.goRoot, "misc", "wasm", "wasm_exec.js")
		wasmExecFile, err := os.Open(wasmExecPath)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "opening wasm_exec.js: %v", err)
			http.Error(writer, fmt.Sprintf("opening wasm_exec.js: %v", err), http.StatusBadGateway)
			return
		}
		defer func() { _ = wasmExecFile.Close() }()
		if _, err := io.Copy(writer, wasmExecFile); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error copying binary bytes: %v", err)
			return
		}
	})
	r.PathPrefix("/").Handler(http.StripPrefix("", http.FileServer(http.Dir(cfg.StaticPath))))
	return r
}

type binaryHandler struct {
	mainPackage string
	goRoot      string
}

func (bh *binaryHandler) build() ([]byte, error) {
	tempFile, err := ioutil.TempFile(os.TempDir(), "warp-bootstrap-*.wasm")
	if err != nil {
		return nil, fmt.Errorf("opening tempfile: %w", err)
	}
	defer func() { _ = os.Remove(tempFile.Name()) }()
	if err := tempFile.Close(); err != nil {
		return nil, fmt.Errorf("closing tempfile: %w", err)
	}
	cmd := exec.Command(filepath.Join(bh.goRoot, "bin", "go"), "build", "-o", tempFile.Name(), bh.mainPackage)
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	out, err := cmd.CombinedOutput()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", out)
		return nil, fmt.Errorf("running build command (%v): %w", cmd, err)
	}
	data, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("reading binary: %w", err)
	}
	return data, nil
}

func (bh *binaryHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	data, err := bh.build()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "building binary: %v\n", err)
		http.Error(writer, fmt.Sprintf("building binary: %v", err), http.StatusBadGateway)
		return
	}
	writer.Header().Set("Content-Type", "application/wasm")
	if _, err := io.Copy(writer, bytes.NewBuffer(data)); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error copying binary bytes: %w\n", err)
		return
	}
}
