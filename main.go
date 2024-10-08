package main

//todo:
// clear old side lines when a cube no longer becomes visible
//rum remove diff when becomes no longer visible
//figure out why it leaves trail.

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type cube struct {
	sides [6]*side
}

func (c *cube) getVisibleSides(vector point) (visible, toCleanUp []*side) {
	visible = []*side{}
	toCleanUp = []*side{}
	for _, s := range c.sides {
		// fmt.Println("checking side: ", s.side_color, s.was_visible)
		if s.center.dot(vector) > 0 {
			s.setVisible(true)
			visible = append(visible, s)
		} else if s.setVisible(false) {
			toCleanUp = append(toCleanUp, s)
		}
	}
	return
}

func main() {

	side1 := newSide(Green, point{1, 0, 0})
	side2 := newSide(Red, point{0, 1, 0})
	side3 := newSide(White, point{0, 0, 1})
	side4 := newSide(Blue, point{-1, 0, 0})
	side5 := newSide(Orange, point{0, -1, 0})
	side6 := newSide(Yellow, point{0, 0, -1})

	fmt.Println(side1.getEdges())

	fmt.Println(byte(Orange.hex>>24), byte(Orange.hex>>16), byte(Orange.hex>>8), byte(Orange.hex))

	cube1 := cube{[6]*side{&side1, &side2, &side3, &side4, &side5, &side6}}
	view := viewingPlane{point{1, 1, 1}, point{-1, 1, 0}, point{1, 1, -2}}
	view.normalize()

	// Initialize SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatalf("Failed to initialize SDL: %s\n", err)
	}
	defer sdl.Quit()

	// Create a window
	window, err := sdl.CreateWindow("2D Array Background", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("Failed to create window: %s\n", err)
	}
	defer window.Destroy()

	// Create a renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("Failed to create renderer: %s\n", err)
	}
	defer renderer.Destroy()

	// Create a texture to render the 2D array as a background
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, winWidth, winHeight)
	if err != nil {
		log.Fatalf("Failed to create texture: %s\n", err)
	}
	defer texture.Destroy()

	// Convert the 2D array into a 1D byte array for SDL
	pixels := make([]byte, winWidth*winHeight*4) // 4 bytes per pixel (RGBA)
	updatePixels(pixels, 0)

	running := true

	visible, _ := cube1.getVisibleSides(view.normal)
	for _, s := range visible {
		s.setNewLines(view)
		s.draw(pixels)
	}
	texture.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*4)
	renderer.Copy(texture, nil, nil)
	renderer.Present()

	for counter := 0; running; counter++ {
		update := false
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false

			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN { // Check for key press (not release)
					switch e.Keysym.Sym {
					case sdl.K_UP:
						update = true
						view.rotate(false, true)
						// Add logic for up arrow press here
					case sdl.K_DOWN:
						update = true
						view.rotate(false, false)
						// Add logic for down arrow press here
					case sdl.K_LEFT:
						update = true
						view.rotate(true, true)
						// Add logic for left arrow press here
					case sdl.K_RIGHT:
						update = true
						view.rotate(true, false)
						// Add logic for right arrow press here
					}
				}
			}
		}
		if update {
			visible, toCleanUp := cube1.getVisibleSides(view.normal)
			updatePixels(pixels, counter)
			for _, s := range toCleanUp {
				s.cleanup(pixels)
			}
			for _, s := range visible {
				s.setNewLines(view)
				s.draw(pixels)
				s.side_lines = s.new_lines
			}
			texture.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*4)

		}
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		//sdl.Delay(100)
	}
}

func updatePixels(pixels []byte, offsetVal int) {
	for y := 0; y < winHeight; y++ {
		for x := 0; x < winWidth; x++ {
			index := ((offsetVal + y*winWidth + x) * 4) % (winHeight * winWidth * 4)
			pixels[index+3] = byte(Grey.hex >> 24)
			pixels[index+2] = byte(Grey.hex >> 16)
			pixels[index+1] = byte(Grey.hex >> 8)
			pixels[index] = 255 // Alpha

		}
	}
}
