package main

import (
	"math"
)

const (
	winWidth  = 800
	winHeight = 600
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

type pixel struct {
	x int
	y int
}

// called point, more used as vector

type point struct {
	x float64
	y float64
	z float64
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
