package meter

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/go-audio/wav"
	"github.com/nfnt/resize"
)

type VUMeter struct {
	Width   int
	Height  int
	Bars    int
	BPM     int
	Decoder *wav.Decoder
}

func (vumeter *VUMeter) Render() (out image.Image, err error) {

	if !vumeter.Decoder.IsValidFile() {
		err = errors.New("The file is invalid")

		return
	}

	/*dur, durErr := vumeter.Decoder.Duration()
	if durErr != nil {
		return
	}

	fmt.Println(dur)*/

	barWidth := int(math.Ceil(float64(vumeter.Width) / (float64(vumeter.Bars) * 2)))

	imageWidth := vumeter.Bars * barWidth * 2
	imageHeight := vumeter.Height * 2

	im := image.NewRGBA(image.Rectangle{Max: image.Point{X: imageWidth, Y: imageHeight}})

	for bar := 0; bar < vumeter.Bars; bar++ {
		barStart := bar * barWidth * 2

		barHeight := barWidth

		barY := vumeter.Height - barHeight

		draw.Draw(im, image.Rect(barStart, barY, barStart+barWidth, barY+(barHeight*2)), image.White, image.ZP, draw.Src)
	}

	for x := 0; x < imageWidth; x++ {

		for y := 0; y < vumeter.Height; y++ {
			r, g, b, a := im.At(x, y).RGBA()

			im.SetRGBA(x, y, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
		}
	}

	out = resize.Resize(uint(vumeter.Width), uint(vumeter.Height), im, resize.Lanczos3)

	return
}
