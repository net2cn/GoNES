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
}

type cartridgeHeader struct {
	name         [4]byte
	prgROMChunks uint8
	chrROMChunks uint8
	mapper1      uint8
	mapper2      uint8
	prgRAMSize   uint8
	tvSystem1    uint8
	tvSystem2    uint8
	unused       [5]byte
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

	// Unused.
	if header.mapper1&0x04 != 0 {
		trainer := make([]byte, 512)
		if _, err := io.ReadFull(file, trainer); err != nil {
			fmt.Printf("Failed to read cartridge: %s\n", err)
			return nil, err
		}
	}

	// Determine mapper ID
	cart.mapperID = ((header.mapper2 >> 4) << 4) | (header.mapper1 >> 4)

	var fileType uint8 = 1
	switch fileType {
	case 0:
	case 1:
		cart.prgBanks = header.prgROMChunks
		cart.prgMemory = make([]uint8, int(cart.prgBanks)*16384)
		if _, err := io.ReadFull(file, cart.prgMemory); err != nil {
			fmt.Printf("Failed to read PRG banks: %s\n", err)
			return nil, err
		}

		cart.chrBanks = header.chrROMChunks
		cart.chrMemory = make([]uint8, int(cart.chrBanks)*8192)
		if _, err := io.ReadFull(file, cart.chrMemory); err != nil {
			fmt.Printf("Failed to read CHR banks: %s\n", err)
			return nil, err
		}
	case 2:
	}
	// Close file once done.
	file.Close()
	return &cart, nil
}

func (cart *Cartridge) CPUWrite(addr uint16, data uint8) bool {
	return false

}

func (cart *Cartridge) CPURead(addr uint16, data *uint8) bool {
	return false
}

func (cart *Cartridge) PPUWrite(addr uint16, data uint8) bool {
	return false

}

func (cart *Cartridge) PPURead(addr uint16, data *uint8) bool {
	return false
}
