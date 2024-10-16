package main

//todo:
// clear old side lines when a cube no longer becomes visible
//rum remove diff when becomes no longer visible
//figure out why it leaves trail.

import (
	"fmt"
	"log"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// type cube struct {
// 	sides [6]*side
// }

func main() {
	center_pieces := getCenterPieces()
	corner_pieces := getCornerPieces(center_pieces)
	edge_pieces := getEdgePieces(center_pieces)
	rubiks_cube := rubiks_cube{corner_pieces, edge_pieces, center_pieces}
	view := viewingPlane{point{1, 1, 1}, point{1, 1, -2}, point{-1, 1, 0}}
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

	visible, _ := rubiks_cube.getVisibleSides(view.normal)
	for _, s := range visible {
		s.draw(pixels, view)
	}
	texture.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*4)
	renderer.Copy(texture, nil, nil)
	renderer.Present()

	// for counting fps
	frames := 1
	updates := 1
	last_time := time.Now()
	key_states := [4]bool{false, false, false, false}

	for counter := 0; running; counter++ {
		update := false
		frames += 1
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN { // Check for key press (not release)
					switch e.Keysym.Sym {
					case sdl.K_UP:
						key_states[0] = true
					case sdl.K_DOWN:
						key_states[1] = true
					case sdl.K_LEFT:
						key_states[2] = true
					case sdl.K_RIGHT:
						key_states[3] = true
					}
				}
				if e.Type == sdl.KEYUP { // Check for key press (not release)
					switch e.Keysym.Sym {
					case sdl.K_UP:
						key_states[0] = false
					case sdl.K_DOWN:
						key_states[1] = false
					case sdl.K_LEFT:
						key_states[2] = false
					case sdl.K_RIGHT:
						key_states[3] = false
					}
				}
			}
		}
		if key_states[0] {
			view.rotate(false, true)
			update = true
		}
		if key_states[1] {
			view.rotate(false, false)
			update = true
		}
		if key_states[2] {
			view.rotate(true, true)
			update = true
		}
		if key_states[3] {
			view.rotate(true, false)
			update = true
		}
		if update {
			visible, _ := rubiks_cube.getVisibleSides(view.normal)
			updatePixels(pixels, counter)
			for _, s := range visible {
				s.draw(pixels, view)
			}
			texture.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*4)
			updates += 1

		}
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		sdl.Delay(10)
		if frames%100 == 0 {
			timeDiff := time.Since(last_time).Seconds()
			last_time = time.Now()
			fmt.Printf("looped  %f fps \n", 100/timeDiff)
			fmt.Printf("updated the canvase %f fps \n", float64(updates)/timeDiff)
			frames = 0
			updates = 0
		}
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
