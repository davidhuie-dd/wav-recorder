package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"flag"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gordonklaus/portaudio"
)

const (
	framesPerBuffer = 256
	bitDepth        = 32
	channels        = 1
)

func printDevices() error {
	if err := portaudio.Initialize(); err != nil {
		return fmt.Errorf("error initializing portaudio: %w", err)
	}

	devs, err := portaudio.Devices()
	if err != nil {
		return fmt.Errorf("error listing audio devices: %w", err)
	}
	for i, d := range devs {
		fmt.Printf("Device #%d: %s\n", i, d.Name)
	}

	return nil
}

func writeAudioToFile(devNum int, dest string, length time.Duration) error {
	if err := portaudio.Initialize(); err != nil {
		return fmt.Errorf("error initializing portaudio: %w", err)
	}

	devs, err := portaudio.Devices()
	if err != nil {
		return fmt.Errorf("error listing portaudio devices: %w", err)
	}

	dev := devs[devNum]
	buf := make([]int32, framesPerBuffer)
	sampleRate := dev.DefaultSampleRate

	log.Printf("Using device %d: %s", devNum, dev.Name)

	stream, err := portaudio.OpenStream(portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   dev,
			Channels: channels,
			Latency:  dev.DefaultLowInputLatency,
		},
		SampleRate:      dev.DefaultSampleRate,
		FramesPerBuffer: framesPerBuffer,
	}, buf)
	if err != nil {
		return fmt.Errorf("error creating portaudio stream: %w", err)
	}
	defer stream.Stop()

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Printf("error closing file: %s", err)
		}
	}()

	enc := wav.NewEncoder(out,
		int(sampleRate),
		bitDepth,
		channels,
		1,
	)
	defer func() {
		if err := enc.Close(); err != nil {
			log.Printf("error encoder file: %s", err)
		}
	}()

	if err := stream.Start(); err != nil {
		return fmt.Errorf("error starting portaudio stream: %w", err)
	}

	log.Printf("Recording for the following %s...", length.String())

	format := audio.Format{
		NumChannels: channels,
		SampleRate:  int(sampleRate),
	}

	data := make([]int, len(buf))
	samples := (int(length.Seconds()) * int(sampleRate)) / framesPerBuffer

	for i := 0; i < samples; i++ {
		if err := stream.Read(); err != nil {
			return fmt.Errorf("error reading audio stream: %w", err)
		}

		for i, v := range buf {
			data[i] = int(v)
		}

		b := &audio.IntBuffer{
			Format:         &format,
			SourceBitDepth: bitDepth,
			Data:           data,
		}
		if err := enc.Write(b); err != nil {
			return fmt.Errorf("error writing stream: %w", err)
		}
	}

	log.Println("...done!")

	return nil
}

func main() {
	dev := flag.Int("dev", 0, "the device to read from")
	listDev := flag.Bool("list", false, "list the available audio devices")
	dest := flag.String("dest", "out.wav", "where to place the output wav file")
	recLength := flag.Duration("len", 10*time.Second, "the amount of time to record")
	flag.Parse()

	if *listDev {
		if err := printDevices(); err != nil {
			panic(fmt.Errorf("error printing devices: %w", err))
		}
	} else {
		if err := writeAudioToFile(*dev, *dest, *recLength); err != nil {
			panic(fmt.Errorf("error writing audio to file: %w", err))
		}
	}
}
