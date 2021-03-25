package meter

import (
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"

	"github.com/go-audio/wav"
	"github.com/nfnt/resize"
)

type VUMeter struct {
	Width     int
	Height    int
	Bars      int
	BPM       float64
	FrameRate float64
}

type PeakData struct {
	BPM             float64 `json:"bpm"`
	TotalChannels   int     `json:"channels"`
	SampleRate      int     `json:"sample_rate"`
	TotalSamples    int     `json:"total_samples"`
	SamplesPerFrame int     `json:"samples_per_frame"`
	BarsData        []int   `json:"bars_data"`
}

func (vumeter *VUMeter) ReadPeaksData(fileName string) (peakData PeakData, err error) {
	file, fileErr := ioutil.ReadFile(fileName)
	if fileErr != nil {
		err = fileErr
		return
	}

	err = json.Unmarshal(file, &peakData)

	return
}

func (vumeter *VUMeter) GeneratePeaksData(decoder *wav.Decoder) (peakData PeakData, err error) {
	if !decoder.IsValidFile() {
		err = errors.New("The file is invalid")

		return
	}

	totalChannels := decoder.NumChans
	sampleRate := decoder.SampleRate

	maxValue := math.Exp2(float64(decoder.BitDepth)) / 2

	totalSamples := float64(totalChannels) * float64(sampleRate)
	samplesPerFrame := int(totalSamples / float64(vumeter.FrameRate))

	var barPeaks []int

	buf, pcmErr := decoder.FullPCMBuffer()
	if pcmErr != nil {
		err = pcmErr
		return
	}

	var peakValue int

	for i, s := range buf.Data {

		if s > peakValue {
			peakValue = s
		}

		if i > 0 && i%samplesPerFrame == 0 {

			peakPercent := int((float64(peakValue) / maxValue) * 100)

			barPeaks = append(barPeaks, peakPercent)
			peakValue = 0
		}
	}

	peakData.BPM = vumeter.BPM
	peakData.TotalChannels = int(totalChannels)
	peakData.SampleRate = int(sampleRate)
	peakData.SamplesPerFrame = samplesPerFrame
	peakData.BarsData = barPeaks

	return
}

func (vumeter *VUMeter) Render(peakData PeakData, frame int) (out image.Image, err error) {

	barWidth := int(math.Ceil(float64(vumeter.Width) / (float64(vumeter.Bars) * 2)))

	imageWidth := vumeter.Bars * barWidth * 2
	imageHeight := vumeter.Height * 2

	im := image.NewRGBA(image.Rectangle{Max: image.Point{X: imageWidth, Y: imageHeight}})

	var fullBars int

	if frame < len(peakData.BarsData) {
		fullBars = int(float64(vumeter.Bars) * (float64(peakData.BarsData[frame]) / 100.0))
	}

	barsOffset := (vumeter.Bars - fullBars) / 2

	anglePerBar := math.Pi / float64(fullBars)

	for bar := 0; bar < vumeter.Bars; bar++ {
		barStart := bar * barWidth * 2

		var barHeight int

		//barHeight := barWidth
		if bar > barsOffset && bar < (barsOffset+fullBars) {

			barAngle := float64((bar - barsOffset)) * anglePerBar

			barHeight = int(float64(vumeter.Height) * math.Sin(barAngle))

			if barHeight < barWidth {
				barHeight = barWidth
			}
		} else {
			barHeight = barWidth
		}

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
