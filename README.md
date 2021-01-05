![Gopherbrick](gopherbrick.png)
# evb provides EVB-specific functions for the Go ev3dev interface

[![Build Status](https://travis-ci.org/ev3go/evb.svg?branch=master)](https://travis-ci.org/ev3go/evb) [![Coverage Status](https://coveralls.io/repos/ev3go/evb/badge.svg?branch=master&service=github)](https://coveralls.io/github/ev3go/evb?branch=master) [![GoDoc](https://godoc.org/github.com/ev3go/evb?status.svg)](https://godoc.org/github.com/ev3go/evb)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fev3go%2Fevb.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fev3go%2Fevb?ref=badge_shield)

github.com/ev3go/evb depends on an ev3dev kernel 4.4.9-ti-rt-r25-11-ev3dev-bb.org or better.

## Example code

```
package main

import (
	"image"
	"image/draw"
	"time"

	"github.com/ev3go/ev3dev/fb"
	"github.com/ev3go/evb"
)

func main() {
	evb.LCD.Init(true)
	defer evb.LCD.Close()

	// Render the gopherbrick to the screen.
	draw.Draw(evb.LCD, evb.LCD.Bounds(), gopher, gopher.Bounds().Min, draw.Src)

	time.Sleep(10 * time.Second)
}

var gopher = &fb.RGB565{
	Pix: []uint8{
		<snip image data bytes>
	},
	Stride: 440,
	Rect: image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: 220,
			Y: 176,
		},
	},
}
```



## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fev3go%2Fevb.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fev3go%2Fevb?ref=badge_large)