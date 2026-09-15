package display

import "image/color"

type Display interface {
	Size() (x, y int16)

	SetPixel(x, y int16, c color.RGBA)

	Display() error
	ClearDisplay()
}
