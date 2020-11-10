package mdl

import "fmt"

type Model struct {
	Vertices     []Vertex
	VertexItems  int
	TextureItems int
	NormalItems  int
}

type Vertex struct {
	Vs [4]float32
	Ts [4]float32
	Ns [4]float32
}

func NewVertex(vertexItems, textureItems, normalItems []float32) (Vertex, error) {
	if len(vertexItems) > 4 {
		return Vertex{}, fmt.Errorf("more than 4 vertex items: %d", len(vertexItems))
	}
	if len(textureItems) > 4 {
		return Vertex{}, fmt.Errorf("more than 4 texture items: %d", len(textureItems))
	}
	if len(normalItems) > 4 {
		return Vertex{}, fmt.Errorf("more than 4 normal items: %d", len(normalItems))
	}
	var v Vertex
	copy(v.Vs[:], vertexItems)
	copy(v.Ts[:], textureItems)
	copy(v.Ns[:], normalItems)
	return v, nil
}

func (model *Model) Interleaved() (vs []float32, is []uint16, err error) {
	if len(model.Vertices) == 0 {
		return nil, nil, nil
	}
	is = make([]uint16, 0, len(model.Vertices))
	vs = make([]float32, 0, len(is))
	vertices := model.Vertices
	vertexLocation := 0
	vertexLocations := make(map[Vertex]int)
	for _, vertex := range vertices {
		if vertexLocation > 0xffff {
			return nil, nil, fmt.Errorf("index does not fit in uint16: %d", vertexLocation)
		}
		vIndex, ok := vertexLocations[vertex]
		if ok {
			is = append(is, uint16(vIndex))
			continue
		}
		vertexLocations[vertex] = vertexLocation

		vs = append(vs, vertex.Vs[:model.VertexItems]...)
		vs = append(vs, vertex.Ts[:model.TextureItems]...)
		vs = append(vs, vertex.Ns[:model.NormalItems]...)
		is = append(is, uint16(vertexLocation))
		vertexLocation++
	}
	return vs, is, nil
}
