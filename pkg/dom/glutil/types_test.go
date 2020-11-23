package glutil

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"testing"
)

func TestBuild(t *testing.T) {
	type Struct struct {
		Stuff int32
	}
	type testType struct {
		F  float32
		V3 mgl32.Vec3
		M4 mgl32.Mat4
		S  Struct
	}
	tv := &testType{
		F:  1.0,
		V3: mgl32.Vec3{1.0, 2.0, 3.0},
		M4: mgl32.Ident4(),
		S: Struct{
			Stuff: 1,
		},
	}
	built, err := build(tv)
	if err != nil {
		t.Fatalf("building: %v", err)
	}
	size := built.Size()
	if want := 4+4*4+4*4*4+4; want != size {
		t.Fatalf("invalid size, want %d got %d", want, size)
	}
	fmt.Printf("safe!\n")
}
