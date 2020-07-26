package nes

type Mapper0 struct {
	prgBanks uint8
	chrBanks uint8
}

func NewMapper0(prgBanks uint8, chrBanks uint8) Mapper {
	mapper := Mapper0{}
	mapper.prgBanks = prgBanks
	mapper.chrBanks = chrBanks
	return &mapper
}

func (mapper *Mapper0) CPUMapRead(addr uint16, mappedAddr *uint32) bool {
	var tmp uint32
	if addr >= 0x8000 && addr <= 0xFFFF {
		if mapper.prgBanks > 1 {
			tmp = uint32(addr & 0x7FFF)
			mappedAddr = &tmp
		} else {
			tmp = uint32(addr & 0x3FFF)
			mappedAddr = &tmp
		}
		return true
	}

	return false
}

func (mapper *Mapper0) CPUMapWrite(addr uint16, mappedAddr *uint32) bool {
	if addr >= 0x8000 && addr <= 0xFFFF {
		return true
	}

	return false
}

func (mapper *Mapper0) PPUMapRead(addr uint16, mappedAddr *uint32) bool {
	var tmp uint32
	if addr >= 0x0000 && addr <= 0x1FFF {
		tmp = uint32(addr)
		mappedAddr = &tmp
		return true
	}

	return false
}

func (mapper *Mapper0) PPUMapWrite(addr uint16, mappedAddr *uint32) bool {
	// if addr >= 0x0000 && addr <= 0x1FFF {
	// 	return true
	// }

	return false
}
