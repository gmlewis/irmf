// -*- compile-command: "go run aprbfem.go -n 36"; -*-

// aprbfem creates an axial-plus-radial-bi-filar-electro-magnet
// STL file using Go. It is based on 30x30x132mm-vert.irmf.
//
// Usage:
//   go run aprbfem.go -out aprbfem.stl -n 120
package main

import (
	"flag"
	"log"
	"math"

	"github.com/gmlewis/go3d/vec3"
	"github.com/gmlewis/irmf-slicer/stl"
)

var (
	filename = flag.String("out", "aprbfem.stl", "Output filename")
	numDiv   = flag.Int("n", 120, "Number of divisions per rotation")
)

// TriWriter is a writer that accepts STL triangles.
type TriWriter interface {
	Write(t *stl.Tri) error
}

func main() {
	flag.Parse()

	w, err := stl.New(*filename)
	if err != nil {
		log.Fatalf("stl.New: %v", err)
	}

	m := &arBifilarElectromagnet{
		numPairs:    11,
		innerRadius: 3.0,
		size:        0.85,
		singleGap:   0.15,
		numTurns:    61,
		w:           w,
	}

	m.render()

	if err := w.Close(); err != nil {
		log.Fatalf("w.Close: %v", err)
	}

	log.Printf("Done.")
}

type arBifilarElectromagnet struct {
	// initializers:
	numPairs    int
	innerRadius float64
	size        float64
	singleGap   float64
	numTurns    int
	w           TriWriter

	// calculated:
	nTurns          float64
	inc             float64
	connectorRadius float64
	doubleGap       float64
}

func (m *arBifilarElectromagnet) render() {
	m.nTurns = float64(m.numTurns)
	m.inc = math.Pi / float64(m.numPairs)
	m.connectorRadius = m.innerRadius + float64(m.numPairs)*(m.size+m.singleGap)
	m.doubleGap = m.size + 2.0*m.singleGap

	for i := 1; i <= m.numPairs; i++ {
		m.coilPlusConnectorWires(1, i)
		m.coilPlusConnectorWires(2, i)
	}
}

func (m *arBifilarElectromagnet) coilPlusConnectorWires(wireNum, coilNum int) {
	radiusOffset := float64(coilNum - 1)
	// spacingAngle := float64(m.numPairs-4) * m.inc * math.Atan(2.0*radiusOffset/float64(m.numPairs-1))
	coilRadius := radiusOffset + m.innerRadius
	trimStartAngle := 0.05

	m.coilSquareFace(wireNum, coilRadius, trimStartAngle, 0.0)
}

func (m *arBifilarElectromagnet) coilSquareFace(wireNum int, radius, trimStartAngle, trimEndAngle float64) {
	delta := 2.0 * math.Pi / float64(*numDiv)
	angle := 0.0
	endAngle := float64(m.numTurns)*2.0*math.Pi + angle

	ri := radius - 0.5*m.size
	ro := radius + 0.5*m.size
	for ; angle <= endAngle-delta; angle += delta {
		m.coilWireSegment(wireNum, angle, angle+delta, ri, ro)
	}
}

func (m *arBifilarElectromagnet) coilWireSegment(wireNum int, a1, a2, ri, ro float64) {
	z1 := a1 / math.Pi
	z2 := a2 / math.Pi
	if wireNum%2 == 0 {
		a1 += math.Pi
		a2 += math.Pi
	}

	pu := func(r, a, z float64) *vec3.T {
		return &vec3.T{float32(r * math.Cos(a)), float32(r * math.Sin(a)), float32(z - 0.5*m.size)}
	}
	pd := func(r, a, z float64) *vec3.T {
		return &vec3.T{float32(r * math.Cos(a)), float32(r * math.Sin(a)), float32(z + 0.5*m.size)}
	}
	cp := func(v *vec3.T) *vec3.T { return &vec3.T{v[0], v[1], v[2]} }

	p1uo := pu(ro, a1, z1)
	// p1ui := pu(ri, a1, z1)
	p1do := pd(ro, a1, z1)
	// p1di := pd(ri, a1, z1)
	p2uo := pu(ro, a2, z2)
	// p2ui := pu(ri, a2, z2)
	p2do := pd(ro, a2, z2)
	// p2di := pd(ri, a2, z2)

	n1o := vec3.T{float32(math.Cos(a1)), float32(math.Sin(a1)), 0}
	n1u := vec3.Cross(pu(ro, a1, z1-1.0).Sub(p1uo), cp(p2uo).Sub(p1uo))
	n1u.Normalize()
	// n1i := vec3.T{float32(math.Cos(a1 + math.Pi)), float32(math.Sin(a1 + math.Pi)), 0}
	n1d := vec3.Cross(pd(ro, a1, z1+1.0).Sub(p1do), cp(p2do).Sub(p1do))
	n1d.Normalize()

	// n2o := vec3.T{float32(math.Cos(a2)), float32(math.Sin(a2)), 0}
	n2u := vec3.Cross(pu(ro, a2, z2-1.0).Sub(p2uo), cp(p2uo).Sub(p2uo))
	n2u.Normalize()
	// n2i := vec3.T{float32(math.Cos(a2 + math.Pi)), float32(math.Sin(a2 + math.Pi)), 0}
	n2d := vec3.Cross(pd(ro, a2, z2+1.0).Sub(p2do), cp(p2do).Sub(p2do))
	n2d.Normalize()

	write := func(n, v1, v2, v3 *vec3.T) {
		m.w.Write(&stl.Tri{
			N:  [3]float32{n[0], n[1], n[2]},
			V1: [3]float32{v1[0], v1[1], v1[2]},
			V2: [3]float32{v2[0], v2[1], v2[2]},
			V3: [3]float32{v3[0], v3[1], v3[2]},
		})
	}
	write(&n1o, p1uo, p2do, p1do)
	write(&n1o, p1uo, p2uo, p2do)
}
