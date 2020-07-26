package nes

// Bus The main bus of a NES.
type Bus struct {
	CPU       *CPU
	CPURAM    []uint8
	PPU       *PPU
	cartridge *Cartridge

	systemClockCounter uint32
}

// NewBus Create a NES main bus with device attached to it.
func NewBus() *Bus {
	// Init RAM space.
	ram := make([]uint8, 64*1024)
	bus := Bus{nil, ram, nil, nil, 0}

	// Connect CPU to bus
	bus.CPU = ConnectCPU(&bus)
	return &bus
}

func (bus *Bus) CPUWrite(addr uint16, data uint8) {
	if bus.cartridge.CPUWrite(addr, data) {

	} else if addr >= 0x0000 && addr <= 0x1FFF {
		bus.CPURAM[addr&0x07FF] = data // addr&0x07FF yields back the geniune value after mirroring
	} else if addr >= 0x2000 && addr < 0x3FFF {
		bus.PPU.CPUWrite(addr&0x0007, data)
	}
}

func (bus *Bus) CPURead(addr uint16, readOnly ...bool) uint8 {
	var data uint8 = 0x00
	bReadOnly := false
	if len(readOnly) > 0 {
		bReadOnly = readOnly[0]
	}

	if bus.cartridge.CPURead(addr, &data) {

	} else if addr >= 0x0000 && addr <= 0x1FFF {
		data = bus.CPURAM[addr&0x07FF] // addr&0x07FF yields back the geniune value after mirroring
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		data = bus.PPU.CPURead(addr&0x0007, bReadOnly)
	}

	return data
}

// NES interface

func (bus *Bus) InsertCartridge(cart *Cartridge) {
	bus.cartridge = cart
	bus.PPU.ConnectCartridge(cart)
}

func (bus *Bus) Reset() {
	bus.CPU.Reset()
	bus.systemClockCounter = 0
}

func (bus *Bus) Clock() {

}
