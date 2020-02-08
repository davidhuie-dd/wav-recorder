package file

import (
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"

	"github.com/cryptix/wav"
)

type SpeechToTextArgs struct {
	Buffer     []int16
	BufferSize uint
}

type SpeechToTextResp struct {
	Decoding string
}

type FileProcessor struct {
	client *rpc.Client
}

func NewFileProcessor(client *rpc.Client) *FileProcessor {
	return &FileProcessor{client: client}
}

func (fp *FileProcessor) Transcribe(path string) error {
	i, err := os.Stat(path)
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	r, err := wav.NewReader(f, i.Size())
	if err != nil {
		return err
	}

	var d []int16
	for {
		s, err := r.ReadSample()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		d = append(d, int16(s))
	}

	log.Printf("Transcribing file: %s", path)

	args := &SpeechToTextArgs{
		Buffer:     d,
		BufferSize: uint(len(d)),
	}
	resp := &SpeechToTextResp{}
	if err := fp.client.Call("DeepSpeechServer.SpeechToText", args, resp); err != nil {
		return fmt.Errorf("error feeding audio content: %w", err)
	}

	log.Printf("Transcription: %s", resp.Decoding)

	return nil
}
