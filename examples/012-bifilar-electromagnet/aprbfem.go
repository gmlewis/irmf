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
	leadLen  = flag.Float64("lead_len", 5.0, "Length of two external leads")
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
		leadLen:     *leadLen,
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
	leadLen     float64
	size        float64
	singleGap   float64
	numTurns    int
	w           TriWriter

	// calculated:
	inc             float64
	connectorRadius float64
	doubleGap       float64
	height          float32

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
	m.height = float32((m.size + m.singleGap) * float64((m.numTurns*2 + 1)))

	for i := 1; i <= m.numPairs; i++ {
		m.coilPlusConnectorWires(1, i)
		m.coilPlusConnectorWires(2, i)
	}
}

func (m *arBifilarElectromagnet) radiusOffset(coilNum int) float64 {
	return (m.size + m.singleGap) * float64(coilNum-1)
}

func (m *arBifilarElectromagnet) spacingAngle(coilNum int) float64 {
	return float64(m.numPairs-4) * m.inc * math.Atan(2.0*m.radiusOffset(coilNum)/float64(m.numPairs-1))
}

func (m *arBifilarElectromagnet) coilRadius(coilNum int) float64 {
	return m.radiusOffset(coilNum) + m.innerRadius
}

func (m *arBifilarElectromagnet) endAngle(wireNum, coilNum int) float64 {
	radius := m.coilRadius(coilNum)
	spacingAngle := m.spacingAngle(coilNum)
	endAngle := float64(m.numTurns) * 2.0 * math.Pi

	nextCoilNum := coilNum + 1
	nextSpacingAngle := m.spacingAngle(nextCoilNum)
	if nextCoilNum > *numPairs {
		nextSpacingAngle = m.spacingAngle(1) + 2.0*math.Pi
	}

	// Account for the edge of the wire connector
	ro := radius + 0.5*m.size
	angleEnd := m.size / ro

	result := endAngle + nextSpacingAngle - math.Pi - spacingAngle - angleEnd
	return result
}

func (m *arBifilarElectromagnet) coilPlusConnectorWires(wireNum, coilNum int) {
	radius := m.coilRadius(coilNum)
	spacingAngle := m.spacingAngle(coilNum)

	ri := radius - 0.5*m.size
	ro := radius + 0.5*m.size

	angle := 0.5 * m.size / ro // Start at the edge of the wire connector
	endAngle := m.endAngle(wireNum, coilNum) + angle
	delta := (endAngle - angle) / float64(*numDivs**numTurns)

	// The first segment and the last segment are special cases because they connect
	// up to the wire segments that pair up the coils in the correct sequence.
	m.coilWireSegment(true, false, wireNum, coilNum, spacingAngle, angle+spacingAngle, ri, ro)

	for i := 0; i < *numDivs**numTurns; i++ {
		lastFace := i == *numDivs**numTurns-1
		m.coilWireSegment(false, lastFace, wireNum, coilNum, angle+spacingAngle, angle+delta+spacingAngle, ri, ro)
		angle += delta
	}
}

func (m *arBifilarElectromagnet) coilWireSegment(firstFace, lastFace bool, wireNum, coilNum int, origA1, origA2, ri, ro float64) {
	a1, a2 := origA1, origA2
	z1 := (m.size + m.singleGap) * a1 / math.Pi
	z2 := (m.size + m.singleGap) * a2 / math.Pi
	if wireNum == 2 {
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

	nu := vec3.Cross(cp(p2ui).Sub(p1uo), cp(p2uo).Sub(p1uo))
	nu.Normalize()
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
	quad := func(v1, v2, v3, v4 *vec3.T) *vec3.T {
		v31 := cp(v3).Sub(v1)
		n1 := vec3.Cross(cp(v2).Sub(v1), v31)
		n1.Normalize()
		write(&n1, v1, v2, v3)
		n2 := vec3.Cross(v31, cp(v4).Sub(v1))
		n2.Normalize()
		write(&n2, v1, v3, v4)
		return &n2
	}

	if firstFace {
		// a1 is the center of the connector
		da := m.size / ro
		a0 := a1 - 0.5*da
		adja1 := a1 + 0.5*da
		z0 := (m.size + m.singleGap) * (origA1 - 0.5*da) / math.Pi
		adjz1 := (m.size + m.singleGap) * (origA1 + 0.5*da) / math.Pi

		p0uo := pu(ro, a0, z0)
		p0ui := pu(ri, a0, z0)
		p0do := pd(ro, a0, z0)
		p0di := pd(ri, a0, z0)

		adjP1uo := pu(ro, adja1, adjz1)
		adjP1ui := pu(ri, adja1, adjz1)
		adjP1do := pd(ro, adja1, adjz1)
		adjP1di := pd(ri, adja1, adjz1)

		n := vec3.Cross(cp(p0di).Sub(p0do), cp(p0ui).Sub(p0do))
		n.Normalize()

		quad(p0do, p0di, p0ui, p0uo)       // end-cap
		quad(p0uo, p0ui, adjP1ui, adjP1uo) // upward
		quad(p0ui, p0di, adjP1di, adjP1ui) // inner
		quad(p0do, adjP1do, adjP1di, p0di) // downward

		ni01 := vec3.T{float32(math.Cos(a1 + math.Pi)), float32(math.Sin(a1 + math.Pi)), 0}

		vlen := m.connectorRadius + 0.5*m.size - ro

		outP0uo := cp(&ni01).Scale(-float32(vlen)).Add(p0uo)
		outP0ui := cp(&ni01).Scale(-float32(vlen - m.size)).Add(p0uo)
		outP0do := cp(&ni01).Scale(-float32(vlen)).Add(p0do)
		outP0di := cp(&ni01).Scale(-float32(vlen - m.size)).Add(p0do)
		outP1uo := cp(&ni01).Scale(-float32(vlen)).Add(adjP1uo)
		outP1ui := cp(&ni01).Scale(-float32(vlen - m.size)).Add(adjP1uo)
		outP1do := cp(&ni01).Scale(-float32(vlen)).Add(adjP1do)
		outP1di := cp(&ni01).Scale(-float32(vlen - m.size)).Add(adjP1do)

		zu := 0.5 * (outP0ui[2] + outP1ui[2])
		zd := 0.5 * (outP0di[2] + outP1di[2])
		outP0uo[2] = zu
		outP0ui[2] = zu
		outP0do[2] = zd
		outP0di[2] = zd
		outP1uo[2] = zu
		outP1ui[2] = zu
		outP1do[2] = zd
		outP1di[2] = zd

		quad(outP0do, outP0di, outP0ui, outP0uo) // end-cap
		quad(outP0di, p0do, p0uo, outP0ui)       // end-cap connector
		quad(outP0uo, outP1uo, outP1do, outP0do) // outer
		quad(outP0ui, p0uo, adjP1uo, outP1ui)    // upward connector
		quad(outP0uo, outP0ui, outP1ui, outP1uo) // upward
		quad(outP0di, outP1di, adjP1do, p0do)    // downward
		quad(outP1do, outP1uo, outP1ui, outP1di) // backface
		quad(outP1di, outP1ui, adjP1uo, adjP1do) // backface connector

		h := m.height
		if coilNum == 1 && wireNum == 1 {
			h += float32(*leadLen + m.size)
		}

		botP0uo := cp(&vec3.UnitZ).Scale(h).Add(outP0uo)
		botP0ui := cp(&vec3.UnitZ).Scale(h).Add(outP0ui)
		botP0do := cp(&vec3.UnitZ).Scale(h).Add(outP0do)
		botP0di := cp(&vec3.UnitZ).Scale(h).Add(outP0di)
		botP1uo := cp(&vec3.UnitZ).Scale(h).Add(outP1uo)
		botP1ui := cp(&vec3.UnitZ).Scale(h).Add(outP1ui)
		botP1do := cp(&vec3.UnitZ).Scale(h).Add(outP1do)
		botP1di := cp(&vec3.UnitZ).Scale(h).Add(outP1di)

		// axial connector
		quad(botP0uo, botP0ui, outP0di, outP0do) // forward (was end-cap)
		quad(outP0do, outP1do, botP1uo, botP0uo) // outer
		quad(botP1uo, outP1do, outP1di, botP1ui) // backface
		quad(botP0ui, botP1ui, outP1di, outP0di) // inner

		if coilNum == 1 && wireNum == 1 {
			quad(botP1uo, botP1ui, botP0ui, botP0uo) // end cap of exit wire
		} else {
			// angle connector
			quad(botP0do, botP0di, botP0ui, botP0uo) // forward (was end-cap)
			quad(botP1do, botP0do, botP0uo, botP1uo) // outer
			quad(botP1do, botP1uo, botP1ui, botP1di) // backface
			quad(botP1do, botP1di, botP0di, botP0do) // end cap

			// connector to inner start of next coil
			nextRing := coilNum - 1
			if nextRing < 1 {
				nextRing = *numPairs
			}
			nextRadius := m.coilRadius(nextRing)

			conlen := m.connectorRadius - nextRadius
			conP0uo := cp(&ni01).Scale(float32(conlen)).Add(outP0uo)
			conP0do := cp(&ni01).Scale(float32(conlen)).Add(outP0do)
			conP1uo := cp(&ni01).Scale(float32(conlen)).Add(outP1uo)
			conP1do := cp(&ni01).Scale(float32(conlen)).Add(outP1do)
			conP0uo[2] = botP0uo[2]
			conP0do[2] = botP0do[2]
			conP1uo[2] = botP1uo[2]
			conP1do[2] = botP1do[2]

			// radial connector
			quad(botP0do, botP0di, botP0ui, botP0uo) // end-cap
			quad(botP0di, conP0do, conP0uo, botP0ui) // end-cap connector
			quad(botP0uo, botP1uo, botP1do, botP0do) // outer
			quad(botP0ui, conP0uo, conP1uo, botP1ui) // upward connector
			quad(botP0uo, botP0ui, botP1ui, botP1uo) // upward
			quad(botP0di, botP1di, conP1do, conP0do) // downward
			quad(botP1do, botP1uo, botP1ui, botP1di) // backface
			quad(botP1di, botP1ui, conP1uo, conP1do) // backface connector

			uc := &connector{
				inwardN: &ni01,
				center:  pu(0.5*(ro+ri), a1, 0.5*(z0+z1-m.size)),
				p1:      conP0uo,
				p2:      conP0do,
				p3:      conP1uo,
				p4:      conP1do,
			}
			m.upperConnectors = append(m.upperConnectors, uc)
		}
		return
	}

	quad(p1uo, p2uo, p2do, p1do) // outer-facing
	quad(p1uo, p1ui, p2ui, p2uo) // upward-facing
	quad(p1ui, p1di, p2di, p2ui) // inner-facing
	quad(p1do, p2do, p2di, p1di) // downward-facing

	if lastFace {
		// a2 is the end of the spiral
		da := m.size / ro
		a3 := a2 + da
		z3 := (m.size + m.singleGap) * (origA2 + da) / math.Pi

		p3uo := pu(ro, a3, z3)
		p3ui := pu(ri, a3, z3)
		p3do := pd(ro, a3, z3)
		p3di := pd(ri, a3, z3)

		if coilNum == *numPairs && wireNum == 2 { // exit wire
			quad(p3di, p3do, p3uo, p3ui) // end-cap
			quad(p2uo, p3uo, p3do, p2do) // outer
			quad(p2ui, p2di, p3di, p3ui) // inner
			quad(p3ui, p3uo, p2uo, p2ui) // upward

			h := float32(*leadLen)
			botP3uo := cp(&vec3.UnitZ).Scale(h).Add(p3uo)
			botP3ui := cp(&vec3.UnitZ).Scale(h).Add(p3ui)
			botP2uo := cp(&vec3.UnitZ).Scale(h).Add(p2uo)
			botP2ui := cp(&vec3.UnitZ).Scale(h).Add(p2ui)

			quad(botP3ui, botP3uo, p3do, p3di) // forward
			quad(botP3uo, botP2uo, p2do, p3do) // outer
			quad(botP2uo, botP2ui, p2di, p2do) // backface
			quad(botP2ui, botP3ui, p3di, p2di) // inner

			quad(botP3uo, botP3ui, botP2ui, botP2uo) // end cap of exit wire
			return
		}

		if coilNum == *numPairs && wireNum == 1 { // special loop-back case
			quad(p3di, p3do, p3uo, p3ui) // end-cap
			quad(p2ui, p2di, p3di, p3ui) // inner
			quad(p3ui, p3uo, p2uo, p2ui) // upward
			quad(p2di, p2do, p3do, p3di) // downward
			return
		}

	}
}
