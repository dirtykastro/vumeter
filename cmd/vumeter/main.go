package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/go-audio/wav"

	"github.com/dirtykastro/graphicutils"
	meter "github.com/dirtykastro/vumeter"
)

func main() {
	audioFile := flag.String("audio", "", "audio file path")
	outputFile := flag.String("file", "", "destination file name")

	flag.Parse()

	if *audioFile == "" {
		fmt.Println("the audio file is required")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *outputFile == "" {
		fmt.Println("the destination file is required")
		flag.PrintDefaults()
		os.Exit(0)
	}

	audioPeakFile := *audioFile + ".pk"

	vumeter := &meter.VUMeter{Width: 200, Height: 50, Bars: 60, BPM: 88.0}

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

	im, err := vumeter.Render(peakData)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out, err := os.Create(*outputFile)
	if err != nil {
		fmt.Println("Error:", err)
	}

	defer out.Close()

	png.Encode(out, im)

}
