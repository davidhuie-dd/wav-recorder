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

### Example

This command records from device #2 for one minute, storing the output
in `my-song.wav`.

```
$ wav-recorder -dest my-song.wav -dev 2 -len 1m

2020/02/06 15:31:42 Using device 2: Logitech Webcam C930e
2020/02/06 15:31:42 Recording for the following 1m0s...
2020/02/06 15:32:43 ...done!
```
