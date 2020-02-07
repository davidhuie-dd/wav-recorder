package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/davidhuie-dd/mic2text/ds"
)

func main() {
	model := flag.String("model", "output_graph.pbmm", "the path to the model to used")
	beamWidth := flag.Int("beamWidth", 1024, "the deepspeech beam width")
	lm := flag.String("lm", "lm.binary", "the path to the language model")
	trie := flag.String("trie", "trie", "the path to the trie")
	lmWeight := flag.Float64("lmWeigth", .75, "the language model weight")
	validWordCountWeight := flag.Float64("validWordCountWeight", 1.85, "the valid word count weight")
	addr := flag.String("addr", ":3000", "the address to listen on")
	flag.Parse()

	serv := ds.NewDeepSpeechServer(*model, *beamWidth, *lm, *trie,
		*lmWeight, *validWordCountWeight)

	rpc.Register(serv)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", *addr)
	if e != nil {
		log.Fatal("error opening listener:", e)
	}

	log.Printf("Starting mic2text server on %s", l.Addr().String())

	http.Serve(l, nil)
}
