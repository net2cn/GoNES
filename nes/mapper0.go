package nes

// Mapper0 Mapper0 struct
type Mapper0 struct {
	prgBanks uint8
	chrBanks uint8
}

// NewMapper0 Creates a new mapper of mapper0.
func NewMapper0(prgBanks uint8, chrBanks uint8) Mapper {
	mapper := Mapper0{}
	mapper.prgBanks = prgBanks
	mapper.chrBanks = chrBanks
	return &mapper
}

// CPUMapRead Mapper0's CPUMapRead implementation.
func (mapper *Mapper0) CPUMapRead(addr uint16, mappedAddr *uint32) bool {
	if addr >= 0x8000 && addr <= 0xFFFF {
		if mapper.prgBanks > 1 {
			*mappedAddr = uint32(addr & 0x7FFF)
		} else {
			*mappedAddr = uint32(addr & 0x3FFF)
		}
		return true
	}

	return false
}

// CPUMapWrite Mapper0's CPUMapWrite implementation.
func (mapper *Mapper0) CPUMapWrite(addr uint16, mappedAddr *uint32) bool {
	if addr >= 0x8000 && addr <= 0xFFFF {
		if mapper.prgBanks > 1 {
			*mappedAddr = uint32(addr & 0x7FFF)
		} else {
			*mappedAddr = uint32(addr & 0x3FFF)
		}

		return true
	}

	return false
}

// PPUMapRead Mapper0's PPUMapRead implementation.
func (mapper *Mapper0) PPUMapRead(addr uint16, mappedAddr *uint32) bool {
	if addr >= 0x0000 && addr <= 0x1FFF {
		*mappedAddr = uint32(addr)
		return true
	}

	return false
}

// PPUMapWrite Mapper0's PPUMapWrite implementation.
func (mapper *Mapper0) PPUMapWrite(addr uint16, mappedAddr *uint32) bool {
	if addr >= 0x0000 && addr <= 0x1FFF {
		if mapper.chrBanks == 0 {
			*mappedAddr = uint32(addr)
			return true
		}
	}

	return false
}
