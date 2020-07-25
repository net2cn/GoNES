package main

import (
	"fmt"
	"github.com/net2cn/GoNES/nes"
	"github.com/veandco/go-sdl2/sdl"
)

var windowTitle string = "GoNES SDL2"
var windowWidth, windowHeight int32 = 680, 480

type demoCPU struct {
}

func (cpu *demoCPU) Construct(width int32, height int32) error {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var src, dst sdl.Rect
	var err error

	window, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println("Failed to create window: %s", err)
		return err
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("Failed to create renderer: %s", err)
		return err
	}
	defer renderer.Destroy()

	src = sdl.Rect{0, 0, 512, 512}
	dst = sdl.Rect{100, 50, 512, 512}

	renderer.Clear()
	renderer.SetDrawColor(255, 0, 0, 0)
	renderer.FillRect(&src)
	renderer.SetDrawColor(0, 255, 0, 0)
	renderer.FillRect(&dst)
	renderer.Present()

	sdl.PollEvent()
	sdl.Delay(2000)

	return nil
}

func (cpu *demoCPU) Start() {

}

func main() {
	fmt.Println("Hello world.")
	demo := demoCPU{}
	err := demo.Construct(windowWidth, windowHeight)
	if err != nil {
		return
	}
	demo.Start()
}
