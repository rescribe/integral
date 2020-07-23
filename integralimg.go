// Copyright 2019 Nick White.
// Use of this source code is governed by the GPLv3
// license that can be found in the LICENSE file.

// integralimg is a package for processing integral images, aka
// summed area tables. These are structures which precompute the
// sum of pixels to the left and above each pixel, which can make
// several common image processing operations much faster.
package integralimg

import (
	"image"
	"image/color"
	"math"
)

// I is the Integral Image
type I [][]uint64

// SqImage is a Square integral Image.
// A squared integral image is an integral image for which the square of
// each pixel is saved; this is useful for efficiently calculating
// Standard Deviation.
type SqImage [][]uint64

func (i I) ColorModel() color.Model { return color.Gray16Model }

func (i I) Bounds() image.Rectangle {
	return image.Rectangle {image.Point{0, 0}, image.Point{len(i[0]), len(i)}}
}

// at64 is used to return the raw uint64 for a given pixel. Accessing
// this separately to a (potentially lossy) conversion to a Gray16 is
// necessary for SqImage to function accurately.
func (i I) at64(x, y int) uint64 {
	if !(image.Point{x, y}.In(i.Bounds())) {
		return 0
	}

	var prevx, prevy, prevxy uint64
	prevx, prevy, prevxy = 0, 0, 0
	if x > 0 {
		prevx = i[y][x-1]
	}
	if y > 0 {
		prevy = i[y-1][x]
	}
	if x > 0 && y > 0 {
		prevxy = i[y-1][x-1]
	}
	orig := i[y][x] + prevxy - prevx - prevy
	return orig
}

func (i I) At(x, y int) color.Color {
	c := i.at64(x, y)
	return color.Gray16{uint16(c)}
}

func (i I) set64(x, y int, c uint64) {
	var prevx, prevy, prevxy uint64
	prevx, prevy, prevxy = 0, 0, 0
	if x > 0 {
		prevx = i[y][x-1]
	}
	if y > 0 {
		prevy = i[y-1][x]
	}
	if x > 0 && y > 0 {
		prevxy = i[y-1][x-1]
	}
	final := c + prevx + prevy - prevxy
	i[y][x] = final
}

func (i I) Set(x, y int, c color.Color) {
	gray := color.Gray16Model.Convert(c).(color.Gray16).Y
	i.set64(x, y, uint64(gray))
}

// NewImage returns a new integral Image with the given bounds.
func NewImage(r image.Rectangle) *I {
	w, h := r.Dx(), r.Dy()
	var rows I
	for i := 0; i < h; i++ {
		col := make([]uint64, w)
		rows = append(rows, col)
	}
	return &rows
}

func (i SqImage) ColorModel() color.Model { return I(i).ColorModel() }

func (i SqImage) Bounds() image.Rectangle {
	return I(i).Bounds()
}

func (i SqImage) At(x, y int) color.Color {
	c := I(i).at64(x, y)
	rt := math.Sqrt(float64(c))
	return color.Gray16{uint16(rt)}
}

func (i SqImage) Set(x, y int, c color.Color) {
	gray := uint64(color.Gray16Model.Convert(c).(color.Gray16).Y)
	I(i).set64(x, y, gray * gray)
}

// NewSqImage returns a new squared integral Image with the given bounds.
func NewSqImage(r image.Rectangle) *SqImage {
	i := NewImage(r)
	s := SqImage(*i)
	return &s
}

// Window is a part of an Integral Image
type Window struct {
	topleft uint64
	topright uint64
	bottomleft uint64
	bottomright uint64
	width int
	height int
}

// ToIntegralImg creates an integral image
func ToIntegralImg(img *image.Gray) I {
	var integral I
	var oldy, oldx, oldxy uint64
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		newrow := []uint64{}
		for x := b.Min.X; x < b.Max.X; x++ {
			oldx, oldy, oldxy = 0, 0, 0
			if x > 0 {
				oldx = newrow[x-1]
			}
			if y > 0 {
				oldy = integral[y-1][x]
			}
			if x > 0 && y > 0 {
				oldxy = integral[y-1][x-1]
			}
			pixel := uint64(img.GrayAt(x, y).Y)
			i := pixel + oldx + oldy - oldxy
			newrow = append(newrow, i)
		}
		integral = append(integral, newrow)
	}
	return integral
}

// ToSqIntegralImg creates an integral image of the square of all
// pixel values
func ToSqIntegralImg(img *image.Gray) I {
	var integral I
	var oldy, oldx, oldxy uint64
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		newrow := []uint64{}
		for x := b.Min.X; x < b.Max.X; x++ {
			oldx, oldy, oldxy = 0, 0, 0
			if x > 0 {
				oldx = newrow[x-1]
			}
			if y > 0 {
				oldy = integral[y-1][x]
			}
			if x > 0 && y > 0 {
				oldxy = integral[y-1][x-1]
			}
			pixel := uint64(img.GrayAt(x, y).Y)
			i := pixel * pixel + oldx + oldy - oldxy
			newrow = append(newrow, i)
		}
		integral = append(integral, newrow)
	}
	return integral
}

// GetWindow gets the values of the corners of a square part of an
// Integral Image, plus the dimensions of the part, which can
// be used to quickly calculate the mean of the area
func (i I) GetWindow(x, y, size int) Window {
	step := size / 2

	minx, miny := 0, 0
	maxy := len(i)-1
	maxx := len(i[0])-1

	if y > (step+1) {
		miny = y - step - 1
	}
	if x > (step+1) {
		minx = x - step - 1
	}

	if maxy > (y + step) {
		maxy = y + step
	}
	if maxx > (x + step) {
		maxx = x + step
	}

	return Window { i[miny][minx], i[miny][maxx], i[maxy][minx], i[maxy][maxx], maxx-minx, maxy-miny}
}

func (i SqImage) GetWindow(x, y, size int) Window {
	return I(i).GetWindow(x, y, size)
}

// GetVerticalWindow gets the values of the corners of a vertical
// slice of an Integral Image, starting at x
func (i I) GetVerticalWindow(x, width int) Window {
	maxy := len(i) - 1
	maxx := x + width
	if maxx > len(i[0])-1 {
		maxx = len(i[0]) - 1
	}

	return Window { i[0][x], i[0][maxx], i[maxy][x], i[maxy][maxx], width, maxy }
}

func (i SqImage) GetVerticalWindow(x, width int) Window {
	return I(i).GetVerticalWindow(x, width)
}

// Sum returns the sum of all pixels in a Window
func (w Window) Sum() uint64 {
	return w.bottomright + w.topleft - w.topright - w.bottomleft
}

// Size returns the total size of a Window
func (w Window) Size() int {
	return w.width * w.height
}

// Mean returns the average value of pixels in a Window
func (w Window) Mean() float64 {
	return float64(w.Sum()) / float64(w.Size())
}

// Proportion returns the proportion of pixels which are on
func (w Window) Proportion() float64 {
	area := w.width * w.height
	// divide by 255 as each on pixel has the value of 255
	sum := float64(w.Sum()) / 255
	return float64(area) / sum - 1
}

// MeanStdDevWindow calculates the mean and standard deviation of
// a section on an Integral Image, using the corresponding Square
// Integral Image.
func MeanStdDevWindow(i I, sq SqImage, x, y, size int) (float64, float64) {
	imean := i.GetWindow(x, y, size).Mean()
	smean := sq.GetWindow(x, y, size).Mean()

	variance := smean - (imean * imean)

	return imean, math.Sqrt(variance)
}
