package imgconv

import (
	"image"
	"image/color"
	"math"
	"math/rand"
)

// uiImage has few colours and sharp edges, like charts and screenshots.
func uiImage(w, h int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := color.NRGBA{240, 240, 240, 255}
			if (x/32+y/32)%2 == 0 {
				c = color.NRGBA{31, 58, 95, 255}
			}
			if y%40 < 12 && x%100 < 75 {
				c = color.NRGBA{76, 155, 232, 255}
			}
			img.SetNRGBA(x, y, c)
		}
	}
	return img
}

// photoImage has smooth gradients plus noise, like a photograph.
func photoImage(w, h int) *image.NRGBA {
	rng := rand.New(rand.NewSource(1))
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r := 128 + 100*math.Sin(float64(x)/90) + rng.Float64()*20
			g := 128 + 100*math.Cos(float64(y)/70) + rng.Float64()*20
			b := 128 + 80*math.Sin(float64(x+y)/120) + rng.Float64()*20
			img.SetNRGBA(x, y, color.NRGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}
	return img
}
