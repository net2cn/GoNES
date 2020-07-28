package nes

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type Cartridge struct {
	prgMemory []uint8
	chrMemory []uint8

	mapperID uint8
	prgBanks uint8
	chrBanks uint8

	mapper Mapper
}

type cartridgeHeader struct {
	Name         [4]byte
	PRGROMChunks byte
	CHRROMChunks byte
	Mapper1      byte
	Mapper2      byte
	PRGRAMSize   byte

	// Unused
	TVSystem1 byte
	TVSystem2 byte
	Unused    [5]byte
}

func NewCartridge(filePath string) (*Cartridge, error) {
	cart := Cartridge{}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to load cartridge: %s\n", err)
		return nil, err
	}
	defer file.Close()

	header := cartridgeHeader{}
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		fmt.Printf("Failed to to read file header: %s\n", err)
		return nil, err
	}

	if header.Name != [4]byte{0x4E, 0x45, 0x53,0x1A}{
		fmt.Printf("Failed to load cartridge: This is not a valid file.")
		return nil,err
	}

	// Unused.
	if header.Mapper1&0x04 != 0 {
		trainer := make([]byte, 512)
		if _, err := io.ReadFull(file, trainer); err != nil {
			fmt.Printf("Failed to read cartridge: %s\n", err)
			return nil, err
		}
	}

	// Determine mapper ID
	cart.mapperID = ((header.Mapper2 >> 4) << 4) | (header.Mapper1 >> 4)

	var fileType uint8 = 1
	switch fileType {
	case 0:
	case 1:
		cart.prgBanks = header.PRGROMChunks
		cart.prgMemory = make([]uint8, int(cart.prgBanks)*16384)
		if _, err := io.ReadFull(file, cart.prgMemory); err != nil {
			fmt.Printf("Failed to read PRG banks: %s\n", err)
			return nil, err
		}

		cart.chrBanks = header.CHRROMChunks
		cart.chrMemory = make([]uint8, int(cart.chrBanks)*8192)
		if _, err := io.ReadFull(file, cart.chrMemory); err != nil {
			fmt.Printf("Failed to read CHR banks: %s\n", err)
			return nil, err
		}
	case 2:
	}
	// Close file once done.
	file.Close()

	switch cart.mapperID {
	case 0:
		cart.mapper = NewMapper0(cart.prgBanks, cart.chrBanks)
	}

	return &cart, nil
}

// CPU IO

func (cart *Cartridge) CPURead(addr uint16, data *uint8) bool {
	var mappedAddr uint32 = 0
	if cart.mapper.CPUMapRead(addr, &mappedAddr) {
		*data = cart.prgMemory[mappedAddr]
		return true
	}

	return false
}

func (cart *Cartridge) CPUWrite(addr uint16, data uint8) bool {
	var mappedAddr uint32 = 0
	if cart.mapper.CPUMapWrite(addr, &mappedAddr) {
		cart.prgMemory[mappedAddr] = data
		return true
	}

	return false
}

// PPU IO

func (cart *Cartridge) PPURead(addr uint16, data *uint8) bool {
	var mappedAddr uint32 = 0
	if cart.mapper.PPUMapRead(addr, &mappedAddr) {
		*data = cart.chrMemory[mappedAddr]
		return true
	}

	return false
}

func (cart *Cartridge) PPUWrite(addr uint16, data uint8) bool {
	var mappedAddr uint32 = 0
	if cart.mapper.PPUMapWrite(addr, &mappedAddr) {
		cart.chrMemory[mappedAddr] = data
		return true
	}

	return false
}
