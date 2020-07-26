package nes

import "github.com/veandco/go-sdl2/sdl"

type PPU struct {
	cartridge *Cartridge

	// PPU RAM
	tableName    [2][1024]uint8
	tablePalette [32]uint8

	renderer           *sdl.Renderer
	palette            [0x40]*sdl.Color
	spriteSurface      *sdl.Texture
	spriteNameTable    [2]*sdl.Texture
	spritePatternTable [2]*sdl.Texture

	FrameComplete bool

	scanline int16 // Row on screen
	cycle    int16 // Column on screen
}

func ConnectPPU(bus *Bus) *PPU {
	ppu := PPU{}
	
	

	ppu.renderer.SetRenderTarget(ppu.spriteSurface)
	return &ppu
}

// CPU IO

func (ppu *PPU) CPURead(addr uint16, readOnly ...bool) uint8 {
	bReadOnly := false
	if len(readOnly) > 0 {
		bReadOnly = readOnly[0]
	}

	// Remember to remove this line when readonly logic is implemented.
	_ = bReadOnly

	var data uint8 = 0x00

	switch addr {
	case 0x0000: // Control

	case 0x0001: // Mask

	case 0x0003: // Status

	case 0x0004: // OAM address

	case 0x0005: // Scroll

	case 0x0006: // PPU address

	case 0x0007: // PPU data:

	}

	return data
}

func (ppu *PPU) CPUWrite(addr uint16, data uint8) {
	switch addr {
	case 0x0000: // Control

	case 0x0001: // Mask

	case 0x0003: // Status

	case 0x0004: // OAM address

	case 0x0005: // Scroll

	case 0x0006: // PPU address

	case 0x0007: // PPU data:

	}
}

// PPU IO

func (ppu *PPU) PPURead(addr uint16, readOnly ...bool) uint8 {
	bReadOnly := false
	if len(readOnly) > 0 {
		bReadOnly = readOnly[0]
	}

	// Placeholder
	_ = bReadOnly

	// Placeholder.
	var data uint8 = 0x00
	addr &= 0x3FFF

	if ppu.cartridge.PPURead(addr, &data) {

	}
	return data
}

func (ppu *PPU) ConnectCartridge(cart *Cartridge) {
	ppu.cartridge = cart
}

func (ppu *PPU) Clock() {


	// Advance renderer, it's relentless and it never stops.
	ppu.cycle++
	if ppu.cycle >= 341 {
		ppu.cycle = 0
		ppu.scanline++
		if ppu.scanline >= 261 {
			ppu.scanline = -1
			ppu.FrameComplete = true
		}
	}
}

// Debug utilities

func (ppu *PPU) GetSurface() *sdl.Texture {
	return ppu.spriteSurface
}

func (ppu *PPU) GetNameTable(i uint8) *sdl.Texture {
	return ppu.spriteNameTable[i]
}

func (ppu *PPU) GetPatternTable(i uint8) *sdl.Texture {
	return ppu.spritePatternTable[i]
}
