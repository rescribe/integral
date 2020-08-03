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

	integral := NewImage(b)
	draw.Draw(integral, b, img, b.Min, draw.Src)

	if !imgsequal(img, integral) {
		t.Errorf("Read png image differs to integral image\n")
	}
}

func TestSqFromPNG(t *testing.T) {
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

	integral := NewSqImage(b)
	draw.Draw(integral, b, img, b.Min, draw.Src)

	if !imgsequal(img, integral) {
		t.Errorf("Read png image differs to square integral image\n")
	}
}

func TestSum(t *testing.T) {
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

	imgplus := newGray16Plus(b)
	integral := NewImage(b)

	draw.Draw(imgplus, b, img, b.Min, draw.Src)
	draw.Draw(integral, b, img, b.Min, draw.Src)

	cases := []struct {
		name string
		r    image.Rectangle
	}{
		{"fullimage", b},
		{"small", image.Rect(1, 1, 5, 5)},
		{"toobig", image.Rect(0, 0, 2000, b.Dy())},
		{"toosmall", image.Rect(-1, -1, 4, 5)},
		{"small2", image.Rect(0, 0, 4, 4)},
	}

	for _, c := range cases{
		t.Run(c.name, func(t *testing.T) {
			sumimg := imgplus.sum(c.r)
			sumint := integral.Sum(c.r)
			if sumimg != sumint {
				t.Errorf("Sum of integral image differs to regular image: regular: %d, integral: %d\n", sumimg, sumint)
			}
		})
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

type grayPlus struct {
	image.Gray16
}

func newGray16Plus(r image.Rectangle) *grayPlus {
	var g grayPlus
	g.Gray16 = *image.NewGray16(r)
	return &g
}

func (i grayPlus) sum(r image.Rectangle) uint64 {
	var sum uint64
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			c := i.Gray16At(x, y).Y
			sum += uint64(c)
		}
	}
	return sum
}
