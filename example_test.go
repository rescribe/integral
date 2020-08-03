// Copyright 2020 Nick White.
// Use of this source code is governed by the GPLv3
// license that can be found in the LICENSE file.

package integralimg_test

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"log"
	"strings"

	"rescribe.xyz/integralimg"
)

const smallImg = `iVBORw0KGgoAAAANSUhEUgAAAAYAAAADCAAAAACVaiEnAAAGqHpUWHRSYXcgcHJvZmlsZSB0eXBlIGV4aWYAAHjarVhbsusoDPxnFbMEEA/BcsRDVbODWf60sJOTOM69d6omrhgHZJC6pYZz3Prnb3V/4UPRF5cy19JK8fiklhoJHqo/Pm3fg0/7fnzK2Yb3fufjOUDosufzN8tpL+jPPy881gj9vd/Vc4TqOdFj5XPCaCsTHuark+inoz+kc6K2To9b5VdXOx3tOA23K+c38p76OYn9dq8diYHSzLCKRCuG6Pe9Hh7E4yv41n2HUxgteE4xu6M5PQEgb+E9Wu9fAXoH/3xyV/QfYF7BJzkt4gXL8mCt3A+EfA/+hvhl4fj0iN4HRvf8Ec75VZ1VdR3RSSpAtJwZ5d0DHXsHhh2Qx/1awcX4Zjzzvhqu6sUPkDP98B3XCC0QWFEXUphBgoa12xEGXEy0iNESDYq7r0amRiMaT8muoMSxxQkGKQ5aLkZ009OXsNdte70RKlaeAaYUMFnAK18v96vB/3I51WEQBQOzHljBL7LMhRvGnN1hBUKCnrzlDfDjOun3L4mFVAWDecNcEaD4fkzRc/jJrbh5jrDLaI+qCI7nOQEgwtoZzoQIBnwJMYcSPBNxCMCxgiCB5xQTdTAQcqYJJynFWMgxVbK18Q6HbUuZClk3tMnKB9XE4KZFAVkpZeQPp4ockhxzyjmXzLm63LKUWFLJpRQuJnLCkRNnLsxcubHUWFPNtVSutbYqjVqEBuZWGrfaWhMhJ1hIMJfAXtDTqceeeu6lc6+9dRlIn5FGHmXwqKMNmTTjhEzMMnnW2aas4BaUYqWVV1m86mpLFLmmUZNmLcpatak8WTtZ/bj+A2vhZI02U2bHT9bQ65gfUwSTk2ycgTFKAYyzMYCEJuPM15ASGXPGmW+EosgEJ7Nx42YwxkBhWoGyhid3P8z9EW8u1z/ijX7HnDPq/g/mHKj75O2GtWn73NiMHVVomPqI6tMyhKrD13vcjna0maNi4sw6R9AxZum5pCVNLS26ruV1lVKxEc6Qsm1A2VqXfBhL61DMvZJXk0kedvdZei2yFo/CquQBl66+PO6VZh62+ey5rHXng9Xl2fVo81irT1sAnjVr88BU4HyGzLEL4h4MocTKFJx5j60jv/n5vQ3K5WPIQneX2DviIc5JDeSsCtGHD5ggFYNR7JBwtIvfQnZvMTPHFaNq83HEHYz63aYAsIK88nK0sdbjyb0MYY287LWqK1krK0WFSKwqvc3J8Bu/wWpvJ3WR1umgo4vH11Z7J5tHS6o6c1ccFtZoQaXzYot5L6Pu+/gNVqCEJ6pVc9EBA1Rah01H7TqaZsTKhmidKMGunVAsUA4tsMoFfkn7MrBSPyB3N2nGUrS2nmYtK+I+wxSwnJCmOEJhVwbwyP/uUWa6WNLUVp2B6jeb4HfVfhnHsMSbvo8X3Muov31nL6NUhDqwQnAgzDL9ucjR455dqO36G+N32zNbjlwR90iVxR/mSJoBcVutTJGUKnStRl/jLvA+prGTr0UbvrSMqv2oZiix2GSDJ+1Ja1QHbybLqBMKBB1qTx2yorzWw44ciFUNFs5C1JgUsHdBHiELVhymZ10s4IJyuEvxyJd5EXw4syhUt7NoDssiU+Ax2c+zcEc6Cneto3CpPFi/KSZ3V126rrZWDdkYkp9yET3B2lA5NaQAlCXRBVXsT9hoTSEXHQo5ExQxQjonJ0x0AF7FNNKhcIA0lPNGDVf5ZGxENoFH1gRvE0Fcdtq66Zv03p642V5T34Db28Nvldg9pVg/bLHx4m+Xnhp7wWkOAhGGYlMstmf1voNFAkzjeLiJhO4qMk0bWDM0my0Ik5nbgRKw/5ZIo0HKkFqFx/DcHOEQOQvCNYmxmEENONEV8uQi7QB7gEk9kox2Zu76YysoBlUwcvLRZZl+b/zddhTHa4gu6QkamiG6ReY00HL83cBB4JqmFzjVhrWL7Jnbq7avFBvDWL7RDXzul8Dm0CgUpBmlunpddgDqkCE4hAzpEyftR+m5l9ree9pHYQoOW8gq5HTvmUmwwW07LLrWtL1zJ6v7yVqc4/eRZpqU6E74klAlx5HmStx1RRdH+9gbcZrRsIkftRD+osHBw+5hb+oPCG0v2iBuCJ0+EJyXJcjj4KIZUjKXSMsFhxyQhxdZRD7E1B1qilCv1ZnWxmwWEMa2k32URX0NxdkDDqJih6Z+KPDaAop9AsfMG9RuhL3Ro2h/KvUhkPZfkC2Rh0CaPK5VvqPnbEKI0fsSKbzlCnIayi7YrBdyLN0JurrZIEr/AhzRmfaPMzpJAAABJGlDQ1BJQ0MgcHJvZmlsZQAAeJydkLFKw1AUhr+0okUUB8WhOGRwtOBiJ5eqEASFGCsYndIkxWISQ5JSfAPfRB+mgyD4BD6BgrP/jQ4OZvHC4f84nPP/915o2UmYlgu7kGZV4XgD/9K/spfe6NBmjR26QVjmA9c9ofF8vmIZfekZr+a5P89iFJehdK7KwryowNoX92dVbljFxu3QOxQ/iO0ozSLxk3g7SiPDZtdLk2n442lusxJnF+emr9rC4ZhTXGxGTJmQUNGTZuoc0WdP6lAQcE9JKE2I1ZtppuJGVMrJ4UA0FOk2DXndOs9VykgeE3mZhDtSeZo8zP9+r32c1ZvW5jwPiqButVWt8RjeH2HVh/VnWL5uyOr8flvDTL+e+ecbvwDmi1BkPgPvaAAAAAlwSFlzAAALEwAACxMBAJqcGAAAAAd0SU1FB+QIAw8tLQ1tVrwAAAAdSURBVAjXY34TxMnwlIFJWMhMlYGBiUFIVUiIAQAu0QMPFiP3IQAAAABJRU5ErkJggg==`

func smallImgPNG() io.Reader { return base64.NewDecoder(base64.StdEncoding, strings.NewReader(smallImg)) }

func ExampleImage_Sum() {
	img, _, err := image.Decode(smallImgPNG())
	if err != nil {
		log.Fatal(err)
	}
	b := img.Bounds()
	integral := integralimg.NewImage(b)
	draw.Draw(integral, b, img, b.Min, draw.Src)
	fmt.Printf("Sum: %d\n", integral.Sum(b))
	// Output:
	// Sum: 551008
}

func ExampleImage_Mean() {
	img, _, err := image.Decode(smallImgPNG())
	if err != nil {
		log.Fatal(err)
	}
	b := img.Bounds()
	integral := integralimg.NewImage(b)
	draw.Draw(integral, b, img, b.Min, draw.Src)
	fmt.Printf("Mean: %f\n", integral.Mean(b))
	// Output:
	// Mean: 30611.555556
}

func ExampleMeanStdDev() {
	img, _, err := image.Decode(smallImgPNG())
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
	// Mean: 30611.555556, Standard Deviation: 25372.651592
}
