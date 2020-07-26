package nes

type Mapper interface {
	CPUMapRead(addr uint16, mappedAddr *uint32) bool
	CPUMapWrite(addr uint16, mappedAddr *uint32) bool
	PPUMapRead(addr uint16, mappedAddr *uint32) bool
	PPUMapWrite(addr uint16, mappedAddr *uint32) bool
}
