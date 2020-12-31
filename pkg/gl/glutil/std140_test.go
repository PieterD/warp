package glutil

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"testing"
)

func TestStd140(t *testing.T) {
	type UniformBlock struct {
		F    float32
		S    int `std140:"sampler2D"`
		Vec2 mgl32.Vec2
		Vec3 mgl32.Vec3
		Vec4 mgl32.Vec4
		Mat4 mgl32.Mat4
	}
	ub := UniformBlock{
		F:    1.0,
		S:    5,
		Vec2: mgl32.Vec2{1, 2},
		Vec3: mgl32.Vec3{11, 22, 33},
		Vec4: mgl32.Vec4{111, 222, 333, 444},
		Mat4: mgl32.Ident4(),
	}
	d, err := Std140Data(&ub)
	if err != nil {
		t.Fatalf("Std140Data failed: %v", err)
	}
	_ = d
	s, err := Std140Uniform(&ub)
	if err != nil {
		t.Fatalf("Std140Uniform failed: %v", err)
	}
	fmt.Printf("%s\n", s)
}
