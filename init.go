package wdesdl

import (
	"github.com/banthar/Go-SDL/sdl"
	"github.com/skelterjohn/geom"
	"github.com/skelterjohn/go.wde"
	"runtime"
)

var (
	done    chan struct{}
	sdlWrap = newSdlWrap()
)

func initsdl() {
	go sdlWrap.wrap()
	// wait until the init is completed
	<-sdlWrap.InitDone
}

func init() {
	initsdl()
	wde.BackendNewWindow = newWindow
	wde.BackendRun = run
	wde.BackendStop = stop
	done = make(chan struct{}, 0)
}

func run() {
	<-done
}

func stop() {
	sdlWrap.Quit <- struct{}{}
	done <- struct{}{}
}

// Wrap all interactions with the SDL backend using channels
type SdlWrap struct {
	Size     chan *geom.Coord
	Title    chan string
	Quit     chan struct{}
	Events   chan interface{}
	InitDone chan struct{}
	Surface  chan *sdl.Surface
}

// Create a default SdlWrap
func newSdlWrap() (s *SdlWrap) {
	s = &SdlWrap{
		Size:     make(chan *geom.Coord, 0),
		Title:    make(chan string, 0),
		Quit:     make(chan struct{}, 0),
		Events:   make(chan interface{}, 1),
		InitDone: make(chan struct{}, 0),
		Surface:  make(chan *sdl.Surface, 1)}
	return
}

// Wrap the SDL into one single thread
func (s SdlWrap) wrap() {
	runtime.LockOSThread()

	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		panic("Unable to init SDL. Cause: " + sdl.GetError())
	}
	sdl.EnableUNICODE(1)

	s.InitDone <- struct{}{}

	for {
		select {
		case t := <-s.Title:
			sdl.WM_SetCaption(t, t)
		case <-s.Quit:
			close(s.Events)
			sdl.Quit()
		case sz := <-s.Size:
			s.Surface <- sdl.SetVideoMode(int(sz.X), int(sz.Y), 32, sdl.HWSURFACE|sdl.DOUBLEBUF|sdl.RESIZABLE)
		default:
			e := sdl.PollEvent()
			if e != nil {
				s.Events <- e
			}
		}
	}
}
