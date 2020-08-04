// Copyright 2020 Nick White.
// Use of this source code is governed by the GPLv3
// license that can be found in the LICENSE file.

package integral_test

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"

	"rescribe.xyz/integral"
)

func ExampleImage_Sum() {
	f, err := os.Open("testdata/in.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	b := img.Bounds()
	in := integral.NewImage(b)
	draw.Draw(in, b, img, b.Min, draw.Src)
	fmt.Printf("Sum: %d\n", in.Sum(b))
	// Output:
	// Sum: 601340165
}

func ExampleImage_Mean() {
	f, err := os.Open("testdata/in.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	b := img.Bounds()
	in := integral.NewImage(b)
	draw.Draw(in, b, img, b.Min, draw.Src)
	fmt.Printf("Mean: %f\n", in.Mean(b))
	// Output:
	// Mean: 54677.229042
}

func ExampleMeanStdDev() {
	f, err := os.Open("testdata/in.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	b := img.Bounds()
	in := integral.NewImage(b)
	sq := integral.NewSqImage(b)
	draw.Draw(in, b, img, b.Min, draw.Src)
	draw.Draw(sq, b, img, b.Min, draw.Src)
	mean, stddev := integral.MeanStdDev(*in, *sq, b)
	fmt.Printf("Mean: %f, Standard Deviation: %f\n", mean, stddev)
	// Output:
	// Mean: 54677.229042, Standard Deviation: 21643.721672
}
