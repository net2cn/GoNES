package nes

type Bus struct {
	CPU *CPU
	RAM []uint8
}

func NewBus() *Bus {
	// Init RAM space.
	ram := make([]uint8, 64*128)
	bus := Bus{nil, ram}
	bus.CPU = ConnectCPU(&bus)
	return &bus
}

func (bus *Bus) Write(addr uint16, data uint8) {
	if addr >= 0x0000 && addr <= 0xFFFF {
		bus.RAM[addr] = data
	}
}

func (bus *Bus) Read(addr uint16) uint8 {
	if addr >= 0x0000 && addr <= 0xFFFF {
		return bus.RAM[addr]
	}

	return 0x00
}
