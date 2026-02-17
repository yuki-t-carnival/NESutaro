package cpu

import (
	"fmt"
)

const TraceSize = 256

type Tracer struct {
	buf   [TraceSize]TraceEntry
	index int // The index of the next buffer to use.
}

type TraceEntry struct {
	a, x, y, s, p byte
	pc            uint16
	op            byte
	opName        string
}

func NewTracer(c *CPU) *Tracer {
	t := &Tracer{}
	t.Record(c)
	return t
}

// The Record Saves the current CPU Registers state in a ring buffer.
func (t *Tracer) Record(c *CPU) {
	op := c.read(c.pc)
	var opName string
	opName = opTable[op].Name
	t.buf[t.index] = TraceEntry{
		pc:     c.pc,
		a:      c.a,
		x:      c.x,
		y:      c.y,
		s:      c.s,
		p:      c.p,
		op:     op,
		opName: opName,
	}
	t.index = (t.index + 1) % TraceSize // ring buffer
}

// The Dump Outputs all CPU log buffers to console.
func (t *Tracer) Dump() {
	for i := range TraceSize {
		idx := (t.index + i) % TraceSize // Output from the oldest dump.
		buf := t.buf[idx]
		fmt.Printf(
			"PC:%04X "+
				"A:%02X "+
				"X:%02X "+
				"Y:%02X "+
				"S:%02X "+
				"P:%02X "+
				"Op:%02X "+
				"Fn:%s\n",
			buf.pc,
			buf.a,
			buf.x,
			buf.y,
			buf.s,
			buf.p,
			buf.op,
			buf.opName,
		)
	}
}

// The GetCPUInfo Gets CPU status strings for the debug screen.
func (t *Tracer) GetCPUInfo() []string {
	var idx int
	if t.index == 0 {
		idx = (TraceSize - 1) % TraceSize
	} else {
		idx = (t.index - 1) % TraceSize
	}
	buf := t.buf[idx]
	var str []string
	str = append(str, fmt.Sprintf("PC:%04X", buf.pc))
	str = append(str, fmt.Sprintf(" A:%02X", buf.a))
	str = append(str, fmt.Sprintf(" X:%02X", buf.x))
	str = append(str, fmt.Sprintf(" Y:%02X", buf.y))
	str = append(str, fmt.Sprintf(" S:%02X", buf.s))
	str = append(str, fmt.Sprintf(" P:%02X", buf.p))
	str = append(str, "")
	str = append(str, fmt.Sprintf("Op:%02X", buf.op))
	str = append(str, fmt.Sprintf("Fn:%s", buf.opName))
	return str
}
