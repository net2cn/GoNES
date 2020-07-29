package nes

// Bus The main bus of a NES.
type Bus struct {
	CPU        *CPU
	CPURAM     []uint8
	PPU        *PPU
	cartridge  *Cartridge
	Controller []uint8

	controllerState []uint8

	systemClockCounter uint32
}

// NewBus Create a NES main bus with device attached to it.
func NewBus() *Bus {
	// Init RAM space.
	ram := make([]uint8, 64*1024)
	bus := Bus{CPU: nil,
		CPURAM:          ram,
		PPU:             nil,
		cartridge:       nil,
		Controller:      make([]uint8, 2),
		controllerState: make([]uint8, 2)}

	// Connect CPU to bus
	bus.CPU = ConnectCPU(&bus)
	bus.PPU = ConnectPPU(&bus)
	return &bus
}

// CPU IO

// CPURead Allow CPU read from bus.
func (bus *Bus) CPURead(addr uint16, readOnly ...bool) uint8 {
	var data uint8 = 0x00
	bReadOnly := false
	if len(readOnly) > 0 {
		bReadOnly = readOnly[0]
	}

	if bus.cartridge.CPURead(addr, &data) {
		// Cartridge address range
	} else if addr >= 0x0000 && addr <= 0x1FFF {
		data = bus.CPURAM[addr&0x07FF] // addr&0x07FF yields back the geniune value after mirroring
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		data = bus.PPU.CPURead(addr&0x0007, bReadOnly)
	} else if addr >= 0x4016 && addr <= 0x4017 {
		data = 0
		if (bus.controllerState[addr&0x0001] & 0x80) > 0 {
			data = 1
		}
		bus.controllerState[addr&0x0001] <<= 1
	}

	return data
}

// CPUWrite Allow CPU write to bus.
func (bus *Bus) CPUWrite(addr uint16, data uint8) {
	if bus.cartridge.CPUWrite(addr, data) {

	} else if addr >= 0x0000 && addr <= 0x1FFF {
		bus.CPURAM[addr&0x07FF] = data // addr&0x07FF yields back the geniune value after mirroring
	} else if addr >= 0x2000 && addr < 0x3FFF {
		bus.PPU.CPUWrite(addr&0x0007, data)
	} else if addr >= 0x4016 && addr <= 0x4017 {
		bus.controllerState[addr&0x0001] = bus.Controller[addr&0x0001]
	}
}

// NES interface

// InsertCartridge Connects game cartridge to bus.
func (bus *Bus) InsertCartridge(cart *Cartridge) {
	bus.cartridge = cart
	bus.PPU.ConnectCartridge(cart)
}

// Reset Reset whole bus and the devices attached to it.
func (bus *Bus) Reset() {
	bus.CPU.Reset()
	bus.systemClockCounter = 0
}

// Clock Clock bus once.
func (bus *Bus) Clock() {
	bus.PPU.Clock()

	if bus.systemClockCounter%3 == 0 {
		bus.CPU.Clock()
	}

	if bus.PPU.NMI {
		bus.PPU.NMI = false
		bus.CPU.NMI()
	}
	bus.systemClockCounter++
}
