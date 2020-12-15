// -*- compile-command: "go run aprbfem.go"; -*-

// aprbfem creates an axial-plus-radial-bi-filar-electro-magnet
// STL file using Go. It is based on 30x30x132mm-vert.irmf.
//
// Usage:
//   go run aprbfem.go -h
//   go run aprbfem.go -out aprbfem.stl
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
	innerR   = flag.Float64("inner_radius", 3.0, "Inner radius in millimeters")
	numDivs  = flag.Int("num_divs", 36, "Number of divisions per rotation")
	numPairs = flag.Int("num_pairs", 11, "Number of coil pairs")
	numTurns = flag.Int("num_turns", 61, "Total number of turns per coil")
	wireGap  = flag.Float64("wire_gap", 0.15, "Gap between wires in millimeters")
	wireSize = flag.Float64("wire_size", 0.85, "Width of (square) wire in millimeters")
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
		numPairs:    *numPairs,
		innerRadius: *innerR,
		size:        *wireSize,
		singleGap:   *wireGap,
		numTurns:    *numTurns,
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
	inc             float64
	connectorRadius float64
	doubleGap       float64
	height          float32

	lowerConnectors []*connector
	upperConnectors []*connector
}

type connector struct {
	inwardN        *vec3.T
	center         *vec3.T
	p1, p2, p3, p4 *vec3.T
}

func (m *arBifilarElectromagnet) render() {
	m.inc = math.Pi / float64(m.numPairs)
	m.connectorRadius = m.innerRadius + float64(m.numPairs)*(m.size+m.singleGap)
	m.doubleGap = m.size + 2.0*m.singleGap
	m.height = float32(m.numTurns*2 + 1)

	for i := 1; i <= m.numPairs; i++ {
		m.coilPlusConnectorWires(1, i)
		m.coilPlusConnectorWires(2, i)
	}
}

func (m *arBifilarElectromagnet) coilPlusConnectorWires(wireNum, coilNum int) {
	radiusOffset := float64(coilNum - 1)
	spacingAngle := float64(m.numPairs-4) * m.inc * math.Atan(2.0*radiusOffset/float64(m.numPairs-1))
	coilRadius := radiusOffset + m.innerRadius
	trimStartAngle := 0.05

	m.coilSquareFace(wireNum, coilRadius, trimStartAngle, 0.0, spacingAngle)

	m.coilConnectorWires(wireNum, coilNum, coilRadius, trimStartAngle, 0.0, spacingAngle)
}

func (m *arBifilarElectromagnet) coilConnectorWires(wireNum, coilNum int, coilRadius, trimStartAngle, trimEndAngle, spacingAngle float64) {
	if coilNum == m.numPairs {
		if wireNum == 1 {
			m.innerExitWire(wireNum, coilNum, coilRadius, trimStartAngle, 0.0, spacingAngle)
			return
		}
		m.outerExitWire(wireNum, coilNum, coilRadius, trimStartAngle, 0.0, spacingAngle)
		return
	}

	m.coilConnectorWire(wireNum, coilNum, coilRadius, trimStartAngle, 0.0, spacingAngle)
}

func (m *arBifilarElectromagnet) innerExitWire(wireNum, coilNum int, coilRadius, trimStartAngle, trimEndAngle, spacingAngle float64) {
}

func (m *arBifilarElectromagnet) outerExitWire(wireNum, coilNum int, coilRadius, trimStartAngle, trimEndAngle, spacingAngle float64) {
}

func (m *arBifilarElectromagnet) coilConnectorWire(wireNum, coilNum int, coilRadius, trimStartAngle, trimEndAngle, spacingAngle float64) {
}

func (m *arBifilarElectromagnet) coilSquareFace(wireNum int, radius, trimStartAngle, trimEndAngle, spacingAngle float64) {
	delta := 2.0 * math.Pi / float64(*numDivs)
	angle := 0.0
	endAngle := float64(m.numTurns)*2.0*math.Pi + angle

	ri := radius - 0.5*m.size
	ro := radius + 0.5*m.size
	firstFace := true
	for ; angle <= endAngle-delta; angle += delta {
		lastFace := (angle+delta > endAngle-delta)
		m.coilWireSegment(firstFace, lastFace, wireNum, angle+spacingAngle, angle+delta+spacingAngle, ri, ro)
		firstFace = false
	}
}

func (m *arBifilarElectromagnet) coilWireSegment(firstFace, lastFace bool, wireNum int, origA1, origA2, ri, ro float64) {
	a1, a2 := origA1, origA2
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
	p1ui := pu(ri, a1, z1)
	p1do := pd(ro, a1, z1)
	p1di := pd(ri, a1, z1)
	p2uo := pu(ro, a2, z2)
	p2ui := pu(ri, a2, z2)
	p2do := pd(ro, a2, z2)
	p2di := pd(ri, a2, z2)

	ma := 0.5 * (a1 + a2)

	no := vec3.T{float32(math.Cos(ma)), float32(math.Sin(ma)), 0}
	nu := vec3.Cross(cp(p2ui).Sub(p1uo), cp(p2uo).Sub(p1uo))
	nu.Normalize()
	ni := vec3.T{float32(math.Cos(ma + math.Pi)), float32(math.Sin(ma + math.Pi)), 0}
	nd := vec3.Cross(cp(p2do).Sub(p1do), cp(p2di).Sub(p1do))
	nd.Normalize()

	write := func(n, v1, v2, v3 *vec3.T) {
		m.w.Write(&stl.Tri{
			N:  [3]float32{n[0], n[1], n[2]},
			V1: [3]float32{v1[0], v1[1], v1[2]},
			V2: [3]float32{v2[0], v2[1], v2[2]},
			V3: [3]float32{v3[0], v3[1], v3[2]},
		})
	}
	quad := func(n, v1, v2, v3, v4 *vec3.T) {
		write(n, v1, v2, v3)
		write(n, v1, v3, v4)
	}

	if firstFace {
		da := m.size / ro
		a0 := a1 - da
		z0 := (origA1 - da) / math.Pi
		p0uo := pu(ro, a0, z0)
		p0ui := pu(ri, a0, z0)
		p0do := pd(ro, a0, z0)
		p0di := pd(ri, a0, z0)
		n := vec3.Cross(cp(p0di).Sub(p0do), cp(p0ui).Sub(p0do))
		n.Normalize()
		nb := cp(&n).Invert()

		quad(&n, p0do, p0di, p0ui, p0uo) // end-cap
		// quad(&no, p0uo, p1uo, p1do, p0do) // outer
		quad(&nu, p0uo, p0ui, p1ui, p1uo) // upward
		quad(&ni, p0ui, p0di, p1di, p1ui) // inner
		quad(&nd, p0do, p1do, p1di, p0di) // downward

		ma01 := 0.5 * (a1 + a0)
		ni01 := vec3.T{float32(math.Cos(ma01 + math.Pi)), float32(math.Sin(ma01 + math.Pi)), 0}

		vlen := m.connectorRadius + 0.5*m.size - ro
		outP0uo := cp(&ni01).Scale(-float32(vlen)).Add(p0uo)
		outP0ui := cp(&ni01).Scale(-float32(vlen - m.size)).Add(p0uo)
		outP0do := cp(&ni01).Scale(-float32(vlen)).Add(p0do)
		outP0di := cp(&ni01).Scale(-float32(vlen - m.size)).Add(p0do)
		outP1uo := cp(&ni01).Scale(-float32(vlen)).Add(p1uo)
		outP1ui := cp(&ni01).Scale(-float32(vlen - m.size)).Add(p1uo)
		outP1do := cp(&ni01).Scale(-float32(vlen)).Add(p1do)
		outP1di := cp(&ni01).Scale(-float32(vlen - m.size)).Add(p1do)

		quad(&n, outP0do, outP0di, outP0ui, outP0uo)  // end-cap
		quad(&n, outP0di, p0do, p0uo, outP0ui)        // end-cap connector
		quad(&no, outP0uo, outP1uo, outP1do, outP0do) // outer
		quad(&nu, outP0ui, p0uo, p1uo, outP1ui)       // upward connector
		quad(&nu, outP0uo, outP0ui, outP1ui, outP1uo) // upward
		quad(&nd, outP0di, outP1di, p1do, p0do)       // downward
		quad(nb, outP1do, outP1uo, outP1ui, outP1di)  // backface
		quad(nb, outP1di, outP1ui, p1uo, p1do)        // backface connector

		botP0uo := cp(&vec3.UnitZ).Scale(m.height).Add(outP0uo)
		botP0ui := cp(&vec3.UnitZ).Scale(m.height).Add(outP0ui)
		//	botP0do := cp(&vec3.UnitZ).Scale(m.height).Add(outP0do)
		//	botP0di := cp(&vec3.UnitZ).Scale(m.height).Add(outP0di)
		botP1uo := cp(&vec3.UnitZ).Scale(m.height).Add(outP1uo)
		botP1ui := cp(&vec3.UnitZ).Scale(m.height).Add(outP1ui)
		//	botP1do := cp(&vec3.UnitZ).Scale(m.height).Add(outP1do)
		//	botP1di := cp(&vec3.UnitZ).Scale(m.height).Add(outP1di)

		// axial connector
		quad(&n, botP0uo, botP0ui, outP0di, outP0do)    // forward (end-cap)
		quad(&no, outP0do, outP1do, botP1uo, botP0uo)   // outer
		quad(nb, botP1uo, outP1do, outP1di, botP1ui)    // backface
		quad(&ni01, botP0ui, botP1ui, outP1di, outP0di) // inner

		uc := &connector{
			inwardN: &ni01,
			center:  pu(0.5*(ro+ri), ma01, 0.5*(z0+z1-m.size)),
			p1:      p0uo,
			p2:      p0ui,
			p3:      p1ui,
			p4:      p1uo,
		}
		m.upperConnectors = append(m.upperConnectors, uc)
	}

	// outer-facing
	quad(&no, p1uo, p2uo, p2do, p1do)
	// upward-facing
	quad(&nu, p1uo, p1ui, p2ui, p2uo)
	// inner-facing
	quad(&ni, p1ui, p1di, p2di, p2ui)
	// downward-facing
	quad(&nd, p1do, p2do, p2di, p1di)
	if lastFace {
		n := vec3.Cross(cp(p2ui).Sub(p2do), cp(p2di).Sub(p2do))
		n.Normalize()
		edge := cp(&n).Scale(float32(m.size))
		p3uo := cp(p2uo).Add(edge)
		p3ui := cp(p2ui).Add(edge)
		p3do := cp(p2do).Add(edge)
		p3di := cp(p2di).Add(edge)
		// quad(&n, p2do, p2uo, p2ui, p2di)
		quad(&n, p3di, p3do, p3uo, p3ui) // end-cap
		quad(&n, p2uo, p3uo, p3do, p2do) // outer
		quad(&n, p2ui, p2di, p3di, p3ui) // inner
		quad(&n, p3ui, p3uo, p2uo, p2ui) // upward
		lc := &connector{
			inwardN: &ni,
			center:  cp(p3uo).Add(p2ui).Scale(0.5),
			p1:      p3ui,
			p2:      p3uo,
			p3:      p2uo,
			p4:      p2ui,
		}
		m.lowerConnectors = append(m.lowerConnectors, lc)
	}
}
