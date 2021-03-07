package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"

	"github.com/go-audio/wav"

	//"github.com/dirtykastro/graphicutils"
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

	f, err := os.Open(*audioFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var wavDecoder *wav.Decoder

	wavDecoder = wav.NewDecoder(f)

	vumeter := &meter.VUMeter{Width: 200, Height: 50, Bars: 60, Decoder: wavDecoder}

	im, err := vumeter.Render()

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
