package nes

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

type PPU struct {
	cartridge *Cartridge

	// PPU RAM
	tableName    [2][1024]uint8
	tablePattern [2][4096]uint8
	tablePalette [32]uint8

	palette            [][]uint8
	sprite             *sdl.Surface
	spriteNameTable    []*sdl.Surface
	spritePatternTable []*sdl.Surface

	FrameComplete bool

	scanline int32 // Row on screen
	cycle    int32 // Column on screen
}

func ConnectPPU(bus *Bus) *PPU {
	var err error
	ppu := PPU{}

	// Initialize PPU spritesheet
	ppu.sprite, err = sdl.CreateRGBSurfaceWithFormat(0, 256*2, 240*2, 8, sdl.PIXELFORMAT_RGB888)
	if err != nil {
		fmt.Printf("Failed to create sprite: %s\n", err)
		panic(err)
	}

	// Initialize PPU palette
	ppu.palette = ppu.initializePalette()

	// // Initialize PPU sprite name table
	// ppu.spriteNameTable = make([]*sdl.Texture, 2)

	// // Initialize PPU sprite pattern table
	// ppu.spritePatternTable = make([]*sdl.Texture, 2)
	// for i := range ppu.spritePatternTable {
	// 	ppu.spritePatternTable[i] = new(sdl.Texture)
	// }

	return &ppu
}

func (ppu *PPU) initializePalette() [][]uint8 {
	return [][]uint8{{84, 84, 84, 0},
		{0, 30, 116, 0},
		{8, 16, 144, 0},
		{48, 0, 136, 0},
		{68, 0, 100, 0},
		{92, 0, 48, 0},
		{84, 4, 0, 0},
		{60, 24, 0, 0},
		{32, 42, 0, 0},
		{8, 58, 0, 0},
		{0, 64, 0, 0},
		{0, 60, 0, 0},
		{0, 50, 60, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},

		{152, 150, 152, 0},
		{8, 76, 196, 0},
		{48, 50, 236, 0},
		{92, 30, 228, 0},
		{136, 20, 176, 0},
		{160, 20, 100, 0},
		{152, 34, 32, 0},
		{120, 60, 0, 0},
		{84, 90, 0, 0},
		{40, 114, 0, 0},
		{8, 124, 0, 0},
		{0, 118, 40, 0},
		{0, 102, 120, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},

		{236, 238, 236, 0},
		{76, 154, 236, 0},
		{120, 124, 236, 0},
		{176, 98, 236, 0},
		{228, 84, 236, 0},
		{236, 88, 180, 0},
		{236, 106, 100, 0},
		{212, 136, 32, 0},
		{160, 170, 0, 0},
		{116, 196, 0, 0},
		{76, 208, 32, 0},
		{56, 204, 108, 0},
		{56, 180, 204, 0},
		{60, 60, 60, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},

		{236, 238, 236, 0},
		{168, 204, 236, 0},
		{188, 188, 236, 0},
		{212, 178, 236, 0},
		{236, 174, 236, 0},
		{236, 174, 212, 0},
		{236, 180, 176, 0},
		{228, 196, 144, 0},
		{204, 210, 120, 0},
		{180, 222, 120, 0},
		{168, 226, 144, 0},
		{152, 226, 180, 0},
		{160, 214, 228, 0},
		{160, 162, 160, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}
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
	var pixelColor []uint8
	// Draw old-fashioned static noise.
	if rand.Int()%2 != 0 {
		pixelColor = ppu.palette[0x3F]
	} else {
		pixelColor = ppu.palette[0x30]
	}

	ppu.sprite.Set(int(ppu.cycle-1+1), int(ppu.scanline+1), color.RGBA{pixelColor[0], pixelColor[1], pixelColor[2], pixelColor[3]})

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

func (ppu *PPU) GetSprite() *sdl.Surface {
	return ppu.sprite
}

func (ppu *PPU) GetNameTable(i uint8) *sdl.Surface {
	return ppu.spriteNameTable[i]
}

func (ppu *PPU) GetPatternTable(i uint8) *sdl.Surface {
	return ppu.spritePatternTable[i]
}
