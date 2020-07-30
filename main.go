package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
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

type debugger struct {
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
func (debug *debugger) drawString(x int, y int, str string, color *sdl.Color) {
	if len(str) <= 0 {
		return
	}

	// Create text
	text, err := debug.font.RenderUTF8Blended(str, *color)
	if err != nil {
		fmt.Printf("Failed to create text: %s\n", err)
		return
	}
	defer text.Free()

	// Draw text, noted that you should always draw on buffer instead of directly draw on screen and blow yourself up.
	if err = text.Blit(nil, debug.buffer, &sdl.Rect{X: int32(x), Y: int32(y)}); err != nil {
		fmt.Printf("Failed to draw text: %s\n", err)
		return
	}
}

func (debug *debugger) drawSprite(x int, y int, sprite *sdl.Surface) {
	sprite.Blit(nil, debug.buffer, &sdl.Rect{X: int32(x), Y: int32(y)})
}

func (debug *debugger) drawPartialSprite(dstX int, dstY int, sprite *sdl.Surface, srcX int, srcY int, w int, h int) {
	dstRect := sdl.Rect{X: int32(dstX), Y: int32(dstY), W: int32(w), H: int32(h)}
	srcRect := sdl.Rect{X: int32(srcX), Y: int32(srcY), W: int32(w), H: int32(h)}
	sprite.Blit(&srcRect, debug.buffer, &dstRect)
}

// Draw out our RAM layouts.
func (debug *debugger) drawRAM(x int, y int, nAddr uint16, nRows int, nColumns int) {
	var nRAMX, nRAMY int = x, y
	for row := 0; row < nRows; row++ {
		sOffset := "$" + nes.ConvertToHex(nAddr, 4) + ":"
		for col := 0; col < nColumns; col++ {
			sOffset += " " + nes.ConvertToHex(uint16(debug.bus.CPURead(nAddr, true)), 2)
			nAddr++
		}
		debug.drawString(nRAMX, nRAMY, sOffset, &sdl.Color{R: 0, G: 255, B: 0, A: 0})
		nRAMY += 10
	}
}

// Draw CPU's internal state.
func (debug *debugger) drawCPU(x int, y int) {
	var color *sdl.Color = &sdl.Color{R: 0, G: 255, B: 0, A: 0}
	debug.drawString(x, y, "STATUS:", color)
	debug.drawString(x+64, y, "C", debug.getFlagColor(debug.bus.CPU.Status&(1<<0)))
	debug.drawString(x+80, y, "Z", debug.getFlagColor(debug.bus.CPU.Status&(1<<1)))
	debug.drawString(x+96, y, "I", debug.getFlagColor(debug.bus.CPU.Status&(1<<2)))
	debug.drawString(x+112, y, "D", debug.getFlagColor(debug.bus.CPU.Status&(1<<3)))
	debug.drawString(x+128, y, "B", debug.getFlagColor(debug.bus.CPU.Status&(1<<4)))
	debug.drawString(x+144, y, "U", debug.getFlagColor(debug.bus.CPU.Status&(1<<5)))
	debug.drawString(x+160, y, "O", debug.getFlagColor(debug.bus.CPU.Status&(1<<6)))
	debug.drawString(x+178, y, "N", debug.getFlagColor(debug.bus.CPU.Status&(1<<7)))
	debug.drawString(x, y+10, "PC: $"+nes.ConvertToHex(debug.bus.CPU.PC, 4), color)
	// It looks really sucks.
	debug.drawString(x, y+20, "A:  $"+nes.ConvertToHex(uint16(debug.bus.CPU.A), 2)+" ["+nes.ConvertUint8ToString(debug.bus.CPU.A)+"]", color)
	debug.drawString(x, y+30, "X:  $"+nes.ConvertToHex(uint16(debug.bus.CPU.X), 2)+" ["+nes.ConvertUint8ToString(debug.bus.CPU.X)+"]", color)
	debug.drawString(x, y+40, "Y:  $"+nes.ConvertToHex(uint16(debug.bus.CPU.Y), 2)+" ["+nes.ConvertUint8ToString(debug.bus.CPU.Y)+"]", color)
	debug.drawString(x, y+50, "SP: $"+nes.ConvertToHex(uint16(debug.bus.CPU.SP), 4), color)
}

// Draw disassembled code.
func (debug *debugger) drawASM(x int, y int, nLines int) {
	// I hate it without sorted map. Lots of hacky stuffs happen here.
	itA := debug.mapASM[debug.bus.CPU.PC]
	var nLineY int = (nLines>>1)*10 + y
	var idx int = sort.SearchInts(debug.mapKeys, int(debug.bus.CPU.PC))
	if itA != debug.mapASM[uint16(debug.mapKeys[len(debug.mapKeys)-1])] {
		debug.drawString(x, nLineY, itA, &sdl.Color{R: 255, G: 255, B: 0, A: 0})
		for nLineY < (nLines*10)+y {
			nLineY += 10
			idx++
			if idx != len(debug.mapKeys)-1 {
				debug.drawString(x, nLineY, debug.mapASM[uint16(debug.mapKeys[idx])], &sdl.Color{R: 0, G: 255, B: 0, A: 0})
			}
		}
	}
	idx = sort.SearchInts(debug.mapKeys, int(debug.bus.CPU.PC))
	nLineY = (nLines>>1)*10 + y
	if itA != debug.mapASM[uint16(debug.mapKeys[len(debug.mapKeys)-1])] {
		debug.drawString(x, nLineY, itA, &sdl.Color{R: 255, G: 255, B: 0, A: 0})
		for nLineY > y {
			nLineY -= 10
			idx--
			// Check if our index is out of range since we're subtracting it.
			if idx != len(debug.mapKeys)-1 && idx > 0 {
				debug.drawString(x, nLineY, debug.mapASM[uint16(debug.mapKeys[idx])], &sdl.Color{R: 0, G: 255, B: 0, A: 0})
			}
		}
	}
}

func (debug *debugger) drawNameTable(x int, y int, nameTable *sdl.Surface) {
	for v := 0; v < 30; v++ {
		for h := 0; h < 32; h++ {
			var id uint8 = uint8(uint32(debug.bus.PPU.TableName[0][v*32+h]))

			debug.drawPartialSprite(x+h*8, y+v*8, nameTable,
				int(id&0x0F<<3), int((id>>4)&0x0F<<3),
				8, 8)
		}
	}
}

// Test what color should we use based on flag's state.
func (debug *debugger) getFlagColor(flag uint8) *sdl.Color {
	if flag == 0 {
		return &sdl.Color{R: 0, G: 255, B: 0, A: 0}
	}
	return &sdl.Color{R: 255, G: 0, B: 0, A: 0}
}

// Construct our debug.
func (debug *debugger) Construct(filePath string, width int32, height int32) error {
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
	if debug.font, err = ttf.OpenFont(fontPath, fontSize); err != nil {
		fmt.Printf("Failed to load font: %s\n", err)
		panic(err)
	}

	// Create window
	debug.window, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Printf("Failed to create window: %s\n", err)
		panic(err)
	}

	// Create draw surface and draw buffer
	if debug.surface, err = debug.window.GetSurface(); err != nil {
		fmt.Printf("Failed to get window surface: %s\n", err)
		panic(err)
	}

	if debug.buffer, err = debug.surface.Convert(debug.surface.Format, debug.window.GetFlags()); err != nil {
		fmt.Printf("Failed to create buffer: %s\n", err)
	}

	// Create renderer
	debug.renderer, err = sdl.CreateRenderer(debug.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("Failed to create renderer: %s\n", err)
		panic(err)
	}

	// Init our NES.
	debug.bus = nes.NewBus()

	// Load cartridge
	debug.cart, err = nes.NewCartridge(filePath)
	if err != nil {
		return err
	}

	debug.emulationRun = false
	debug.residualTime = 0.0

	// Insert cartridge
	debug.bus.InsertCartridge(debug.cart)

	// Disassemble ASM
	debug.mapASM = debug.bus.CPU.Disassemble(0x0000, 0xFFFF)

	// Create a collection of keys so that we can iter over.
	debug.mapKeys = make([]int, 0)
	for k := range debug.mapASM {
		debug.mapKeys = append(debug.mapKeys, int(k))
	}
	sort.Ints(debug.mapKeys)

	debug.bus.Reset()

	// Get inputLock ready for user input.
	debug.inputLock = false

	return nil
}

func (debug *debugger) Update(elapsedTime int64) bool {
	// Use double buffering technique to prevent flickering.

	// Get user inputs.
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
			if !debug.inputLock {
				switch t.Keysym.Sym {
				// NES controller
				// Bugged...
				case sdl.K_s:
					debug.bus.Controller[0] |= 0x10
				case sdl.K_UP:
					debug.bus.Controller[0] |= 0x08
				case sdl.K_DOWN:
					debug.bus.Controller[0] |= 0x04
				case sdl.K_LEFT:
					debug.bus.Controller[0] |= 0x02
				case sdl.K_RIGHT:
					debug.bus.Controller[0] |= 0x01
				// Debug utility.
				case sdl.K_c:
					// Golang's do while.
					for done := true; done; done = debug.bus.CPU.Complete() != true {
						debug.bus.Clock()
					}
					for done := true; done; done = debug.bus.CPU.Complete() == true {
						debug.bus.Clock()
					}
				case sdl.K_f:
					for done := true; done; done = debug.bus.PPU.FrameComplete != true {
						debug.bus.Clock()
					}
					for done := true; done; done = debug.bus.CPU.Complete() != true {
						debug.bus.Clock()
					}
					debug.bus.PPU.FrameComplete = false
				case sdl.K_SPACE:
					debug.emulationRun = !debug.emulationRun
				case sdl.K_r:
					debug.bus.CPU.Reset()
				case sdl.K_p:
					debug.selectedPalette++
					debug.selectedPalette &= 0x07
				case sdl.K_d:
					if !nes.IsPathExists("./debug") {
						os.Mkdir("debug", os.ModePerm)
					}
					img.SavePNG(debug.bus.PPU.GetScreen(), "./debug/sprite.png")
				}
			}

			// Anti-jittering
			if t.Repeat > 0 {
				debug.inputLock = false
			} else {
				if t.State == sdl.RELEASED {
					debug.inputLock = false
				} else if t.State == sdl.PRESSED {
					debug.inputLock = true
				}
			}
		}
	}

	// Check if we have reach the end of a frame
	if debug.emulationRun {
		if debug.residualTime > 0 {
			debug.residualTime -= elapsedTime
		} else {
			debug.residualTime += 1000/60 - elapsedTime
			// Golang's do while.
			for done := true; done; done = debug.bus.PPU.FrameComplete != true {
				debug.bus.Clock()
			}
			debug.bus.PPU.FrameComplete = false
		}
	}

	// Render stuffs.
	// Always remember to draw on buffer.

	// Draw screen, sprites.
	debug.drawSprite(0, 0, debug.bus.PPU.GetScreen())
	// Quick hack to render background tiles
	// nameTable := debug.bus.PPU.GetPatternTable(0, debug.selectedPalette)
	// debug.drawNameTable(0, 0, nameTable)

	debug.drawSprite(416, 349, debug.bus.PPU.GetPatternTable(0, debug.selectedPalette))
	debug.drawSprite(416+132, 349, debug.bus.PPU.GetPatternTable(1, debug.selectedPalette))

	// Draw selected palettes border.
	switchSize := 6
	debug.buffer.FillRect(
		&sdl.Rect{X: int32(419 + int(debug.selectedPalette)*(switchSize*5) - 3),
			Y: 337,
			W: int32((switchSize * 5)),
			H: int32(switchSize * 2)},
		0x00FFFF00,
	)

	// Draw palettes.
	for p := 0; p < 8; p++ {
		for s := 0; s < 4; s++ {
			debug.buffer.FillRect(
				&sdl.Rect{X: int32(419 + p*(switchSize*5) + s*switchSize),
					Y: 340,
					W: int32(switchSize),
					H: int32(switchSize)},
				nes.ConvertColorToUint32(
					debug.bus.PPU.GetColorFromPaletteRAM(uint8(p), uint8(s))))
		}
	}

	// Draw CPU & ASM
	debug.drawCPU(416, 2)
	// debug.drawASM(416, 72, 25)

	// Draw DMA
	for i := 0; i < 26; i++ {
		s := nes.ConvertToHex(uint16(i), 2) + ": (" + strconv.Itoa(int(debug.bus.PPU.OAM[i*4+3])) + "," +
			strconv.Itoa(int(debug.bus.PPU.OAM[i*4+0])) + ")" +
			"ID: " + nes.ConvertToHex(uint16(debug.bus.PPU.OAM[i*4+1]), 2) +
			"AT: " + nes.ConvertToHex(uint16(debug.bus.PPU.OAM[i*4+2]), 2)
		debug.drawString(416, 72+i*10, s, &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	}

	// Draw key hints.
	debug.drawString(0, 362, "SPACE - Run/stop", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(0, 372, "R - Reset", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(0, 382, "F - Step one frame", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(0, 392, "C - Step one instruction", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(0, 402, "D - Dump screen", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(0, 412, "P - Change palette", &sdl.Color{R: 0, G: 255, B: 0, A: 0})

	fps := 0
	if elapsedTime != 0 {
		fps = int(1000 / elapsedTime)
	}
	debug.drawString(612, 2, "FPS: "+strconv.Itoa(fps), &sdl.Color{R: 0, G: 255, B: 0, A: 0})

	// Swap buffer and present our rendered content.
	debug.window.UpdateSurface()
	debug.buffer.Blit(nil, debug.surface, nil)

	// Clear out buffer for next render round.
	debug.buffer.FillRect(nil, 0xFF000000)
	debug.renderer.Clear()

	return true
}

func (debug *debugger) Start() {
	elapsedTime := startTime.Sub(endTime).Milliseconds()

	running := true
	for running {
		startTime = time.Now()
		running = debug.Update(elapsedTime)
		endTime = time.Now()
		elapsedTime = endTime.Sub(startTime).Milliseconds()
	}

}

func main() {
	fmt.Println(windowTitle)

	// I really enjoy its graphics. I mean the anime movie.
	fmt.Println("HELLO WORLD -ALLTALE-")
	fmt.Println("With programming we have god's hand.")
	debug := debugger{}
	err := debug.Construct("./roms/dk.nes", windowWidth, windowHeight)
	if err != nil {
		return
	}
	debug.Start()
}
