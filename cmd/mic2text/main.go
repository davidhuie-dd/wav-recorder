package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"time"

	"github.com/davidhuie-dd/mic2text/mic"
	"github.com/davidhuie-dd/mic2text/transcribe"
)

func main() {
	dev := flag.Int("dev", 0, "the device to read from")
	listDev := flag.Bool("list", false, "list the devices")
	addr := flag.String("addr", "localhost:3000", "the address to the mic2text server")
	pollInterval := flag.Duration("poll", 5*time.Second, "how often to poll for a transcription")
	flag.Parse()

	if *listDev {
		if err := mic.PrintDevices(); err != nil {
			panic(fmt.Errorf("error printing devices: %w", err))
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
