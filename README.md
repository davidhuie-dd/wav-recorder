# wav-recorder

Record wav audio from a system device using [PortAudio](http://www.portaudio.com/).

## Installation

A [recent version](https://golang.org/) of Go is required.

Then, PortAudio must also be installed on your system:
```bash
$ brew install portaudio
```

To install wav-recorder, run:
```bash
$ go get github.com/davidhuie-dd/wav-recorder
```

## Usage

```
Usage of wav-recorder:
  -dest string
    	where to place the output wav file (default "out.wav")
  -dev int
    	the device to read from
  -len duration
    	the amount of time to record (default 10s)
  -list
    	list the available audio devices
```
