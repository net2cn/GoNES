package nes

import (
	"strings"
)

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
	_ = iota
	flagCarryBit
	flagZero
	flagDisableInterrupts
	flagDecimalMode // Redundant
	flagBreak
	flagUnused
	flagOverflow
	flagNegative
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

type Instruction struct {
	name    string
	operate uint8
	addr    uint8
	cycles  uint8
}

type CPU struct {
	// Public
	Bus *Bus

	A      uint8  // Accumulator register
	X      uint8  // X register
	Y      uint8  // Y register
	SP     uint8  // Stack pointer
	PC     uint16 // Program pointer
	Status uint8  // Status register

	// Private
	fetched    uint8
	temp       uint16
	addrAbs    uint16
	addrRel    uint16
	opcode     uint8
	cycles     uint8
	clockCount uint32

	table [256]func() uint8
}

func (cpu *CPU) createTable() {
	cpu.table = [256]func() uint8{
		cpu.brk, cpu.ora, cpu.kil, cpu.slo, cpu.nop, cpu.ora, cpu.asl, cpu.slo,
		cpu.php, cpu.ora, cpu.asl, cpu.anc, cpu.nop, cpu.ora, cpu.asl, cpu.slo,
		cpu.bpl, cpu.ora, cpu.kil, cpu.slo, cpu.nop, cpu.ora, cpu.asl, cpu.slo,
		cpu.clc, cpu.ora, cpu.nop, cpu.slo, cpu.nop, cpu.ora, cpu.asl, cpu.slo,
		cpu.jsr, cpu.and, cpu.kil, cpu.rla, cpu.bit, cpu.and, cpu.rol, cpu.rla,
		cpu.plp, cpu.and, cpu.rol, cpu.anc, cpu.bit, cpu.and, cpu.rol, cpu.rla,
		cpu.bmi, cpu.and, cpu.kil, cpu.rla, cpu.nop, cpu.and, cpu.rol, cpu.rla,
		cpu.sec, cpu.and, cpu.nop, cpu.rla, cpu.nop, cpu.and, cpu.rol, cpu.rla,
		cpu.rti, cpu.eor, cpu.kil, cpu.sre, cpu.nop, cpu.eor, cpu.lsr, cpu.sre,
		cpu.pha, cpu.eor, cpu.lsr, cpu.alr, cpu.jmp, cpu.eor, cpu.lsr, cpu.sre,
		cpu.bvc, cpu.eor, cpu.kil, cpu.sre, cpu.nop, cpu.eor, cpu.lsr, cpu.sre,
		cpu.cli, cpu.eor, cpu.nop, cpu.sre, cpu.nop, cpu.eor, cpu.lsr, cpu.sre,
		cpu.rts, cpu.adc, cpu.kil, cpu.rra, cpu.nop, cpu.adc, cpu.ror, cpu.rra,
		cpu.pla, cpu.adc, cpu.ror, cpu.arr, cpu.jmp, cpu.adc, cpu.ror, cpu.rra,
		cpu.bvs, cpu.adc, cpu.kil, cpu.rra, cpu.nop, cpu.adc, cpu.ror, cpu.rra,
		cpu.sei, cpu.adc, cpu.nop, cpu.rra, cpu.nop, cpu.adc, cpu.ror, cpu.rra,
		cpu.nop, cpu.sta, cpu.nop, cpu.sax, cpu.sty, cpu.sta, cpu.stx, cpu.sax,
		cpu.dey, cpu.nop, cpu.txa, cpu.xaa, cpu.sty, cpu.sta, cpu.stx, cpu.sax,
		cpu.bcc, cpu.sta, cpu.kil, cpu.ahx, cpu.sty, cpu.sta, cpu.stx, cpu.sax,
		cpu.tya, cpu.sta, cpu.txs, cpu.tas, cpu.shy, cpu.sta, cpu.shx, cpu.ahx,
		cpu.ldy, cpu.lda, cpu.ldx, cpu.lax, cpu.ldy, cpu.lda, cpu.ldx, cpu.lax,
		cpu.tay, cpu.lda, cpu.tax, cpu.lax, cpu.ldy, cpu.lda, cpu.ldx, cpu.lax,
		cpu.bcs, cpu.lda, cpu.kil, cpu.lax, cpu.ldy, cpu.lda, cpu.ldx, cpu.lax,
		cpu.clv, cpu.lda, cpu.tsx, cpu.las, cpu.ldy, cpu.lda, cpu.ldx, cpu.lax,
		cpu.cpy, cpu.cmp, cpu.nop, cpu.dcp, cpu.cpy, cpu.cmp, cpu.dec, cpu.dcp,
		cpu.iny, cpu.cmp, cpu.dex, cpu.axs, cpu.cpy, cpu.cmp, cpu.dec, cpu.dcp,
		cpu.bne, cpu.cmp, cpu.kil, cpu.dcp, cpu.nop, cpu.cmp, cpu.dec, cpu.dcp,
		cpu.cld, cpu.cmp, cpu.nop, cpu.dcp, cpu.nop, cpu.cmp, cpu.dec, cpu.dcp,
		cpu.cpx, cpu.sbc, cpu.nop, cpu.isc, cpu.cpx, cpu.sbc, cpu.inc, cpu.isc,
		cpu.inx, cpu.sbc, cpu.nop, cpu.sbc, cpu.cpx, cpu.sbc, cpu.inc, cpu.isc,
		cpu.beq, cpu.sbc, cpu.kil, cpu.isc, cpu.nop, cpu.sbc, cpu.inc, cpu.isc,
		cpu.sed, cpu.sbc, cpu.nop, cpu.isc, cpu.nop, cpu.sbc, cpu.inc, cpu.isc,
	}
}

func ConnectCPU(bus *Bus) *CPU {
	cpu := CPU{}
	cpu.Bus = bus
	return &cpu
}

// IO
// Read reads one byte from bus and return a word value.
func (cpu *CPU) read(addr uint16) uint8 {
	return cpu.Bus.Read(addr)
}

func (cpu *CPU) write(addr uint16, data uint8) {
	cpu.Bus.Write(addr, data)
}

func (cpu *CPU) getFlag(f uint8) uint8 {
	if (cpu.Status & f) > 0 {
		return 1
	}

	return 0
}

func (cpu *CPU) setFlag(f uint8, v bool) {
	if v {
		cpu.Status |= f
	} else {
		cpu.Status &= ^f
	}
}

func (cpu *CPU) Clock() {
	if cpu.cycles == 0 {
		cpu.opcode = cpu.read(cpu.PC)

		cpu.setFlag(flagUnused, true)

		cpu.PC++

		// Get starting number of cycles
		cpu.cycles = instructionCycles[cpu.opcode]

		mode := instructionModes[cpu.opcode]

		// TODO: refactor this switch
		switch mode {
		case modeImplied:
			cpu.fetched = cpu.A
		case modeImmediate:
			cpu.PC++
			cpu.addrAbs = cpu.PC
		case modeZeroPage:
			cpu.addrAbs = uint16(cpu.read(cpu.PC))
			cpu.PC++
			cpu.addrAbs &= 0x00FF
		case modeZeroPageX:
			cpu.addrAbs = uint16(cpu.read(cpu.PC) + cpu.X)
			cpu.PC++
			cpu.addrAbs &= 0x00FF
		case modeZeroPageY:
			cpu.addrAbs = uint16(cpu.read(cpu.PC) + cpu.Y)
			cpu.PC++
			cpu.addrAbs &= 0x00FF
		case modeRelative:
			cpu.addrRel = uint16(cpu.read(cpu.PC))
			cpu.PC++
			if (cpu.addrRel & 0x80) != 0 {
				cpu.addrRel |= 0xFF00
			}
		case modeAbsolute:
			var lo uint16 = uint16(cpu.read(cpu.PC))
			cpu.PC++
			var hi uint16 = uint16(cpu.read(cpu.PC))
			cpu.PC++

			cpu.addrAbs = uint16((hi << 8) | lo)
		case modeAbsoluteX:
			var lo uint16 = uint16(cpu.read(cpu.PC))
			cpu.PC++
			var hi uint16 = uint16(cpu.read(cpu.PC))
			cpu.PC++

			cpu.addrAbs = uint16((hi << 8) | lo)
			cpu.addrAbs += uint16(cpu.X)

			// If page change happens, we add an additional cycles
			// according to the MOS6502 spec. (a caveat)
			if (cpu.addrAbs & 0xFF00) != uint16((hi << 8)) {
				cpu.cycles++
			}
		case modeAbosluteY:
			var lo uint16 = uint16(cpu.read(cpu.PC))
			cpu.PC++
			var hi uint16 = uint16(cpu.read(cpu.PC))
			cpu.PC++

			cpu.addrAbs = uint16((hi << 8) | lo)
			cpu.addrAbs += uint16(cpu.Y)

			if (cpu.addrAbs & 0xFF00) != uint16((hi << 8)) {
				cpu.cycles++
			}
		case modeIndirect:
			var ptrLo uint16 = uint16(cpu.read(cpu.PC))
			cpu.PC++
			var ptrHi uint16 = uint16(cpu.read(cpu.PC))
			cpu.PC++

			var ptr uint16 = uint16((ptrHi << 8) | ptrLo)

			// Notice that accroding to wiki, there's a page boundary
			// hardware bug.
			if ptrLo == 0x00FF {
				cpu.addrAbs = uint16((cpu.read(ptr&0xFF00) << 8) | cpu.read(ptr+0))
			} else {
				cpu.addrAbs = uint16((cpu.read(ptr+1) << 8) | cpu.read(ptr+0))
			}
		case modeIndirectX:
			var t uint16 = uint16(cpu.read(cpu.PC))
			cpu.PC++

			var lo uint16 = uint16(cpu.read((t + uint16(cpu.X)) & 0x00FF))
			var hi uint16 = uint16(cpu.read((t + uint16(cpu.X) + 1) & 0x00FF))

			cpu.addrAbs = (hi << 8) | lo
		case modeIndirectY:
			var t uint16 = uint16(cpu.read(cpu.PC))
			cpu.PC++

			var lo uint16 = uint16(cpu.read(t & 0x00FF))
			var hi uint16 = uint16(cpu.read((t + 1) & 0x00FF))

			cpu.addrAbs = (hi << 8) | lo
			cpu.addrAbs += uint16(cpu.Y)

			// We may cross the page boundary so we may need extra
			// cycles.
			if (cpu.addrAbs & 0xFF00) != (hi << 8) {
				cpu.cycles++
			}
		}

		var additionalCycles uint8 = cpu.table[cpu.opcode]()
		cpu.cycles += additionalCycles

		cpu.setFlag(flagUnused, true)
	}
	cpu.clockCount++
	cpu.cycles--
}

func (cpu CPU) reset() {
	cpu.addrAbs = 0xFFFC
	var lo uint16 = uint16(cpu.read(cpu.addrAbs + 0))
	var hi uint16 = uint16(cpu.read(cpu.addrAbs + 1))

	cpu.PC = (hi << 8) | lo

	cpu.A = 0
	cpu.X = 0
	cpu.Y = 0
	cpu.SP = 0xFD
	cpu.Status = 0x00 | flagUnused

	cpu.addrAbs = 0x0000
	cpu.addrRel = 0x0000
	cpu.fetched = 0x00

	// Interrupt reset need cycles.
	cpu.cycles = 8
}

// Interrupts
// Interrupt request
func (cpu *CPU) irq() {
	if cpu.getFlag(flagDisableInterrupts) == 0 {
		cpu.write(0x0100+uint16(cpu.SP), uint8((cpu.PC>>8)&0x00FF))
		cpu.SP--
		cpu.write(0x0100+uint16(cpu.SP), uint8(cpu.PC&0x00FF))
		cpu.SP--

		cpu.setFlag(flagBreak, false)
		cpu.setFlag(flagUnused, true)
		cpu.setFlag(flagDisableInterrupts, true)
		cpu.write(0x0100+uint16(cpu.SP), cpu.Status)
		cpu.SP--

		cpu.addrAbs = 0xFFFE
		var lo uint16 = uint16(cpu.read(cpu.addrAbs + 0))
		var hi uint16 = uint16(cpu.read(cpu.addrAbs + 1))
		cpu.PC = (hi << 8) | lo

		cpu.cycles = 7
	}
}

// Non-Maskable interrupt
func (cpu *CPU) nmi() {
	cpu.write(0x0100+uint16(cpu.SP), uint8((cpu.PC>>8)&0x00FF))
	cpu.SP--
	cpu.write(0x0100+uint16(cpu.SP), uint8(cpu.PC&0x00FF))
	cpu.SP--

	cpu.setFlag(flagBreak, false)
	cpu.setFlag(flagUnused, true)
	cpu.setFlag(flagDisableInterrupts, true)
	cpu.write(0x0100+uint16(cpu.SP), cpu.Status)
	cpu.SP--

	cpu.addrAbs = 0xFFFA
	var lo uint16 = uint16(cpu.read(cpu.addrAbs + 0))
	var hi uint16 = uint16(cpu.read(cpu.addrAbs + 1))
	cpu.PC = (hi << 8) | lo

	cpu.cycles = 8
}

// Instructions
func (cpu *CPU) fetch() uint8 {
	var fetched uint8 = 0
	if instructionModes[cpu.opcode] != modeImplied {
		fetched = cpu.read(cpu.addrAbs)
	}
	return fetched
}

// Legal instructions
func (cpu *CPU) adc() uint8 {
	cpu.fetch()

	cpu.temp = uint16(cpu.A) + uint16(cpu.fetched) + uint16(cpu.getFlag(flagCarryBit))

	cpu.setFlag(flagCarryBit, cpu.temp > 255)
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0)
	// V = (A ^ R) & ~(A ^ M)
	cpu.setFlag(flagOverflow, (^(uint16(cpu.A)^uint16(cpu.fetched))&(uint16(cpu.A)^uint16(cpu.temp)))&0x0080 == 0)
	cpu.setFlag(flagNegative, cpu.temp&0x80 == 0)
	cpu.A = uint8(cpu.temp & 0x00FF)

	return 1
}

func (cpu *CPU) sbc() uint8 {
	cpu.fetch()

	var value uint16 = uint16(cpu.fetched) ^ 0x00FF

	cpu.temp = uint16(cpu.A) + value + uint16(cpu.getFlag(flagCarryBit))
	cpu.setFlag(flagCarryBit, cpu.temp&0xFF00 == 0)
	cpu.setFlag(flagZero, cpu.temp&0x00FF == 0)
	cpu.setFlag(flagOverflow, (cpu.temp^uint16(cpu.A)&(cpu.temp^value)&0x0080) == 0)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 == 0)
	cpu.A = uint8(cpu.temp & 0x00FF)

	return 1
}

func (cpu *CPU) and() uint8 {
	cpu.fetch()
	cpu.A = cpu.A & cpu.fetched
	cpu.setFlag(flagZero, cpu.A == 0x00)
	cpu.setFlag(flagNegative, cpu.A&0x80 == 0)
	return 1
}

func (cpu *CPU) asl() uint8 {
	cpu.fetch()
	cpu.temp = uint16(cpu.fetched) << 1
	cpu.setFlag(flagCarryBit, (cpu.temp&0xFF00) > 0)
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x00)
	cpu.setFlag(flagNegative, cpu.temp&0x80 == 0)

	if instructionModes[cpu.opcode] == modeImplied {
		cpu.A = uint8(cpu.temp & 0x00FF)
	} else {
		cpu.write(cpu.addrAbs, uint8(cpu.temp&0x00FF))
	}

	return 0
}

func (cpu *CPU) bcc() uint8 {
	if cpu.getFlag(flagCarryBit) == 0 {
		cpu.cycles++
		cpu.addrAbs = cpu.PC + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.cycles++
		}

		cpu.PC = cpu.addrAbs
	}
	return 0
}

func (cpu *CPU) bcs() uint8 {
	if cpu.getFlag(flagCarryBit) == 1 {
		cpu.cycles++
		cpu.addrAbs = cpu.PC + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.cycles++
		}

		cpu.PC = cpu.addrAbs
	}
	return 0
}

func (cpu *CPU) beq() uint8 {
	if cpu.getFlag(flagZero) == 1 {
		cpu.cycles++
		cpu.addrAbs = cpu.PC + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.cycles++
		}

		cpu.PC = cpu.addrAbs
	}

	return 0
}

func (cpu *CPU) bit() uint8 {
	cpu.fetch()
	cpu.temp = uint16(cpu.A & cpu.fetched)
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x00)
	cpu.setFlag(flagNegative, cpu.fetched&(1<<7) == 0)
	cpu.setFlag(flagOverflow, cpu.fetched&(1<<6) == 0)
	return 0
}

func (cpu *CPU) bmi() uint8 {
	if cpu.getFlag(flagNegative) == 1 {
		cpu.cycles++
		cpu.addrAbs = cpu.PC + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.cycles++
		}

		cpu.PC = cpu.addrAbs
	}

	return 0
}

func (cpu *CPU) bne() uint8 {
	if cpu.getFlag(flagZero) == 0 {
		cpu.cycles++
		cpu.addrAbs = cpu.PC + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.cycles++
		}

		cpu.PC = cpu.addrAbs
	}

	return 0
}

func (cpu *CPU) bpl() uint8 {
	if cpu.getFlag(flagNegative) == 0 {
		cpu.cycles++
		cpu.addrAbs = cpu.PC + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.cycles++
		}

		cpu.PC = cpu.addrAbs
	}

	return 0
}

func (cpu *CPU) brk() uint8 {
	cpu.PC++

	cpu.setFlag(flagDisableInterrupts, true)
	cpu.write(0x0100+uint16(cpu.SP), uint8((cpu.PC>>8)&0x00FF))
	cpu.SP--
	cpu.write(0x0100+uint16(cpu.SP), uint8(cpu.PC&0x00FF))
	cpu.SP--

	cpu.setFlag(flagBreak, true)
	cpu.write(0x0100+uint16(cpu.SP), cpu.Status)
	cpu.SP--
	cpu.setFlag(flagBreak, false)

	cpu.PC = uint16(cpu.read(0xFFFE)) | (uint16(cpu.read(0xFFFF)) << 8)
	return 0
}

func (cpu *CPU) bvc() uint8 {
	if cpu.getFlag(flagOverflow) == 0 {
		cpu.cycles++
		cpu.addrAbs = cpu.PC + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.cycles++
		}

		cpu.PC = cpu.addrAbs
	}

	return 0
}

func (cpu *CPU) bvs() uint8 {
	if cpu.getFlag(flagOverflow) == 1 {
		cpu.cycles++
		cpu.addrAbs = cpu.PC + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.cycles++
		}

		cpu.PC = cpu.addrAbs
	}

	return 0
}

func (cpu *CPU) clc() uint8 {
	cpu.setFlag(flagCarryBit, false)
	return 0
}

func (cpu *CPU) cld() uint8 {
	cpu.setFlag(flagDecimalMode, false)
	return 0
}

func (cpu *CPU) cli() uint8 {
	cpu.setFlag(flagDisableInterrupts, false)
	return 0
}

func (cpu *CPU) clv() uint8 {
	cpu.setFlag(flagOverflow, false)
	return 0
}

func (cpu *CPU) cmp() uint8 {
	cpu.fetch()
	cpu.temp = uint16(cpu.A) - uint16(cpu.fetched)
	cpu.setFlag(flagCarryBit, cpu.A >= cpu.fetched)
	cpu.setFlag(flagZero, (cpu.temp*0x00FF) == 0x0000)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 == 0)
	return 1
}

func (cpu *CPU) cpx() uint8 {
	cpu.fetch()
	cpu.temp = uint16(cpu.X) - uint16(cpu.fetched)
	cpu.setFlag(flagCarryBit, cpu.X >= cpu.fetched)
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x0000)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 == 0)
	return 0
}

func (cpu *CPU) cpy() uint8 {
	cpu.fetch()
	cpu.temp = uint16(cpu.Y) - uint16(cpu.fetched)
	cpu.setFlag(flagCarryBit, cpu.Y >= cpu.fetched)
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x0000)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 == 0)
	return 0
}

func (cpu *CPU) dec() uint8 {
	cpu.fetch()
	cpu.temp = uint16(cpu.fetched - 1)
	cpu.write(cpu.addrAbs, uint8(cpu.temp&0x00FF))
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x0000)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 == 0)
	return 0
}

func (cpu *CPU) dex() uint8 {
	cpu.X--
	cpu.setFlag(flagZero, cpu.X == 0x00)
	cpu.setFlag(flagNegative, cpu.X&0x80 == 0)
	return 0
}

func (cpu *CPU) dey() uint8 {
	cpu.Y--
	cpu.setFlag(flagZero, cpu.Y == 0x00)
	cpu.setFlag(flagNegative, cpu.Y&0x80 == 0)
	return 0
}

func (cpu *CPU) eor() uint8 {
	cpu.fetch()
	cpu.A = cpu.A ^ cpu.fetched
	cpu.setFlag(flagZero, cpu.A == 0x00)
	cpu.setFlag(flagNegative, cpu.A&0x80 == 0)
	return 1
}

func (cpu *CPU) inc() uint8 {
	cpu.fetch()
	cpu.temp = uint16(cpu.fetched + 1)
	cpu.write(cpu.addrAbs, uint8(cpu.temp&0x00FF))
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x0000)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 == 0)
	return 0
}

func (cpu *CPU) inx() uint8 {
	cpu.X++
	cpu.setFlag(flagZero, cpu.X == 0x00)
	cpu.setFlag(flagNegative, cpu.X&0x80 == 0)
	return 0
}

func (cpu *CPU) iny() uint8 {
	cpu.Y++
	cpu.setFlag(flagZero, cpu.Y == 0x00)
	cpu.setFlag(flagNegative, cpu.Y&0x80 == 0)
	return 0
}

func (cpu *CPU) jmp() uint8 {
	cpu.PC = cpu.addrAbs
	return 0
}

func (cpu *CPU) jsr() uint8 {
	cpu.PC--

	cpu.write(0x0100+uint16(cpu.SP), uint8((cpu.PC>>8)&0x00FF))
	cpu.SP--
	cpu.write(0x0100+uint16(cpu.SP), uint8(cpu.PC&0x00FF))
	cpu.SP--

	cpu.PC = cpu.addrAbs
	return 0
}

func (cpu *CPU) lda() uint8 {
	cpu.fetch()
	cpu.A = cpu.fetched
	cpu.setFlag(flagZero, cpu.A == 0x00)
	cpu.setFlag(flagNegative, cpu.A&0x80 == 0)
	return 1
}

func (cpu *CPU) ldx() uint8 {
	cpu.fetch()
	cpu.X = cpu.fetched
	cpu.setFlag(flagZero, cpu.X == 0x00)
	cpu.setFlag(flagNegative, cpu.X&0x80 == 0)
	return 1
}

func (cpu *CPU) ldy() uint8 {
	cpu.fetch()
	cpu.setFlag(flagZero, cpu.Y == 0x00)
	cpu.setFlag(flagNegative, cpu.Y&0x80 == 0)
	return 1
}

func (cpu *CPU) lsr() uint8 {
	cpu.fetch()
	cpu.setFlag(flagCarryBit, cpu.fetched&0x0001 == 0)
	cpu.temp = uint16(cpu.fetched) >> 1
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x0000)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 == 0)
	if instructionModes[cpu.opcode] == modeImplied {
		cpu.A = uint8(cpu.temp & 0x00FF)
	} else {
		cpu.write(cpu.addrAbs, uint8(cpu.temp&0x00FF))
	}
	return 0
}

func (cpu *CPU) nop() uint8 {
	// Note that not all NOPs are equal. There're indeed some illegal ones.
	switch cpu.opcode {
	case 0x1C:
	case 0x3C:
	case 0x5C:
	case 0x7C:
	case 0xDC:
	case 0xFC:
		return 1
	}
	return 0
}

func (cpu *CPU) ora() uint8 {
	cpu.fetch()
	cpu.A = cpu.A | cpu.fetched
	cpu.setFlag(flagZero, cpu.A == 0x00)
	cpu.setFlag(flagNegative, cpu.A&0x80 == 0)
	return 1
}

func (cpu *CPU) pha() uint8 {
	cpu.write(0x0100+uint16(cpu.SP), cpu.A)
	cpu.SP--
	return 0
}

func (cpu *CPU) php() uint8 {
	cpu.write(0x0100+uint16(cpu.SP), cpu.Status|flagBreak|flagUnused)
	cpu.setFlag(flagBreak, false)
	cpu.setFlag(flagUnused, false)
	cpu.SP--
	return 0
}

func (cpu *CPU) pla() uint8 {
	cpu.SP++
	cpu.A = cpu.read(0x0100 + uint16(cpu.SP))
	cpu.setFlag(flagZero, cpu.A == 0x00)
	cpu.setFlag(flagNegative, cpu.A&0x80 == 0)
	return 0
}

func (cpu *CPU) plp() uint8 {
	cpu.SP++
	cpu.Status = cpu.read(0x0100 + uint16(cpu.SP))
	cpu.setFlag(flagUnused, true)
	return 0
}

func (cpu *CPU) rol() uint8 {
	cpu.fetch()
	cpu.temp = uint16(uint16(cpu.fetched) << 1)
	cpu.setFlag(flagCarryBit, cpu.temp&0xFF00 == 0)
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x0000)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 == 0)
	if instructionModes[cpu.opcode] == modeImplied {
		cpu.A = uint8(cpu.temp & 0x00FF)
	} else {
		cpu.write(cpu.addrAbs, uint8(cpu.temp&0x00FF))
	}
	return 0
}

func (cpu *CPU) ror() uint8 {
	cpu.fetch()
	cpu.temp = uint16(cpu.getFlag(flagCarryBit)<<7 | (cpu.fetched >> 1))
	cpu.setFlag(flagCarryBit, cpu.fetched&0x01 == 0)
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x00)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 == 0)
	if instructionModes[cpu.opcode] == modeImplied {
		cpu.A = uint8(cpu.temp & 0x00FF)
	} else {
		cpu.write(cpu.addrAbs, uint8(cpu.temp&0x00FF))
	}
	return 0
}

func (cpu *CPU) rti() uint8 {
	cpu.SP++
	cpu.Status = cpu.read(0x0100 + uint16(cpu.SP))
	cpu.Status &= ^uint8(flagBreak)
	cpu.Status &= ^uint8(flagUnused)

	cpu.SP++
	cpu.PC = uint16(cpu.read(0x0100 + uint16(cpu.SP)))
	cpu.SP++
	cpu.PC |= uint16(cpu.read(0x0100+uint16(cpu.SP))) << 8
	return 0
}

func (cpu *CPU) rts() uint8 {
	cpu.SP++
	cpu.PC = uint16(cpu.read(0x0100 + uint16(cpu.SP)))
	cpu.SP++
	cpu.PC |= uint16(cpu.read(0x0100+uint16(cpu.SP))) << 8

	cpu.PC++
	return 0
}

func (cpu *CPU) sec() uint8 {
	cpu.setFlag(flagCarryBit, true)
	return 0
}

func (cpu *CPU) sed() uint8 {
	cpu.setFlag(flagDecimalMode, true)
	return 0
}

func (cpu *CPU) sei() uint8 {
	cpu.setFlag(flagDisableInterrupts, true)
	return 0
}

func (cpu *CPU) sta() uint8 {
	cpu.write(cpu.addrAbs, cpu.A)
	return 0
}

func (cpu *CPU) stx() uint8 {
	cpu.write(cpu.addrAbs, cpu.X)
	return 0
}

func (cpu *CPU) sty() uint8 {
	cpu.write(cpu.addrAbs, cpu.Y)
	return 0
}

func (cpu *CPU) tax() uint8 {
	cpu.X = cpu.A
	cpu.setFlag(flagZero, cpu.X == 0x00)
	cpu.setFlag(flagNegative, cpu.X&0x80 == 0)
	return 0
}

func (cpu *CPU) tay() uint8 {
	cpu.Y = cpu.A
	cpu.setFlag(flagZero, cpu.Y == 0x00)
	cpu.setFlag(flagNegative, cpu.Y&0x80 == 0)
	return 0
}

func (cpu *CPU) tsx() uint8 {
	cpu.X = cpu.SP
	cpu.setFlag(flagZero, cpu.X == 0x00)
	cpu.setFlag(flagNegative, cpu.X&0x80 == 0)
	return 0
}

func (cpu *CPU) txa() uint8 {
	cpu.A = cpu.X
	cpu.setFlag(flagZero, cpu.A == 0x00)
	cpu.setFlag(flagNegative, cpu.A&0x80 == 0)
	return 0
}

func (cpu *CPU) txs() uint8 {
	cpu.SP = cpu.X
	return 0
}

func (cpu *CPU) tya() uint8 {
	cpu.A = cpu.Y
	cpu.setFlag(flagZero, cpu.A == 0x00)
	cpu.setFlag(flagNegative, cpu.A&0x80 == 0)
	return 0
}

// Illegal insturctions
// Since these instructions are not used in NES, I'll leave them blank.
func (cpu *CPU) ahx() uint8 {
	return 0
}

func (cpu *CPU) alr() uint8 {
	return 0
}

func (cpu *CPU) anc() uint8 {
	return 0
}

func (cpu *CPU) arr() uint8 {
	return 0
}

func (cpu *CPU) axs() uint8 {
	return 0
}

func (cpu *CPU) dcp() uint8 {
	return 0
}

func (cpu *CPU) isc() uint8 {
	return 0
}

func (cpu *CPU) kil() uint8 {
	return 0
}

func (cpu *CPU) las() uint8 {
	return 0
}

func (cpu *CPU) lax() uint8 {
	return 0
}

func (cpu *CPU) rla() uint8 {
	return 0
}

func (cpu *CPU) rra() uint8 {
	return 0
}

func (cpu *CPU) sax() uint8 {
	return 0
}

func (cpu *CPU) shx() uint8 {
	return 0
}

func (cpu *CPU) shy() uint8 {
	return 0
}

func (cpu *CPU) slo() uint8 {
	return 0
}

func (cpu *CPU) sre() uint8 {
	return 0
}

func (cpu *CPU) tas() uint8 {
	return 0
}

func (cpu *CPU) xaa() uint8 {
	return 0
}

// Helper functions
func (cpu *CPU) Complete() bool {
	return cpu.cycles == 0
}

func (cpu *CPU) disassemble(nStart uint16, nStop uint16) {
	var addr uint16 = nStart
	var value, hi, lo uint8 = 0x00, 0x00, 0x00

	var mapLines map[uint16]string
	var lineAddr uint16 = 0

	for addr <= nStop {
		lineAddr = uint16(addr)

		sInst := "$" + hex(addr, 4) + ": "

		opcode := cpu.Bus.Read(addr)
		addr++
		sInst += instructionNames[opcode] + " "

		mode := instructionModes[opcode]

		switch mode {
		case modeImplied:
			sInst += " {IMP}"
		case modeImmediate:
			value = cpu.Bus.Read(addr)
			sInst += "#$" + hex(uint16(value), 2) + " {IMM}"
		case modeZeroPage:
			lo = cpu.Bus.Read(addr)
			addr++
			hi = 0x00
			sInst += "$" + hex(uint16(lo), 2) + " {ZP0}"
		case modeZeroPageX:
			lo = cpu.Bus.Read(addr)
			addr++
			hi = 0x00
			sInst += "$" + hex(uint16(lo), 2) + ", X {ZPX}"
			// TODO: Finish this simple disassembler
		}
	}
}

func replaceAtIndex(str string, replacement byte, index int) string {
	return str[:index] + string(replacement) + str[index+1:]
}

func hex(n uint16, d uint8) string {
	var s string = strings.Repeat("0", int(d))
	for i := int(d) - 1; i >= 0; i, n = i-1, n>>4 {
		s = replaceAtIndex(s, "0123456789ABCDEF"[n&0xF], i)
	}
	return s
}
