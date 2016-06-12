// Copyright ©2016 The ev3go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package evb

import (
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"os"
	"sync"
	"syscall"

	"github.com/ev3go/ev3dev"
)

const (
	// LCDWidth is the width of the LCD screen in pixel565.
	LCDWidth = 220

	// LCDHeight is the height of the LCD screen in pixel565.
	LCDHeight = 176

	// LCDStride is the width of the LCD screen memory in bytes.
	LCDStride = 440
)

// LCD is the draw image used draw directly to the evb LCD screen.
// Drawing operations are safe for concurrent use. It must be
// initialized before use.
var LCD ev3dev.FrameBuffer = new(lcd)

// lcd is a reader/writer locked draw.Image.
type lcd struct {
	mu  sync.RWMutex
	img *RGB565
	f   *os.File
}

func (p *lcd) Init(zero bool) error {
	p.mu.RLock()
	if p.f == nil {
		p.mu.RUnlock()
		return p.frameBuffer("/dev/fb0", zero)
	}
	p.mu.RUnlock()
	if zero {
		p.mu.Lock()
		for i := 0; i < LCDHeight*LCDStride; i++ {
			p.img.Pix[i] = 0
		}
		p.mu.Unlock()
	}
	return nil
}

func (p *lcd) Close() (err error) {
	defer func() {
		p.mu.Unlock()
		_err := p.f.Close()
		p.f = nil
		if err == nil {
			err = _err
		}
	}()
	p.mu.Lock()
	return syscall.Munmap(p.img.Pix)
}

func (p *lcd) frameBuffer(path string, zero bool) error {
	defer p.mu.Unlock()
	p.mu.Lock()
	var err error
	p.f, err = os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	fbdev, err := syscall.Mmap(int(p.f.Fd()), 0, LCDHeight*LCDStride, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return err
	}
	if zero {
		for i := 0; i < LCDHeight*LCDStride; i++ {
			fbdev[i] = 0
		}
	}
	p.img, err = newRGB565With(fbdev, image.Rect(0, 0, LCDWidth, LCDHeight))
	return err
}

func (p *lcd) ColorModel() color.Model { return p.img.ColorModel() }
func (p *lcd) Bounds() image.Rectangle { return p.img.Bounds() }
func (p *lcd) At(x, y int) color.Color {
	defer p.mu.RUnlock()
	p.mu.RLock()
	if p.f == nil {
		return nil
	}
	return p.img.At(x, y)
}
func (p *lcd) Set(x, y int, c color.Color) {
	p.mu.RLock()
	if p.f == nil {
		p.mu.RUnlock()
		return
	}
	p.mu.RUnlock()
	p.mu.Lock()
	p.img.Set(x, y, c)
	p.mu.Unlock()
}

// NewRGB565 returns a new RGB565 image with the given bounds.
func NewRGB565(r image.Rectangle) *RGB565 {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 2*w*h)
	return &RGB565{Pix: pix, Rect: r}
}

// newRGB565With returns a new RGB565 image with the given bounds,
// backed by the []byte, pix. If the length of pix does not equal
// 2*w*h, a error is returned.
func newRGB565With(pix []byte, r image.Rectangle) (*RGB565, error) {
	w, h := r.Dx(), r.Dy()
	if len(pix) != 2*w*h {
		return nil, errors.New("ev3dev: bad pixel buffer length")
	}
	return &RGB565{Pix: pix, Rect: r}, nil
}

// RGB565 is an in-memory image whose At method returns Pixel565 values.
type RGB565 struct {
	// Pix holds the image's pixels, as RGB565 values.
	// The Pixel565 at (x, y) is the pair of bytes at
	// Pix[2*(x-Rect.Min.X) + (y-Rect.Min.Y)*2*Rect.Dx].
	// Pixel565 values are encoded little endian in Pix.
	Pix []uint8
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// ColorModel returns the RGB565 color model.
func (p *RGB565) ColorModel() color.Model { return RGB565Model }

// Bounds returns the bounding rectangle for the image.
func (p *RGB565) Bounds() image.Rectangle { return p.Rect }

// At returns the color of the pixel565 at (x, y).
func (p *RGB565) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return Pixel565(0)
	}
	i := p.pixOffset(x, y)
	return Pixel565(binary.LittleEndian.Uint16(p.Pix[i : i+2]))
}

// Set sets the color of the pixel565 at (x, y) to c.
func (p *RGB565) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.pixOffset(x, y)
	binary.LittleEndian.PutUint16(p.Pix[i:i+2], uint16(RGB565Model.Convert(c).(Pixel565)))
}

// pixOffset returns the index into p.Pix for the first byte
// containing the pixel at (x, y).
func (p *RGB565) pixOffset(x, y int) int {
	return 2*(x-p.Rect.Min.X) + (y-p.Rect.Min.Y)*2*p.Rect.Dx()
}

// Pixel565 is an RGB565 pixel.
type Pixel565 uint16

// RGBA returns the RGBA values for the receiver.
func (c Pixel565) RGBA() (r, g, b, a uint32) {
	r = uint32(c&0xf800) >> (11 - 3) // Shift to align high bit to bit 7.
	r |= r >> 5                      // Adjust by highest 3 bits.
	r |= r << 8

	g = uint32(c&0x7e0) >> (5 - 2) // Shift to align high bit to bit 7.
	g |= g >> 6                    // Adjust by highest 2 bits.
	g |= g << 8

	b = uint32(c & 0x1f)
	b <<= 3     // Shift to align high bit to bit 7.
	b |= b >> 5 // Adjust by highest 3 bits.
	b |= b << 8

	return r, g, b, 0xffff
}

// RGB565Model is the color model for RGB565 images.
var RGB565Model color.Model = color.ModelFunc(rgb565Model)

func rgb565Model(c color.Color) color.Color {
	if _, ok := c.(Pixel565); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	r >>= 3
	g >>= 2
	b >>= 3
	return Pixel565((r&0x1f)<<11 | (g&0x3f)<<5 | b&0x1f)
}
