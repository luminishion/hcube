package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/luminishion/hcube/draw"
)

const windowWidth = 480
const windowHeight = 480

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "HCube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	previousTime := glfw.GetTime()

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth) / windowHeight, 0.01, 50.0)
	view := mgl32.LookAtV(mgl32.Vec3{2, 2, 2}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	model := mgl32.Ident4()

	mvp := projection
	mvp = mvp.Mul4(view)
	mvp = mvp.Mul4(model)

	hcube := draw.NewHcube(4)
	defer hcube.Close()
	hcube.SetColor(mgl32.Vec3{1, 1, 1})
	hcube.SetMVP(mvp)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0, 0, 0, 1)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		hcube.Rotate4d(elapsed)
		hcube.Render(2.5)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
