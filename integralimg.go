// Copyright 2020 Nick White.
// Use of this source code is governed by the GPLv3
// license that can be found in the LICENSE file.

// integralimg is a package for processing integral images, aka
// summed area tables. These are structures which precompute the
// sum of pixels to the left and above each pixel, which can make
// several common image processing operations much faster.
//
// integralimg.Image and integralimg.SqImage fully implement the
// image.Image and image/draw.Draw interfaces, and hence can be
// used like so:
//
//     img, _, err := image.Decode(f)
//     integral := integralimg.NewImage(b)
//     draw.Draw(integral, b, img, b.Min, draw.Src)
//
// The Sum(), Mean() and MeanStdDev() functions provided for the
// integral versions of Images significantly speed up many common
// image processing operations.
package integralimg

import (
	"image"
	"image/color"
	"math"
)

// Image is an integral Image
type Image [][]uint64

// SqImage is a Square integral Image.
// A squared integral image is an integral image for which the square of
// each pixel is saved; this is useful for efficiently calculating
// Standard Deviation.
type SqImage [][]uint64

func (i Image) ColorModel() color.Model { return color.Gray16Model }

func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(i[0]), len(i))
}

// at64 is used to return the raw uint64 for a given pixel. Accessing
// this separately to a (potentially lossy) conversion to a Gray16 is
// necessary for SqImage to function accurately.
func (i Image) at64(x, y int) uint64 {
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

func (i Image) At(x, y int) color.Color {
	c := i.at64(x, y)
	return color.Gray16{uint16(c)}
}

func (i Image) set64(x, y int, c uint64) {
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

func (i Image) Set(x, y int, c color.Color) {
	gray := color.Gray16Model.Convert(c).(color.Gray16).Y
	i.set64(x, y, uint64(gray))
}

// NewImage returns a new integral Image with the given bounds.
func NewImage(r image.Rectangle) *Image {
	w, h := r.Dx(), r.Dy()
	var rows Image
	for i := 0; i < h; i++ {
		col := make([]uint64, w)
		rows = append(rows, col)
	}
	return &rows
}

func (i SqImage) ColorModel() color.Model { return Image(i).ColorModel() }

func (i SqImage) Bounds() image.Rectangle {
	return Image(i).Bounds()
}

func (i SqImage) At(x, y int) color.Color {
	c := Image(i).at64(x, y)
	rt := math.Sqrt(float64(c))
	return color.Gray16{uint16(rt)}
}

func (i SqImage) Set(x, y int, c color.Color) {
	gray := uint64(color.Gray16Model.Convert(c).(color.Gray16).Y)
	Image(i).set64(x, y, gray * gray)
}

// NewSqImage returns a new squared integral Image with the given bounds.
func NewSqImage(r image.Rectangle) *SqImage {
	i := NewImage(r)
	s := SqImage(*i)
	return &s
}

func lowest(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func highest(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (i Image) topLeft(r image.Rectangle) uint64 {
	x := highest(r.Min.X, 0)
	y := highest(r.Min.Y, 0)
	return i[y][x]
}

func (i Image) topRight(r image.Rectangle) uint64 {
	x := lowest(r.Max.X, i.Bounds().Dx() - 1)
	y := highest(r.Min.Y, 0)
	return i[y][x]
}

func (i Image) bottomLeft(r image.Rectangle) uint64 {
	x := highest(r.Min.X, 0)
	y := lowest(r.Max.Y, i.Bounds().Dy() - 1)
	return i[y][x]
}

func (i Image) bottomRight(r image.Rectangle) uint64 {
	x := lowest(r.Max.X, i.Bounds().Dx() - 1)
	y := lowest(r.Max.Y, i.Bounds().Dy() - 1)
	return i[y][x]
}

// Sum returns the sum of all pixels in a rectangle
func (i Image) Sum(r image.Rectangle) uint64 {
	return i.bottomRight(r) + i.topLeft(r) - i.topRight(r) - i.bottomLeft(r)
}

// Mean returns the average value of pixels in a rectangle
func (i Image) Mean(r image.Rectangle) float64 {
	in := r.Intersect(i.Bounds())
	return float64(i.Sum(r)) / float64(in.Dx() * in.Dy())
}

// Sum returns the sum of all pixels in a rectangle
func (i SqImage) Sum(r image.Rectangle) uint64 {
	return Image(i).Sum(r)
}

// Mean returns the average value of pixels in a rectangle
func (i SqImage) Mean(r image.Rectangle) float64 {
	return Image(i).Mean(r)
}

// MeanStdDev calculates the mean and standard deviation of a
// section of an image, using the corresponding regular and square
// integral images.
func MeanStdDev(i Image, sq SqImage, r image.Rectangle) (float64, float64) {
	imean := i.Mean(r)
	smean := sq.Mean(r)

	variance := smean - (imean * imean)

	return imean, math.Sqrt(variance)
}
