package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/net2cn/GoNES/nes"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// SDL2 variables
var windowTitle string = "GoNES SDL2"
var fontPath string = "./ui/assets/UbuntuMono-R.ttf" // Man, do not use a variable-width font! It looks too ugly with that!
var fontSize int = 15
var windowWidth, windowHeight int32 = 680, 480

var startTime time.Time = time.Now()
var endTime time.Time = time.Now()

type demoCPU struct {
	bus  *nes.Bus
	cart *nes.Cartridge

	emulationRun bool
	residualTime int64

	selectedPalette uint8

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
}

func (demo *demoCPU) drawSprite(x int, y int, sprite *sdl.Surface) {
	sprite.Blit(nil, demo.buffer, &sdl.Rect{X: int32(x), Y: int32(y)})
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

	// Init sdl2
	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		fmt.Printf("Failed to init sdl2: %s\n", err)
		panic(err)
	}

	// Initialize font
	if err = ttf.Init(); err != nil {
		fmt.Printf("Failed to init font: %s\n", err)
		panic(err)
	}

	// Load the font for our text
	if demo.font, err = ttf.OpenFont(fontPath, fontSize); err != nil {
		fmt.Printf("Failed to load font: %s\n", err)
		panic(err)
	}

	// Create window
	demo.window, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Printf("Failed to create window: %s\n", err)
		panic(err)
	}

	// Create draw surface and draw buffer
	if demo.surface, err = demo.window.GetSurface(); err != nil {
		fmt.Printf("Failed to get window surface: %s\n", err)
		panic(err)
	}

	if demo.buffer, err = demo.surface.Convert(demo.surface.Format, demo.window.GetFlags()); err != nil {
		fmt.Printf("Failed to create buffer: %s\n", err)
	}

	// Create renderer
	demo.renderer, err = sdl.CreateRenderer(demo.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("Failed to create renderer: %s\n", err)
		panic(err)
	}

	// Init our NES.
	demo.bus = nes.NewBus()

	// Load cartridge
	demo.cart, err = nes.NewCartridge("./test/nestest.nes")
	if err != nil {
		return err
	}

	demo.emulationRun = false
	demo.residualTime = 0.0

	// Insert cartridge
	demo.bus.InsertCartridge(demo.cart)

	// Disassemble ASM
	demo.mapASM = demo.bus.CPU.Disassemble(0x0000, 0xFFFF)

	// Create a collection of keys so that we can iter over.
	demo.mapKeys = make([]int, 0)
	for k := range demo.mapASM {
		demo.mapKeys = append(demo.mapKeys, int(k))
	}
	sort.Ints(demo.mapKeys)

	demo.bus.Reset()

	// Get inputLock ready for user input.
	demo.inputLock = false

	return nil
}

func (demo *demoCPU) Update(elapsedTime int64) bool {
	// Use double buffering technique to prevent flickering.

	// Check if we have reach the end of a frame
	if demo.emulationRun {
		if demo.residualTime > 0 {
			demo.residualTime -= elapsedTime
		} else {
			demo.residualTime += 1000/60 - elapsedTime
			// Golang's do while.
			for done := true; done; done = demo.bus.PPU.FrameComplete != true {
				demo.bus.Clock()
			}
			demo.bus.PPU.FrameComplete = false
		}
	}

	//Get user inputs.
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
			if !demo.inputLock {
				switch t.Keysym.Sym {
				case sdl.K_c:
					// Golang's do while.
					for done := true; done; done = demo.bus.CPU.Complete() != true {
						demo.bus.Clock()
					}
					for done := true; done; done = demo.bus.CPU.Complete() == true {
						demo.bus.Clock()
					}
				case sdl.K_f:
					for done := true; done; done = demo.bus.PPU.FrameComplete != true {
						demo.bus.Clock()
					}
					for done := true; done; done = demo.bus.CPU.Complete() != true {
						demo.bus.Clock()
					}
					demo.bus.PPU.FrameComplete = false
				case sdl.K_SPACE:
					demo.emulationRun = !demo.emulationRun
				case sdl.K_r:
					demo.bus.CPU.Reset()
				case sdl.K_p:
					demo.selectedPalette++
					demo.selectedPalette &= 0x07
				case sdl.K_d:
					if !nes.IsPathExists("./debug") {
						os.Mkdir("debug", os.ModePerm)
					}
					img.SavePNG(demo.bus.PPU.GetScreen(), "./debug/sprite.png")
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
	// Always remember to draw on buffer.

	// Draw screen, sprites.
	demo.drawSprite(0, 0, demo.bus.PPU.GetScreen())
	demo.drawSprite(416, 348, demo.bus.PPU.GetPatternTable(0, demo.selectedPalette))
	demo.drawSprite(416+132, 348, demo.bus.PPU.GetPatternTable(1, demo.selectedPalette))

	// Draw selected palettes border.
	switchSize := 6
	demo.buffer.FillRect(
		&sdl.Rect{X: int32(419 + int(demo.selectedPalette)*(switchSize*5) - 3),
			Y: 337,
			W: int32((switchSize * 5)),
			H: int32(switchSize * 2)},
		0x00FFFF00,
	)
	// Draw palettes.
	for p := 0; p < 8; p++ {
		for s := 0; s < 4; s++ {
			demo.buffer.FillRect(
				&sdl.Rect{X: int32(419 + p*(switchSize*5) + s*switchSize),
					Y: 340,
					W: int32(switchSize),
					H: int32(switchSize)},
				nes.ConvertColorToUint32(
					demo.bus.PPU.GetColorFromPaletteRAM(uint8(p), uint8(s))))
		}
	}

	// Draw CPU & ASM
	demo.drawCPU(416, 2)
	demo.drawASM(416, 72, 25)

	// Draw key hints.
	demo.drawString(0, 362, "SPACE - Run/stop", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	demo.drawString(0, 372, "R - Reset", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	demo.drawString(0, 382, "F - Step one frame", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	demo.drawString(0, 392, "C - Step one instruction", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	demo.drawString(0, 402, "D - Dump screen", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	demo.drawString(0, 412, "P - Change palette", &sdl.Color{R: 0, G: 255, B: 0, A: 0})

	// Swap buffer and present our rendered content.
	demo.window.UpdateSurface()
	demo.buffer.Blit(nil, demo.surface, nil)

	// Clear out buffer for next render round.
	demo.buffer.FillRect(nil, 0xFF000000)
	demo.renderer.Clear()

	return true
}

func (demo *demoCPU) Start() {
	elapsedTime := startTime.Sub(endTime).Milliseconds()

	running := true
	for running {
		startTime = time.Now()
		running = demo.Update(elapsedTime)
		endTime = time.Now()
		elapsedTime = endTime.Sub(startTime).Milliseconds()
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
