portaudio-go
============

The package provides Go bindings for [PortAudio](http://www.portaudio.com). <br/>
All the code has automatically been generated with rules defined in [portaudio.yml](/portaudio.yml).

### Usage

```
$ brew install portaudio
(or use your package manager)

$ go get github.com/xlab/portaudio-go/portaudio
```

See [example#1](https://github.com/xlab/alac-go/blob/master/cmd/alac-player/main.go).

### Rebuilding the package

You will need to get the [cgogen](https://git.io/cgogen) tool installed first.

```
$ git clone https://github.com/xlab/portaudio-go && cd portaudio-go
$ make clean
$ make
```
