package hardware

import (
	"blinky/display"
	"machine"

	"tinygo.org/x/drivers/sh1106"
)

type Hardware struct {
	ButtonUp   *Button
	ButtonDown *Button
	ButtonOk   *Button
	Led        *Led
	Display    display.Display
}

type HardwareConfig struct {
	ButtonUpPin   machine.Pin
	ButtonDownPin machine.Pin
	ButtonOkPin   machine.Pin
	LedPin        machine.Pin
	SCLPin        machine.Pin
	SDAPin        machine.Pin
}

func NewHardware(c HardwareConfig) *Hardware {
	return &Hardware{
		ButtonUp:   newButton(c.ButtonUpPin),
		ButtonDown: newButton(c.ButtonDownPin),
		ButtonOk:   newButton(c.ButtonOkPin),
		Led:        newLed(c.LedPin),
		Display:    newDisplay(c.SCLPin, c.SDAPin),
	}
}

type Button struct {
	pin machine.Pin
}

func newButton(pin machine.Pin) *Button {
	b := &Button{
		pin: pin,
	}
	pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	return b
}

func (b *Button) Pressed() bool {
	return !b.pin.Get()
}

type Led struct {
	pin machine.Pin
}

func newLed(pin machine.Pin) *Led {
	led := &Led{
		pin: pin,
	}
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return led
}

func (l *Led) Set(on bool) {
	if on {
		l.pin.High()
	} else {
		l.pin.Low()
	}
}

func (l *Led) On() {
	l.Set(true)
}

func (l *Led) Off() {
	l.Set(false)
}

func newDisplay(sclPin, sdaPin machine.Pin) display.Display {
	i2c := machine.I2C0
	i2c.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
		SCL:       sclPin,
		SDA:       sdaPin,
	})
	display := sh1106.NewI2C(i2c)
	display.Configure(sh1106.Config{})

	display.ClearDisplay()

	return &display
	//return newDirtyTrackingDisplay(&display)
	//return newDirtyPagesDisplay(&display)
	//return newDirtyTrackingDisplaySecond(&display)
}
