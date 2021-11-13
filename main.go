package main

// typedef unsigned char Uint8;
// void SoundOut(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"runtime/pprof"
	"sort"
	"strconv"
	"time"
	"unsafe"

	"github.com/net2cn/GoNES/nes"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// SDL2 variables
var windowTitle string = "GoNES SDL2"
var fontPath string = "./assets/UbuntuMono-R.ttf" // Man, do not use a variable-width font! It looks too ugly with that!
var fontSize int = 15
var windowWidth, windowHeight int32 = 680, 480

// Timer.
var startTime time.Time = time.Now()
var endTime time.Time = time.Now()
var elapsedTime int64 = 0

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

// Draw a full sprite.
func (debug *debugger) drawSprite(x int, y int, sprite *sdl.Surface) {
	sprite.Blit(nil, debug.buffer, &sdl.Rect{X: int32(x), Y: int32(y)})
}

// Draw a part of sprite.
func (debug *debugger) drawPartialSprite(dstX int, dstY int, sprite *sdl.Surface, srcX int, srcY int, w int, h int) {
	dstRect := sdl.Rect{X: int32(dstX), Y: int32(dstY), W: int32(w), H: int32(h)}
	srcRect := sdl.Rect{X: int32(srcX), Y: int32(srcY), W: int32(w), H: int32(h)}
	sprite.Blit(&srcRect, debug.buffer, &dstRect)
}

// Draw out our RAM layouts.
func (debug *debugger) drawRAM(x int, y int, addr uint16, rows int, columns int) {
	var nRAMX, nRAMY int = x, y
	for row := 0; row < rows; row++ {
		sOffset := "$" + nes.ConvertToHex(addr, 4) + ":"
		for col := 0; col < columns; col++ {
			sOffset += " " + nes.ConvertToHex(uint16(debug.bus.CPURead(addr, true)), 2)
			addr++
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
func (debug *debugger) drawASM(x int, y int, lines int) {
	// I hate it without sorted map. Lots of hacky stuffs happen here.
	// Draw backwards.
	itA := debug.mapASM[debug.bus.CPU.PC]
	var nLineY int = (lines>>1)*10 + y
	var idx int = sort.SearchInts(debug.mapKeys, int(debug.bus.CPU.PC))
	if itA != debug.mapASM[uint16(debug.mapKeys[len(debug.mapKeys)-1])] {
		debug.drawString(x, nLineY, itA, &sdl.Color{R: 255, G: 255, B: 0, A: 0})
		for nLineY < (lines*10)+y {
			nLineY += 10
			idx++
			if idx != len(debug.mapKeys)-1 {
				debug.drawString(x, nLineY, debug.mapASM[uint16(debug.mapKeys[idx])], &sdl.Color{R: 0, G: 255, B: 0, A: 0})
			}
		}
	}
	// Draw forwards.
	idx = sort.SearchInts(debug.mapKeys, int(debug.bus.CPU.PC))
	nLineY = (lines>>1)*10 + y
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

// Draw DMA state.
func (debug *debugger) drawDMA(x int, y int, lines int) {
	// Draw DMA
	for i := 0; i < lines; i++ {
		s := nes.ConvertToHex(uint16(i), 2) + ": (" + strconv.Itoa(int(debug.bus.PPU.OAM[i*4+3])) + "," +
			strconv.Itoa(int(debug.bus.PPU.OAM[i*4+0])) + ")" +
			" ID: " + nes.ConvertToHex(uint16(debug.bus.PPU.OAM[i*4+1]), 2) +
			" AT: " + nes.ConvertToHex(uint16(debug.bus.PPU.OAM[i*4+2]), 2)
		debug.drawString(x, y+i*10, s, &sdl.Color{R: 0, G: 255, B: 0, A: 0})
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
		fmt.Printf("Failed to init sdl2 video: %s\n", err)
		panic(err)
	}

	if err = sdl.Init(sdl.INIT_AUDIO); err != nil {
		fmt.Printf("Failed to init sdl2 audio: %s\n", err)
		panic(err)
	}

	spec := &sdl.AudioSpec{
		Freq:     48000,
		Format:   sdl.AUDIO_U8,
		Channels: 1,
		Samples:  4096,
		Callback: sdl.AudioCallback(C.SoundOut),
	}

	if err := sdl.OpenAudio(spec, nil); err != nil {
		fmt.Printf("Failed to open sdl2 audio: %s\n", err)
		panic(err)
	}

	sdl.PauseAudio(true)

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

	// Create renderer
	debug.renderer, err = sdl.CreateRenderer(debug.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("Failed to create renderer: %s\n", err)
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
	// This is because Golang does not provide a sortable map. (They use hashmap instead of rb tree, which truly ummmms me.)
	debug.mapKeys = make([]int, 0)
	for k := range debug.mapASM {
		debug.mapKeys = append(debug.mapKeys, int(k))
	}
	// Sort our map keys to make sure it is in order so that we can consturct our tricky iterator.
	sort.Ints(debug.mapKeys)

	debug.bus.Reset()

	// Get inputLock ready for user input.
	debug.inputLock = false

	return nil
}

//export SoundOut
func SoundOut(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))

	var phase float64
	for i := 0; i < n; i += 2 {
		phase += 2 * math.Pi * 440 / 48000
		sample := C.Uint8((math.Sin(phase) + 0.999999) * 128)
		buf[i] = sample
		buf[i+1] = sample
	}
}

func (debug *debugger) Update(elapsedTime int64) bool {
	// Use double buffering technique to prevent flickering.
	// Render sequence:
	// Acquire inputs,
	// Emulate,
	// Rneder to buffer,
	// Swap screen and buffer (present),
	// Clear buffer for next round.

	// Get NES controller inputs.
	keyState := sdl.GetKeyboardState()

	debug.bus.Controller[0] = 0x00
	if keyState[sdl.SCANCODE_X] != 0 {
		debug.bus.Controller[0] |= 0x80
	} else {
		debug.bus.Controller[0] |= 0x00
	}

	if keyState[sdl.SCANCODE_Z] != 0 {
		debug.bus.Controller[0] |= 0x40
	} else {
		debug.bus.Controller[0] |= 0x00
	}

	if keyState[sdl.SCANCODE_A] != 0 {
		debug.bus.Controller[0] |= 0x20
	} else {
		debug.bus.Controller[0] |= 0x00
	}

	if keyState[sdl.SCANCODE_S] != 0 {
		debug.bus.Controller[0] |= 0x10
	} else {
		debug.bus.Controller[0] |= 0x00
	}

	if keyState[sdl.SCANCODE_UP] != 0 {
		debug.bus.Controller[0] |= 0x08
	} else {
		debug.bus.Controller[0] |= 0x00
	}

	if keyState[sdl.SCANCODE_DOWN] != 0 {
		debug.bus.Controller[0] |= 0x04
	} else {
		debug.bus.Controller[0] |= 0x00
	}

	if keyState[sdl.SCANCODE_LEFT] != 0 {
		debug.bus.Controller[0] |= 0x02
	} else {
		debug.bus.Controller[0] |= 0x00
	}

	if keyState[sdl.SCANCODE_RIGHT] != 0 {
		debug.bus.Controller[0] |= 0x01
	} else {
		debug.bus.Controller[0] |= 0x00
	}

	// Get debugger inputs.
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
			if !debug.inputLock {
				switch t.Keysym.Sym {
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
					sdl.PauseAudio(!debug.emulationRun)
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
				// Held.
				debug.inputLock = false
			} else {
				// Pressed once.
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
			debug.residualTime += 1000000/60 - elapsedTime
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
	debug.drawDMA(416, 72, 25)

	// Draw key hints.
	debug.drawString(2, 306, "NES", &sdl.Color{R: 255, G: 255, B: 0, A: 0})
	debug.drawString(2, 316, "Z - A", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(2, 326, "X - B", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(2, 336, "A - Select", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(2, 346, "S - Start", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(2, 356, "UP - Up", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(2, 366, "DOWN - Down", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(2, 376, "LEFT - Left", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(2, 386, "RIGHT - Right", &sdl.Color{R: 0, G: 255, B: 0, A: 0})

	debug.drawString(2, 406, "Debugger", &sdl.Color{R: 255, G: 255, B: 0, A: 0})
	debug.drawString(2, 416, "SPACE - Run/stop", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(2, 426, "R - Reset", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(2, 436, "F - Step one frame", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(2, 446, "C - Step one instruction", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(2, 456, "D - Dump screen", &sdl.Color{R: 0, G: 255, B: 0, A: 0})
	debug.drawString(2, 466, "P - Change palette", &sdl.Color{R: 0, G: 255, B: 0, A: 0})

	// Swap buffer and present our rendered content.
	debug.surface, debug.buffer = debug.buffer, debug.surface
	debug.window.UpdateSurface()

	// Clear out buffer for next render round.
	debug.buffer.FillRect(nil, 0xFF000000)
	debug.renderer.Clear()

	return true
}

func (debug *debugger) Start() {
	var passedTime int64 = 0
	var passedFrame int64 = 0

	running := true
	for running {
		startTime = time.Now()
		running = debug.Update(elapsedTime)
		endTime = time.Now()
		elapsedTime = endTime.Sub(startTime).Microseconds()

		passedTime += elapsedTime
		passedFrame++
		if passedTime >= 1000000 {
			debug.window.SetTitle(windowTitle + " FPS: " + strconv.Itoa(int(1000000/(passedTime/passedFrame))))
			passedTime = 0
			passedFrame = 0
		}
	}
}

func main() {
	fmt.Println(windowTitle)
	// I really enjoy its graphics. I mean the anime movie.
	fmt.Println("HELLO WORLD -ALLTALE-")
	fmt.Println("With programming we have god's hand.")

	// Read flags
	var file *string
	if len(os.Args) > 2 { // Use flag to parse arguments.
		var file = flag.String("file", "", "NES ROM file")
		var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to file")

		flag.Parse()

		// Handle flags
		if *file == "" {
			fmt.Println("Please specify a NES ROM file.")
			os.Exit(1)
		}

		if *cpuprofile != "" {
			f, err := os.Create(*cpuprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}
	} else if len(os.Args) < 2 {
		fmt.Println("Please specify a NES ROM file. Usage: GoNES.exe [path to NES ROM file].")
		os.Exit(1)
	} else { // Drag-n-drop support.
		file = &os.Args[1]
	}

	// Construct a debugger instance.
	debug := debugger{}
	err := debug.Construct(*file, windowWidth, windowHeight)
	if err != nil {
		return
	}

	// Start debugger.
	debug.Start()
}
