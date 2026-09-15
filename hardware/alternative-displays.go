package hardware

import (
	"blinky/display"
	"image/color"

	"tinygo.org/x/drivers/sh1106"
)

// ##################### Dirty Tracking Display 1
type dirtyTrackingDisplay struct {
	device    *sh1106.Device
	width     int16
	height    int16
	pages     int16
	buffer    [1024]byte
	dirtyMinX [8]int16
	dirtyMaxX [8]int16
}

func newDirtyTrackingDisplay(device *sh1106.Device) display.Display {
	w, h := device.Size()
	pages := h / 8
	var buffer [1024]byte  // width * pages
	var dirtyMinX [8]int16 // height / 8
	var dirtyMaxX [8]int16 // height / 8
	d := &dirtyTrackingDisplay{
		device:    device,
		width:     w,
		height:    h,
		pages:     pages,
		buffer:    buffer,
		dirtyMinX: dirtyMinX,
		dirtyMaxX: dirtyMaxX,
	}
	d.clearDirty()

	d.invalidateAll()
	return d
}

func (d *dirtyTrackingDisplay) clearDirty() {
	var i int16
	for i = 0; i < d.pages; i++ {
		d.dirtyMinX[i] = 127
		d.dirtyMaxX[i] = 0
	}
}

func (d *dirtyTrackingDisplay) invalidateAll() {
	for i := range d.dirtyMinX {
		d.dirtyMinX[i] = 0
		d.dirtyMaxX[i] = 127
	}
}

func (d *dirtyTrackingDisplay) Size() (x, y int16) {
	return d.device.Size()
}

func (d *dirtyTrackingDisplay) SetPixel(x, y int16, c color.RGBA) {
	if x < 0 || x >= d.width || y < 0 || y >= d.height {
		return
	}

	page := y / 8
	index := page*d.width + x
	bit := uint8(1 << (y % 8))

	color := false
	if c.R+c.G+c.B > 0 {
		color = true
	}

	if color {
		d.buffer[index] |= bit
	} else {
		d.buffer[index] &^= bit
	}

	// Dirty markieren
	if int16(x) < d.dirtyMinX[page] {
		d.dirtyMinX[page] = int16(x)
	}
	if int16(x) > d.dirtyMaxX[page] {
		d.dirtyMaxX[page] = int16(x)
	}
}

func (d *dirtyTrackingDisplay) Display() error {
	return d.displayPartial()
}

func (d *dirtyTrackingDisplay) ClearDisplay() {
	d.device.ClearDisplay()
}

func (d *dirtyTrackingDisplay) displayPartial() error {
	var page int16
	for page = 0; page < d.pages; page++ {

		minX := d.dirtyMinX[page]
		maxX := d.dirtyMaxX[page]

		if minX > maxX {
			continue // nichts geändert
		}

		// Page setzen
		d.device.Command(0xB0 | uint8(page))

		// Spalten setzen (SH1106 Offset beachten!)
		col := int(minX) + 2 // oft +2 Offset beim SH1106!
		d.device.Command(uint8(0x00 | (col & 0x0F)))
		d.device.Command(uint8(0x10 | (col >> 4)))

		// Daten senden
		start := page*d.width + int16(minX)
		end := page*d.width + int16(maxX) + 1

		data := d.buffer[start:end]
		d.device.Tx(data, false)
	}

	d.clearDirty()
	return nil
}

// ##################### Dirty Tracking Display 2
type dirtyTrackingDisplaySecond struct {
	device    *sh1106.Device
	width     int16
	height    int16
	pages     int16
	buffer    [1024]byte
	dirtyMinX [8]int16
	dirtyMaxX [8]int16
}

func newDirtyTrackingDisplaySecond(device *sh1106.Device) display.Display {
	w, h := device.Size()
	pages := h / 8
	var buffer [1024]byte  // width * pages
	var dirtyMinX [8]int16 // height / 8
	var dirtyMaxX [8]int16 // height / 8
	d := &dirtyTrackingDisplaySecond{
		device:    device,
		width:     w,
		height:    h,
		pages:     pages,
		buffer:    buffer,
		dirtyMinX: dirtyMinX,
		dirtyMaxX: dirtyMaxX,
	}
	d.clearDirty()

	d.invalidateAll()
	return d
}

func (d *dirtyTrackingDisplaySecond) clearDirty() {
	var i int16
	for i = 0; i < d.pages; i++ {
		d.dirtyMinX[i] = 127
		d.dirtyMaxX[i] = 0
	}
}

func (d *dirtyTrackingDisplaySecond) invalidateAll() {
	for i := range d.dirtyMinX {
		d.dirtyMinX[i] = 0
		d.dirtyMaxX[i] = 127
	}
}

func (d *dirtyTrackingDisplaySecond) Size() (x, y int16) {
	return d.device.Size()
}

func (d *dirtyTrackingDisplaySecond) SetPixel(x, y int16, c color.RGBA) {
	if x < 0 || x >= d.width || y < 0 || y >= d.height {
		return
	}

	byteIndex := x + (y/8)*d.width

	if c.R != 0 || c.G != 0 || c.B != 0 {
		d.buffer[byteIndex] |= 1 << uint8(y%8)
	} else {
		d.buffer[byteIndex] &^= 1 << uint8(y%8)
	}

	page := y / 8

	if x < d.dirtyMinX[page] {
		d.dirtyMinX[page] = x
	}

	if x > d.dirtyMaxX[page] {
		d.dirtyMaxX[page] = x
	}
}

func (d *dirtyTrackingDisplaySecond) Display() error {
	return d.displayPartial()
}

func (d *dirtyTrackingDisplaySecond) ClearDisplay() {
	d.device.ClearDisplay()
}

func (d *dirtyTrackingDisplaySecond) displayPartial() error {

	for page := int16(0); page < 8; page++ {

		minX := d.dirtyMinX[page]
		maxX := d.dirtyMaxX[page]

		// nichts geändert
		if minX > maxX {
			continue
		}

		// SH1106 Column Offset
		col := minX + 2

		// Page setzen
		d.device.Command(0xB0 | uint8(page))

		// Column setzen
		d.device.Command(uint8(0x00 | (col & 0x0F)))
		d.device.Command(uint8(0x10 | ((col >> 4) & 0x0F)))

		start := int(page)*128 + int(minX)
		end := int(page)*128 + int(maxX) + 1

		// EIN zusammenhängender Transfer
		d.device.Tx(d.buffer[start:end], false)

		// dirty zurücksetzen
		d.dirtyMinX[page] = 127
		d.dirtyMaxX[page] = 0
	}

	return nil
}

// ############### Update only dirte pages
type dirtyPagesDisplay struct {
	device     *sh1106.Device
	width      int16
	height     int16
	pages      int16
	buffer     [1024]byte
	dirtyMinX  [8]int16
	dirtyMaxX  [8]int16
	dirtyPages [8]bool
}

func newDirtyPagesDisplay(device *sh1106.Device) display.Display {
	w, h := device.Size()
	pages := h / 8
	var buffer [1024]byte  // width * pages
	var dirtyMinX [8]int16 // height / 8
	var dirtyMaxX [8]int16 // height / 8
	d := &dirtyPagesDisplay{
		device:    device,
		width:     w,
		height:    h,
		pages:     pages,
		buffer:    buffer,
		dirtyMinX: dirtyMinX,
		dirtyMaxX: dirtyMaxX,
	}
	d.clearDirty()

	d.invalidateAllPages()
	return d
}

func (d *dirtyPagesDisplay) clearDirty() {
	var i int16
	for i = 0; i < d.pages; i++ {
		d.dirtyMinX[i] = 127
		d.dirtyMaxX[i] = 0
	}
}

func (d *dirtyPagesDisplay) invalidateAllPages() {
	for i := range d.dirtyPages {
		d.dirtyPages[i] = true
	}
}

func (d *dirtyPagesDisplay) Size() (x, y int16) {
	return d.device.Size()
}

func (d *dirtyPagesDisplay) SetPixel(x, y int16, c color.RGBA) {
	if x < 0 || x >= d.width || y < 0 || y >= d.height {
		return
	}

	byteIndex := x + (y/8)*d.width

	if c.R != 0 || c.G != 0 || c.B != 0 {
		d.buffer[byteIndex] |= 1 << uint8(y%8)
	} else {
		d.buffer[byteIndex] &^= 1 << uint8(y%8)
	}

	page := y / 8

	if x < d.dirtyMinX[page] {
		d.dirtyMinX[page] = x
	}

	if x > d.dirtyMaxX[page] {
		d.dirtyMaxX[page] = x
	}
}

func (d *dirtyPagesDisplay) Display() error {
	return d.displayPartial()
}

func (d *dirtyPagesDisplay) ClearDisplay() {
	d.device.ClearDisplay()
}

func (d *dirtyPagesDisplay) displayPartial() error {

	var page int16
	for page = int16(0); page < 8; page++ {

		if !d.dirtyPages[page] {
			continue
		}

		// SH1106 Page auswählen
		d.device.Command(0xB0 | uint8(page))

		// SH1106 hat typischerweise +2 Column Offset
		d.device.Command(0x02)
		d.device.Command(0x10)

		start := page * 128
		end := start + 128

		d.device.Tx(d.buffer[start:end], false)

		d.dirtyPages[page] = false
	}

	return nil
}
