package main

import (
	"math"
)

type side struct {
	side_color  color
	center      point
	edges       [4]point
	was_visible bool
}

func newSide(side_color color, center point) side {
	return side{
		side_color:  side_color,
		center:      center,
		was_visible: false,
		edges:       getEdges(center),
	}
}

// func (s *side) setVisible(visible bool) bool {
// 	if s.was_visible == visible {
// 		return false
// 	}
// 	s.was_visible = !s.was_visible
// 	return !visible
// }

func (s *side) displace(d_vector point) {
	for i := 0; i < 4; i++ {
		s.edges[i].x += d_vector.x
		s.edges[i].y += d_vector.y
		s.edges[i].z += d_vector.z
	}
}

func get_lines(view viewingPlane, edges [4]point) []lineSeg {
	var vertices [4]pixel
	for i := 0; i < 4; i++ {
		x, y := mapToPixel(edges[i], view, 8, 6)
		vertices[i] = pixel{x, y}
	}
	_, _, lines := getEquations(vertices)
	return lines
}

func (s *side) draw(pixels []byte, view viewingPlane) {

	lines := get_lines(view, s.edges)

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
}

func getEdges(center point) [4]point {
	result := [4]point{}
	fixX := false
	fixY := false
	if center.x != 0 {
		fixX = true
	} else if center.y != 0 {
		fixY = true
	}

	for i := 0; i < 4; i++ {
		if fixX {
			result[i].x = center.x
			result[i].y = math.Pow(-1, float64(i/2))
			result[i].z = math.Pow(-1, float64(i%2))
		} else if fixY {
			result[i].y = center.y
			result[i].x = math.Pow(-1, float64(i/2))
			result[i].z = math.Pow(-1, float64(i%2))
		} else {
			result[i].z = center.z
			result[i].x = math.Pow(-1, float64(i/2))
			result[i].y = math.Pow(-1, float64(i%2))
		}
	}
	return result
}
