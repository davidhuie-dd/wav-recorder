package ds

import (
	"log"
	"sync"

	"github.com/asticode/go-astideepspeech"
)

type DeepSpeechServer struct {
	mu        sync.Mutex
	modelPath string
	beamWidth int
	model     *astideepspeech.Model
	stream    *astideepspeech.Stream
}

func NewDeepSpeechServer(modelPath string, beamWidth int, lmPath, triePath string,
	lmWeight, validWordCountWeight float64) *DeepSpeechServer {

	model := astideepspeech.New(modelPath, beamWidth)
	model.EnableDecoderWithLM(lmPath, triePath, lmWeight, validWordCountWeight)

	return &DeepSpeechServer{
		model: astideepspeech.New(modelPath, beamWidth),
	}
}

type CreateStreamArgs struct{}

type CreateStreamResp struct{}

// func (d *DeepSpeechServer) CreateStream(_ CreateStreamArgs, _ *CreateStreamResp) error {
// 	d.mu.Lock()
// 	defer d.mu.Unlock()

// 	if d.stream != nil {
// 		return errors.New("error: stream already exists")
// 	}

// 	log.Printf("Starting new stream")

// 	d.stream = astideepspeech.CreateStream(d.model)

// 	return nil
// }

type IntermediateDecodeArgs struct{}

type IntermediateDecodeResp struct {
	Decoding string
}

func (d *DeepSpeechServer) IntermediateDecode(_ IntermediateDecodeArgs, resp *IntermediateDecodeResp) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.stream == nil {
		d.stream = astideepspeech.CreateStream(d.model)
	}

	log.Printf("Processing intermediate decoding")

	*resp = IntermediateDecodeResp{
		Decoding: d.stream.IntermediateDecode(),
	}

	return nil
}

type FinishStreamArgs struct{}

type FinishStreamResp struct {
	Decoding string
}

func (d *DeepSpeechServer) FinishStream(_ FinishStreamArgs, resp *FinishStreamResp) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.stream == nil {
		*resp = FinishStreamResp{}
		return nil
	}

	log.Printf("Finishing stream")

	*resp = FinishStreamResp{
		Decoding: d.stream.FinishStream(),
	}

	d.stream = nil

	return nil
}

type FeedAudioContentArgs struct {
	Buffer     []int16
	BufferSize uint
}

type FeedAudioContentResp struct {
}

func (d *DeepSpeechServer) FeedAudioContent(args FeedAudioContentArgs, _ *FeedAudioContentResp) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.stream == nil {
		d.stream = astideepspeech.CreateStream(d.model)
	}

	log.Printf("Feeding audio content")

	d.stream.FeedAudioContent(args.Buffer, args.BufferSize)

	return nil
}

type SpeechToTextArgs struct {
	Buffer     []int16
	BufferSize uint
}

type SpeechToTextResp struct {
	Decoding string
}

func (d *DeepSpeechServer) SpeechToText(args SpeechToTextArgs, resp *SpeechToTextResp) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	log.Printf("Converting speech to text")

	text := d.model.SpeechToText(args.Buffer, args.BufferSize)
	*resp = SpeechToTextResp{
		Decoding: text,
	}

	return nil
}
