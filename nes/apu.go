package nes

type APU struct {
	pulse1Enable bool
	pulse1Sampe  float32
}

func ConnectAPU(bus *Bus) *APU {
	apu := APU{}

	return &apu
}

func (apu *APU) CPUWrite(addr uint16, data uint8) {
	switch addr {
	case 0x4000:

	case 0x4001:

	case 0x4002:

	case 0x4003:

	case 0x4005:

	case 0x4006:

	case 0x4007:

	case 0x4008:

	case 0x400C:

	case 0x400E:

	case 0x4015:

	case 0x400F:

	}
}

func (apu *APU) CPURead(addr uint16) uint8 {
	return 0x00
}

func (apu *APU) GetOutputSample() float32 {
	return apu.pulse1Sampe
}

func (apu *APU) Clock() {

}

func (apu *APU) Reset() {

}
