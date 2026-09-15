package game

import (
	"blinky/display"
	"blinky/hardware"
	"image/color"
	"strconv"
	"time"

	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

func NewGameOverController(hw hardware.Hardware, points int16) *GamveOverScreenController {
	_, h := hw.Display.Size()

	pointsText := strconv.FormatInt(int64(points), 10)

	return &GamveOverScreenController{
		screen: &gameOverScreen{
			Model: &gameOverScreenModel{
				Line1: "Game",
				Line2: "Over: " + pointsText,
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

type gameOverScreen struct {
	Model               *gameOverScreenModel
	Display             display.Display
	ColorOn             color.RGBA
	ColorOff            color.RGBA
	LastControlTime     time.Time
	DirectionChangeTime time.Time
}

func (s *gameOverScreen) Render(c color.RGBA) {
	tinyfont.WriteLine(s.Display, &freemono.BoldOblique9pt7b, 0, 14, s.Model.Line1, c)
	tinyfont.WriteLine(s.Display, &freemono.Bold12pt7b, 0, 30, s.Model.Line2, c)

	tinydraw.Line(s.Display, 30, 40, 110, 40, c)
	tinydraw.Line(s.Display, 35, 45, 120, 45, c)
	tinydraw.Line(s.Display, 30, 50, 110, 50, c)
	tinydraw.Line(s.Display, 35, 55, 120, 55, c)

	tinydraw.FilledRectangle(
		s.Display,
		s.Model.spaceShip.x,
		s.Model.spaceShip.y,
		s.Model.spaceShip.w,
		s.Model.spaceShip.h,
		c,
	)
}

type gameOverScreenModel struct {
	Line1      string
	Line2      string
	spaceShip  *spaceShip
	screenDone bool
}

type GamveOverScreenController struct {
	screen   *gameOverScreen
	hardware hardware.Hardware
}

func (c *GamveOverScreenController) Control() bool {
	if time.Since(c.screen.LastControlTime) < time.Millisecond*200 {
		return c.checkForContinue()
	}
	c.screen.LastControlTime = time.Now()

	c.screen.Render(c.screen.ColorOff)

	c.screen.Model.spaceShip.y += c.screen.Model.spaceShip.deltaY

	if time.Since(c.screen.DirectionChangeTime) > 2000*time.Millisecond {
		c.screen.Model.spaceShip.deltaY = -c.screen.Model.spaceShip.deltaY
		c.screen.DirectionChangeTime = time.Now()
	}

	if c.hardware.ButtonOk.Pressed() {
		c.screen.Model.screenDone = true
	}

	c.screen.Render(c.screen.ColorOn)

	return c.checkForContinue()
}

func (c *GamveOverScreenController) checkForContinue() bool {
	if c.screen.Model.screenDone {
		c.screen.Render(c.screen.ColorOff)
	}
	return c.screen.Model.screenDone
}
