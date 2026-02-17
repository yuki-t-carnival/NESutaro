package emulator

import (
	"nesutaro/internal/cartridge"
	"nesutaro/internal/cpu"
	cbus "nesutaro/internal/cpu/bus"
	"nesutaro/internal/joypad"
	"nesutaro/internal/ppu"
	pbus "nesutaro/internal/ppu/bus"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CyclesPerFrame float64 = 1789773 / 60.0
)

type Emulator struct {
	CPU       *cpu.CPU
	cpuCycles float64

	IsPaused    bool
	IsPauseMode bool

	isKeyP       bool
	isKeyS       bool
	isKeyEsc     bool
	isPrevKeyP   bool
	isPrevKeyS   bool
	isPrevKeyEsc bool
}

func NewEmulator(rom /* , sav */ []byte) *Emulator {
	cart := cartridge.NewCartridge(rom /* , sav */)
	pbus := pbus.NewBus(cart)
	p := ppu.NewPPU(pbus)
	j := joypad.NewJoypad()
	cbus := cbus.NewBus(cart, p, j)
	c := cpu.NewCPU(cbus)
	c.Tracer = cpu.NewTracer(c)

	e := &Emulator{
		CPU:         c,
		IsPauseMode: false,
		IsPaused:    false,
	}

	return e
}

func (e *Emulator) RunFrame() int {
	maxCycles := CyclesPerFrame
	e.CPU.Bus.Joypad.Update()
	for e.cpuCycles < maxCycles {
		e.updateEbitenKeys()
		e.updateEmuMode()
		if e.CPU.IsPanic || e.isKeyEsc { // for debug
			e.panicDump()
			return -1
		} else if e.IsPaused {
			return 0
		}
		var c int
		c = e.CPU.Step()
		//e.CPU.Bus.Timer.Step(c, e.CPU.IsStopped)
		e.CPU.Bus.PPU.Step(c)
		//e.CPU.Bus.APU.Step(c / cpuSpeed)
		e.CPU.Tracer.Record(e.CPU)
		e.cpuCycles += float64(c)
	}
	e.cpuCycles -= maxCycles
	return 0
}

// KeyP: Toggle Run/Pause Mode
// KeyS: Run a single step
func (e *Emulator) updateEmuMode() {
	if e.isKeyP {
		e.IsPauseMode = !e.IsPauseMode
	}
	e.IsPaused = e.IsPauseMode && !e.isKeyS
}

// In case of Panic, CPU status is output to the console.
func (e *Emulator) panicDump() {
	e.CPU.Tracer.Dump()
}

func (e *Emulator) updateEbitenKeys() {
	isP := ebiten.IsKeyPressed(ebiten.KeyP)
	isS := ebiten.IsKeyPressed(ebiten.KeyS)
	isEsc := ebiten.IsKeyPressed(ebiten.KeyEscape)
	e.isKeyP = !e.isPrevKeyP && isP
	e.isKeyS = !e.isPrevKeyS && isS
	e.isKeyEsc = !e.isPrevKeyEsc && isEsc
	e.isPrevKeyP = isP
	e.isPrevKeyS = isS
	e.isPrevKeyEsc = isEsc
}

func (e *Emulator) GetDebugLog() []string {
	var state string
	if e.IsPaused {
		state = "PAUSED"
	} else {
		state = "PLAY  "
	}
	strs := []string{}
	strs = append(strs, state)
	strs = append(strs, "")
	strs = append(strs, e.CPU.Tracer.GetCPUInfo()...)
	strs = append(strs, "")
	//strs = append(strs, e.CPU.Bus.APU.GetAPUInfo()...)
	return strs
}
