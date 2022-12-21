package draw

import (
	"log"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Hcube struct {
	edges             []edge
	vertices, project [][]float64
	line              *Line
}

type edge struct {
	v1, v2 int
}

func NewHcube(dimension int) *Hcube {
	if dimension < 3 {
		log.Fatalln("dimension must be >2")
	}

	verticesCount := int(math.Pow(2, float64(dimension)))
	vertices := make([][]float64, verticesCount)
	project := make([][]float64, verticesCount)

	for i := 0; i < verticesCount; i++ {
		v := make([]float64, dimension)

		for n := range v {
			v[n] = float64((i >> n) % 2) - 0.5
		}

		project[i] = make([]float64, dimension)
		vertices[i] = v
	}

	edgeCount := verticesCount / 2 * dimension
	edges := make([]edge, edgeCount)

	count := 0
	for j := 0; j < dimension; j++ {
		shift := int(math.Pow(2, float64(j)))

		for start := 0; start < verticesCount; start += shift * 2 {
			for i := 0; i < shift; i++ {
				idx := start + i
			
				edges[count] = edge{
					v1: idx,
					v2: idx + shift,
				}
				
				count++
			}
		}
	}

	l := NewLine()

	h := &Hcube{
		edges:    edges,
		vertices: vertices,
		line:     l,
		project:  project,
	}

	return h
}

func (h *Hcube) SetColor(color mgl32.Vec3) {
	h.line.SetColor(color)
}

func (h *Hcube) SetMVP(mvp mgl32.Mat4) {
	h.line.SetMVP(mvp)
}

func (h *Hcube) Rotate4d(delta float64) {
	for _, vert := range h.vertices {
		a := vert[2]
		b := vert[3]

		vert[2] = math.Cos(delta)*a - math.Sin(delta)*b
		vert[3] = math.Sin(delta)*a + math.Cos(delta)*b
	}
}

func (h *Hcube) Render(delta float64) {
	project := h.project
	line := h.line
	vertices := h.vertices

	for i, p := range project {
		copy(p, vertices[i])
	}

	for _, p := range project {
		for ind := len(p) - 1; ind > 2; ind-- {
			last := p[ind]

			for i := 0; i < ind; i++ {
				p[i] = delta * p[i] / (last + delta)
			}
		}
	}

	for _, edge := range h.edges {
		p1 := project[edge.v1]
		p2 := project[edge.v2]

		vec1 := mgl32.Vec3{
			float32(p1[0]),
			float32(p1[1]),
			float32(p1[2]),
		}

		vec2 := mgl32.Vec3{
			float32(p2[0]),
			float32(p2[1]),
			float32(p2[2]),
		}

		line.SetPos(vec1, vec2)
		line.Render()

	}
}

func (h *Hcube) Close() {
	h.line.Close()
}
