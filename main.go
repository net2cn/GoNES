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
var fontSize int = 14
var windowWidth, windowHeight int32 = 680, 480

type demoCPU struct {
	bus     *nes.Bus
	mapASM  map[uint16]string
	mapKeys []int

	window   *sdl.Window
	surface  *sdl.Surface
	buffer   *sdl.Surface
	renderer *sdl.Renderer
	font     *ttf.Font
}

func (demo *demoCPU) drawString(x int, y int, str string, color *sdl.Color) {
	// Create text
	text, err := demo.font.RenderUTF8Blended(str, *color)
	if err != nil {
		fmt.Printf("Failed to create text: %s\n", err)
		return
	}
	defer text.Free()

	// Draw text
	if err = text.Blit(nil, demo.buffer, &sdl.Rect{X: int32(x), Y: int32(y)}); err != nil {
		fmt.Printf("Failed to draw text: %s\n", err)
		return
	}

	demo.window.UpdateSurface()
}

// Private draw functions
func (demo *demoCPU) drawRAM(x int, y int, nAddr uint16, nRows int, nColumns int) {
	var nRAMX, nRAMY int = x, y
	for row := 0; row < nRows; row++ {
		sOffset := "$" + nes.ConvertToHex(nAddr, 4) + ":"
		for col := 0; col < nColumns; col++ {
			sOffset += " " + nes.ConvertToHex(uint16(demo.bus.Read(nAddr, true)), 2)
			nAddr++
		}
		demo.drawString(nRAMX, nRAMY, sOffset, &sdl.Color{0, 255, 0, 0})
		nRAMY += 10
	}
}

func (demo *demoCPU) drawCPU(x int, y int) {
	var color *sdl.Color = &sdl.Color{0, 255, 0, 0}
	demo.drawString(x, y, "STATUS:", color)
	demo.drawString(x+64, y, "C", demo.getFlagColor(demo.bus.CPU.Status&1))
	demo.drawString(x+80, y, "Z", demo.getFlagColor(demo.bus.CPU.Status&2))
	demo.drawString(x+96, y, "I", demo.getFlagColor(demo.bus.CPU.Status&3))
	demo.drawString(x+112, y, "D", demo.getFlagColor(demo.bus.CPU.Status&4))
	demo.drawString(x+128, y, "B", demo.getFlagColor(demo.bus.CPU.Status&5))
	demo.drawString(x+144, y, "U", demo.getFlagColor(demo.bus.CPU.Status&6))
	demo.drawString(x+160, y, "O", demo.getFlagColor(demo.bus.CPU.Status&7))
	demo.drawString(x+178, y, "N", demo.getFlagColor(demo.bus.CPU.Status&8))
	demo.drawString(x, y+10, "PC: $"+nes.ConvertToHex(demo.bus.CPU.PC, 4), color)
	demo.drawString(x, y+20, "A:  $"+nes.ConvertToHex(uint16(demo.bus.CPU.A), 2), color)
	demo.drawString(x, y+30, "X:  $"+nes.ConvertToHex(uint16(demo.bus.CPU.X), 2), color)
	demo.drawString(x, y+40, "Y:  $"+nes.ConvertToHex(uint16(demo.bus.CPU.Y), 2), color)
	demo.drawString(x, y+50, "SP: $"+nes.ConvertToHex(uint16(demo.bus.CPU.SP), 4), color)
}

func (demo *demoCPU) drawASM(x int, y int, nLines int) {
	itA := demo.mapASM[demo.bus.CPU.PC]
	var nLineY int = (nLines>>1)*10 + y
	var idx int = sort.SearchInts(demo.mapKeys, int(demo.bus.CPU.PC))
	if itA != demo.mapASM[uint16(demo.mapKeys[len(demo.mapKeys)-1])] {
		demo.drawString(x, nLineY, itA, &sdl.Color{0, 0, 255, 0})
		for nLineY < (nLines*10)+y {
			nLineY += 10
			idx++
			if idx != len(demo.mapKeys)-1 {
				demo.drawString(x, nLineY, demo.mapASM[uint16(demo.mapKeys[idx])], &sdl.Color{0, 255, 0, 0})
			}
		}
	}
	idx = sort.SearchInts(demo.mapKeys, int(demo.bus.CPU.PC))
	nLineY = (nLines>>1)*10 + y
	if itA != demo.mapASM[uint16(demo.mapKeys[len(demo.mapKeys)-1])] {
		demo.drawString(x, nLineY, itA, &sdl.Color{0, 0, 255, 0})
		for nLineY > y {
			nLineY -= 10
			idx--
			if idx != len(demo.mapKeys)-1 && idx > 0 {
				demo.drawString(x, nLineY, demo.mapASM[uint16(demo.mapKeys[idx])], &sdl.Color{0, 255, 0, 0})
			}
		}
	}
}

func (demo *demoCPU) getFlagColor(flag uint8) *sdl.Color {
	if flag == 0 {
		return &sdl.Color{0, 255, 0, 0}
	}
	return &sdl.Color{255, 0, 0, 0}
}

func (demo *demoCPU) Construct(width int32, height int32) error {
	var err error

	demo.bus = nes.NewBus()

	// Init sdl2
	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		fmt.Printf("Failed to init sdl2: %s\n", err)
		return err
	}
	// defer sdl.Quit()

	// Initialize font
	if err = ttf.Init(); err != nil {
		fmt.Printf("Failed to init font: %s\n", err)
		return err
	}
	// defer ttf.Quit()

	// Create window
	demo.window, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Printf("Failed to create window: %s\n", err)
		return err
	}
	// defer demo.window.Destroy()

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
	// defer demo.renderer.Destroy()

	// Load the font for our text
	if demo.font, err = ttf.OpenFont(fontPath, fontSize); err != nil {
		fmt.Printf("Failed to load font: %s\n", err)
		return err
	}
	// defer demo.font.Close()

	// running := true
	// for running {
	// 	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
	// 		switch event.(type) {
	// 		case *sdl.QuitEvent:
	// 			running = false
	// 		}
	// 	}

	// 	sdl.Delay(16)
	// }

	return nil
}

func (demo *demoCPU) Update() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
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
	}

	demo.drawRAM(2, 2, 0x0000, 16, 16)
	demo.drawRAM(2, 182, 0x8000, 16, 16)
	demo.drawCPU(448, 2)
	demo.drawASM(448, 72, 26)

	// Swap buffer
	demo.buffer.Blit(nil, demo.surface, nil)
	demo.buffer.FillRect(nil, 0)

	return true
}

func (demo *demoCPU) Start() {
	var binStr string = "A2 0A 8E 00 00 A2 03 8E 01 00 AC 00 00 A9 00 18 6D 01 00 88 D0 FA 8D 02 00 EA EA EA"
	binStr = strings.ReplaceAll(binStr, " ", "")
	bin, err := hex.DecodeString(binStr)
	if err != nil {
		fmt.Printf("Exception occurred when decoding string.")
	}

	nOffset := 0x8000

	for _, b := range bin {
		demo.bus.RAM[nOffset] = b
		nOffset++
	}

	demo.bus.RAM[0xFFFC] = 0x00
	demo.bus.RAM[0xFFFD] = 0x80

	demo.mapASM = demo.bus.CPU.Disassemble(0x0000, 0xFFFF)
	demo.mapKeys = make([]int, 0)
	for k := range demo.mapASM {
		demo.mapKeys = append(demo.mapKeys, int(k))
	}
	sort.Ints(demo.mapKeys)

	demo.bus.CPU.Reset()

	running := true
	for running {
		running = demo.Update()
		sdl.Delay(16)
	}

}

func main() {
	fmt.Println(windowTitle)

	// I really enjoy its graphics.
	fmt.Println("HELLO WORLD -ALLTALE-")
	fmt.Println("With programming we have god's hand.")
	demo := demoCPU{}
	err := demo.Construct(windowWidth, windowHeight)
	if err != nil {
		return
	}
	demo.Start()
}
