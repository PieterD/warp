# RAW TODO

## vector texture using texelFetch

`vec4 color = texelFetch(someTexture, ivec2(pixelX, pixelY), 0);`

## Point sprites

Programmatic point size:
```
glEnable(GL_PROGRAM_POINT_SIZE)
void main()
{
    gl_Position = projection * view * model * vec4(aPos, 1.0);    
    gl_PointSize = gl_Position.z;    
}  
```

alternatively, fixed point size: `glPointSize()`

In fragment shader:
`gl_PointCoord` gives the location of the point.
GL_POINT_SPRITE_COORD_ORIGIN: GL_LOWER_LEFT, GL_UPPER_LEFT (default)

# TODO

- glVertexAttrib* to set fixed value for disabled attributes
- figure out Multi-Draw extensions and default fallback
- transform feedback: glTransformFeedbackVaryings
- multiple render targets
- separate driver.Buffer's AsUint16Array and AsFloat32Array, and expand.

# Notes

- Mind coordinates: Opengl to top right, texture coordinates to bottom right.

# Links

## Webgl

https://learnopengl.com/
https://webglfundamentals.org/
http://webglsamples.org/WebGL2Samples

Spec:

https://www.khronos.org/registry/webgl/specs/latest/1.0/
https://www.khronos.org/registry/OpenGL/specs/es/2.0/GLSL_ES_Specification_1.00.pdf

https://www.khronos.org/registry/webgl/specs/latest/2.0/
https://www.khronos.org/registry/OpenGL/specs/es/3.0/GLSL_ES_Specification_3.00.pdf

## Css

https://cssbattle.dev/
https://www.youtube.com/watch?v=ECsvqHoFZm8

## Models

Heart model by printable_models: https://free3d.com/3d-model/heart-v1--539992.html
