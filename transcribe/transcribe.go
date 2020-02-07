package transcribe

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

type Transcriber struct {
	interval time.Duration
	client   *rpc.Client
}

func NewTranscriber(client *rpc.Client, interval time.Duration) *Transcriber {
	return &Transcriber{
		client:   client,
		interval: interval,
	}
}

type IntermediateDecodeArgs struct{}

type IntermediateDecodeResp struct {
	Decoding string
}

func (t *Transcriber) getTranscription() (string, error) {
	resp := &IntermediateDecodeResp{}
	if err := t.client.Call("DeepSpeechServer.FinishStream", &IntermediateDecodeArgs{}, resp); err != nil {
		return "", fmt.Errorf("error processing intermediate decode: %w", err)
	}

	return resp.Decoding, nil
}

func (t *Transcriber) Start() {
	tick := time.NewTicker(t.interval)
	for {
		<-tick.C

		trans, err := t.getTranscription()
		if err != nil {
			log.Printf("error getting transcription: %s", err)
		}

		log.Printf("transcription: %s", trans)
	}
}
