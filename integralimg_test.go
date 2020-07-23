// Copyright 2020 Nick White.
// Use of this source code is governed by the GPLv3
// license that can be found in the LICENSE file.

package integralimg

import (
	"image"
	"image/draw"
	_ "image/png"
	"os"
	"testing"
)

func TestFromPNG(t *testing.T) {
	f, err := os.Open("testdata/in.png")
	if err != nil {
		t.Fatalf("Could not open file %s: %v\n", "testdata/in.png", err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("Could not decode image: %v\n", err)
	}
	b := img.Bounds()
	gray := image.NewGray(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(gray, b, img, b.Min, draw.Src)

	integral := ToIntegralImg(gray)

	if !imgsequal(img, integral) {
		t.Errorf("Read png image differs to integral\n")
	}
}

func imgsequal(img1, img2 image.Image) bool {
	b := img1.Bounds()
	if !b.Eq(img2.Bounds()) {
		return false
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r0, g0, b0, a0 := img1.At(x, y).RGBA()
			r1, g1, b1, a1 := img2.At(x, y).RGBA()
			if r0 != r1 {
				return false
			}
			if g0 != g1 {
				return false
			}
			if b0 != b1 {
				return false
			}
			if a0 != a1 {
				return false
			}
		}
	}
	return true
}
