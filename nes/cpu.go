package nes

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
	modeAbsoluteY
	modeIndirect
	modeIndirectX
	modeIndirectY
)

const (
	flagCarryBit          = (1 << 0)
	flagZero              = (1 << 1)
	flagDisableInterrupts = (1 << 2)
	flagDecimalMode       = (1 << 3) // Redundant
	flagBreak             = (1 << 4)
	flagUnused            = (1 << 5)
	flagOverflow          = (1 << 6)
	flagNegative          = (1 << 7)
)

var instructionModes = [256]byte{
	2, 11, 1, 1, 1, 3, 3, 1, 1, 2, 1, 1, 1, 7, 7, 1,
	6, 12, 1, 1, 1, 4, 4, 1, 1, 9, 1, 1, 1, 8, 8, 1,
	7, 11, 1, 1, 3, 3, 3, 1, 1, 2, 1, 1, 7, 7, 7, 1,
	6, 12, 1, 1, 1, 4, 4, 1, 1, 9, 1, 1, 1, 8, 8, 1,
	1, 11, 1, 1, 1, 3, 3, 1, 1, 2, 1, 1, 7, 7, 7, 1,
	6, 12, 1, 1, 1, 4, 4, 1, 1, 9, 1, 1, 1, 8, 8, 1,
	1, 11, 1, 1, 1, 3, 3, 1, 1, 2, 1, 1, 10, 7, 7, 1,
	6, 12, 1, 1, 1, 4, 4, 1, 1, 9, 1, 1, 1, 8, 8, 1,
	1, 11, 1, 1, 3, 3, 3, 1, 1, 1, 1, 1, 7, 7, 7, 1,
	6, 12, 1, 1, 4, 4, 5, 1, 1, 9, 1, 1, 1, 8, 1, 1,
	2, 11, 2, 1, 3, 3, 3, 1, 1, 2, 1, 1, 7, 7, 7, 1,
	6, 12, 1, 1, 4, 4, 5, 1, 1, 9, 1, 1, 8, 8, 9, 1,
	2, 11, 1, 1, 3, 3, 3, 1, 1, 2, 1, 1, 7, 7, 7, 1,
	6, 12, 1, 1, 1, 4, 4, 1, 1, 9, 1, 1, 1, 8, 8, 1,
	2, 11, 1, 1, 3, 3, 3, 1, 1, 2, 1, 1, 7, 7, 7, 1,
	6, 12, 1, 1, 1, 4, 4, 1, 1, 9, 1, 1, 1, 8, 8, 1,
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

// "XXX" means illegal instructions.
var instructionNames = [256]string{
	"BRK", "ORA", "XXX", "XXX", "NOP", "ORA", "ASL", "XXX", "PHP", "ORA", "ASL", "XXX", "NOP", "ORA", "ASL", "XXX",
	"BPL", "ORA", "XXX", "XXX", "NOP", "ORA", "ASL", "XXX", "CLC", "ORA", "NOP", "XXX", "NOP", "ORA", "ASL", "XXX",
	"JSR", "AND", "XXX", "XXX", "BIT", "AND", "ROL", "XXX", "PLP", "AND", "ROL", "XXX", "BIT", "AND", "ROL", "XXX",
	"BMI", "AND", "XXX", "XXX", "NOP", "AND", "ROL", "XXX", "SEC", "AND", "NOP", "XXX", "NOP", "AND", "ROL", "XXX",
	"RTI", "EOR", "XXX", "XXX", "NOP", "EOR", "LSR", "XXX", "PHA", "EOR", "LSR", "XXX", "JMP", "EOR", "LSR", "XXX",
	"BVC", "EOR", "XXX", "XXX", "NOP", "EOR", "LSR", "XXX", "CLI", "EOR", "NOP", "XXX", "NOP", "EOR", "LSR", "XXX",
	"RTS", "ADC", "XXX", "XXX", "NOP", "ADC", "ROR", "XXX", "PLA", "ADC", "ROR", "XXX", "JMP", "ADC", "ROR", "XXX",
	"BVS", "ADC", "XXX", "XXX", "NOP", "ADC", "ROR", "XXX", "SEI", "ADC", "NOP", "XXX", "NOP", "ADC", "ROR", "XXX",
	"NOP", "STA", "NOP", "XXX", "STY", "STA", "STX", "XXX", "DEY", "NOP", "TXA", "XXX", "STY", "STA", "STX", "XXX",
	"BCC", "STA", "XXX", "XXX", "STY", "STA", "STX", "XXX", "TYA", "STA", "TXS", "XXX", "NOP", "STA", "XXX", "XXX",
	"LDY", "LDA", "LDX", "XXX", "LDY", "LDA", "LDX", "XXX", "TAY", "LDA", "TAX", "XXX", "LDY", "LDA", "LDX", "XXX",
	"BCS", "LDA", "XXX", "XXX", "LDY", "LDA", "LDX", "XXX", "CLV", "LDA", "TSX", "XXX", "LDY", "LDA", "LDX", "XXX",
	"CPY", "CMP", "NOP", "XXX", "CPY", "CMP", "DEC", "XXX", "INY", "CMP", "DEX", "XXX", "CPY", "CMP", "DEC", "XXX",
	"BNE", "CMP", "XXX", "XXX", "NOP", "CMP", "DEC", "XXX", "CLD", "CMP", "NOP", "XXX", "NOP", "CMP", "DEC", "XXX",
	"CPX", "SBC", "NOP", "XXX", "CPX", "SBC", "INC", "XXX", "INX", "SBC", "NOP", "SBC", "CPX", "SBC", "INC", "XXX",
	"BEQ", "SBC", "XXX", "XXX", "NOP", "SBC", "INC", "XXX", "SED", "SBC", "NOP", "XXX", "NOP", "SBC", "INC", "XXX",
}

// type Instruction struct {
// 	name    string
// 	operate uint8
// 	addr    uint8
// 	cycles  uint8
// }

// CPU MOS 6502 CPU
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

	opcodeTable [256]func() uint8
	modeTable   [256]func() uint8
}

func (cpu *CPU) createTable() {
	cpu.opcodeTable = [256]func() uint8{
		cpu.brk, cpu.ora, cpu.xxx, cpu.xxx, cpu.nop, cpu.ora, cpu.asl, cpu.xxx, cpu.php, cpu.ora, cpu.asl, cpu.xxx, cpu.nop, cpu.ora, cpu.asl, cpu.xxx,
		cpu.bpl, cpu.ora, cpu.xxx, cpu.xxx, cpu.nop, cpu.ora, cpu.asl, cpu.xxx, cpu.clc, cpu.ora, cpu.nop, cpu.xxx, cpu.nop, cpu.ora, cpu.asl, cpu.xxx,
		cpu.jsr, cpu.and, cpu.xxx, cpu.xxx, cpu.bit, cpu.and, cpu.rol, cpu.xxx, cpu.plp, cpu.and, cpu.rol, cpu.xxx, cpu.bit, cpu.and, cpu.rol, cpu.xxx,
		cpu.bmi, cpu.and, cpu.xxx, cpu.xxx, cpu.nop, cpu.and, cpu.rol, cpu.xxx, cpu.sec, cpu.and, cpu.nop, cpu.xxx, cpu.nop, cpu.and, cpu.rol, cpu.xxx,
		cpu.rti, cpu.eor, cpu.xxx, cpu.xxx, cpu.nop, cpu.eor, cpu.lsr, cpu.xxx, cpu.pha, cpu.eor, cpu.lsr, cpu.xxx, cpu.jmp, cpu.eor, cpu.lsr, cpu.xxx,
		cpu.bvc, cpu.eor, cpu.xxx, cpu.xxx, cpu.nop, cpu.eor, cpu.lsr, cpu.xxx, cpu.cli, cpu.eor, cpu.nop, cpu.xxx, cpu.nop, cpu.eor, cpu.lsr, cpu.xxx,
		cpu.rts, cpu.adc, cpu.xxx, cpu.xxx, cpu.nop, cpu.adc, cpu.ror, cpu.xxx, cpu.pla, cpu.adc, cpu.ror, cpu.xxx, cpu.jmp, cpu.adc, cpu.ror, cpu.xxx,
		cpu.bvs, cpu.adc, cpu.xxx, cpu.xxx, cpu.nop, cpu.adc, cpu.ror, cpu.xxx, cpu.sei, cpu.adc, cpu.nop, cpu.xxx, cpu.nop, cpu.adc, cpu.ror, cpu.xxx,
		cpu.nop, cpu.sta, cpu.nop, cpu.xxx, cpu.sty, cpu.sta, cpu.stx, cpu.xxx, cpu.dey, cpu.nop, cpu.txa, cpu.xxx, cpu.sty, cpu.sta, cpu.stx, cpu.xxx,
		cpu.bcc, cpu.sta, cpu.xxx, cpu.xxx, cpu.sty, cpu.sta, cpu.stx, cpu.xxx, cpu.tya, cpu.sta, cpu.txs, cpu.xxx, cpu.nop, cpu.sta, cpu.xxx, cpu.xxx,
		cpu.ldy, cpu.lda, cpu.ldx, cpu.xxx, cpu.ldy, cpu.lda, cpu.ldx, cpu.xxx, cpu.tay, cpu.lda, cpu.tax, cpu.xxx, cpu.ldy, cpu.lda, cpu.ldx, cpu.xxx,
		cpu.bcs, cpu.lda, cpu.xxx, cpu.xxx, cpu.ldy, cpu.lda, cpu.ldx, cpu.xxx, cpu.clv, cpu.lda, cpu.tsx, cpu.xxx, cpu.ldy, cpu.lda, cpu.ldx, cpu.xxx,
		cpu.cpy, cpu.cmp, cpu.nop, cpu.xxx, cpu.cpy, cpu.cmp, cpu.dec, cpu.xxx, cpu.iny, cpu.cmp, cpu.dex, cpu.xxx, cpu.cpy, cpu.cmp, cpu.dec, cpu.xxx,
		cpu.bne, cpu.cmp, cpu.xxx, cpu.xxx, cpu.nop, cpu.cmp, cpu.dec, cpu.xxx, cpu.cld, cpu.cmp, cpu.nop, cpu.xxx, cpu.nop, cpu.cmp, cpu.dec, cpu.xxx,
		cpu.cpx, cpu.sbc, cpu.nop, cpu.xxx, cpu.cpx, cpu.sbc, cpu.inc, cpu.xxx, cpu.inx, cpu.sbc, cpu.nop, cpu.sbc, cpu.cpx, cpu.sbc, cpu.inc, cpu.xxx,
		cpu.beq, cpu.sbc, cpu.xxx, cpu.xxx, cpu.nop, cpu.sbc, cpu.inc, cpu.xxx, cpu.sed, cpu.sbc, cpu.nop, cpu.xxx, cpu.nop, cpu.sbc, cpu.inc, cpu.xxx,
	}

	cpu.modeTable = [256]func() uint8{
		cpu.imm, cpu.izx, cpu.imp, cpu.imp, cpu.imp, cpu.zp0, cpu.zp0, cpu.imp, cpu.imp, cpu.imm, cpu.imp, cpu.imp, cpu.imp, cpu.abs, cpu.abs, cpu.imp,
		cpu.rel, cpu.izy, cpu.imp, cpu.imp, cpu.imp, cpu.zpx, cpu.zpx, cpu.imp, cpu.imp, cpu.aby, cpu.imp, cpu.imp, cpu.imp, cpu.abx, cpu.abx, cpu.imp,
		cpu.abs, cpu.izx, cpu.imp, cpu.imp, cpu.zp0, cpu.zp0, cpu.zp0, cpu.imp, cpu.imp, cpu.imm, cpu.imp, cpu.imp, cpu.abs, cpu.abs, cpu.abs, cpu.imp,
		cpu.rel, cpu.izy, cpu.imp, cpu.imp, cpu.imp, cpu.zpx, cpu.zpx, cpu.imp, cpu.imp, cpu.aby, cpu.imp, cpu.imp, cpu.imp, cpu.abx, cpu.abx, cpu.imp,
		cpu.imp, cpu.izx, cpu.imp, cpu.imp, cpu.imp, cpu.zp0, cpu.zp0, cpu.imp, cpu.imp, cpu.imm, cpu.imp, cpu.imp, cpu.abs, cpu.abs, cpu.abs, cpu.imp,
		cpu.rel, cpu.izy, cpu.imp, cpu.imp, cpu.imp, cpu.zpx, cpu.zpx, cpu.imp, cpu.imp, cpu.aby, cpu.imp, cpu.imp, cpu.imp, cpu.abx, cpu.abx, cpu.imp,
		cpu.imp, cpu.izx, cpu.imp, cpu.imp, cpu.imp, cpu.zp0, cpu.zp0, cpu.imp, cpu.imp, cpu.imm, cpu.imp, cpu.imp, cpu.ind, cpu.abs, cpu.abs, cpu.imp,
		cpu.rel, cpu.izy, cpu.imp, cpu.imp, cpu.imp, cpu.zpx, cpu.zpx, cpu.imp, cpu.imp, cpu.aby, cpu.imp, cpu.imp, cpu.imp, cpu.abx, cpu.abx, cpu.imp,
		cpu.imp, cpu.izx, cpu.imp, cpu.imp, cpu.zp0, cpu.zp0, cpu.zp0, cpu.imp, cpu.imp, cpu.imp, cpu.imp, cpu.imp, cpu.abs, cpu.abs, cpu.abs, cpu.imp,
		cpu.rel, cpu.izy, cpu.imp, cpu.imp, cpu.zpx, cpu.zpx, cpu.zpy, cpu.imp, cpu.imp, cpu.aby, cpu.imp, cpu.imp, cpu.imp, cpu.abx, cpu.imp, cpu.imp,
		cpu.imm, cpu.izx, cpu.imm, cpu.imp, cpu.zp0, cpu.zp0, cpu.zp0, cpu.imp, cpu.imp, cpu.imm, cpu.imp, cpu.imp, cpu.abs, cpu.abs, cpu.abs, cpu.imp,
		cpu.rel, cpu.izy, cpu.imp, cpu.imp, cpu.zpx, cpu.zpx, cpu.zpy, cpu.imp, cpu.imp, cpu.aby, cpu.imp, cpu.imp, cpu.abx, cpu.abx, cpu.aby, cpu.imp,
		cpu.imm, cpu.izx, cpu.imp, cpu.imp, cpu.zp0, cpu.zp0, cpu.zp0, cpu.imp, cpu.imp, cpu.imm, cpu.imp, cpu.imp, cpu.abs, cpu.abs, cpu.abs, cpu.imp,
		cpu.rel, cpu.izy, cpu.imp, cpu.imp, cpu.imp, cpu.zpx, cpu.zpx, cpu.imp, cpu.imp, cpu.aby, cpu.imp, cpu.imp, cpu.imp, cpu.abx, cpu.abx, cpu.imp,
		cpu.imm, cpu.izx, cpu.imp, cpu.imp, cpu.zp0, cpu.zp0, cpu.zp0, cpu.imp, cpu.imp, cpu.imm, cpu.imp, cpu.imp, cpu.abs, cpu.abs, cpu.abs, cpu.imp,
		cpu.rel, cpu.izy, cpu.imp, cpu.imp, cpu.imp, cpu.zpx, cpu.zpx, cpu.imp, cpu.imp, cpu.aby, cpu.imp, cpu.imp, cpu.imp, cpu.abx, cpu.abx, cpu.imp,
	}
}

// ConnectCPU Initialize a CPU and connect it to the bus
func ConnectCPU(bus *Bus) *CPU {
	cpu := CPU{}
	cpu.Bus = bus
	cpu.createTable()
	return &cpu
}

// Reset Reset CPU
func (cpu *CPU) Reset() {
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

// IO
// Read reads one byte from bus and return a word value.
func (cpu *CPU) read(addr uint16) uint8 {
	return cpu.Bus.CPURead(addr, false)
}

func (cpu *CPU) write(addr uint16, data uint8) {
	cpu.Bus.CPUWrite(addr, data)
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

// Clock Tick CPU once.
func (cpu *CPU) Clock() {
	if cpu.cycles == 0 {
		cpu.opcode = cpu.read(cpu.PC)

		cpu.setFlag(flagUnused, true)

		cpu.PC++

		// Get starting number of cycles
		cpu.cycles = instructionCycles[cpu.opcode]

		additionalCycles1 := cpu.modeTable[cpu.opcode]()
		additionalCycles2 := cpu.opcodeTable[cpu.opcode]()
		cpu.cycles += (additionalCycles1 & additionalCycles2)

		cpu.setFlag(flagUnused, true)
	}
	cpu.clockCount++
	cpu.cycles--
}

// Interrupts

// Irq Interrupt request
func (cpu *CPU) Irq() {
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

// Nmi Non-maskable interrupt
func (cpu *CPU) Nmi() {
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

// Address modes
func (cpu *CPU) imp() uint8 {
	cpu.fetched = cpu.A
	return 0
}

func (cpu *CPU) imm() uint8 {
	cpu.addrAbs = cpu.PC
	cpu.PC++
	return 0
}

func (cpu *CPU) zp0() uint8 {
	cpu.addrAbs = uint16(cpu.read(cpu.PC))
	cpu.PC++
	cpu.addrAbs &= 0x00FF
	return 0
}

func (cpu *CPU) zpx() uint8 {
	cpu.addrAbs = uint16(cpu.read(cpu.PC) + cpu.X)
	cpu.PC++
	cpu.addrAbs &= 0x00FF
	return 0
}

func (cpu *CPU) zpy() uint8 {
	cpu.addrAbs = uint16(cpu.read(cpu.PC) + cpu.Y)
	cpu.PC++
	cpu.addrAbs &= 0x00FF
	return 0
}

func (cpu *CPU) rel() uint8 {
	cpu.addrRel = uint16(cpu.read(cpu.PC))
	cpu.PC++
	if cpu.addrRel&0x80 != 0 {
		cpu.addrRel |= 0xFF00
	}
	return 0
}

func (cpu *CPU) abs() uint8 {
	var lo uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var hi uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++

	cpu.addrAbs = (hi << 8) | lo
	return 0
}

func (cpu *CPU) abx() uint8 {
	var lo uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var hi uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++

	cpu.addrAbs = (hi << 8) | lo
	cpu.addrAbs += uint16(cpu.X)

	// If page change happens, we add an additional cycles
	// according to the MOS6502 spec. (a caveat)
	if (cpu.addrAbs & 0xFF00) != (hi << 8) {
		return 1
	}
	return 0
}

func (cpu *CPU) aby() uint8 {
	var lo uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var hi uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++

	cpu.addrAbs = (hi << 8) | lo
	cpu.addrAbs += uint16(cpu.Y)

	if (cpu.addrAbs & 0xFF00) != (hi << 8) {
		return 1
	}
	return 0
}

func (cpu *CPU) ind() uint8 {
	var ptrLo uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var ptrHi uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++

	var ptr uint16 = (ptrHi << 8) | ptrLo

	// Notice that accroding to wiki, there's a page boundary
	// hardware bug.
	if ptrLo == 0x00FF {
		cpu.addrAbs = (uint16(cpu.read(ptr&0xFF00)) << 8) | uint16(cpu.read(ptr+0))
	} else {
		cpu.addrAbs = (uint16(cpu.read(ptr+1)) << 8) | uint16(cpu.read(ptr+0))
	}
	return 0
}
func (cpu *CPU) izx() uint8 {
	var t uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++

	var lo uint16 = uint16(cpu.read(uint16(t+uint16(cpu.X)) & 0x00FF))
	var hi uint16 = uint16(cpu.read(uint16(t+uint16(t+uint16(cpu.X)+1)) & 0x00FF))

	cpu.addrAbs = (hi << 8) | lo

	return 0
}

func (cpu *CPU) izy() uint8 {
	var t uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++

	var lo uint16 = uint16(cpu.read(t & 0x00FF))
	var hi uint16 = uint16(cpu.read((t + 1) & 0x00FF))

	cpu.addrAbs = (hi << 8) | lo
	cpu.addrAbs += uint16(cpu.Y)

	// We may cross the page boundary so we may need extra
	// cycles.
	if (cpu.addrAbs & 0xFF00) != (hi << 8) {
		return 1
	}
	return 0
}

// Instructions
func (cpu *CPU) fetch() uint8 {
	if instructionModes[cpu.opcode] != modeImplied {
		cpu.fetched = cpu.read(cpu.addrAbs)
	}
	return cpu.fetched
}

// Legal instructions
func (cpu *CPU) adc() uint8 {
	cpu.fetch()

	cpu.temp = uint16(cpu.A) + uint16(cpu.fetched) + uint16(cpu.getFlag(flagCarryBit))

	cpu.setFlag(flagCarryBit, cpu.temp > 255)
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) != 0)
	// V = (A ^ R) & ~(A ^ M)
	cpu.setFlag(flagOverflow, (^(uint16(cpu.A)^uint16(cpu.fetched))&(uint16(cpu.A)^uint16(cpu.temp)))&0x0080 != 0)
	cpu.setFlag(flagNegative, cpu.temp&0x80 != 0)
	cpu.A = uint8(cpu.temp & 0x00FF)

	return 1
}

func (cpu *CPU) sbc() uint8 {
	cpu.fetch()

	var value uint16 = uint16(cpu.fetched) ^ 0x00FF

	cpu.temp = uint16(cpu.A) + value + uint16(cpu.getFlag(flagCarryBit))
	cpu.setFlag(flagCarryBit, cpu.temp&0xFF00 != 0)
	cpu.setFlag(flagZero, cpu.temp&0x00FF != 0)
	cpu.setFlag(flagOverflow, (cpu.temp^uint16(cpu.A)&(cpu.temp^value)&0x0080) != 0)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 != 0)
	cpu.A = uint8(cpu.temp & 0x00FF)

	return 1
}

func (cpu *CPU) and() uint8 {
	cpu.fetch()
	cpu.A = cpu.A & cpu.fetched
	cpu.setFlag(flagZero, cpu.A == 0x00)
	cpu.setFlag(flagNegative, cpu.A&0x80 != 0)
	return 1
}

func (cpu *CPU) asl() uint8 {
	cpu.fetch()
	cpu.temp = uint16(cpu.fetched) << 1
	cpu.setFlag(flagCarryBit, (cpu.temp&0xFF00) > 0)
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x00)
	cpu.setFlag(flagNegative, cpu.temp&0x80 != 0)

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
	cpu.setFlag(flagNegative, cpu.fetched&(1<<7) != 0)
	cpu.setFlag(flagOverflow, cpu.fetched&(1<<6) != 0)
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
	cpu.setFlag(flagNegative, cpu.temp&0x0080 != 0)
	return 1
}

func (cpu *CPU) cpx() uint8 {
	cpu.fetch()
	cpu.temp = uint16(cpu.X) - uint16(cpu.fetched)
	cpu.setFlag(flagCarryBit, cpu.X >= cpu.fetched)
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x0000)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 != 0)
	return 0
}

func (cpu *CPU) cpy() uint8 {
	cpu.fetch()
	cpu.temp = uint16(cpu.Y) - uint16(cpu.fetched)
	cpu.setFlag(flagCarryBit, cpu.Y >= cpu.fetched)
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x0000)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 != 0)
	return 0
}

func (cpu *CPU) dec() uint8 {
	cpu.fetch()
	cpu.temp = uint16(cpu.fetched - 1)
	cpu.write(cpu.addrAbs, uint8(cpu.temp&0x00FF))
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x0000)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 != 0)
	return 0
}

func (cpu *CPU) dex() uint8 {
	cpu.X--
	cpu.setFlag(flagZero, cpu.X == 0x00)
	cpu.setFlag(flagNegative, cpu.X&0x80 != 0)
	return 0
}

func (cpu *CPU) dey() uint8 {
	cpu.Y--
	cpu.setFlag(flagZero, cpu.Y == 0x00)
	cpu.setFlag(flagNegative, cpu.Y&0x80 != 0)
	return 0
}

func (cpu *CPU) eor() uint8 {
	cpu.fetch()
	cpu.A = cpu.A ^ cpu.fetched
	cpu.setFlag(flagZero, cpu.A == 0x00)
	cpu.setFlag(flagNegative, cpu.A&0x80 != 0)
	return 1
}

func (cpu *CPU) inc() uint8 {
	cpu.fetch()
	cpu.temp = uint16(cpu.fetched + 1)
	cpu.write(cpu.addrAbs, uint8(cpu.temp&0x00FF))
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x0000)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 != 0)
	return 0
}

func (cpu *CPU) inx() uint8 {
	cpu.X++
	cpu.setFlag(flagZero, cpu.X == 0x00)
	cpu.setFlag(flagNegative, cpu.X&0x80 != 0)
	return 0
}

func (cpu *CPU) iny() uint8 {
	cpu.Y++
	cpu.setFlag(flagZero, cpu.Y == 0x00)
	cpu.setFlag(flagNegative, cpu.Y&0x80 != 0)
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
	cpu.setFlag(flagNegative, cpu.A&0x80 != 0)
	return 1
}

func (cpu *CPU) ldx() uint8 {
	cpu.fetch()
	cpu.X = cpu.fetched
	cpu.setFlag(flagZero, cpu.X == 0x00)
	cpu.setFlag(flagNegative, cpu.X&0x80 != 0)
	return 1
}

func (cpu *CPU) ldy() uint8 {
	cpu.fetch()
	cpu.Y = cpu.fetched
	cpu.setFlag(flagZero, cpu.Y == 0x00)
	cpu.setFlag(flagNegative, cpu.Y&0x80 != 0)
	return 1
}

func (cpu *CPU) lsr() uint8 {
	cpu.fetch()
	cpu.setFlag(flagCarryBit, cpu.fetched&0x0001 != 0)
	cpu.temp = uint16(cpu.fetched) >> 1
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x0000)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 != 0)
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
	cpu.setFlag(flagNegative, cpu.A&0x80 != 0)
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
	cpu.setFlag(flagNegative, cpu.A&0x80 != 0)
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
	cpu.setFlag(flagCarryBit, cpu.temp&0xFF00 != 0)
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x0000)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 != 0)
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
	cpu.setFlag(flagCarryBit, cpu.fetched&0x01 != 0)
	cpu.setFlag(flagZero, (cpu.temp&0x00FF) == 0x00)
	cpu.setFlag(flagNegative, cpu.temp&0x0080 != 0)
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
	cpu.setFlag(flagNegative, cpu.X&0x80 != 0)
	return 0
}

func (cpu *CPU) tay() uint8 {
	cpu.Y = cpu.A
	cpu.setFlag(flagZero, cpu.Y == 0x00)
	cpu.setFlag(flagNegative, cpu.Y&0x80 != 0)
	return 0
}

func (cpu *CPU) tsx() uint8 {
	cpu.X = cpu.SP
	cpu.setFlag(flagZero, cpu.X == 0x00)
	cpu.setFlag(flagNegative, cpu.X&0x80 != 0)
	return 0
}

func (cpu *CPU) txa() uint8 {
	cpu.A = cpu.X
	cpu.setFlag(flagZero, cpu.A == 0x00)
	cpu.setFlag(flagNegative, cpu.A&0x80 != 0)
	return 0
}

func (cpu *CPU) txs() uint8 {
	cpu.SP = cpu.X
	return 0
}

func (cpu *CPU) tya() uint8 {
	cpu.A = cpu.Y
	cpu.setFlag(flagZero, cpu.A == 0x00)
	cpu.setFlag(flagNegative, cpu.A&0x80 != 0)
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

// This method represents all illegal instructions for convenience.
func (cpu *CPU) xxx() uint8 {
	return 0
}

// Helper functions

// Complete Check if CPU has done its job
func (cpu *CPU) Complete() bool {
	return cpu.cycles == 0
}

// Disassemble Converts 6502 binary to human readable 6502 assembly. It also returns the last value of the map.
func (cpu *CPU) Disassemble(nStart uint16, nStop uint16) map[uint16]string {
	var addr uint32 = uint32(nStart)
	var value, hi, lo uint8 = 0x00, 0x00, 0x00

	// Initialize our map
	var mapLines map[uint16]string = make(map[uint16]string)

	var lineAddr uint16 = 0

	for addr <= uint32(nStop) {
		lineAddr = uint16(addr)

		sInst := "$" + ConvertToHex(uint16(addr), 4) + ": "

		opcode := cpu.Bus.CPURead(uint16(addr))
		addr++
		sInst += instructionNames[opcode] + " "

		mode := instructionModes[opcode]

		switch mode {
		case modeImplied:
			sInst += " {IMP}"
		case modeImmediate:
			value = cpu.Bus.CPURead(uint16(addr))
			addr++
			sInst += "#$" + ConvertToHex(uint16(value), 2) + " {IMM}"
		case modeZeroPage:
			lo = cpu.Bus.CPURead(uint16(addr))
			addr++
			hi = 0x00
			sInst += "$" + ConvertToHex(uint16(lo), 2) + " {ZP0}"
		case modeZeroPageX:
			lo = cpu.Bus.CPURead(uint16(addr))
			addr++
			hi = 0x00
			sInst += "$" + ConvertToHex(uint16(lo), 2) + ", X {ZPX}"
		case modeZeroPageY:
			lo = cpu.Bus.CPURead(uint16(addr))
			addr++
			hi = 0x00
			sInst += "$" + ConvertToHex(uint16(lo), 2) + ", Y {ZPY}"
		case modeIndirectX:
			lo = cpu.Bus.CPURead(uint16(addr))
			addr++
			hi = 0x00
			sInst += "($" + ConvertToHex(uint16(lo), 2) + "), X {IZX}"
		case modeIndirectY:
			lo = cpu.Bus.CPURead(uint16(addr))
			addr++
			hi = 0x00
			sInst += "($" + ConvertToHex(uint16(lo), 2) + "), Y {IZY}"
		case modeAbsolute:
			lo = cpu.Bus.CPURead(uint16(addr))
			addr++
			hi = cpu.Bus.CPURead(uint16(addr))
			addr++
			sInst += "$" + ConvertToHex(uint16(hi)<<8|uint16(lo), 4) + " {ABS}"
		case modeAbsoluteX:
			lo = cpu.Bus.CPURead(uint16(addr))
			addr++
			hi = cpu.Bus.CPURead(uint16(addr))
			addr++
			sInst += "$" + ConvertToHex(uint16(hi)<<8|uint16(lo), 4) + ", X {ABX}"
		case modeAbsoluteY:
			lo = cpu.Bus.CPURead(uint16(addr))
			addr++
			hi = cpu.Bus.CPURead(uint16(addr))
			addr++
			sInst += "$" + ConvertToHex(uint16(hi)<<8|uint16(lo), 4) + ", Y {ABY}"
		case modeIndirect:
			lo = cpu.Bus.CPURead(uint16(addr))
			addr++
			hi = cpu.Bus.CPURead(uint16(addr))
			addr++
			sInst += "($" + ConvertToHex(uint16(hi)<<8|uint16(lo), 4) + ") {IND}"
		case modeRelative:
			value = cpu.Bus.CPURead(uint16(addr))
			addr++
			sInst += "$" + ConvertToHex(uint16(value), 2) + "[$" + ConvertToHex(uint16(addr+uint32(int8(value))), 4) + "] {REL}" // Make value signed to have a correct relative address
		}
		mapLines[lineAddr] = sInst
	}

	return mapLines
}
