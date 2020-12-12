package gl

import (
	"fmt"
	"testing"
)

func TestDataCoupling(t *testing.T) {
	allAttributeNames := []string{"Coordinates", "TexCoord", "Normal"}
	dc, err := NewDataCoupling(DataCouplingConfig{
		Vertices: []VertexConfig{
			{
				Name:    "Coordinates",
				Type:    Vec3,
				Buffer:  "vertex",
				Padding: 0,
			},
			{
				Name:    "TexCoord",
				Type:    Vec2,
				Buffer:  "vertex",
				Padding: 0,
			},
			{
				Name:    "Normal",
				Type:    Vec3,
				Buffer:  "vertex",
				Padding: 0,
			},
		},
	})
	if err != nil {
		t.Fatalf("creating data coupling: %v", err)
	}
	dc.ProgramConfig(allAttributeNames)
	dc.VertexArrayConfig(allAttributeNames, map[string]*Buffer{"vertex": {}})
	fmt.Printf("done %#v\n", dc)
}
