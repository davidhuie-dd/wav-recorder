package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"time"

	"github.com/davidhuie-dd/mic2text/file"
	"github.com/davidhuie-dd/mic2text/mic"
	"github.com/davidhuie-dd/mic2text/transcribe"
)

func main() {
	dev := flag.Int("dev", 0, "the device to read from")
	listDev := flag.Bool("list", false, "list the devices")
	addr := flag.String("addr", "localhost:3000", "the address to the mic2text server")
	pollInterval := flag.Duration("poll", 5*time.Second, "how often to poll for a transcription")
	seconds := flag.Int("seconds", 0, "listen to the microphone for this number of seconds instead of streaming")
	path := flag.String("file", "", "a wav file to transcribe")
	flag.Parse()

	if *listDev {
		if err := mic.PrintDevices(); err != nil {
			panic(fmt.Errorf("error printing devices: %w", err))
		}
	} else if *seconds > 0 {
		client, err := rpc.DialHTTP("tcp", *addr)
		if err != nil {
			log.Fatal("error dialing mic2text server:", err)
		}

		m := mic.NewMicDevice(client)

		if err := m.ReadFromDeviceForTime(*dev, *seconds); err != nil {
			log.Fatalf("error reading from microphone: %s", err)
		}
	} else if *path != "" {
		client, err := rpc.DialHTTP("tcp", *addr)
		if err != nil {
			log.Fatal("error dialing mic2text server:", err)
		}

		fp := file.NewFileProcessor(client)
		if err := fp.Transcribe(*path); err != nil {
			log.Fatal("error transcribing:", err)
		}
	} else {
		client, err := rpc.DialHTTP("tcp", *addr)
		if err != nil {
			log.Fatal("error dialing mic2text server:", err)
		}

		c := make(chan struct{})
		m := mic.NewMicDevice(client)

		go func() {
			if err := m.ReadFromDevice(*dev, c); err != nil {
				log.Fatal("error reading from microphone:", err)
			}
		}()

		trans := transcribe.NewTranscriber(client, *pollInterval)
		trans.Start()
	}
}
