package mdl

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func FromObj(r io.Reader) (*Model, error) {
	p := &objParser{}
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		if err := p.parseLine(line); err != nil {
			return nil, fmt.Errorf("parsing line %d: %w", lineNum, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file: %w", err)
	}
	model, err := p.toModel()
	if err != nil {
		return nil, fmt.Errorf("converting obj to Model: %w", err)
	}
	return model, nil
}

type objParser struct {
	itemsPerVertex  int
	vertices        []float32
	itemsPerTexture int
	textures        []float32
	itemsPerNormal  int
	normals         []float32
	itemsPerFace    int
	faces           [][3]int
}

func (p *objParser) parseLine(line string) error {
	if strings.HasPrefix(line, "v ") {
		err := p.parseV(strings.TrimSpace(line[2:]))
		if err != nil {
			return fmt.Errorf("parsing vertex : %w", err)
		}
		return nil
	}
	if strings.HasPrefix(line, "vn ") {
		err := p.parseVN(strings.TrimSpace(line[3:]))
		if err != nil {
			return fmt.Errorf("parsing normal: %w", err)
		}
		return nil
	}
	if strings.HasPrefix(line, "vt ") {
		err := p.parseVT(strings.TrimSpace(line[3:]))
		if err != nil {
			return fmt.Errorf("parsing texture coord: %w", err)
		}
		return nil
	}
	if strings.HasPrefix(line, "f ") {
		err := p.parseF(strings.TrimSpace(line[2:]))
		if err != nil {
			return fmt.Errorf("parsing texture coord: %w", err)
		}
		return nil
	}
	return nil
}

func (p *objParser) parseV(raw string) (err error) {
	vertex, err := parseVertexItems(raw)
	if err != nil {
		return fmt.Errorf("parsing vertex items: %w", err)
	}
	if p.itemsPerVertex == 0 {
		p.itemsPerVertex = len(vertex)
	} else if p.itemsPerVertex != len(vertex) {
		return fmt.Errorf("expected %d items per vertex, but found one with %d", p.itemsPerVertex, len(vertex))
	}
	p.vertices = append(p.vertices, vertex[:]...)
	return nil
}

func (p *objParser) parseVN(raw string) (err error) {
	normal, err := parseVertexItems(raw)
	if err != nil {
		return fmt.Errorf("parsing normal items: %w", err)
	}
	if p.itemsPerNormal == 0 {
		p.itemsPerNormal = len(normal)
	} else if p.itemsPerNormal != len(normal) {
		return fmt.Errorf("expected %d items per normal, but found one with %d", p.itemsPerNormal, len(normal))
	}
	p.normals = append(p.normals, normal[:]...)
	return nil
}

func (p *objParser) parseVT(raw string) (err error) {
	texture, err := parseVertexItems(raw)
	if err != nil {
		return fmt.Errorf("parsing texture items: %w", err)
	}
	if p.itemsPerTexture == 0 {
		p.itemsPerTexture = len(texture)
	} else if p.itemsPerTexture != len(texture) {
		return fmt.Errorf("expected %d items per texture, but found one with %d", p.itemsPerTexture, len(texture))
	}
	p.textures = append(p.textures, texture[:]...)
	return nil
}

func (p *objParser) parseF(raw string) (err error) {
	rawItems := strings.Split(raw, " ")
	if len(rawItems) == 0 {
		return fmt.Errorf("no index componenets")
	}
	if p.itemsPerFace == 0 {
		p.itemsPerFace = len(rawItems)
	} else if p.itemsPerFace != len(rawItems) {
		return fmt.Errorf("expected %d items per face, but found one with %d", p.itemsPerFace, len(rawItems))
	}
	for _, rawItem := range rawItems {
		var indices []int
		for _, rawIndex := range strings.Split(rawItem, "/") {
			if rawIndex == "" {
				indices = append(indices, -1)
				continue
			}
			i64, err := strconv.ParseInt(rawIndex, 10, 16)
			if err != nil {
				return fmt.Errorf("parsing raw index: %w", err)
			}
			indices = append(indices, int(i64))
		}
		if len(indices) > 3 {
			return fmt.Errorf("too many indices: %d", len(indices))
		}
		for len(indices) < 3 {
			indices = append(indices, -1)
		}
		var f [3]int
		copy(f[:], indices)
		p.faces = append(p.faces, f)
	}
	return nil
}

func parseVertexItems(raw string) (vertex []float32, err error) {
	split := strings.Split(raw, " ")
	if len(split) == 0 {
		return nil, fmt.Errorf("no vertex components")
	}
	vertex = make([]float32, len(split))
	for i := 0; i < len(split); i++ {
		f64, err := strconv.ParseFloat(split[i], 32)
		if err != nil {
			return nil, fmt.Errorf("parsing float %d: %w", i, err)
		}
		vertex[i] = float32(f64)
	}
	return vertex, nil
}

func (p *objParser) toModel() (*Model, error) {
	//fmt.Printf("vertices[%d]:%d, normals[%d]:%d, textures[%d]:%d faces[%d]:%d\n", p.itemsPerVertex, len(p.vertices), p.itemsPerNormal, len(p.normals), p.itemsPerTexture, len(p.textures), p.itemsPerFace, len(p.faces))
	var vertices []Vertex
	for faceIndex := 0; faceIndex < len(p.faces); faceIndex += p.itemsPerFace {
		face := p.faces[faceIndex : faceIndex+p.itemsPerFace]
		for vertexIndex, faceIndexGroup := range face {
			vStart := (faceIndexGroup[0] - 1) * p.itemsPerVertex
			vEnd := vStart + p.itemsPerVertex
			vertexItems := p.vertices[vStart:vEnd]

			tStart := (faceIndexGroup[1] - 1) * p.itemsPerTexture
			tEnd := tStart + p.itemsPerTexture
			textureItems := p.textures[tStart:tEnd]

			nStart := (faceIndexGroup[2] - 1) * p.itemsPerNormal
			nEnd := nStart + p.itemsPerNormal
			normalItems := p.normals[nStart:nEnd]

			vertex, err := NewVertex(vertexItems, textureItems, normalItems)
			if err != nil {
				return nil, fmt.Errorf("vertex %d:%d: creating new model vertex: %w", faceIndex, vertexIndex, err)
			}
			vertices = append(vertices, vertex)
		}
	}
	switch p.itemsPerFace {
	case 3:
	case 4:
		var err error
		vertices, err = quadsToTriangles(vertices)
		if err != nil {
			return nil, fmt.Errorf("transforming quads to triangles: %w", err)
		}
	default:
		return nil, fmt.Errorf("unhandled amount of items per face: %d", p.itemsPerFace)
	}
	return &Model{
		Vertices:     vertices,
		VertexItems:  p.itemsPerVertex,
		TextureItems: p.itemsPerTexture,
		NormalItems:  p.itemsPerNormal,
	}, nil
}

func quadsToTriangles(quadVertices []Vertex) (triVertices []Vertex, err error) {
	if len(quadVertices)%4 != 0 {
		return nil, fmt.Errorf("number of quad vertices not divisible by 4")
	}
	for i := 0; i < len(quadVertices); i += 4 {
		vertices := quadVertices[i : i+4]
		triVertices = append(triVertices,
			vertices[0],
			vertices[1],
			vertices[3],
			vertices[3],
			vertices[1],
			vertices[2],
		)
	}
	return triVertices, nil
}
