
//	A continuation on the 'basic' gomobile example (golang.org/x/mobile/example/basic)
//  Draws a background color, that can be changed with touch.

package main

import (
	"log"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

var (
	images   *glutil.Images
	fps      *debug.FPS
	program  gl.Program
	buf      gl.Buffer

	// 'green' is used as a variable for the color green, obviously.
	green  float32
	// 'isTouching' is used to signal that the user is actively touching the screen.
	isTouching bool
)

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context
		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop(glctx)
					glctx = nil
				}
			case size.Event:
				sz = e

			case paint.Event:
				if glctx == nil || e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}

				onPaint(glctx, sz)
				a.Publish()
				// Drive the animation by preparing to paint the next frame
				// after this one is shown.
				a.Send(paint.Event{})
			case touch.Event:
				// Switch control to determine what type of touch the program is getting.
				// Only TypeEnd (if releasing touch) will turn isTouching false.
				switch e.Type {
				case touch.TypeBegin:
					isTouching = true
				case touch.TypeMove:
					isTouching = true
				case touch.TypeEnd:
					isTouching = false
				}
			}
		}
	})
}

func onPaint(glctx gl.Context, sz size.Event) {

	// If a touch-input is registered, the float32 green will increase by 0.01 per 'tick'.
	// It will paint the screen with this colour on that condition.

	// If it happens that 'isTouching' is NOT true, the screen will be painted red, and sets 'green' to 0.
	if (isTouching) {
		green += 0.01
		if green > 1 {
			green = 0
		}
		glctx.ClearColor(0, green, 0, 1)
	} else {
		glctx.ClearColor(1, 0, 0, 1)
		green = 0
	}

	glctx.Clear(gl.COLOR_BUFFER_BIT)
	glctx.UseProgram(program)
	fps.Draw(sz)
}

func onStart(glctx gl.Context) {
	var err error
	program, err = glutil.CreateProgram(glctx, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	buf = glctx.CreateBuffer()
	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)

	images = glutil.NewImages(glctx)
	fps = debug.NewFPS(images)
}
func onStop(glctx gl.Context) {
	glctx.DeleteProgram(program)
	glctx.DeleteBuffer(buf)
	fps.Release()
	images.Release()
}


// -----------------------------------------------------------------------------------//

const vertexShader = `#version 100
uniform vec2 offset;

attribute vec4 position;
void main() {
	// offset comes in with x/y values between 0 and 1.
	// position bounds are -1 to 1.
	vec4 offset4 = vec4(2.0*offset.x-1.0, 1.0-2.0*offset.y, 0, 0);
	gl_Position = position + offset4;
}`

const fragmentShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`
