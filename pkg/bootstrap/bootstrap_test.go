package bootstrap

import (
	"fmt"
	"testing"
)

func TestBinaryHandler(t *testing.T) {
	bh := binaryHandler{
		mainPackage: "github.com/PieterD/warp/cmd/gltest",
		goRoot:      `C:\dev\go1.15.2`,
	}
	data, err := bh.build()
	if err != nil {
		t.Fatalf("building binary: %v", err)
	}
	fmt.Printf("%d\n", len(data))
}
