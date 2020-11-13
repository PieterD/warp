package mdl

import (
	"fmt"
	"os"
	"testing"
)

func TestFromObj(t *testing.T) {
	//h, err := os.Open("../../misc/static/models/12190_Heart_v1_L3.obj")
	//h, err := os.Open("../../misc/static/models/square.obj")
	h, err := os.Open("../../misc/static/models/cube.obj")
	if err != nil {
		t.Fatalf("opening object file: %v", err)
	}
	defer func() { _ = h.Close() }()

	model, err := FromObj(h)
	if err != nil {
		t.Fatalf("reading object: %v", err)
	}
	vs, is, err := model.Interleaved()
	_ = vs
	_ = is
	fmt.Printf("safe\n")
}
