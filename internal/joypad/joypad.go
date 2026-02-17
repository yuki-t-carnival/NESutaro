package joypad

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Joypad struct {
	// Ebiten Inputs
	inputs     byte
	snapInputs byte

	// Others
	isGamepadEnabled bool   // From config.toml
	gamepadBind      [8]int // From config.toml

	setIndex  byte
	isPolling bool
}

func NewJoypad() *Joypad {
	return &Joypad{}
}

func (j *Joypad) Update() {
	j.updateEbitenInputs()
}

func (j *Joypad) updateEbitenInputs() {
	if j.isGamepadEnabled {
		j.inputs = j.updateEbitenKeys() | j.updateEbitenGamepadButtons()
	} else {
		j.inputs = j.updateEbitenKeys()
	}
}

func (j *Joypad) updateEbitenKeys() byte {
	inputs := [8]bool{
		ebiten.IsKeyPressed(ebiten.KeyZ),         // A
		ebiten.IsKeyPressed(ebiten.KeyX),         // B
		ebiten.IsKeyPressed(ebiten.KeyShiftLeft), // SELECT
		ebiten.IsKeyPressed(ebiten.KeyEnter),     // START
		ebiten.IsKeyPressed(ebiten.KeyUp),        // UP
		ebiten.IsKeyPressed(ebiten.KeyDown),      // DOWN
		ebiten.IsKeyPressed(ebiten.KeyLeft),      // LEFT
		ebiten.IsKeyPressed(ebiten.KeyRight),     // RIGHT
	}
	var keys byte
	for i, b := range inputs {
		if b {
			keys |= 1 << i
		}
	}
	return keys
}

func (j *Joypad) updateEbitenGamepadButtons() byte {
	id := ebiten.GamepadID(0)
	var inputs [8]bool
	for i, v := range j.gamepadBind {
		inputs[i] = ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton(v))
	}
	var gamepad byte
	for i, b := range inputs {
		if b {
			gamepad |= 1 << i
		}
	}
	return gamepad
}

// From config.toml
func (j *Joypad) SetIsGamepadEnabled(b bool) {
	j.isGamepadEnabled = b
}

// From config.toml
func (j *Joypad) SetIsGamepadBind(bind [8]int) {
	j.gamepadBind = bind
}

func (j *Joypad) Read4016() byte {
	if j.isPolling {
		return j.snapInputs >> 0 & 1
	} else {
		if j.setIndex == 8 {
			return 1
		} else {
			input := j.snapInputs >> j.setIndex & 1
			j.setIndex += 1
			return input
		}
	}
}

func (j *Joypad) Write4016(val byte) {
	if j.isPolling {
		if val&1 == 0 {
			j.isPolling = false
		}
	} else {
		if val&1 == 1 {
			j.isPolling = true
			j.snapInputs = j.inputs
			j.setIndex = 0
		}
	}
}
