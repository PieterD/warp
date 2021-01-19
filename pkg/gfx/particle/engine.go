package particle

import "github.com/PieterD/warp/pkg/gl"

type Config struct {
	MaxParticles int
}

/*
vertex attributes:
 - lifecycle (0.0 - 1.0)
 - location vec3
 - momentum vec3
 - perturbation (0.0 - 1.0)

*/

type Engine struct {
}

func New(cfg Config) (*Engine, error) {
	return &Engine{}, nil
}

func (e *Engine) Destroy() {

}

// Advance advances the engine by the given lifecycle increment.
func (e *Engine) Advance(glx *gl.Context, lifeCycleIncrement float32) {

}

// Draw draws the engine in its current state to the given framebuffer, or the default one if nil.
func (e *Engine) Draw(glx *gl.Context, fb *gl.FramebufferObject) {

}
