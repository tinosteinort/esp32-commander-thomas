package main

import (
	"blinky/game"
	"blinky/hardware"
	"image/color"
	"machine"
	"runtime"
	"time"
)

var white = color.RGBA{255, 255, 255, 255}
var black = color.RGBA{0, 0, 0, 255}

func main() {

	hardware := hardware.NewHardware(hardware.HardwareConfig{
		ButtonUpPin:   machine.GPIO15,
		ButtonDownPin: machine.GPIO14,
		ButtonOkPin:   machine.GPIO25,
		LedPin:        machine.LED,
		SCLPin:        machine.GPIO22,
		SDAPin:        machine.GPIO21,
	})

	game := game.New(*hardware)

	hardware.Led.Off()

	lastDisplayTime := time.Now()
	printSystemInfo := false
	lastSystemInfoTime := time.Now()

	for {

		if hardware.ButtonUp.Pressed() {
			//deltaY = -1
			hardware.Led.Set(true)
		} else if hardware.ButtonDown.Pressed() {
			//deltaY = 1
			hardware.Led.Set(true)
		} else if hardware.ButtonOk.Pressed() {
			//deltaY = 1
			hardware.Led.Set(true)
		} else {
			//deltaY = 0
			hardware.Led.Set(false)
		}
		time.Sleep(time.Millisecond * 1)

		game.Control()

		if time.Since(lastDisplayTime) > time.Millisecond*50 {
			err := hardware.Display.Display()
			if err != nil {
				println(err)
			}
			lastDisplayTime = time.Now()
		}

		if printSystemInfo && time.Since(lastSystemInfoTime) > time.Second*5 {
			printMemoryStats()
			lastSystemInfoTime = time.Now()
		}
	}
}

func printMemoryStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	println("HeapAlloc:", m.HeapAlloc)
	println("HeapInuse:", m.HeapInuse)
	println("HeapReleased:", m.HeapReleased)
	println("Mallocs:", m.Mallocs)
	println("Frees:", m.Frees)
	println("TotalAlloc:", m.TotalAlloc)
	println("==================================")
}
