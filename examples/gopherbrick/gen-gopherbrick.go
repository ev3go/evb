// Copyright ©2016 The ev3go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/kortschak/utter"

	"github.com/ev3go/ev3dev/fb"
)

func main() {
	f, err := os.Open("gopherbrick-screen.png")
	if err != nil {
		log.Fatalf("failed to open gopherbrick image file: %v", err)
	}
	defer f.Close()

	src, err := png.Decode(f)
	if err != nil {
		log.Fatalf("failed to decode gopherbrick image file: %v", err)
	}

	dst := fb.NewRGB565(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)

	utter.Config.ElideType = true
	utter.Config.CommentBytes = false
	utter.Config.BytesWidth = 16
	utter.Config.Indent = "\t"
	gopher := utter.Sdump(dst)

	out, err := os.Create("gopherbrick.go")
	if err != nil {
		log.Fatalf("failed to create gopherbrick.go file: %v", err)
	}
	defer out.Close()
	fmt.Fprintf(out, `// generated by gen-gopherbrick; DO NOT EDIT

// Copyright ©2016 The ev3go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The Go gopher was designed by Renee French and is
// licensed under the Creative Commons Attributions 3.0.

//go:generate go run gen-gopherbrick.go

// gopherbrick demonstrates use of the evb screen.
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

var gopher = %s`, gopher)
}
