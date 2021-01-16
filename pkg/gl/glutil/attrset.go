package glutil

import (
	"fmt"
	"strings"

	"github.com/PieterD/warp/pkg/gl"
)

type AttrSet struct {
	enabled map[string]struct{}
	buffers map[string][]string
	attrs   []asAttribute
	vao     gl.VertexArrayObject
}

type asAttribute struct {
	attrName        string
	attrType        gl.Type
	attrIndex       int
	bufferName      string
	instanceDivisor int
}

type AttrConfig struct {
	Name   string
	Type   gl.Type
	Buffer string
	// Padding is the amount of bytes between the end of this attribute, and the beginning of the next.
	// 0 means the next attribute begins at the next byte in memory, with no gap of size 0.
	// >0 means there will be a gap after the attribute with a number of bytes equal to Padding.
	Padding int
	// Instancedivisor sets the glVertexAttribDivisor, required for instanced rendering.
	// 0 means the attribute advances once per vertex, as normal.
	// >0 means the attribute advances once every InstanceDivisor instances.
	InstanceDivisor int
}

func NewAttrSet(cfgs ...AttrConfig) (*AttrSet, error) {
	as := &AttrSet{
		enabled: make(map[string]struct{}),
	}
	index := 0
	for i, cfg := range cfgs {
		if cfg.Name == "" {
			return nil, fmt.Errorf("empty name is not allowed: attribute %d", i)
		}
		if cfg.Type == 0 {
			return nil, fmt.Errorf("empty Type is not allowed: attribute %d", i)
		}
		if _, ok := as.enabled[cfg.Name]; ok {
			return nil, fmt.Errorf("non-unique attribute name: %s", cfg.Name)
		}
		as.enabled[cfg.Name] = struct{}{}
		as.attrs = append(as.attrs, asAttribute{
			attrName:        cfg.Name,
			attrType:        cfg.Type,
			attrIndex:       index,
			bufferName:      cfg.Buffer,
			instanceDivisor: cfg.InstanceDivisor,
		})
		as.buffers[cfg.Buffer] = append(as.buffers[cfg.Buffer], cfg.Name)
		//TODO: different advancement for very large attrs (4 for mat4 I think)
		index++
	}
	return as, nil
}

func (as *AttrSet) Enable(names ...string) {
	enabled := make(map[string]struct{})
	for _, name := range names {
		enabled[name] = struct{}{}
	}
	as.enabled = enabled
}

func (as *AttrSet) ShaderCode() string {
	buf := &strings.Builder{}
	for _, attr := range as.attrs {
		buf.WriteString(fmt.Sprintf("layout (location = %d) in %s %s;\n",
			attr.attrIndex, attr.attrType.GLSL(), attr.attrName))
	}
	return buf.String()
}

func (as *AttrSet) Buffers(buffers map[string]gl.BufferObject) {
	// Mark VAO for regeneration.
	panic("not implemented")
}

// Content sets the amount of content (number of vertices) present in the current buffers.
func (as *AttrSet) Content(vertices int) {
	panic("not implemented")
}

func (as *AttrSet) VAO(glx *gl.Context) (gl.VertexArrayObject, error) {
	// Generate VAO the first time it is requested.
	// Update it each next time it is requested,
	// even though it returns the same VAO value.

	// Return an error is not all required buffers are set.
	// A buffer is only required if it has currently enabled attrs.

	// Attr 0 must be enabled.
	panic("not implemented")
}
