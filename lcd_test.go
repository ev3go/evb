// Copyright ©2016 The ev3go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package evb

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var genGolden = flag.Bool("gen.golden", false, "generate golden image files")

func decodeImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

var testImages = []string{
	"gopherbrick",
	"corner",
	"black",
}

func TestRGB565(t *testing.T) {
	for _, test := range testImages {
		golden := filepath.FromSlash("testdata/" + test + "-565.png")

		src, err := decodeImage(filepath.FromSlash("testdata/" + test + ".png"))
		if err != nil {
			t.Fatalf("failed to read src image file %v.png: %v", test, err)
		}

		got := NewRGB565(src.Bounds())
		draw.Draw(got, got.Bounds(), src, src.Bounds().Min, draw.Src)

		if *genGolden {
			f, err := os.Create(golden)
			if err != nil {
				t.Fatalf("failed to create golden image file %v-565.png: %v", test, err)
			}
			defer f.Close()
			err = png.Encode(f, got)
			if err != nil {
				t.Fatalf("failed to encode golden image %v-565.png: %v", test, err)
			}
			continue
		}

		gol, err := decodeImage(golden)
		if err != nil {
			t.Fatalf("failed to read golden image file %v-565.png: %v", test, err)
		}
		want := NewRGB565(gol.Bounds())
		draw.Draw(want, want.Bounds(), gol, gol.Bounds().Min, draw.Src)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("RGB565 from source does not match expected image for %v test", test)
		}
	}
}

var pixelTests = []struct {
	rgb    color.RGBA
	rgb565 Pixel565
}{
	{rgb: color.RGBA{R: 0x05, G: 0x0a, B: 0x0b, A: 0xff}, rgb565: 0x0041},
	{rgb: color.RGBA{R: 0x0e, G: 0x21, B: 0x26, A: 0xff}, rgb565: 0x0904},
	{rgb: color.RGBA{R: 0x5a, G: 0xda, B: 0xff, A: 0xff}, rgb565: 0x5edf},
}

func TestRGB565Model(t *testing.T) {
	for _, test := range pixelTests {
		got := RGB565Model.Convert(test.rgb)
		want := test.rgb565
		if got != want {
			t.Errorf("unexpected RGB565 value for %+v: got: %016b, want: %016b", test.rgb, got, want)
		}
	}
}

func TestPixel565RGBA(t *testing.T) {
	for _, test := range pixelTests {
		got := color.RGBAModel.Convert(test.rgb565).(color.RGBA)
		got.R &= 0xf8
		got.G &= 0xfc
		got.B &= 0xf8
		want := test.rgb
		want.R &= 0xf8
		want.G &= 0xfc
		want.B &= 0xf8
		if got != want {
			t.Errorf("unexpected RGBA value for %016b: got: %+v, want: %+v", test.rgb565, got, want)
		}
	}
}
