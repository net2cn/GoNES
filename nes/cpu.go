package nes

const (
	_ = iota
	interruptNone
	interruptNMI
	interruptIRQ
)

const (
	_ = iota
	modeImplied
	modeImmediate
	modeZeroPage
	modeZeroPageX
	modeZeroPageY
	modeRelative
	modeAbsolute
	modeAbsoluteX
	modeAbosluteY
	modeIndirect
	modeIndirectX
	modeIndirectY
)

const (
	flagCarryBit          = (1 << 0)
	flagZero              = (2 << 0)
	flagDisableInterrupts = (3 << 0)
	flagDecimalMode       = (4 << 0) // Redundant
	flagBreak             = (5 << 0)
	flagUnused            = (6 << 0)
	flagOverflow          = (7 << 0)
	flagNegative          = (8 << 0)
)

var instructionModes = [256]byte{
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	1, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 8, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 13, 13, 6, 3, 6, 3, 2, 2, 3, 3,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 13, 13, 6, 3, 6, 3, 2, 2, 3, 3,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
}

var instructionSizes = [256]byte{
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	3, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	1, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	1, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 0, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 0, 3, 0, 0,
	2, 2, 2, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
}

var instructionCycles = [256]byte{
	7, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 3, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 5, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 6, 2, 6, 4, 4, 4, 4, 2, 5, 2, 5, 5, 5, 5, 5,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 5, 2, 5, 4, 4, 4, 4, 2, 4, 2, 4, 4, 4, 4, 4,
	2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
}

var instructionPageCycles = [256]byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
}

var instructionNames = [256]string{
	"BRK", "ORA", "KIL", "SLO", "NOP", "ORA", "ASL", "SLO",
	"PHP", "ORA", "ASL", "ANC", "NOP", "ORA", "ASL", "SLO",
	"BPL", "ORA", "KIL", "SLO", "NOP", "ORA", "ASL", "SLO",
	"CLC", "ORA", "NOP", "SLO", "NOP", "ORA", "ASL", "SLO",
	"JSR", "AND", "KIL", "RLA", "BIT", "AND", "ROL", "RLA",
	"PLP", "AND", "ROL", "ANC", "BIT", "AND", "ROL", "RLA",
	"BMI", "AND", "KIL", "RLA", "NOP", "AND", "ROL", "RLA",
	"SEC", "AND", "NOP", "RLA", "NOP", "AND", "ROL", "RLA",
	"RTI", "EOR", "KIL", "SRE", "NOP", "EOR", "LSR", "SRE",
	"PHA", "EOR", "LSR", "ALR", "JMP", "EOR", "LSR", "SRE",
	"BVC", "EOR", "KIL", "SRE", "NOP", "EOR", "LSR", "SRE",
	"CLI", "EOR", "NOP", "SRE", "NOP", "EOR", "LSR", "SRE",
	"RTS", "ADC", "KIL", "RRA", "NOP", "ADC", "ROR", "RRA",
	"PLA", "ADC", "ROR", "ARR", "JMP", "ADC", "ROR", "RRA",
	"BVS", "ADC", "KIL", "RRA", "NOP", "ADC", "ROR", "RRA",
	"SEI", "ADC", "NOP", "RRA", "NOP", "ADC", "ROR", "RRA",
	"NOP", "STA", "NOP", "SAX", "STY", "STA", "STX", "SAX",
	"DEY", "NOP", "TXA", "XAA", "STY", "STA", "STX", "SAX",
	"BCC", "STA", "KIL", "AHX", "STY", "STA", "STX", "SAX",
	"TYA", "STA", "TXS", "TAS", "SHY", "STA", "SHX", "AHX",
	"LDY", "LDA", "LDX", "LAX", "LDY", "LDA", "LDX", "LAX",
	"TAY", "LDA", "TAX", "LAX", "LDY", "LDA", "LDX", "LAX",
	"BCS", "LDA", "KIL", "LAX", "LDY", "LDA", "LDX", "LAX",
	"CLV", "LDA", "TSX", "LAS", "LDY", "LDA", "LDX", "LAX",
	"CPY", "CMP", "NOP", "DCP", "CPY", "CMP", "DEC", "DCP",
	"INY", "CMP", "DEX", "AXS", "CPY", "CMP", "DEC", "DCP",
	"BNE", "CMP", "KIL", "DCP", "NOP", "CMP", "DEC", "DCP",
	"CLD", "CMP", "NOP", "DCP", "NOP", "CMP", "DEC", "DCP",
	"CPX", "SBC", "NOP", "ISC", "CPX", "SBC", "INC", "ISC",
	"INX", "SBC", "NOP", "SBC", "CPX", "SBC", "INC", "ISC",
	"BEQ", "SBC", "KIL", "ISC", "NOP", "SBC", "INC", "ISC",
	"SED", "SBC", "NOP", "ISC", "NOP", "SBC", "INC", "ISC",
}

type instruction struct {
	name    string
	operate uint8
	addr    uint8
	cycles  uint8
}

type CPU struct {
	Bus *Bus

	A      uint8  // Accumulator register
	X      uint8  // X register
	Y      uint8  // Y register
	SP     uint8  // Stack pointer
	PC     uint16 // Program pointer
	Status uint8  // Status register

	Fetched uint8

	AddrAbs uint16
	AddrRel uint16
	Opcode  uint8
	Cycles  uint8
}

func ConnectCPU(bus *Bus) *CPU {
	cpu := CPU{}
	cpu.Bus = bus
	return &cpu
}

// Read reads one byte from bus and return a word value.
func (cpu *CPU) Read(addr uint16) uint8 {
	return cpu.Bus.Read(addr)
}

func (cpu *CPU) Write(addr uint16, data uint8) {
	cpu.Bus.Write(addr, data)
}

func (cpu *CPU) GetFlag(f uint8) uint8 {
	if (cpu.Status & f) > 0 {
		return 1
	}

	return 0
}

func (cpu *CPU) SetFlag(f uint8, v bool) {
	if v {
		cpu.Status |= f
	} else {
		cpu.Status &= ^f
	}
}

func (cpu *CPU) Clock() {
	if cpu.Cycles == 0 {
		cpu.Opcode = cpu.Read(cpu.PC)

		cpu.SetFlag(flagUnused, true)

		cpu.PC++

		// Get starting nmber of cycles
		cpu.Cycles = instructionCycles[cpu.Opcode]

		mode := instructionModes[cpu.Opcode]

		switch mode {
		case modeImplied:
			cpu.Fetched = cpu.A
		case modeImmediate:
			cpu.PC++
			cpu.AddrAbs = cpu.PC
		case modeZeroPage:
			cpu.AddrAbs = uint16(cpu.Read(cpu.PC))
			cpu.PC++
			cpu.AddrAbs &= 0x00FF
		case modeZeroPageX:
			cpu.AddrAbs = uint16(cpu.Read(cpu.PC) + cpu.X)
			cpu.PC++
			cpu.AddrAbs &= 0xFF
		case modeZeroPageY:
			cpu.AddrAbs = uint16(cpu.Read(cpu.PC) + cpu.Y)
			cpu.PC++
			cpu.AddrAbs &= 0x00FF
		case modeRelative:
			cpu.AddrRel = uint16(cpu.Read(cpu.PC))
			cpu.PC++
			if (cpu.AddrRel & 0x80) != 0 {
				cpu.AddrRel |= 0xFF00
			}
		case modeAbsolute:
			var lo uint16 = uint16(cpu.Read(cpu.PC))
			cpu.PC++
			var hi uint16 = uint16(cpu.Read(cpu.PC))
			cpu.PC++

			cpu.AddrAbs = uint16((hi << 8) | lo)
		case modeAbsoluteX:
			var lo uint16 = uint16(cpu.Read(cpu.PC))
			cpu.PC++
			var hi uint16 = uint16(cpu.Read(cpu.PC))
			cpu.PC++

			cpu.AddrAbs = uint16((hi << 8) | lo)
			cpu.AddrAbs += uint16(cpu.X)

			if (cpu.AddrAbs & 0xFF00) != uint16((hi << 8)) {
				cpu.Cycles++
			}
		case modeAbosluteY:
			var lo uint16 = uint16(cpu.Read(cpu.PC))
			cpu.PC++
			var hi uint16 = uint16(cpu.Read(cpu.PC))
			cpu.PC++

			cpu.AddrAbs = uint16((hi << 8) | lo)
			cpu.AddrAbs += uint16(cpu.Y)

			if (cpu.AddrAbs & 0xFF00) != uint16((hi << 8)) {
				cpu.Cycles++
			}
		case modeIndirect:
			var ptrLo uint16 = uint16(cpu.Read(cpu.PC))
			cpu.PC++
			var ptrHi uint16 = uint16(cpu.Read(cpu.PC))
			cpu.PC++

			var ptr uint16 = uint16((ptrHi << 8) | ptrLo)

			if ptrLo == 0x00FF {
				cpu.AddrAbs = uint16((cpu.Read(ptr&0xFF00) << 8) | cpu.Read(ptr+0))
			} else {
				cpu.AddrAbs = uint16((cpu.Read(ptr+1) << 8) | cpu.Read(ptr+0))
			}
		case modeIndirectX:
			var t uint16 = uint16(cpu.Read(cpu.PC))
			cpu.PC++

			var lo uint16 = uint16(cpu.Read((t + uint16(cpu.X)) & 0x00FF))
			var hi uint16 = uint16(cpu.Read((t + uint16(cpu.X) + 1) & 0x00FF))

			cpu.AddrAbs = (hi << 8) | lo
		case modeIndirectY:
			var t uint16 = uint16(cpu.Read(cpu.PC))
			cpu.PC++

			var lo uint16 = uint16(cpu.Read(t & 0x00FF))
			var hi uint16 = uint16(cpu.Read((t + 1) & 0x00FF))

			cpu.AddrAbs = (hi << 8) | lo
			cpu.AddrAbs += uint16(cpu.Y)

			if (cpu.AddrAbs & 0xFF00) != (hi << 8) {
				cpu.Cycles++
			}
		}

		// TODO: Implement Operation
	}
	cpu.Cycles--
}

func (cpu *CPU) Fetch() uint8 {
	var fetched uint8 = 0
	if instructionModes[cpu.Opcode] != modeImplied {
		fetched = cpu.Read(cpu.AddrAbs)
	}
	return fetched
}
