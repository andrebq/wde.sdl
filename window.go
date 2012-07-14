package wdesdl

import (
	_ "fmt"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/skelterjohn/geom"
	"github.com/skelterjohn/go.wde"
	"image/draw"
)

// Represent a error detected inside the SDL library
type sdlError string

// Implement the error interface
func (s sdlError) Error() string {
	return "SDL: " + string(s)
}

// A wde window backed by a SDL surface
type sdlWindow struct {
	*sdl.Surface
	events chan interface{}
}

// Change the title of the window.
func (s *sdlWindow) SetTitle(title string) {
	sdlWrap.Title <- title
}

// Resize the window
func (s *sdlWindow) SetSize(w, h int) {
	sdlWrap.Size <- &geom.Coord{float64(w), float64(h)}
}

// Return the current size of the window
func (s *sdlWindow) Size() (w, h int) {
	r := s.Bounds()
	return r.Max.X, r.Max.Y
}

// Display the window
func (s *sdlWindow) Show() {
	return
}

// Return the current window
func (s *sdlWindow) Screen() draw.Image {
	return s
}

// Swap the buffers
func (s *sdlWindow) FlushImage() {
	s.Flip()
}

// Return the event channel.
func (s *sdlWindow) EventChan() <-chan interface{} {
	return s.events
}

// Close the current window
func (s *sdlWindow) Close() (err error) {
	s.Surface.Free()
	return
}

// Pool the events
func (s *sdlWindow) poolSdl() {
	for e := range sdlWrap.Events {
		switch e := e.(type) {
		case *sdl.QuitEvent:
			s.events <- wde.CloseEvent{}
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
				kde := wde.KeyDownEvent{wdekeyFromCode(e)}
				s.events <- kde

				if len(string(rune(e.Keysym.Unicode))) > 0 {
					kte := wde.KeyTypedEvent{KeyEvent: wde.KeyEvent{wdekeyFromCode(e)},
						Glyph: string(rune(e.Keysym.Unicode)),
						Chord: ""}
					s.events <- kte
				}
			} else {
				kde := wde.KeyUpEvent{wdekeyFromCode(e)}
				s.events <- kde
			}
		default:
			_ = e
		}
	}
}

// Convert the sdl keycode to the wde table.
func wdekeyFromCode(e *sdl.KeyboardEvent) string {
	//TODO Implement this
	return ""
}

// Create a new SDL window
func newWindow(width, height int) (wdeWindow wde.Window, err error) {
	w := &sdlWindow{}
	wdeWindow = w
	sdlWrap.Size <- &geom.Coord{float64(width), float64(height)}
	w.Surface = <-sdlWrap.Surface
	w.events = make(chan interface{}, 16)
	go w.poolSdl()

	if w.Surface == nil {
		err = sdlError(sdl.GetError())
		return
	}

	return
}
