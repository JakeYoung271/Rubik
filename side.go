package main

import (
	"math"
)

type side struct {
	side_color  color
	center      point
	was_visible bool
	side_lines  []lineSeg
	new_lines   []lineSeg
}

func newSide(side_color color, center point) side {
	return side{
		side_color:  side_color,
		center:      center,
		was_visible: false,
		side_lines:  []lineSeg{},
		new_lines:   []lineSeg{},
	}
}

func (s *side) setVisible(visible bool) bool {
	if s.was_visible == visible {
		return false
	}
	s.was_visible = !s.was_visible
	return !visible
}
func (s *side) setNewLines(view viewingPlane) {
	edges := s.getEdges()
	var vertices [4]pixel
	for i := 0; i < 4; i++ {
		x, y := mapToPixel(edges[i], view, 3, 2.25)
		vertices[i] = pixel{x, y}
	}
	_, _, lines := getEquations(vertices)
	s.new_lines = lines
}

func (s *side) clearDiff(pixels []byte) {
	if len(s.side_lines) == 0 {
		return
	}

	xMin, xMax := winWidth, 0
	for i := range s.side_lines {
		if s.side_lines[i].xMin < xMin {
			xMin = s.side_lines[i].xMin
		}
		if s.side_lines[i].xMax > xMax {
			xMax = s.side_lines[i].xMax
		}
	}
	for x := xMin + 1; x < xMax; x++ {
		var yVal int
		//xInNewRect := false
		newMin, newMax := winHeight, 0
		oldMin, oldMax := winHeight, 0
		for _, line := range s.new_lines {
			if line.xMin <= x && line.xMax >= x {
				//xInNewRect = true
				yVal = int(line.m*float64(x) + line.b)
				if yVal < newMin {
					newMin = yVal
				}
				if yVal > newMax {
					newMax = yVal
				}
			}
		}
		for _, line := range s.side_lines {
			if line.xMin <= x && line.xMax >= x {
				yVal = int(line.m*float64(x) + line.b)
				if yVal < oldMin {
					oldMin = yVal
				}
				if yVal > oldMax {
					oldMax = yVal
				}
			}
		}
		for y := oldMin; y <= oldMax; y++ {
			if y < newMin || y > newMax {
				index := (y*winWidth + x) * 4
				if pixels[index+3] == byte(s.side_color.hex>>24) && pixels[index+2] == byte(s.side_color.hex>>16) && pixels[index+1] == byte(s.side_color.hex>>8) {
					pixels[index+3] = byte(Grey.hex >> 24)
					pixels[index+2] = byte(Grey.hex >> 16)
					pixels[index+1] = byte(Grey.hex >> 8)
				}
			}
		}
	}
}

func (s *side) addDiff(pixels []byte) {
	xMin, xMax := winWidth, 0
	for _, line := range s.new_lines {
		if line.xMin < xMin {
			xMin = line.xMin
		}
		if line.xMax > xMax {
			xMax = line.xMax
		}
	}
	// fmt.Println("xMin: ", xMin, " xMax: ", xMax)
	for x := xMin; x <= xMax; x++ {
		var yVal int
		//var xInOldRect bool
		newMin, newMax := winHeight, 0
		oldMin, oldMax := winHeight, 0
		for _, line := range s.new_lines {
			if line.xMin <= x && line.xMax >= x {
				yVal = int(line.m*float64(x) + line.b)
				if yVal < newMin {
					newMin = yVal
				}
				if yVal > newMax {
					newMax = yVal
				}
			}
		}
		for _, line := range s.side_lines {
			if line.xMin <= x && line.xMax >= x {
				//xInOldRect = true
				yVal = int(line.m*float64(x) + line.b)
				if yVal < oldMin {
					oldMin = yVal
				}
				if yVal > oldMax {
					oldMax = yVal
				}
			}
		}
		for y := newMin; y <= newMax; y++ {
			if y < oldMin || y > oldMax {
				index := (y*winWidth + x) * 4
				pixels[index+3] = byte(s.side_color.hex >> 24)
				pixels[index+2] = byte(s.side_color.hex >> 16)
				pixels[index+1] = byte(s.side_color.hex >> 8)
			}
		}
	}
}

func (s *side) cleanup(pixels []byte) {
	minX, maxX := winWidth, 0

	for _, line := range s.side_lines {
		if line.xMin < minX {
			minX = line.xMin
		}
		if line.xMax > maxX {
			maxX = line.xMax
		}
	}

	for x := minX; x <= maxX; x++ {
		yMin, yMax := winHeight, 0
		for _, line := range s.side_lines {
			if line.xMin <= x && line.xMax >= x {
				yVal := int(line.m*float64(x) + line.b)
				if yVal < yMin {
					yMin = yVal
				}
				if yVal > yMax {
					yMax = yVal
				}
			}
		}
		for y := yMin; y <= yMax; y++ {
			index := (y*winWidth + x) * 4
			pixels[index+3] = byte(Grey.hex >> 24)
			pixels[index+2] = byte(Grey.hex >> 16)
			pixels[index+1] = byte(Grey.hex >> 8)
		}
	}

	s.side_lines = []lineSeg{}
	s.new_lines = []lineSeg{}
}

func (s *side) draw(pixels []byte) {

	lines := s.new_lines

	minX, maxX := winWidth, 0
	for i := range lines {
		if lines[i].xMin < minX {
			minX = lines[i].xMin
		}
		if lines[i].xMax > maxX {
			maxX = lines[i].xMax
		}
	}

	for x := minX; x <= maxX; x++ {
		yMin, yMax := winHeight, 0
		var yVal int
		for i := range lines {
			if lines[i].xMin <= x && lines[i].xMax >= x {
				yVal = int(lines[i].m*float64(x) + lines[i].b)
				if yVal < yMin {
					yMin = yVal
				}
				if yVal > yMax {
					yMax = yVal
				}
			}
		}
		for y := yMin; y <= yMax; y++ {
			pixels[(y*winWidth+x)*4+3] = byte(s.side_color.hex >> 24)
			pixels[(y*winWidth+x)*4+2] = byte(s.side_color.hex >> 16)
			pixels[(y*winWidth+x)*4+1] = byte(s.side_color.hex >> 8)
		}
	}
	s.side_lines = lines

}

func (s *side) getEdges() [4]point {
	result := [4]point{}
	fixX := false
	fixY := false
	if s.center.x != 0 {
		fixX = true
	} else if s.center.y != 0 {
		fixY = true
	}

	for i := 0; i < 4; i++ {
		if fixX {
			result[i].x = s.center.x
			result[i].y = math.Pow(-1, float64(i/2))
			result[i].z = math.Pow(-1, float64(i%2))
		} else if fixY {
			result[i].y = s.center.y
			result[i].x = math.Pow(-1, float64(i/2))
			result[i].z = math.Pow(-1, float64(i%2))
		} else {
			result[i].z = s.center.z
			result[i].x = math.Pow(-1, float64(i/2))
			result[i].y = math.Pow(-1, float64(i%2))
		}
	}
	return result
}
