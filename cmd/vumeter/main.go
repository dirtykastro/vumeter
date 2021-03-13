package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/go-audio/wav"

	"github.com/dirtykastro/graphicutils"
	meter "github.com/dirtykastro/vumeter"
)

func main() {
	audioFile := flag.String("audio", "", "audio file path")
	outputFolder := flag.String("folder", "", "destination folder")
	width := flag.Int("width", 200, "image width")
	height := flag.Int("height", 50, "image height")
	bars := flag.Int("bars", 60, "total bars")
	bpm := flag.Float64("bpm", 95.0, "song speed BPM")
	frameRate := flag.Float64("frame_rate", 30.0, "video frame rate")
	frames := flag.Int("frames", 0, "total frames")

	flag.Parse()

	if *audioFile == "" {
		fmt.Println("the audio file is required")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *outputFolder == "" {
		fmt.Println("the destination folder is required")
		flag.PrintDefaults()
		os.Exit(0)
	}

	audioPeakFile := *audioFile + ".pk"

	vumeter := &meter.VUMeter{Width: *width, Height: *height, Bars: *bars, BPM: *bpm, FrameRate: *frameRate}

	var peakData meter.PeakData

	if graphicutils.Exists(audioPeakFile) {
		var err error
		peakData, err = vumeter.ReadPeaksData(audioPeakFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	} else {

		f, err := os.Open(*audioFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var wavDecoder *wav.Decoder

		wavDecoder = wav.NewDecoder(f)

		peakData, err = vumeter.GeneratePeaksData(wavDecoder)

		peakFile, _ := json.MarshalIndent(peakData, "", " ")

		_ = ioutil.WriteFile(audioPeakFile, peakFile, 0644)
	}

	for frame := 0; frame < *frames; frame++ {

		im, err := vumeter.Render(peakData, frame)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		out, err := os.Create(*outputFolder + "/vumeter" + strconv.Itoa(10000+frame) + ".png")
		if err != nil {
			fmt.Println("Error:", err)
		}

		defer out.Close()

		png.Encode(out, im)
	}
}
