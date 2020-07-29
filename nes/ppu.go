package nes

import (
	"fmt"
	"image/color"

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

// Bitmask for PPU loppy register.
const (
	loppyCoarseX    = 0x001F
	loppyCoarseY    = 0x03E0
	loppyNameTableX = (1 << 10)
	loppyNameTableY = (1 << 11)
	loppyFineY      = 0x7000
	loppyUnused     = (1 << 15)
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

	vramAddr uint16
	tramAddr uint16

	fineX uint8

	nextTileID   uint8
	nextTileAttr uint8
	nextTileLSB  uint8
	nextTileMSB  uint8

	shifterPatternLo uint16
	shifterPatternHi uint16
	shifterAttribLo  uint16
	shifterAttribHi  uint16
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

// Addr IO
func (ppu *PPU) getLoppyRegister(addr *uint16, m int) uint16 {
	return (*addr & uint16(m)) / uint16(m&-m)
}

func (ppu *PPU) setLoppyRegister(addr *uint16, m int, val uint8) {
	*addr |= uint16(val) * uint16(m&-m)
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
		ppu.ppuDataBuffer = ppu.PPURead(ppu.vramAddr)

		// If CPU is reading address above 0x3F00, we need to instantly return its value
		// instead of delay 1 cylce.
		if ppu.vramAddr >= 0x3F00 {
			data = ppu.ppuDataBuffer
		}

		if ppu.getFlag(&ppu.control, controlIncrementMode) != 0 {
			ppu.vramAddr += 32
		} else {
			ppu.vramAddr++
		}
	}

	return data
}

// CPUWrite CPU write to PPU.
func (ppu *PPU) CPUWrite(addr uint16, data uint8) {
	switch addr {
	case 0x0000: // Control
		ppu.control = data
		ppu.setLoppyRegister(&ppu.tramAddr, loppyNameTableX,
			ppu.getFlag(&ppu.control, controlNameTableX))
		ppu.setLoppyRegister(&ppu.tramAddr, loppyNameTableY,
			ppu.getFlag(&ppu.control, controlNameTableY))
	case 0x0001: // Mask
		ppu.mask = data
	case 0x0002: // Status

	case 0x0003: // OAM address

	case 0x0004: // OAM data

	case 0x0005: // Scroll
		if ppu.addressLatch == 0 {
			ppu.fineX = data & 0x07
			ppu.setLoppyRegister(&ppu.tramAddr, loppyCoarseX, data>>3)
			ppu.addressLatch = 1
		} else {
			ppu.setLoppyRegister(&ppu.tramAddr, loppyFineY, data&0x07)
			ppu.setLoppyRegister(&ppu.tramAddr, loppyCoarseY, data>>3)
			ppu.addressLatch = 0
		}
	case 0x0006: // PPU address
		if ppu.addressLatch == 0 {
			ppu.tramAddr = (uint16(data&0x3F) << 8) | (ppu.tramAddr & 0x00FF)
			ppu.addressLatch = 1

		} else {
			ppu.tramAddr = (ppu.tramAddr & 0xFF00) | uint16(data)
			ppu.vramAddr = ppu.tramAddr
			ppu.addressLatch = 0
		}
	case 0x0007: // PPU data
		ppu.PPUWrite(ppu.vramAddr, data)
		if ppu.getFlag(&ppu.control, controlIncrementMode) != 0 {
			ppu.vramAddr += 32
		} else {
			ppu.vramAddr++
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

		if ppu.getFlag(&ppu.mask, maskGreyscale) != 0 {
			data = ppu.tablePalette[addr] & 0x30
		} else {
			data = ppu.tablePalette[addr] & 0x3F
		}
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
	var incrementScrollX func() = func() {
		if (ppu.getFlag(&ppu.mask, maskRenderBackground) != 0) ||
			(ppu.getFlag(&ppu.mask, maskRenderSprites) != 0) {
			if ppu.getLoppyRegister(&ppu.vramAddr, loppyCoarseX) == 31 {
				ppu.setLoppyRegister(&ppu.vramAddr, loppyCoarseX, 0)
				ppu.setLoppyRegister(&ppu.vramAddr, loppyNameTableX,
					uint8(^ppu.getLoppyRegister(&ppu.vramAddr, loppyNameTableX)))
			} else {
				ppu.setLoppyRegister(&ppu.vramAddr, loppyCoarseX,
					uint8(ppu.getLoppyRegister(&ppu.vramAddr, loppyCoarseX)+1))
			}
		}
	}

	var incrementScrollY func() = func() {
		if (ppu.getFlag(&ppu.mask, maskRenderBackground) != 0) ||
			(ppu.getFlag(&ppu.mask, maskRenderSprites) != 0) {
			if ppu.getLoppyRegister(&ppu.vramAddr, loppyFineY) < 7 {
				ppu.setLoppyRegister(&ppu.vramAddr, loppyFineY,
					uint8(ppu.getLoppyRegister(&ppu.vramAddr, loppyFineY)+1))
			} else {
				ppu.setLoppyRegister(&ppu.vramAddr, loppyFineY, 0)

				if ppu.getLoppyRegister(&ppu.vramAddr, loppyCoarseY) == 29 {
					ppu.setLoppyRegister(&ppu.vramAddr, loppyCoarseY, 0)
					ppu.setLoppyRegister(&ppu.vramAddr, loppyNameTableY,
						uint8(^ppu.getLoppyRegister(&ppu.vramAddr, loppyNameTableY)))
				} else if ppu.getLoppyRegister(&ppu.vramAddr, loppyCoarseY) == 31 {
					ppu.setLoppyRegister(&ppu.vramAddr, loppyCoarseY, 0)
				} else {
					ppu.setLoppyRegister(&ppu.vramAddr, loppyCoarseY,
						uint8(ppu.getLoppyRegister(&ppu.vramAddr, loppyCoarseY)+1))
				}
			}
		}
	}

	var transferAddressX func() = func() {
		if (ppu.getFlag(&ppu.mask, maskRenderBackground) != 0) ||
			(ppu.getFlag(&ppu.mask, maskRenderSprites) != 0) {
			ppu.setLoppyRegister(&ppu.vramAddr, loppyNameTableX,
				uint8(ppu.getLoppyRegister(&ppu.tramAddr, loppyNameTableX)))
			ppu.setLoppyRegister(&ppu.vramAddr, loppyCoarseX,
				uint8(ppu.getLoppyRegister(&ppu.tramAddr, loppyCoarseX)))
		}
	}

	var transferAddressY func() = func() {
		if (ppu.getFlag(&ppu.mask, maskRenderBackground) != 0) ||
			(ppu.getFlag(&ppu.mask, maskRenderSprites) != 0) {
			ppu.setLoppyRegister(&ppu.vramAddr, loppyFineY,
				uint8(ppu.getLoppyRegister(&ppu.tramAddr, loppyFineY)))
			ppu.setLoppyRegister(&ppu.vramAddr, loppyNameTableY,
				uint8(ppu.getLoppyRegister(&ppu.tramAddr, loppyNameTableY)))
			ppu.setLoppyRegister(&ppu.vramAddr, loppyCoarseY,
				uint8(ppu.getLoppyRegister(&ppu.tramAddr, loppyCoarseY)))
		}
	}

	var loadBackgroundShifters func() = func() {
		ppu.shifterPatternLo = (ppu.shifterPatternLo & 0xFF00) | uint16(ppu.nextTileLSB)
		ppu.shifterPatternHi = (ppu.shifterPatternHi & 0xFF00) | uint16(ppu.nextTileMSB)
		if ppu.nextTileAttr&0b01 != 0 {
			ppu.shifterAttribLo = (ppu.shifterAttribLo & 0xFF00) | 0xFF
		} else {
			ppu.shifterAttribLo = (ppu.shifterAttribLo & 0xFF00) | 0x00
		}
		if ppu.nextTileAttr&0b10 != 0 {
			ppu.shifterAttribHi = (ppu.shifterAttribHi & 0xFF00) | 0xFF
		} else {
			ppu.shifterAttribHi = (ppu.shifterAttribHi & 0xFF00) | 0x00
		}
	}

	var updateShifters func() = func() {
		if ppu.getFlag(&ppu.mask, maskRenderBackground) != 0 {
			ppu.shifterPatternLo <<= 1
			ppu.shifterPatternHi <<= 1

			ppu.shifterAttribLo <<= 1
			ppu.shifterAttribHi <<= 1
		}
	}

	if ppu.scanline >= -1 && ppu.scanline < 240 {
		if ppu.scanline == 0 && ppu.cycle == 0 {
			ppu.cycle = 1
		}

		if ppu.scanline == -1 && ppu.cycle == 1 {
			ppu.setFlag(&ppu.status, statusVerticalBlank, false)
		}

		if (ppu.cycle >= 2 && ppu.cycle <= 258) || (ppu.cycle >= 321 && ppu.cycle < 338) {
			updateShifters()

			switch (ppu.cycle - 1) % 8 {
			case 0:
				loadBackgroundShifters()
				ppu.nextTileID = ppu.PPURead(0x2000 | (ppu.vramAddr & 0x0FFF))
			case 2:
				ppu.nextTileAttr = ppu.PPURead(0x23C0 |
					(ppu.getLoppyRegister(&ppu.vramAddr, loppyNameTableY) << 11) |
					(ppu.getLoppyRegister(&ppu.vramAddr, loppyNameTableX) << 10) |
					((ppu.getLoppyRegister(&ppu.vramAddr, loppyCoarseY) >> 2) << 3) |
					(ppu.getLoppyRegister(&ppu.vramAddr, loppyCoarseX) >> 2))
				if ppu.getLoppyRegister(&ppu.vramAddr, loppyCoarseY)&0x02 != 0 {
					ppu.nextTileAttr >>= 4
				}
				if ppu.getLoppyRegister(&ppu.vramAddr, loppyCoarseX)&0x02 != 0 {
					ppu.nextTileAttr >>= 2
				}
				ppu.nextTileAttr &= 0x03
			case 4:
				ppu.nextTileLSB = ppu.PPURead((uint16(ppu.getFlag(&ppu.control, controlPatternBackground)) << 12) +
					uint16(ppu.nextTileID)<<4 +
					ppu.getLoppyRegister(&ppu.vramAddr, loppyFineY) + 0)
			case 6:
				ppu.nextTileMSB = ppu.PPURead((uint16(ppu.getFlag(&ppu.control, controlPatternBackground)) << 12) +
					uint16(ppu.nextTileID)<<4 +
					ppu.getLoppyRegister(&ppu.vramAddr, loppyFineY) + 8)
			case 7:
 				incrementScrollX()
			}
		}

		if ppu.cycle == 256 {
			incrementScrollY()
		}

		if ppu.cycle == 257 {
			loadBackgroundShifters()
			transferAddressX()
		}

		if ppu.cycle == 338 || ppu.cycle == 340 {
			ppu.nextTileID = ppu.PPURead(0x2000 | (ppu.vramAddr & 0x0FFF))
		}

		if ppu.scanline == -1 && ppu.cycle >= 280 && ppu.cycle < 305 {
			transferAddressY()
		}
	}

	if ppu.scanline == 240 {
		// Placeholder
	}

	if ppu.scanline >= 241 && ppu.scanline <= 261 {
		if ppu.scanline == 241 && ppu.cycle == 1 {
			ppu.setFlag(&ppu.status, statusVerticalBlank, true)
			if ppu.getFlag(&ppu.control, controlEnableNMI) != 0 {
				ppu.NMI = true
			}
		}
	}

	var bgPixel uint8 = 0x00
	var bgPalette uint8 = 0x00

	if ppu.getFlag(&ppu.mask, maskRenderBackground) != 0 {
		var bitMux uint16 = 0x8000 >> ppu.fineX

		var p0Pixel uint8 = 0
		if (ppu.shifterPatternLo & bitMux) > 0 {
			p0Pixel = 1
		}
		var p1Pixel uint8 = 0
		if (ppu.shifterPatternHi & bitMux) > 0 {
			p1Pixel = 1
		}
		bgPixel = (p1Pixel << 1) | p0Pixel

		var pal0 uint8 = 0
		if (ppu.shifterAttribLo & bitMux) > 0 {
			pal0 = 1
		}
		var pal1 uint8 = 0
		if (ppu.shifterAttribHi & bitMux) > 0 {
			pal1 = 1
		}
		bgPalette = (pal1 << 1) | pal0
	}

	ppu.screen.Set(int(ppu.cycle-1+1), int(ppu.scanline+1),
		ppu.GetColorFromPaletteRAM(bgPalette, bgPixel))

	// Draw old-fashioned static noise.
	// var pixelColor []uint8
	// if rand.Int()%2 != 0 {
	// 	pixelColor = ppu.palette[0x3F]
	// } else {
	// 	pixelColor = ppu.palette[0x30]
	// }

	// ppu.screen.Set(int(ppu.cycle-1+1), int(ppu.scanline+1),
	// color.RGBA{pixelColor[0],
	// 	pixelColor[1],
	// 	pixelColor[2],
	// 	pixelColor[3]})

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
