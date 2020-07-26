package nes

type PPU struct {
	cartridge *Cartridge

	// PPU RAM
	tableName    [2][1024]uint8
	tablePalette [32]uint8
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

func (ppu *PPU) PPUWrite(addr uint16, data uint8) {
	addr &= 0x3FFF
}

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

}
