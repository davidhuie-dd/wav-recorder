package mic

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/gordonklaus/portaudio"
)

const (
	sampleRate      = 8000
	framesPerBuffer = 4000
	channels        = 1
)

type FeedAudioContentArgs struct {
	Buffer     []int16
	BufferSize uint
}

type FeedAudioContentResp struct {
}

func PrintDevices() error {
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

func NewMicDevice(client *rpc.Client) *MicDevice {
	return &MicDevice{
		client: client,
	}
}

type MicDevice struct {
	client *rpc.Client
}

func (m *MicDevice) ReadFromDevice(devNum int, stop chan struct{}) error {
	if err := portaudio.Initialize(); err != nil {
		return fmt.Errorf("error initializing portaudio: %w", err)
	}

	devs, err := portaudio.Devices()
	if err != nil {
		return fmt.Errorf("error listing portaudio devices: %w", err)
	}

	dev := devs[devNum]

	log.Printf("Using device %d: %s", devNum, dev.Name)

	buf := make([]int16, framesPerBuffer)
	stream, err := portaudio.OpenStream(portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   dev,
			Channels: channels,
			Latency:  dev.DefaultLowInputLatency,
		},
		SampleRate:      sampleRate,
		FramesPerBuffer: framesPerBuffer,
	}, buf)
	if err != nil {
		return fmt.Errorf("error creating portaudio stream: %w", err)
	}

	if err := stream.Start(); err != nil {
		return fmt.Errorf("error starting portaudio stream: %w", err)
	}
	defer stream.Stop()

	for {
		select {
		case <-stop:
			return nil
		default:
		}

		if err := stream.Read(); err != nil {
			return fmt.Errorf("error reading audio stream: %w", err)
		}

		args := &FeedAudioContentArgs{
			Buffer:     buf,
			BufferSize: framesPerBuffer,
		}

		if err := m.client.Call("DeepSpeechServer.FeedAudioContent", args, &FeedAudioContentResp{}); err != nil {
			log.Printf("error feeding audio content: %s", err)
		}
	}
}
