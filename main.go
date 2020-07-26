package main

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/net2cn/GoNES/nes"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// SDL2 variables
var windowTitle string = "GoNES SDL2"
var fontPath string = "./ui/assets/UbuntuMono-R.ttf" // Man, do not use a variable-width font! It looks too ugly with that!
var fontSize int = 15
var windowWidth, windowHeight int32 = 680, 480

type demoCPU struct {
	bus       *nes.Bus
	mapASM    map[uint16]string
	mapKeys   []int
	inputLock bool

	window   *sdl.Window
	surface  *sdl.Surface
	buffer   *sdl.Surface
	renderer *sdl.Renderer
	font     *ttf.Font
}

// Private draw functions
func (demo *demoCPU) drawString(x int, y int, str string, color *sdl.Color) {
	if len(str) <= 0 {
		return
	}

	// Create text
	text, err := demo.font.RenderUTF8Blended(str, *color)
	if err != nil {
		fmt.Printf("Failed to create text: %s\n", err)
		return
	}
	defer text.Free()

	// Draw text, noted that you should always draw on buffer instead of directly draw on screen and blow yourself up.
	if err = text.Blit(nil, demo.buffer, &sdl.Rect{X: int32(x), Y: int32(y)}); err != nil {
		fmt.Printf("Failed to draw text: %s\n", err)
		return
	}

	demo.window.UpdateSurface()
}

// Draw out our RAM layouts.
func (demo *demoCPU) drawRAM(x int, y int, nAddr uint16, nRows int, nColumns int) {
	var nRAMX, nRAMY int = x, y
	for row := 0; row < nRows; row++ {
		sOffset := "$" + nes.ConvertToHex(nAddr, 4) + ":"
		for col := 0; col < nColumns; col++ {
			sOffset += " " + nes.ConvertToHex(uint16(demo.bus.CPURead(nAddr, true)), 2)
			nAddr++
		}
		demo.drawString(nRAMX, nRAMY, sOffset, &sdl.Color{R: 0, G: 255, B: 0, A: 0})
		nRAMY += 10
	}
}

// Draw CPU's internal state.
func (demo *demoCPU) drawCPU(x int, y int) {
	var color *sdl.Color = &sdl.Color{R: 0, G: 255, B: 0, A: 0}
	demo.drawString(x, y, "STATUS:", color)
	demo.drawString(x+64, y, "C", demo.getFlagColor(demo.bus.CPU.Status&(1<<0)))
	demo.drawString(x+80, y, "Z", demo.getFlagColor(demo.bus.CPU.Status&(1<<1)))
	demo.drawString(x+96, y, "I", demo.getFlagColor(demo.bus.CPU.Status&(1<<2)))
	demo.drawString(x+112, y, "D", demo.getFlagColor(demo.bus.CPU.Status&(1<<3)))
	demo.drawString(x+128, y, "B", demo.getFlagColor(demo.bus.CPU.Status&(1<<4)))
	demo.drawString(x+144, y, "U", demo.getFlagColor(demo.bus.CPU.Status&(1<<5)))
	demo.drawString(x+160, y, "O", demo.getFlagColor(demo.bus.CPU.Status&(1<<6)))
	demo.drawString(x+178, y, "N", demo.getFlagColor(demo.bus.CPU.Status&(1<<7)))
	demo.drawString(x, y+10, "PC: $"+nes.ConvertToHex(demo.bus.CPU.PC, 4), color)
	// It looks really sucks.
	demo.drawString(x, y+20, "A:  $"+nes.ConvertToHex(uint16(demo.bus.CPU.A), 2)+" ["+nes.ConvertUint8ToString(demo.bus.CPU.A)+"]", color)
	demo.drawString(x, y+30, "X:  $"+nes.ConvertToHex(uint16(demo.bus.CPU.X), 2)+" ["+nes.ConvertUint8ToString(demo.bus.CPU.X)+"]", color)
	demo.drawString(x, y+40, "Y:  $"+nes.ConvertToHex(uint16(demo.bus.CPU.Y), 2)+" ["+nes.ConvertUint8ToString(demo.bus.CPU.Y)+"]", color)
	demo.drawString(x, y+50, "SP: $"+nes.ConvertToHex(uint16(demo.bus.CPU.SP), 4), color)
}

// Draw disassembled code.
func (demo *demoCPU) drawASM(x int, y int, nLines int) {
	// I hate it without sorted map. Lots of hacky stuffs happen here.
	itA := demo.mapASM[demo.bus.CPU.PC]
	var nLineY int = (nLines>>1)*10 + y
	var idx int = sort.SearchInts(demo.mapKeys, int(demo.bus.CPU.PC))
	if itA != demo.mapASM[uint16(demo.mapKeys[len(demo.mapKeys)-1])] {
		demo.drawString(x, nLineY, itA, &sdl.Color{R: 255, G: 255, B: 0, A: 0})
		for nLineY < (nLines*10)+y {
			nLineY += 10
			idx++
			if idx != len(demo.mapKeys)-1 {
				demo.drawString(x, nLineY, demo.mapASM[uint16(demo.mapKeys[idx])], &sdl.Color{R: 0, G: 255, B: 0, A: 0})
			}
		}
	}
	idx = sort.SearchInts(demo.mapKeys, int(demo.bus.CPU.PC))
	nLineY = (nLines>>1)*10 + y
	if itA != demo.mapASM[uint16(demo.mapKeys[len(demo.mapKeys)-1])] {
		demo.drawString(x, nLineY, itA, &sdl.Color{R: 255, G: 255, B: 0, A: 0})
		for nLineY > y {
			nLineY -= 10
			idx--
			// Check if our index is out of range since we're subtracting it.
			if idx != len(demo.mapKeys)-1 && idx > 0 {
				demo.drawString(x, nLineY, demo.mapASM[uint16(demo.mapKeys[idx])], &sdl.Color{R: 0, G: 255, B: 0, A: 0})
			}
		}
	}
}

// Test what color should we use based on flag's state.
func (demo *demoCPU) getFlagColor(flag uint8) *sdl.Color {
	if flag == 0 {
		return &sdl.Color{R: 0, G: 255, B: 0, A: 0}
	}
	return &sdl.Color{R: 255, G: 0, B: 0, A: 0}
}

// Construct our demo.
func (demo *demoCPU) Construct(width int32, height int32) error {
	var err error

	// Init our NES.
	demo.bus = nes.NewBus()

	// Init sdl2
	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		fmt.Printf("Failed to init sdl2: %s\n", err)
		return err
	}

	// Initialize font
	if err = ttf.Init(); err != nil {
		fmt.Printf("Failed to init font: %s\n", err)
		return err
	}

	// Load the font for our text
	if demo.font, err = ttf.OpenFont(fontPath, fontSize); err != nil {
		fmt.Printf("Failed to load font: %s\n", err)
		return err
	}

	// Create window
	demo.window, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Printf("Failed to create window: %s\n", err)
		return err
	}

	// Create draw surface and draw buffer
	if demo.surface, err = demo.window.GetSurface(); err != nil {
		fmt.Printf("Failed to get window surface: %s\n", err)
		return err
	}

	if demo.buffer, err = demo.surface.Convert(demo.surface.Format, demo.window.GetFlags()); err != nil {
		fmt.Printf("Failed to create buffer: %s\n", err)
	}

	// Create renderer
	demo.renderer, err = sdl.CreateRenderer(demo.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("Failed to create renderer: %s\n", err)
		return err
	}

	// ASM code
	// this snippet of codes calculates 10 * 3
	// Redundant codes are for validating if those instructions are
	// implemented correctly
	// *=$8000
	// LDX #10
	// STX $0000
	// LDX #3
	// STX $0001
	// LDY $0000
	// LDA #0
	// CLC
	// loop
	// ADC $0001
	// DEY
	// BNE loop
	// STA $0002
	// NOP
	// NOP
	// NOP
	var binStr string = "A2 0A 8E 00 00 A2 03 8E 01 00 AC 00 00 A9 00 18 6D 01 00 88 D0 FA 8D 02 00 EA EA EA" // binary represent
	binStr = strings.ReplaceAll(binStr, " ", "")

	bin, err := hex.DecodeString(binStr)
	if err != nil {
		fmt.Printf("Exception occurred when decoding string.")
	}

	nOffset := 0x8000

	for _, b := range bin {
		demo.bus.CPURAM[nOffset] = b
		nOffset++
	}

	demo.bus.CPURAM[0xFFFC] = 0x00
	demo.bus.CPURAM[0xFFFD] = 0x80

	demo.mapASM = demo.bus.CPU.Disassemble(0x0000, 0xFFFF)
	// Create a collection of keys so that we can iter over.
	demo.mapKeys = make([]int, 0)
	for k := range demo.mapASM {
		demo.mapKeys = append(demo.mapKeys, int(k))
	}
	sort.Ints(demo.mapKeys)

	demo.bus.CPU.Reset()

	// Get inputLock ready for user input.
	demo.inputLock = false

	return nil
}

func (demo *demoCPU) Update() bool {
	// Using double buffering technique to prevent flickering.

	//Get user inputs.
	// I should probably split out this part and make this looks neater.
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
			if !demo.inputLock {
				switch t.Keysym.Sym {
				case sdl.K_SPACE:
					// Golang's do while.
					for done := true; done; done = demo.bus.CPU.Complete() != true {
						demo.bus.CPU.Clock()
					}
				case sdl.K_r:
					demo.bus.CPU.Reset()
				case sdl.K_i:
					demo.bus.CPU.Irq()
				case sdl.K_n:
					demo.bus.CPU.Nmi()
				}
			}

			// Anti-jittering
			if t.Repeat > 0 {
				demo.inputLock = false
			} else {
				if t.State == sdl.RELEASED {
					demo.inputLock = false
				} else if t.State == sdl.PRESSED {
					demo.inputLock = true
				}
			}

		}
	}

	// Render stuffs.
	demo.drawRAM(2, 2, 0x0000, 16, 16)
	demo.drawRAM(2, 182, 0x8000, 16, 16)
	demo.drawCPU(448, 2)
	demo.drawASM(448, 72, 26)
	demo.drawString(2, 362, "SPACE - step one, R - reset, I - IRQ, N - NMI", &sdl.Color{R: 0, G: 255, B: 0, A: 0})

	// Swap buffer and present our rendered content.
	demo.buffer.Blit(nil, demo.surface, nil)
	demo.buffer.FillRect(nil, 0xFF000000)

	return true
}

func (demo *demoCPU) Start() {
	running := true
	for running {
		running = demo.Update()
		sdl.Delay(16)
	}

}

func main() {
	fmt.Println(windowTitle)

	// I really enjoy its graphics. I mean the anime movie.
	fmt.Println("HELLO WORLD -ALLTALE-")
	fmt.Println("With programming we have god's hand.")
	demo := demoCPU{}
	err := demo.Construct(windowWidth, windowHeight)
	if err != nil {
		return
	}
	demo.Start()
}
