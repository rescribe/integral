// Copyright 2020 Nick White.
// Use of this source code is governed by the GPLv3
// license that can be found in the LICENSE file.

package integralimg_test

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"

	"rescribe.xyz/integralimg"
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
	integral := integralimg.NewImage(b)
	draw.Draw(integral, b, img, b.Min, draw.Src)
	fmt.Printf("Sum: %d\n", integral.Sum(b))
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
	integral := integralimg.NewImage(b)
	draw.Draw(integral, b, img, b.Min, draw.Src)
	fmt.Printf("Mean: %f\n", integral.Mean(b))
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
	integral := integralimg.NewImage(b)
	sqIntegral := integralimg.NewSqImage(b)
	draw.Draw(integral, b, img, b.Min, draw.Src)
	draw.Draw(sqIntegral, b, img, b.Min, draw.Src)
	mean, stddev := integralimg.MeanStdDev(*integral, *sqIntegral, b)
	fmt.Printf("Mean: %f, Standard Deviation: %f\n", mean, stddev)
	// Output:
	// Mean: 54677.229042, Standard Deviation: 21643.721672
}
