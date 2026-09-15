package game

import (
	"blinky/display"
	"blinky/hardware"
	"image/color"
	"math/rand/v2"
	"slices"
	"strconv"
	"time"

	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

type playScreen struct {
	Model           *playScreenModel
	Display         display.Display
	ColorOn         color.RGBA
	ColorOff        color.RGBA
	LastControlTime time.Time
}

type playScreenModel struct {
	spaceShip         *spaceShip
	points            int16
	bullets           []*bullet
	newBulletAlllowed bool
	enemies           []*enemy
	lives             int16
}

func (m playScreenModel) isGameOver() bool {
	return m.lives < 1
}

type PlayScreenController struct {
	screen   *playScreen
	hardware hardware.Hardware
}

func NewPlayScreenController(hardware hardware.Hardware) *PlayScreenController {
	_, h := hardware.Display.Size()
	hardware.Display.ClearDisplay()

	return &PlayScreenController{
		screen: &playScreen{
			Model: &playScreenModel{
				spaceShip: &spaceShip{
					x:      int16(10),
					y:      h - (h / 2),
					w:      10,
					h:      5,
					deltaY: 0,
				},
				points:            0,
				newBulletAlllowed: true,
				lives:             3,
			},
			ColorOn:         color.RGBA{255, 255, 255, 255},
			ColorOff:        color.RGBA{0, 0, 0, 255},
			LastControlTime: time.Now(),
			Display:         hardware.Display,
		},
		hardware: hardware,
	}
}

func (s *playScreen) Render(c color.RGBA) {

	points := strconv.FormatInt(int64(s.Model.points), 10)
	tinyfont.WriteLine(s.Display, &freemono.Regular9pt7b, 50, 10, points, c)

	var i int16
	for i = 0; i < s.Model.lives; i++ {
		tinydraw.FilledRectangle(
			s.Display,
			i*(s.Model.spaceShip.w+2),
			4,
			s.Model.spaceShip.w,
			s.Model.spaceShip.h,
			c,
		)
	}

	tinydraw.FilledRectangle(
		s.Display,
		s.Model.spaceShip.x,
		s.Model.spaceShip.y,
		s.Model.spaceShip.w,
		s.Model.spaceShip.h,
		c,
	)

	// frame
	w, h := s.Display.Size()
	tinydraw.Rectangle(
		s.Display,
		0,
		16,
		w-2,
		h-16,
		c,
	)

	for _, b := range s.Model.bullets {
		tinydraw.Rectangle(
			s.Display,
			b.x,
			b.y,
			b.w,
			b.h,
			c,
		)
	}

	for _, e := range s.Model.enemies {
		tinydraw.FilledRectangle(
			s.Display,
			e.x,
			e.y,
			e.w,
			e.h,
			c,
		)
	}
}

func (c *PlayScreenController) Control() playScreenResult {
	if time.Since(c.screen.LastControlTime) < time.Millisecond*100 {
		return c.checkScreenIsDone()
	}
	c.screen.LastControlTime = time.Now()

	c.screen.Render(c.screen.ColorOff)

	ship := c.screen.Model.spaceShip
	displayW, displayH := c.screen.Display.Size()

	if c.hardware.ButtonUp.Pressed() {
		ship.deltaY = -2
	} else if c.hardware.ButtonDown.Pressed() {
		ship.deltaY = 2
	} else {
		ship.deltaY = 0
	}

	if c.hardware.ButtonOk.Pressed() {
		if c.screen.Model.newBulletAlllowed {
			c.screen.Model.newBulletAlllowed = false
			c.screen.Model.bullets = append(c.screen.Model.bullets, newBullet(ship))
		}
	} else {
		c.screen.Model.newBulletAlllowed = true
	}

	if c.shouldCreateNewEnemy() {
		c.screen.Model.enemies = append(c.screen.Model.enemies, c.newEnemy())
	}

	// move spaceShip
	if ship.y+ship.deltaY >= 16 &&
		ship.y+ship.h+ship.deltaY < displayH {

		ship.y += c.screen.Model.spaceShip.deltaY
	}

	// move all bullets
	for _, b := range c.screen.Model.bullets {
		b.x += b.deltaX

	}

	// move all enemies
	for _, e := range c.screen.Model.enemies {
		e.x += e.deltaX
	}

	c.checkBulletEnemyCollision()
	c.checkSpaceShipEnemyCollision()

	// then delete all bullets, which are not full in screen
	c.screen.Model.bullets = slices.DeleteFunc(c.screen.Model.bullets, func(b *bullet) bool {
		return b.x+b.w >= displayW
	})
	// then delete all enemies, which are not full in screen
	c.screen.Model.enemies = slices.DeleteFunc(c.screen.Model.enemies, func(e *enemy) bool {
		return e.x <= 0
	})

	c.screen.Render(c.screen.ColorOn)

	return c.checkScreenIsDone()
}

func newBullet(s *spaceShip) *bullet {
	bh := int16(2)
	bw := int16(5)
	return &bullet{
		x:      s.x + s.w,
		y:      s.y + s.h - bh,
		w:      bw,
		h:      bh,
		deltaX: 7,
	}
}

func (c *PlayScreenController) newEnemy() *enemy {
	ew := int16(8)
	eh := int16(4)

	displayW, displayH := c.screen.Display.Size()
	min := int(16)
	max := int(displayH - eh)
	randY := min + rand.IntN(max-min)

	return &enemy{
		x:      displayW - ew,
		y:      int16(randY),
		w:      ew,
		h:      eh,
		deltaX: -5,
	}
}

func (c *PlayScreenController) shouldCreateNewEnemy() bool {
	if c.screen.Model.points < 5 {
		return len(c.screen.Model.enemies) < 2
	} else if c.screen.Model.points < 10 {
		return len(c.screen.Model.enemies) < 3
	} else if c.screen.Model.points < 15 {
		return len(c.screen.Model.enemies) < 4
	}
	return len(c.screen.Model.enemies) < 5
}

func (c *PlayScreenController) checkBulletEnemyCollision() {
	for i := len(c.screen.Model.bullets) - 1; i >= 0; i-- {
		for j := len(c.screen.Model.enemies) - 1; j >= 0; j-- {

			if bulletOverlapsEnemy(c.screen.Model.bullets[i], c.screen.Model.enemies[j]) {
				c.onBulletHitsEnemy()

				c.screen.Model.bullets = append(c.screen.Model.bullets[:i], c.screen.Model.bullets[i+1:]...)
				c.screen.Model.enemies = append(c.screen.Model.enemies[:j], c.screen.Model.enemies[j+1:]...)

				break
			}
		}
	}
}

func (c *PlayScreenController) checkSpaceShipEnemyCollision() {
	for j := len(c.screen.Model.enemies) - 1; j >= 0; j-- {

		if spaceShipOverlapsEnemy(c.screen.Model.spaceShip, c.screen.Model.enemies[j]) {
			c.onEnemyHitsSpaceShip()

			c.screen.Model.enemies = append(c.screen.Model.enemies[:j], c.screen.Model.enemies[j+1:]...)

			break
		}
	}
}

func bulletOverlapsEnemy(b *bullet, e *enemy) bool {
	return b.x < e.x+e.w &&
		b.x+b.w > e.x &&
		b.y < e.y+e.h &&
		b.y+b.h > e.y
}

func spaceShipOverlapsEnemy(s *spaceShip, e *enemy) bool {
	return s.x < e.x+e.w &&
		s.x+s.w > e.x &&
		s.y < e.y+e.h &&
		s.y+s.h > e.y
}

func (c *PlayScreenController) onBulletHitsEnemy() {
	c.screen.Model.points += 1
}

func (c *PlayScreenController) onEnemyHitsSpaceShip() {
	c.screen.Model.lives -= 1
}

func (c *PlayScreenController) checkScreenIsDone() playScreenResult {
	gameOver := c.screen.Model.isGameOver()
	if gameOver {
		c.screen.Render(c.screen.ColorOff)
	}
	return playScreenResult{
		screenDone: gameOver,
		points:     c.screen.Model.points,
	}
}
