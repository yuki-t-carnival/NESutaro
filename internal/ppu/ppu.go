package ppu

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"nesutaro/internal/ppu/bus"
	"os"
)

const (
	VblankFlag      byte = 1 << 7
	VblankNMIEnable byte = 1 << 7
)

var xFlipLUT [256]byte

type PPU struct {
	Bus      *bus.Bus
	viewport [256 * 240]int
	screen   [2]*image.RGBA
	front    int

	oam [256]byte

	cycles int

	v         uint16
	t         uint16
	x         uint8
	w         bool
	ppuctrl   byte
	ppumask   byte
	ppustatus byte
	oamaddr   byte

	hasNMI bool

	ly int

	nesPal [64]color.RGBA
}

func NewPPU(b *bus.Bus) *PPU {
	palFilePath := "nes.pal"
	palFile, err := os.ReadFile(palFilePath)
	if err != nil {
		log.Fatal(err)
	}
	p := &PPU{
		Bus: b,
	}
	p.screen[0] = image.NewRGBA(image.Rect(0, 0, 256, 240))
	p.screen[1] = image.NewRGBA(image.Rect(0, 0, 256, 240))

	for i := 0; i < 64; i++ {
		p.nesPal[i] = color.RGBA{palFile[i*3+0], palFile[i*3+1], palFile[i*3+2], 255}
	}
	return p
}

func (p *PPU) Step(cpuCycles int) {
	p.cycles += cpuCycles * 3
	for p.cycles >= 341 {
		isBGEnabled := p.ppumask>>3&1 == 1
		switch {
		case p.ly < 240: // HBlank: 0 ~ 239
			if isBGEnabled {
				p.drawBGLine()

				// dot 256
				if p.v&0x7000 != 0x7000 {
					p.v += 0x1000 // fineY++
				} else {
					p.v &^= 0x7000 // fineY = 0
					coarseY := p.v & 0x03E0 >> 5
					switch coarseY {
					case 29:
						coarseY = 0
						p.v ^= 0x0800 // tableY reverse
					case 31:
						coarseY = 0
					default:
						coarseY += 1
					}
					p.v = p.v&^0x03E0 + coarseY<<5
				}
				p.v = p.v&^0x041F + p.t&0x041F // p.v = p.t (coarseX, tableX only) dot 257
			}
			if p.ly > 0 {
				isSpriteEnabled := p.ppumask>>4&1 == 1
				if isSpriteEnabled {
					p.drawSpritesLine(p.ly - 1)
				}
			}

		case p.ly == 240: // VBlank: 240 ~ 259
			p.front ^= 1
			backColor := p.nesPal[p.Bus.Read(0x3F00)]
			for y := 0; y < 240; y++ {
				for x := 0; x < 256; x++ {
					p.screen[p.front^1].SetRGBA(x, y, backColor)
				}
			}
			p.ppustatus |= VblankFlag
			if p.ppuctrl&VblankNMIEnable != 0 {
				p.hasNMI = true
			}
		case p.ly == 260: // Post(pre) render scanline: 260, 261
			p.ppustatus &^= VblankFlag
			p.v = p.v&^0x7BE0 + p.t&0x7BE0 // p.v = p.t (fineY, coarseY, tableY only) dot 280 ~ 304
		}

		p.cycles -= 341

		p.ly += 1
		if p.ly >= 262 {
			p.ly = 0
		}

	}
}

// ============================================ BG =================================================

func (p *PPU) drawBGLine() {
	prevCoarseX := -1
	//p.x := int(p.x & 7)

	var tileLine [2]byte
	startX := int((p.ppumask>>1&1 ^ 1) * 8)
	for i := 0; i < 256; i++ {
		if i >= startX {
			currentCoarseX := int(p.v & 0x001F)

			if currentCoarseX != prevCoarseX {
				nameAddr := uint16(0x2000) + p.v&0x0FFF
				tileIndex := p.Bus.Read(nameAddr)
				tileLine = p.getBGTileLine(tileIndex)
			}
			raw := p.getTileLinePixel(&tileLine, int(p.x))
			if raw != 0 {
				rgba := p.resolveBGPalette(raw)
				p.setPixel(i, p.ly, rgba)
			}
			prevCoarseX = currentCoarseX
		}

		if p.x&7 != 7 {
			p.x += 1
		} else {
			p.x = 0
			coarseX := p.v & 0x001F
			if coarseX == 31 {
				coarseX = 0
				p.v ^= 0x0400
			} else {
				coarseX += 1
			}
			p.v = p.v&^0x001F + coarseX
		}
	}
}
func (p *PPU) getBGTileLine(tileIndex byte) [2]byte {
	base := uint16(0x0000)
	if p.ppuctrl>>4&1 == 1 {
		base = 0x1000
	}
	fineY := p.v & 0x7000 >> 12
	var data [2]byte
	data[0] = p.Bus.Read(base + uint16(tileIndex)<<4 + (0 << 3) + fineY)
	data[1] = p.Bus.Read(base + uint16(tileIndex)<<4 + (1 << 3) + fineY)
	return data
}
func (p *PPU) resolveBGPalette(raw int) color.RGBA {
	blockX := p.v & 0x001C >> 2 // = coarseX / 4
	blockY := p.v & 0x0380 >> 7 // = coarseY / 4
	attrAddr := uint16(0x23C0) + p.v&0x0C00 + blockY<<3 + blockX
	attr := p.Bus.Read(attrAddr)

	shift := p.v&0x0040>>4 + p.v&0x0002

	palVal := uint16(attr >> shift & 3)
	colorVal := p.Bus.Read(0x3F00 + palVal<<2 + uint16(raw))
	return p.nesPal[colorVal&0x3F]
}

// ========================================= Sprites ===============================================

func (p *PPU) drawSpritesLine(ly int) {
	tileH := 8 * int(1+p.ppuctrl>>5&1)
	var tileLine [2]byte
	startX := int(p.ppumask>>2&1^1) * 8
	for i := 63; i >= 0; i-- {
		posY := int(p.oam[i<<2+0])
		tileIdx := p.oam[i<<2+1]
		attr := int(p.oam[i<<2+2])
		posX := int(p.oam[i<<2+3])

		isXFlip := attr>>6&1 == 1
		isYFlip := attr>>7&1 == 1

		if posY <= ly && ly < posY+tileH {
			fineY := ly - posY
			tileY := fineY
			if isYFlip {
				tileY = tileH - 1 - tileY
			}
			tileLine = p.getSpriteTileLine(tileIdx, tileY)
			for j := 0; j < 8; j++ {
				if startX <= posX+j && posX+j < 256 {
					tileX := j
					if isXFlip {
						tileX = 7 - tileX
					}
					raw := p.getTileLinePixel(&tileLine, tileX)
					palVal := attr & 3
					if raw != 0 {
						rgba := p.resolveSpritesPalette(palVal, raw)
						p.setPixel(posX+j, ly+1, rgba)
					}
					p.ppustatus |= 1 << 6
				}
			}
		}
	}
}
func (p *PPU) getSpriteTileLine(tileIndex byte, fineY int) [2]byte {
	idx := uint16(tileIndex)
	var start uint16

	tileH := 8 * int(1+p.ppuctrl>>5&1)
	if tileH == 16 {
		start = idx&1<<12 + idx&0xFE<<4
	} else {
		spriteAddrNum := uint16(p.ppuctrl >> 3 & 1)
		start = spriteAddrNum<<12 + idx<<4
	}

	var data [2]byte
	data[0] = p.Bus.Read(start + (0 << 3) + uint16(fineY))
	data[1] = p.Bus.Read(start + (1 << 3) + uint16(fineY))
	return data
}
func (p *PPU) resolveSpritesPalette(palVal int, raw int) color.RGBA {
	colorVal := p.Bus.Read(0x3F10 + uint16(palVal)<<2 + uint16(raw))
	return p.nesPal[colorVal&0x3F]
}

// ======================================== BG/Sprites =============================================

func (p *PPU) getTileLinePixel(tileLine *[2]byte, fineX int) int {
	shift := 7 - fineX
	lo := tileLine[0] & (1 << shift) >> shift
	hi := tileLine[1] & (1 << shift) >> shift
	px := int(hi<<1 | lo)
	return px & 0x03
}
func (p *PPU) setPixel(x, y int, rgba color.RGBA) {
	p.screen[p.front^1].SetRGBA(x, y, rgba)
}

// Get Viewport pixels converted from colorIndex to RGBA
func (p *PPU) GetGameScreen() *image.RGBA {
	return p.screen[p.front]
}

func (p *PPU) ReadOAM(addr uint16) byte {
	return p.oam[addr]
}
func (p *PPU) WriteOAM(addr uint16, val byte) {
	p.oam[addr] = val
}

func (p *PPU) WritePPUCTRL(val byte) {
	p.ppuctrl = val
	p.t = p.t&0x73FF | uint16(val)&0x03<<10
}

func (p *PPU) WritePPUMASK(val byte) {
	p.ppumask = val
}

func (p *PPU) ReadPPUSTATUS() byte {
	p.w = false
	val := p.ppustatus
	//p.ppustatus &^= VblankFlag
	return val
}

func (p *PPU) WriteOAMADDR(val byte) {
	p.oamaddr = val
}

func (p *PPU) ReadOAMDATA() byte {
	return p.oam[p.oamaddr]
}

func (p *PPU) WriteOAMDATA(val byte) {
	p.oam[p.oamaddr] = val
	p.oamaddr += 1
}

func (p *PPU) WritePPUSCROLL(val byte) {
	isFirstWriting := !p.w
	if isFirstWriting {
		coarseX := uint16(val >> 3)
		p.t = p.t&^0x001F + coarseX
		p.x = val & 0x07
	} else {
		fineY := uint16(val&0x07) << 12
		coarseY := uint16(val&0xF8) << 2
		p.t = fineY + p.t&^0x73E0 + coarseY
	}
	p.w = !p.w
}

func (p *PPU) WritePPUADDR(val byte) {
	isFirstWriting := !p.w
	if isFirstWriting {
		// Bit 14 is forced 0 when writting the PPUADDR high byte.
		p.t = uint16(val&0x3F)<<8 + p.t&^0xFF00
	} else {
		p.t = p.t&^0x00FF + uint16(val)
		p.v = p.t
	}
	p.w = !p.w
}

func (p *PPU) ReadPPUDATA() byte {
	val := p.Bus.Read(p.v)
	if p.ppuctrl>>2&1 == 0 {
		p.v += 1
	} else {
		p.v += 32
	}
	return val
}

func (p *PPU) WritePPUDATA(val byte) {
	p.Bus.Write(p.v, val)
	if p.ppuctrl>>2&1 == 0 {
		p.v += 1
	} else {
		p.v += 32
	}
}

func (p *PPU) WriteOAMDMA(data *[0x100]byte) {
	p.oam = *data
}

func (p *PPU) HasNMI() bool {
	return p.hasNMI
}

func (p *PPU) DisableNMI() {
	p.hasNMI = false
}

func (p *PPU) OAMLog() {
	for i, v := range p.oam {
		fmt.Printf("p.oam[%d]=%02X\n", i, v)
	}
}
