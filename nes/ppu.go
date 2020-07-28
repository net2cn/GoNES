package nes

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

// Bitmask for status register.
const (
	statusUnused         = 0x1F // Mask the unused register.
	statusSpriteOverFlow = (1 << 5)
	statusSpriteZeroHit  = (1 << 6)
	statusVerticalBlank  = (1 << 7)
)

// Bitmask for mask register.
const (
	maskGreyscale            = (1 << 0)
	maskRenderBackgroundLeft = (1 << 1)
	maskRenderSpritesLeft    = (1 << 2)
	maskRenderBackground     = (1 << 3)
	maskRenderSprites        = (1 << 4)
	maskEnhanceRed           = (1 << 5)
	maskEnhanceGreen         = (1 << 6)
	maskEnhanceBlue          = (1 << 7)
)

// Bitmask for control register.
const (
	controlNameTableX        = (1 << 0)
	controlNameTableY        = (1 << 1)
	controlIncrementMode     = (1 << 2)
	controlPatternSprite     = (1 << 3)
	controlPatternBackground = (1 << 4)
	controlSpriteSize        = (1 << 5)
	controlSlaveMode         = (1 << 6)
	controlEnableNMI         = (1 << 7)
)

// PPU Nintendo 2C02 PPU struct
type PPU struct {
	cartridge *Cartridge

	// PPU RAM
	TableName    [2][1024]uint8
	tablePattern [2][4096]uint8
	tablePalette [32]uint8

	NMI bool

	palette            [][]uint8
	screen             *sdl.Surface
	spriteNameTable    []*sdl.Surface
	spritePatternTable []*sdl.Surface

	FrameComplete bool

	scanline int32 // Row on screen
	cycle    int32 // Column on screen

	status  uint8 // Control register
	mask    uint8 // Mask register
	control uint8 // Control register

	addressLatch  uint8
	ppuDataBuffer uint8 // Data would delayed by 1 cycle when read.
	ppuAddress    uint16
}

// ConnectPPU Initialize a PPU and connect it to the bus.
func ConnectPPU(bus *Bus) *PPU {
	var err error
	ppu := PPU{}

	// Initialize PPU spritesheet
	// For NES screen, there's 341 cycles and 261 scanlines for each screen,
	// but NES is only generating a frame with 256 cycles and 240 scanlines.
	// So I doubled the frame size in both width and height so that the
	// screen won't overflow.
	ppu.screen, err = sdl.CreateRGBSurfaceWithFormat(0, 256*2, 240*2, 8, sdl.PIXELFORMAT_RGB888)
	if err != nil {
		fmt.Printf("Failed to create sprite: %s\n", err)
		panic(err)
	}

	// Initialize PPU palette
	ppu.palette = ppu.initializePalette()

	// Initialize PPU sprite name table
	ppu.spriteNameTable = make([]*sdl.Surface, 2)
	for i := range ppu.spriteNameTable {
		ppu.spriteNameTable[i], err = sdl.CreateRGBSurfaceWithFormat(0, 256*2, 240*2, 8, sdl.PIXELFORMAT_RGB888)
		if err != nil {
			fmt.Printf("Failed to create sprite name table %d: %s\n", i, err)
			panic(err)
		}
	}

	// Initialize PPU sprite pattern table
	ppu.spritePatternTable = make([]*sdl.Surface, 2)
	for i := range ppu.spritePatternTable {
		ppu.spritePatternTable[i], err = sdl.CreateRGBSurfaceWithFormat(0, 256*2, 240*2, 8, sdl.PIXELFORMAT_RGB888)
		if err != nil {
			fmt.Printf("Failed to create pattern name table %d: %s\n", i, err)
			panic(err)
		}
	}

	ppu.NMI = false

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

// REG IO

func (ppu *PPU) getFlag(reg *uint8, f uint8) uint8 {
	if (*reg & f) > 0 {
		return 1
	}

	return 0
}

func (ppu *PPU) setFlag(reg *uint8, f uint8, v bool) {
	if v {
		*reg |= f
	} else {
		*reg &= ^f
	}
}

// CPU IO

// CPURead CPU read from PPU.
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

	case 0x0002: // Status
		data = (ppu.status & 0xE0) | (ppu.ppuDataBuffer & 0x1F)
		ppu.setFlag(&ppu.status, statusVerticalBlank, false)
		ppu.addressLatch = 0
	case 0x0003: // OAM address

	case 0x0004: // OAM data

	case 0x0005: // Scroll

	case 0x0006: // PPU address

	case 0x0007: // PPU data:
		data = ppu.ppuDataBuffer
		ppu.ppuDataBuffer = ppu.PPURead(ppu.ppuAddress)

		if ppu.ppuAddress >= 0x3F00 {
			data = ppu.ppuDataBuffer
		}
		ppu.ppuAddress++
	}

	return data
}

// CPUWrite CPU write to PPU.
func (ppu *PPU) CPUWrite(addr uint16, data uint8) {
	switch addr {
	case 0x0000: // Control
		ppu.control = data
	case 0x0001: // Mask
		ppu.mask = data
	case 0x0002: // Status

	case 0x0003: // OAM address

	case 0x0004: // OAM data

	case 0x0005: // Scroll

	case 0x0006: // PPU address
		if ppu.addressLatch == 0 {
			ppu.ppuAddress = (ppu.ppuAddress & 0x00FF) | (uint16(data) << 8)
			ppu.addressLatch = 1

		} else {
			ppu.ppuAddress = (ppu.ppuAddress & 0xFF00) | uint16(data)
			ppu.addressLatch = 0
		}
	case 0x0007: // PPU data
		ppu.PPUWrite(ppu.ppuAddress, data)
		if ppu.getFlag(&ppu.control, controlIncrementMode) != 0 {
			ppu.ppuAddress += 32
		} else {
			ppu.ppuAddress++
		}
	}
}

// PPU IO

// PPURead PPU read from PPU.
func (ppu *PPU) PPURead(addr uint16, readOnly ...bool) uint8 {
	bReadOnly := false
	if len(readOnly) > 0 {
		bReadOnly = readOnly[0]
	}

	// Placeholder
	_ = bReadOnly

	var data uint8 = 0x00
	addr &= 0x3FFF

	if ppu.cartridge.PPURead(addr, &data) { // Mapper relocation

	} else if addr >= 0x0000 && addr <= 0x1FFF { // Pattern memory
		data = ppu.tablePattern[(addr&0x1000)>>12][addr&0x0FFF]
	} else if addr >= 0x2000 && addr <= 0x3EFF { // Name Table memory
		addr &= 0x0FFF

		if ppu.cartridge.mirror == mirrorVertical {
			if addr >= 0x0000 && addr <= 0x03FF {
				data = ppu.TableName[0][addr&0x03FF]
			} else if addr >= 0x0400 && addr <= 0x07FF {
				data = ppu.TableName[1][addr&0x03FF]
			} else if addr >= 0x0800 && addr <= 0x03FF {
				data = ppu.TableName[0][addr&0x03FF]
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				data = ppu.TableName[1][addr&0x03FF]
			}
		} else if ppu.cartridge.mirror == mirrorHorizontal {
			if addr >= 0x0000 && addr <= 0x03FF {
				data = ppu.TableName[0][addr&0x03FF]
			} else if addr >= 0x0400 && addr <= 0x07FF {
				data = ppu.TableName[0][addr&0x03FF]
			} else if addr >= 0x0800 && addr <= 0x03FF {
				data = ppu.TableName[1][addr&0x03FF]
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				data = ppu.TableName[1][addr&0x03FF]
			}
		}
	} else if addr >= 0x3F00 && addr <= 0x3FFF { // Palette Memory
		addr &= 0x001F
		if addr == 0x0010 {
			addr = 0x0000
		}
		if addr == 0x0014 {
			addr = 0x0004
		}
		if addr == 0x0018 {
			addr = 0x0008
		}
		if addr == 0x001C {
			addr = 0x000C
		}
		data = ppu.tablePalette[addr]
	}

	return data
}

// PPUWrite PPU write to PPU.
func (ppu *PPU) PPUWrite(addr uint16, data uint8) {
	addr &= 0x3FFF

	if ppu.cartridge.PPUWrite(addr, data) { // Mapper relocation

	} else if addr >= 0x0000 && addr <= 0x1FFF { // Pattern memory
		ppu.tablePattern[(addr&0x1000)>>12][addr&0x0FFF] = data
	} else if addr >= 0x2000 && addr <= 0x3EFF { // Name Table memory
		addr &= 0x0FFF

		if ppu.cartridge.mirror == mirrorVertical {
			if addr >= 0x0000 && addr <= 0x03FF {
				ppu.TableName[0][addr&0x03FF] = data
			} else if addr >= 0x0400 && addr <= 0x07FF {
				ppu.TableName[1][addr&0x03FF] = data
			} else if addr >= 0x0800 && addr <= 0x03FF {
				ppu.TableName[0][addr&0x03FF] = data
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				ppu.TableName[1][addr&0x03FF] = data
			}
		} else if ppu.cartridge.mirror == mirrorHorizontal {
			if addr >= 0x0000 && addr <= 0x03FF {
				ppu.TableName[0][addr&0x03FF] = data
			} else if addr >= 0x0400 && addr <= 0x07FF {
				ppu.TableName[0][addr&0x03FF] = data
			} else if addr >= 0x0800 && addr <= 0x03FF {
				ppu.TableName[1][addr&0x03FF] = data
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				ppu.TableName[1][addr&0x03FF] = data
			}
		}
	} else if addr >= 0x3F00 && addr <= 0x3FFF { // Palette Memory
		addr &= 0x001F
		if addr == 0x0010 {
			addr = 0x0000
		}
		if addr == 0x0014 {
			addr = 0x0004
		}
		if addr == 0x0018 {
			addr = 0x0008
		}
		if addr == 0x001C {
			addr = 0x000C
		}
		ppu.tablePalette[addr] = data
	}
}

// ConnectCartridge Connect cartridge to PPU.
func (ppu *PPU) ConnectCartridge(cart *Cartridge) {
	ppu.cartridge = cart
}

// Clock Clock PPU once.
func (ppu *PPU) Clock() {
	if ppu.scanline == -1 && ppu.cycle == 1 {
		ppu.setFlag(&ppu.status, statusVerticalBlank, false)
	}

	if ppu.scanline == 241 && ppu.cycle == 1 {
		ppu.setFlag(&ppu.status, statusVerticalBlank, true)
		if ppu.getFlag(&ppu.control, controlEnableNMI) != 0 {
			ppu.NMI = true
		}
	}

	var pixelColor []uint8
	// Draw old-fashioned static noise.
	if rand.Int()%2 != 0 {
		pixelColor = ppu.palette[0x3F]
	} else {
		pixelColor = ppu.palette[0x30]
	}

	ppu.screen.Set(int(ppu.cycle-1+1), int(ppu.scanline+1),
		color.RGBA{pixelColor[0],
			pixelColor[1],
			pixelColor[2],
			pixelColor[3]})

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

// GetScreen Return PPU rendered screen.
func (ppu *PPU) GetScreen() *sdl.Surface {
	return ppu.screen
}

// GetNameTable Return PPU internal name table.
func (ppu *PPU) GetNameTable(i uint8) *sdl.Surface {
	return ppu.spriteNameTable[i]
}

// GetColorFromPaletteRAM Get a color from PPU internal palette RAM.
func (ppu *PPU) GetColorFromPaletteRAM(palette uint8, pixel uint8) color.RGBA {
	data := ppu.palette[ppu.PPURead(0x3F00+(uint16(palette)<<2)+uint16(pixel))&0x3F]
	return color.RGBA{data[0], data[1], data[2], data[3]}
}

// GetPatternTable Get PPU internal pattern table.
func (ppu *PPU) GetPatternTable(i uint8, palette uint8) *sdl.Surface {
	for tileY := 0; tileY < 16; tileY++ {
		for tileX := 0; tileX < 16; tileX++ {
			var offset uint16 = uint16(tileY*256 + tileX*16) // Byte offset

			for row := 0; row < 8; row++ {
				var tileLSB uint8 = ppu.PPURead(uint16(uint16(i)*0x1000 + offset + uint16(row) + 0x0000))
				var tileMSB uint8 = ppu.PPURead(uint16(uint16(i)*0x1000 + offset + uint16(row) + 0x0008))

				for col := 0; col < 8; col++ {
					var pixel uint8 = (tileLSB & 0x01) + (tileMSB & 0x01)
					tileLSB >>= 1
					tileMSB >>= 1

					ppu.spritePatternTable[i].Set(
						tileX*8+(7-col),
						tileY*8+row,
						ppu.GetColorFromPaletteRAM(palette, pixel),
					)
				}
			}
		}
	}

	return ppu.spritePatternTable[i]
}
