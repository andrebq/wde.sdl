package main

import (
	_ "github.com/andrebq/wde.sdl"
	"github.com/skelterjohn/go.wde"
	"log"
)

func run() {
	w, err := wde.NewWindow(480, 320)
	if err != nil {
		log.Printf("Unable to create window %v", err)
		wde.Stop()
	}
	w.SetTitle("WDE SDL Backend")
	w.Show()
	for ev := range w.EventChan() {
		switch ev := ev.(type) {
		case wde.CloseEvent:
			log.Printf("Going out. %v", ev)
			wde.Stop()
		default:
			log.Printf("wde: %v", ev)
		}
	}
}

func main() {
	go run()
	wde.Run()
}
