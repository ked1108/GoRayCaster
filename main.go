package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"math"
)

var window *glfw.Window

const WIDTH = 1024
const HEIGHT = 512
const PI = 3.141592653589

type player struct {
	pa float64
	px float64
	py float64
}

var p player
var pdx, pdy float64
var gameMap []int8 = []int8{
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 1, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
}

var mapX, mapY, mapS = 8, 8, 64

func (p player) drawPlayer() {
	gl.Color3f(1, 0, 1)
	gl.PointSize(16)
	gl.Begin(gl.POINTS)
	gl.Vertex2f(linearMapX(float32(p.px)), linearMapY(float32(p.py)))
	gl.End()

	gl.LineWidth(4)
	gl.Begin(gl.LINES)
	gl.Vertex2f(linearMapX(float32(p.px)), linearMapY(float32(p.py)))
	gl.Vertex2f(linearMapX(float32(p.px+5*pdx)), linearMapY(float32(p.py+5*pdy)))
	gl.End()
}

func linearMapX(d float32) float32 {
	slope := float32(2.0 / WIDTH)
	return -1.0 + (slope * d)
}

func linearMapY(d float32) float32 {
	slope := float32(2.0 / HEIGHT)
	return -1 + (slope * d)
}

func drawMap2d() {
	for x := 0; x < mapX; x++ {
		for y := 0; y < mapY; y++ {
			if gameMap[y*mapX+x] == 1 {
				gl.Color3f(0, 0, 0)
			} else {
				gl.Color3f(1, 1, 1)
			}
			var xo = float32(x * mapS)
			var yo = float32(y * mapS)
			gl.Rectf(linearMapX(xo+1), linearMapY(yo+1), linearMapX(xo+float32(mapS-1)), linearMapY(yo+float32(mapS-1)))
		}
	}
}

func (p player) drawRays2d() {
	var ra, ry, rx, xo, yo float64
	var r, mx, my, mp, dof int
	ra = p.pa
	for r := 0; r < 1; r++ {
		aTan := -1 / math.Tan(ra)

		if ra > PI {
			ry = float64((int(p.py)>>6)<<6) - 0.0001
			rx = (p.py-ry)*aTan + p.px
			yo = -64
			xo = -yo * aTan
		}
		if ra < PI {
			ry = float64((int(p.py)>>6)<<6) - 64
			rx = (p.py-ry)*aTan + p.px
			yo = 64
			xo = -yo * aTan
		}
		if ra == 0 || ra == PI {
			rx = p.px
			ry = p.py
			dof = 8
		}

		for dof < 8 {
			mx = int(rx) >> 6
			my = int(ry) >> 6
			mp = my * mapX * mx

			if mp < mapX*mapY && gameMap[mp] == 1 {
				dof = 8
			} else {
				rx += xo
				ry += yo
				dof += 1
			}
		}
	}
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if window.GetKey(glfw.KeyA) == glfw.Press {
		p.pa += 0.1
		if p.pa > 2*PI {
			p.pa -= 2 * PI
		}
		pdx = math.Cos(p.pa) * 5
		pdy = math.Sin(p.pa) * 5
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		p.pa -= 0.1
		if p.pa < 0 {
			p.pa += 2 * PI
		}
		pdx = math.Cos(p.pa) * 5
		pdy = math.Sin(p.pa) * 5
	}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		p.py += pdy
		p.px += pdx
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		p.py -= pdy
		p.px -= pdx
	}
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	glfw.DefaultWindowHints()
	window, err = glfw.CreateWindow(WIDTH, HEIGHT, "RayCaster", nil, nil)
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0.5, 0.5, 0.5, 0.5)
	window.SwapBuffers()

	p = player{0, 200, 200}

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		drawMap2d()
		p.drawPlayer()
		window.SwapBuffers()
		window.SetKeyCallback(keyCallback)
		glfw.PollEvents()
	}
	glfw.Terminate()
}
