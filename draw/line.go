package draw

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const vertexShaderSource = `
#version 330 core

layout (location = 0) in vec3 aPos;

uniform mat4 MVP;

void main()
{
   gl_Position = MVP * vec4(aPos.x, aPos.y, aPos.z, 1.0);
}
` + "\x00"

const fragmentShaderSource = `
#version 330 core

out vec4 FragColor;
uniform vec3 color;

void main()
{
   FragColor = vec4(color, 1.0f);
}
` + "\x00"

type Line struct {
	program, vao, vbo uint32
	vertices          []float32
}

func NewLine() *Line {
	program, err := NewProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		log.Fatalln(err)
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	vertices := make([]float32, 6)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	l := &Line{
		program:  program,
		vao:      vao,
		vbo:      vbo,
		vertices: vertices,
	}

	return l
}

func (l *Line) SetColor(color mgl32.Vec3) {
	gl.UseProgram(l.program)

	gl.Uniform3fv(gl.GetUniformLocation(l.program, gl.Str("color\x00")), 1, &color[0])
}

func (l *Line) SetPos(start, end mgl32.Vec3) {
	vertices := l.vertices

	vertices[0] = start[0]
	vertices[1] = start[1]
	vertices[2] = start[2]
	vertices[3] = end[0]
	vertices[4] = end[1]
	vertices[5] = end[2]

	gl.BindBuffer(gl.ARRAY_BUFFER, l.vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices))
}

func (l *Line) SetMVP(mvp mgl32.Mat4) {
	gl.UseProgram(l.program)

	gl.UniformMatrix4fv(gl.GetUniformLocation(l.program, gl.Str("MVP\x00")), 1, false, &mvp[0])
}

func (l *Line) Render() {
	gl.UseProgram(l.program)

	gl.BindVertexArray(l.vao)
	gl.DrawArrays(gl.LINES, 0, 2)
}

func (l *Line) Close() {
	gl.DeleteVertexArrays(1, &l.vao)
	gl.DeleteBuffers(1, &l.vbo)
	gl.DeleteProgram(l.program)
}
