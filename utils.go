package main

import (
	"math"
)

const (
	winWidth       = 800
	winHeight      = 600
	ROTATION_CONST = math.Pi / 180
	TOLERANCE      = 0.01
)

type color struct {
	hex  uint32
	name string
}

var Green = color{uint32(0x00FF00FF), "green"}
var Red = color{uint32(0xFF0000FF), "red"}
var White = color{uint32(0xFFFFFFFF), "white"}
var Blue = color{uint32(0x0000FFFF), "blue"}
var Yellow = color{uint32(0xFFFF00FF), "yellow"}
var Orange = color{uint32(0xFFA500FF), "orange"}
var Grey = color{uint32(0x808080FF), "grey"}
var Black = color{uint32(0x000000FF), "black"}
var Background = Grey

// called point, more used as vector

type point struct {
	x float64
	y float64
	z float64
}

var SIN_VALUE float64 = math.Sin(ROTATION_CONST)
var COS_VALUE float64 = math.Cos(ROTATION_CONST)

func (lhs point) equals(rhs point) bool {
	return math.Abs(lhs.x-rhs.x) < TOLERANCE && math.Abs(lhs.y-rhs.y) < TOLERANCE && math.Abs(lhs.z-rhs.z) < TOLERANCE
}

func (p point) rotated(axis int, forward bool) point {
	sin_value_local := SIN_VALUE
	cos_value_local := COS_VALUE
	if forward {
		sin_value_local = -sin_value_local

	}
	var x, y, z float64
	switch axis {
	case 0:
		x = p.x
		y = p.y*cos_value_local - p.z*sin_value_local
		z = p.y*sin_value_local + p.z*cos_value_local
	case 1:
		x = p.x*cos_value_local + p.z*sin_value_local
		y = p.y
		z = -p.x*sin_value_local + p.z*cos_value_local
	case 2:
		x = p.x*cos_value_local - p.y*sin_value_local
		y = p.x*sin_value_local + p.y*cos_value_local
		z = p.z
	}
	p.x = x
	p.y = y
	p.z = z
	return p
}

func average(points [4]point) point {
	var result point
	for i := 0; i < 4; i++ {
		result.x += points[i].x
		result.y += points[i].y
		result.z += points[i].z
	}
	result.x /= 4
	result.y /= 4
	result.z /= 4
	return result
}

func scaled(to_scale point, scale float64) point {
	return point{to_scale.x * scale, to_scale.y * scale, to_scale.z * scale}
}

func sq_distance(lhs, rhs point) float64 {
	return (lhs.x-rhs.x)*(lhs.x-rhs.x) + (lhs.y-rhs.y)*(lhs.y-rhs.y) + (lhs.z-rhs.z)*(lhs.z-rhs.z)
}

func (lhs point) dot(rhs point) float64 {
	return lhs.x*rhs.x + lhs.y*rhs.y + lhs.z*rhs.z
}

func (p *point) makeUnit() {
	length := math.Sqrt(p.x*p.x + p.y*p.y + p.z*p.z)
	p.x /= length
	p.y /= length
	p.z /= length
}

type pixel struct {
	x int
	y int
}

func mapToPixel(p point, view viewingPlane, x float64, y float64) (xPixel, yPixel int) {
	xVal := p.dot(view.horiz)
	yVal := p.dot(view.vert)
	xPixel = int((xVal + x) / (2 * x) * winWidth)
	yPixel = int((yVal + y) / (2 * y) * winHeight)
	return
}

// lineSeg related functions

type lineSeg struct {
	xMin int
	xMax int
	m    float64
	b    float64
}

func interpolate(leftPoint, rightPoint pixel) lineSeg {
	if leftPoint.x > rightPoint.x {
		leftPoint, rightPoint = rightPoint, leftPoint
	}
	m := float64(rightPoint.y-leftPoint.y) / float64(rightPoint.x-leftPoint.x)
	b := float64(leftPoint.y) - m*float64(leftPoint.x)
	return lineSeg{leftPoint.x, rightPoint.x, m, b}
}

func getEquations(vertices [4]pixel) (minX, maxX int, lines []lineSeg) {
	vertexList1 := [4]int{0, 0, 1, 2}
	vertexList2 := [4]int{1, 2, 3, 3}

	minX = vertices[0].x
	maxX = vertices[0].x

	for i := 0; i < 4; i++ {
		if vertices[i].x < minX {
			minX = vertices[i].x
		}
		if vertices[i].x > maxX {
			maxX = vertices[i].x
		}
		if vertices[vertexList1[i]].x == vertices[vertexList2[i]].x {
			continue
		}
		lines = append(lines, interpolate(vertices[vertexList1[i]], vertices[vertexList2[i]]))
	}
	return
}

func normalize_number(num float64) float64 {
	if num == 0 {
		return 0
	}
	if num > 0 {
		return 1
	}
	return -1
}

func are_opposites(lhs, rhs point) bool {
	return lhs.x == -rhs.x && lhs.y == -rhs.y && lhs.z == -rhs.z
}
