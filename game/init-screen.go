package game

import (
	"blinky/display"
	"blinky/hardware"
	"image/color"
	"time"

	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

func NewInitScreenController(hw hardware.Hardware) *InitScreenController {
	_, h := hw.Display.Size()
	//hw.Display.ClearDisplay()

	return &InitScreenController{
		screen: &initScreen{
			Model: &initScreenModel{
				Line1: "Commander",
				Line2: "Thomas",
				spaceShip: &spaceShip{
					x:      int16(10),
					y:      h - (h / 4),
					w:      10,
					h:      5,
					deltaY: 1,
				},
				screenDone: false,
			},
			Display:         hw.Display,
			ColorOn:         color.RGBA{255, 255, 255, 255},
			ColorOff:        color.RGBA{0, 0, 0, 255},
			LastControlTime: time.Now(),
		},
		hardware: hw,
	}
}

type initScreen struct {
	Model               *initScreenModel
	Display             display.Display
	ColorOn             color.RGBA
	ColorOff            color.RGBA
	LastControlTime     time.Time
	DirectionChangeTime time.Time
}

func (s *initScreen) Render(c color.RGBA) {
	tinyfont.WriteLine(s.Display, &freemono.BoldOblique9pt7b, 0, 14, s.Model.Line1, s.ColorOn)
	tinyfont.WriteLine(s.Display, &freemono.Bold12pt7b, 0, 30, s.Model.Line2, s.ColorOn)

	tinydraw.Line(s.Display, 30, 40, 110, 40, s.ColorOn)
	tinydraw.Line(s.Display, 35, 45, 120, 45, s.ColorOn)
	tinydraw.Line(s.Display, 30, 50, 110, 50, s.ColorOn)
	tinydraw.Line(s.Display, 35, 55, 120, 55, s.ColorOn)

	tinydraw.FilledRectangle(
		s.Display,
		s.Model.spaceShip.x,
		s.Model.spaceShip.y,
		s.Model.spaceShip.w,
		s.Model.spaceShip.h,
		c,
	)
}

type initScreenModel struct {
	Line1      string
	Line2      string
	spaceShip  *spaceShip
	screenDone bool
}

type InitScreenController struct {
	screen   *initScreen
	hardware hardware.Hardware
}

func (c *InitScreenController) Control() bool {
	if time.Since(c.screen.LastControlTime) < time.Millisecond*200 {
		return c.checkForScreenDone()
	}
	c.screen.LastControlTime = time.Now()

	c.screen.Render(c.screen.ColorOff)

	c.screen.Model.spaceShip.y += c.screen.Model.spaceShip.deltaY

	if time.Since(c.screen.DirectionChangeTime) > 2000*time.Millisecond {
		c.screen.Model.spaceShip.deltaY = -c.screen.Model.spaceShip.deltaY
		c.screen.DirectionChangeTime = time.Now()
	}

	c.screen.Render(c.screen.ColorOn)

	if c.hardware.ButtonOk.Pressed() {
		c.screen.Model.screenDone = true
	}

	return c.checkForScreenDone()
}

func (c *InitScreenController) checkForScreenDone() bool {
	if c.screen.Model.screenDone {
		c.screen.Render(c.screen.ColorOff)
	}
	return c.screen.Model.screenDone
}
