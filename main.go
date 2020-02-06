package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gordonklaus/portaudio"
)

func main() {
	if err := portaudio.Initialize(); err != nil {
		panic(err)
	}

	devs, err := portaudio.Devices()
	if err != nil {
		panic(err)
	}
	for i, d := range devs {
		fmt.Printf("device %d: %s\n", i, d.Name)
		// fmt.Printf("device %d: %#v\n", i, d)
	}

	framesPerBuffer := 256

	dev := devs[2]
	buf := make([]int32, framesPerBuffer)

	fmt.Printf("Using device: %#v\n", dev)

	stream, err := portaudio.OpenStream(portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   dev,
			Channels: 1,
			Latency:  dev.DefaultLowInputLatency,
		},
		SampleRate:      dev.DefaultSampleRate,
		FramesPerBuffer: framesPerBuffer,
	}, buf)
	if err != nil {
		panic(err)
	}

	out, err := os.Create("test.wav")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	enc := wav.NewEncoder(out,
		int(dev.DefaultSampleRate),
		32,
		1,
		1,
	)
	defer func() {
		log.Println(enc.Close())
	}()

	if err := stream.Start(); err != nil {
		panic(err)
	}

	format := audio.Format{
		NumChannels: 1,
		SampleRate:  int(dev.DefaultSampleRate),
	}

	totalSecs := 3
	samples := totalSecs * int(dev.DefaultSampleRate) / 256

	fmt.Println(samples)

	for i := 0; i < samples; i++ {
		if err := stream.Read(); err != nil {
			panic(err)
		}

		d := make([]int, len(buf))
		for i, v := range buf {
			d[i] = int(v)
		}

		b := &audio.IntBuffer{
			Format:         &format,
			SourceBitDepth: 32,
			Data:           d,
		}
		if err := enc.Write(b); err != nil {
			panic(err)
		}
	}
}
