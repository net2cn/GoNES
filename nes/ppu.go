package nes

import (
	"fmt"
	"image/color"

	"github.com/veandco/go-sdl2/sdl"
)

// Bitmask for status register.
const (
	statusUnused         = 0x1F // Mask the unused register.
	statusSpriteOverflow = (1 << 5)
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

const (
	entryY = iota
	entryID
	entryAttribute
	entryX
)

// PPU Nintendo 2C02 PPU struct
type PPU struct {
	cartridge *Cartridge

	// PPU RAM
	TableName    [2][1024]uint8
	tablePattern [2][4096]uint8
	tablePalette [32]uint8
	OAM          [256]uint8

	sprite                 [32]uint8
	spriteCount            uint8
	spriteShifterPatternLo [8]uint8
	spriteShifterPatternHi [8]uint8

	NMI bool

	palette            [64]color.RGBA
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
	oamAddr  uint8

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
	ppu.screen, err = sdl.CreateRGBSurfaceWithFormat(0, 256, 240, 8, sdl.PIXELFORMAT_RGB888)
	if err != nil {
		fmt.Printf("Failed to create sprite: %s\n", err)
		panic(err)
	}

	// Initialize PPU palette
	ppu.initializePalette()

	// Initialize PPU sprite name table
	ppu.spriteNameTable = make([]*sdl.Surface, 2)
	for i := range ppu.spriteNameTable {
		ppu.spriteNameTable[i], err = sdl.CreateRGBSurfaceWithFormat(0, 256, 240, 8, sdl.PIXELFORMAT_RGB888)
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

func (ppu *PPU) initializePalette() {
	colors := []uint32{
		0x666666, 0x002A88, 0x1412A7, 0x3B00A4, 0x5C007E, 0x6E0040, 0x6C0600, 0x561D00,
		0x333500, 0x0B4800, 0x005200, 0x004F08, 0x00404D, 0x000000, 0x000000, 0x000000,
		0xADADAD, 0x155FD9, 0x4240FF, 0x7527FE, 0xA01ACC, 0xB71E7B, 0xB53120, 0x994E00,
		0x6B6D00, 0x388700, 0x0C9300, 0x008F32, 0x007C8D, 0x000000, 0x000000, 0x000000,
		0xFFFEFF, 0x64B0FF, 0x9290FF, 0xC676FF, 0xF36AFF, 0xFE6ECC, 0xFE8170, 0xEA9E22,
		0xBCBE00, 0x88D800, 0x5CE430, 0x45E082, 0x48CDDE, 0x4F4F4F, 0x000000, 0x000000,
		0xFFFEFF, 0xC0DFFF, 0xD3D2FF, 0xE8C8FF, 0xFBC2FF, 0xFEC4EA, 0xFECCC5, 0xF7D8A5,
		0xE4E594, 0xCFEF96, 0xBDF4AB, 0xB3F3CC, 0xB5EBF2, 0xB8B8B8, 0x000000, 0x000000,
	}
	for i, c := range colors {
		r := byte(c >> 16)
		g := byte(c >> 8)
		b := byte(c)
		ppu.palette[i] = color.RGBA{r, g, b, 0xFF}
	}
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

// DO NOT USE!
// func (ppu *PPU) setLoppyRegister(addr *uint16, m int, val uint8) {
// 	*addr |= uint16(val) * uint16(m&-m)
// }

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
		data = ppu.OAM[ppu.oamAddr]
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
		ppu.tramAddr = (ppu.tramAddr & 0xF3FF) | ((uint16(data) & 0x03) << 10)
	case 0x0001: // Mask
		ppu.mask = data
	case 0x0002: // Status

	case 0x0003: // OAM address
		ppu.oamAddr = data
	case 0x0004: // OAM data
		ppu.OAM[ppu.oamAddr] = data
	case 0x0005: // Scroll
		if ppu.addressLatch == 0 {
			ppu.fineX = data & 0x07
			ppu.tramAddr = (ppu.tramAddr & 0xFFE0) | (uint16(data) >> 3)
			ppu.addressLatch = 1
		} else {
			ppu.tramAddr = (ppu.tramAddr & 0x8FFF) | ((uint16(data) & 0x07) << 12)
			ppu.tramAddr = (ppu.tramAddr & 0xFC1F) | ((uint16(data) & 0xF8) << 2)
			ppu.addressLatch = 0
		}
	case 0x0006: // PPU address
		if ppu.addressLatch == 0 {
			ppu.tramAddr = (ppu.tramAddr & 0x80FF) | ((uint16(data) & 0x3F) << 8)
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
			// Vertical
			if addr >= 0x0000 && addr <= 0x03FF {
				data = ppu.TableName[0][addr&0x03FF]
			}
			if addr >= 0x0400 && addr <= 0x07FF {
				data = ppu.TableName[1][addr&0x03FF]
			}
			if addr >= 0x0800 && addr <= 0x0BFF {
				data = ppu.TableName[0][addr&0x03FF]
			}
			if addr >= 0x0C00 && addr <= 0x0FFF {
				data = ppu.TableName[1][addr&0x03FF]
			}
		} else if ppu.cartridge.mirror == mirrorHorizontal {
			// Horizontal
			if addr >= 0x0000 && addr <= 0x03FF {
				data = ppu.TableName[0][addr&0x03FF]
			}
			if addr >= 0x0400 && addr <= 0x07FF {
				data = ppu.TableName[0][addr&0x03FF]
			}
			if addr >= 0x0800 && addr <= 0x0BFF {
				data = ppu.TableName[1][addr&0x03FF]
			}
			if addr >= 0x0C00 && addr <= 0x0FFF {
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
			// Vertical
			if addr >= 0x0000 && addr <= 0x03FF {
				ppu.TableName[0][addr&0x03FF] = data
			}
			if addr >= 0x0400 && addr <= 0x07FF {
				ppu.TableName[1][addr&0x03FF] = data
			}
			if addr >= 0x0800 && addr <= 0x0BFF {
				ppu.TableName[0][addr&0x03FF] = data
			}
			if addr >= 0x0C00 && addr <= 0x0FFF {
				ppu.TableName[1][addr&0x03FF] = data
			}
		} else if ppu.cartridge.mirror == mirrorHorizontal {
			// Horizontal
			if addr >= 0x0000 && addr <= 0x03FF {
				ppu.TableName[0][addr&0x03FF] = data
			}
			if addr >= 0x0400 && addr <= 0x07FF {
				ppu.TableName[0][addr&0x03FF] = data
			}
			if addr >= 0x0800 && addr <= 0x0BFF {
				ppu.TableName[1][addr&0x03FF] = data
			}
			if addr >= 0x0C00 && addr <= 0x0FFF {
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
	// TODO: Refactor these methods.
	var incrementScrollX func() = func() {
		if (ppu.getFlag(&ppu.mask, maskRenderBackground) != 0) ||
			(ppu.getFlag(&ppu.mask, maskRenderSprites) != 0) {
			if ppu.vramAddr&0x001F == 31 {
				ppu.vramAddr &= 0xFFE0
				ppu.vramAddr ^= 0x0400
			} else {
				ppu.vramAddr++
			}
		}
	}

	var incrementScrollY func() = func() {
		if (ppu.getFlag(&ppu.mask, maskRenderBackground) != 0) ||
			(ppu.getFlag(&ppu.mask, maskRenderSprites) != 0) {
			if ppu.vramAddr&0x7000 != 0x7000 {
				ppu.vramAddr += 0x1000
			} else {
				ppu.vramAddr &= 0x8FFF

				y := (ppu.vramAddr & 0x03E0) >> 5
				if y == 29 {
					y = 0
					ppu.vramAddr ^= 0x0800
				} else if y == 31 {
					y = 0
				} else {
					y++
				}
				ppu.vramAddr = (ppu.vramAddr & 0xFC1F) | (y << 5)
			}
		}
	}

	var transferAddressX func() = func() {
		if (ppu.getFlag(&ppu.mask, maskRenderBackground) != 0) ||
			(ppu.getFlag(&ppu.mask, maskRenderSprites) != 0) {
			ppu.vramAddr = (ppu.vramAddr & 0xFBE0) | (ppu.tramAddr & 0x041F)
		}
	}

	var transferAddressY func() = func() {
		if (ppu.getFlag(&ppu.mask, maskRenderBackground) != 0) ||
			(ppu.getFlag(&ppu.mask, maskRenderSprites) != 0) {
			ppu.vramAddr = (ppu.vramAddr & 0x841F) | (ppu.tramAddr & 0x7BE0)
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

		if ppu.getFlag(&ppu.mask, maskRenderSprites) != 0 && ppu.cycle >= 1 && ppu.cycle < 258 {
			for i := uint8(0); i < ppu.spriteCount; i++ {
				if ppu.sprite[i*4+3] > 0 {
					ppu.sprite[i*4+3]--
				} else {
					ppu.spriteShifterPatternLo[i] <<= 1
					ppu.spriteShifterPatternHi[i] <<= 1
				}
			}
		}
	}

	if ppu.scanline >= -1 && ppu.scanline < 240 {
		if ppu.scanline == 0 && ppu.cycle == 0 {
			ppu.cycle = 1
		}

		if ppu.scanline == -1 && ppu.cycle == 1 {
			ppu.setFlag(&ppu.status, statusVerticalBlank, false)
			ppu.setFlag(&ppu.status, statusSpriteOverflow, false)

			for i := 0; i < 8; i++ {
				ppu.spriteShifterPatternLo[i] = 0
				ppu.spriteShifterPatternHi[i] = 0
			}
		}

		if (ppu.cycle >= 2 && ppu.cycle <= 258) || (ppu.cycle >= 321 && ppu.cycle < 338) {
			updateShifters()

			switch (ppu.cycle - 1) % 8 {
			case 0:
				// Fetch next name table.
				loadBackgroundShifters()
				ppu.nextTileID = ppu.PPURead(0x2000 | (ppu.vramAddr & 0x0FFF))
			case 2:
				// Fetch next attribute table.
				addr := (0x23C0 | (ppu.getLoppyRegister(&ppu.vramAddr, loppyNameTableY) << 11) |
					(ppu.getLoppyRegister(&ppu.vramAddr, loppyNameTableX) << 10) |
					((ppu.getLoppyRegister(&ppu.vramAddr, loppyCoarseY) >> 2) << 3) |
					(ppu.getLoppyRegister(&ppu.vramAddr, loppyCoarseX) >> 2))
				ppu.nextTileAttr = ppu.PPURead(addr)
				if ppu.getLoppyRegister(&ppu.vramAddr, loppyCoarseY)&0x02 != 0 {
					ppu.nextTileAttr >>= 4
				}
				if ppu.getLoppyRegister(&ppu.vramAddr, loppyCoarseX)&0x02 != 0 {
					ppu.nextTileAttr >>= 2
				}
				ppu.nextTileAttr &= 0x03
			case 4:
				// Fetch LSB
				ppu.nextTileLSB = ppu.PPURead((uint16(ppu.getFlag(&ppu.control, controlPatternBackground)) << 12) +
					(uint16(ppu.nextTileID) << 4) +
					ppu.getLoppyRegister(&ppu.vramAddr, loppyFineY) + 0)
			case 6:
				// Fetch MSB
				ppu.nextTileMSB = ppu.PPURead((uint16(ppu.getFlag(&ppu.control, controlPatternBackground)) << 12) +
					(uint16(ppu.nextTileID) << 4) +
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

		if ppu.scanline == -1 && ppu.cycle >= 280 && ppu.cycle < 305 {
			transferAddressY()
		}

		if ppu.cycle == 338 || ppu.cycle == 340 {
			ppu.nextTileID = ppu.PPURead(0x2000 | (ppu.vramAddr & 0x0FFF))
		}

		// Sprite Evaluation
		if ppu.cycle == 257 && ppu.scanline >= 0 {
			// memset 0xFF
			ppu.sprite[0] = 0xFF
			for i := 1; i < len(ppu.sprite); i *= 2 {
				copy(ppu.sprite[i:], ppu.sprite[:i])
			}
			ppu.spriteCount = 0

			for i := 0; i < 8; i++ {
				ppu.spriteShifterPatternLo[i] = 0
				ppu.spriteShifterPatternHi[i] = 0
			}

			count := 0
			for count < 64 && ppu.spriteCount < 9 {
				var diff int16 = (int16(ppu.scanline) - int16(ppu.OAM[count*4+entryY]))
				var spriteSize int16 = 16

				if ppu.getFlag(&ppu.control, controlSpriteSize) == 0 {
					spriteSize = 8
				}
				if diff >= 0 && diff < spriteSize {
					if ppu.spriteCount < 8 {
						if count == 0 {

						}
						copy(ppu.sprite[ppu.spriteCount*4:ppu.spriteCount*4+4], ppu.OAM[count*4:count*4+4])
						ppu.spriteCount++
					}
				}
				count++
			}
			if ppu.spriteCount > 8 {
				ppu.setFlag(&ppu.status, statusSpriteOverflow, true)
			}
		}

		if ppu.cycle == 340 {
			for i := uint8(0); i < ppu.spriteCount; i++ {
				var spritePatternBitsLo, spritePatternBitsHi uint8
				var spritePatternAddrLo, spritePatternAddrHi uint16

				if ppu.getFlag(&ppu.control, controlSpriteSize) == 0 {
					// 8x8
					if (ppu.sprite[i*4+entryAttribute] & 0x80) == 0 {
						// Normal
						spritePatternAddrLo =
							(uint16(ppu.getFlag(&ppu.control, controlPatternSprite)) << 12) |
								(uint16(ppu.sprite[i*4+entryID]) << 4) |
								(uint16(ppu.scanline) - uint16(ppu.sprite[i*4+entryY]))

					} else {
						// Flipped
						spritePatternAddrLo =
							(uint16(ppu.getFlag(&ppu.control, controlPatternSprite)) << 12) |
								((uint16(ppu.sprite[i*4+entryID]) + 1) << 4) |
								(7 - uint16(ppu.scanline) - uint16(ppu.sprite[i*4+entryY]))
					}
				} else {
					// 8x16
					if (ppu.sprite[i*4+entryAttribute] & 0x80) == 0 {
						// Normal
						if (uint8(ppu.scanline) - ppu.sprite[i*4+entryY]) < 8 {
							// Top half tile
							spritePatternAddrLo =
								(uint16(ppu.sprite[i*4+entryID]&0x01) << 12) |
									(uint16(ppu.sprite[i*4+entryID]) & 0xFE << 4) |
									(uint16(uint16(ppu.scanline)-uint16(ppu.sprite[i*4+entryY])) & 0x07)
						} else {
							// Bottom half tile
							spritePatternAddrLo =
								(uint16(ppu.sprite[i*4+entryID]&0x01) << 12) |
									((uint16(ppu.sprite[i*4+entryID])&0xFE + 1) << 4) |
									(uint16(uint16(ppu.scanline)-uint16(ppu.sprite[i*4+entryY])) & 0x07)
						}
					} else {
						// Flipped
						if (uint8(ppu.scanline) - ppu.sprite[i*4+0]) < 8 {
							// Top half tile
							spritePatternAddrLo =
								(uint16(ppu.sprite[i*4+entryID]&0x01) << 12) |
									((uint16(ppu.sprite[i*4+entryID]&0xFE) + 1) << 4) |
									((7 - uint16(ppu.scanline) - uint16(ppu.sprite[i*4+entryY])) & 0x07)
						} else {
							// Bottom half tile
							spritePatternAddrLo =
								(uint16(ppu.sprite[i*4+entryID]&0x01) << 12) |
									(uint16(ppu.sprite[i*4+entryID]&0xFE) << 4) |
									((7 - uint16(ppu.scanline) - uint16(ppu.sprite[i*4+entryY])) & 0x07)
						}
					}
				}

				spritePatternAddrHi = spritePatternAddrLo + 8

				spritePatternBitsLo = ppu.PPURead(spritePatternAddrLo)
				spritePatternBitsHi = ppu.PPURead(spritePatternAddrHi)

				if ppu.sprite[i*4+entryAttribute]&0x40 != 0 {
					var flipByte func(b uint8) uint8 = func(b uint8) uint8 {
						b = (b&0xF0)>>4 | (b&0x0F)<<4
						b = (b&0xCC)>>2 | (b&0x33)<<2
						b = (b&0xAA)>>1 | (b&0x55)<<1
						return b
					}

					spritePatternBitsLo = flipByte(spritePatternBitsLo)
					spritePatternBitsHi = flipByte(spritePatternBitsHi)
				}

				ppu.spriteShifterPatternLo[i] = spritePatternBitsLo
				ppu.spriteShifterPatternHi[i] = spritePatternBitsHi
			}
		}
	}

	if ppu.scanline == 240 {
		// Placeholder
	}

	if ppu.scanline >= 241 && ppu.scanline < 261 {
		if ppu.scanline == 241 && ppu.cycle == 1 {
			ppu.setFlag(&ppu.status, statusVerticalBlank, true)
			if ppu.getFlag(&ppu.control, controlEnableNMI) != 0 {
				ppu.NMI = true
			}
		}
	}

	// Render screen
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

	var fgPixel uint8 = 0x00
	var fgPalette uint8 = 0x00
	var fgPriority uint8 = 0x00

	if ppu.getFlag(&ppu.mask, maskRenderSprites) != 0 {
		for i := uint8(0); i < ppu.spriteCount; i++ {
			if ppu.sprite[i*4+entryX] == 0 {
				var fgPixelLo uint8 = 0
				if ppu.spriteShifterPatternLo[i]&0x80 > 0 {
					fgPixelLo = 1
				}
				var fgPixelHi uint8 = 0
				if ppu.spriteShifterPatternHi[i]&0x80 > 0 {
					fgPixelHi = 1
				}
				fgPixel = (fgPixelHi << 1) | fgPixelLo

				fgPalette = (ppu.sprite[i*4+entryAttribute] & 0x03) + 0x04
				if (ppu.sprite[i*4+entryAttribute] & 0x20) == 0 {
					fgPriority = 1
				}

				if fgPixel != 0 {
					break
				}
			}
		}
	}

	var pixel uint8 = 0x00
	var palette uint8 = 0x00

	if bgPixel == 0 && fgPixel == 0 {
		pixel = 0x00
		palette = 0x00
	} else if bgPixel == 0 && fgPixel > 0 {
		pixel = fgPixel
		palette = fgPalette
	} else if bgPixel > 0 && fgPixel == 0 {
		pixel = bgPixel
		palette = bgPalette
	} else if bgPixel > 0 && fgPixel > 0 {
		if fgPriority != 0 {
			pixel = fgPixel
			palette = fgPalette
		} else {
			pixel = bgPixel
			palette = bgPalette
		}
	}

	if ppu.cycle >= 0 && ppu.cycle < 256 && ppu.scanline >= 0 && ppu.scanline < 240 {
		ppu.screen.Set(int(ppu.cycle-1), int(ppu.scanline),
			ppu.GetColorFromPaletteRAM(palette, pixel))
	}

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
	return ppu.palette[ppu.PPURead(0x3F00+(uint16(palette)<<2)+uint16(pixel))&0x3F]
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
					var pixel uint8 = ((tileLSB & 0x01) << 1) | (tileMSB & 0x01)
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
