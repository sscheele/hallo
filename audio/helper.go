package audio

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/gordonklaus/portaudio"
)

//PlayFile is sent a path and plays an audio file
func PlayFile(fileName string) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	id, data, err := readChunk(f)
	if err != nil {
		return err
	}
	if id.String() != "FORM" {
		return errors.New("bad file format")
	}
	_, err = data.Read(id[:])
	if err != nil {
		return err
	}
	if id.String() != "AIFF" {
		return errors.New("bad file format")
	}
	var c commonChunk
	var audio io.Reader
	for {
		id, chunk, err := readChunk(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch id.String() {
		case "COMM":
			if err = binary.Read(chunk, binary.BigEndian, &c); err != nil {
				return err
			}
		case "SSND":
			chunk.Seek(8, 1) //ignore offset and block
			audio = chunk
		}
	}

	//redirect portaudio's annoying garbage to null
	tempErr := os.Stderr
	nullPath := "/dev/null"
	if runtime.GOOS == "windows" {
		nullPath = "nul" //windows is, apparently, magic
	}
	nullF, err := os.Open(nullPath)
	if err != nil {
		return err
	}
	redirStdErr(nullF)
	portaudio.Initialize()
	redirStdErr(tempErr)
	//assume 44100 sample rate, mono, 32 bit
	defer portaudio.Terminate()
	out := make([]int32, 8192)
	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, len(out), &out)
	if err != nil {
		return err
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		return err
	}
	defer stream.Stop()
	for remaining := int(c.NumSamples); remaining > 0; remaining -= len(out) {
		if len(out) > remaining {
			out = out[:remaining]
		}
		err := binary.Read(audio, binary.BigEndian, out)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err = stream.Write(); err != nil {
			return err
		}
		select {
		case <-sig:
			return nil
		default:
		}
	}
	return nil
}

func readChunk(r readerAtSeeker) (id ID, data *io.SectionReader, err error) {
	_, err = r.Read(id[:])
	if err != nil {
		return
	}
	var n int32
	err = binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return
	}
	off, _ := r.Seek(0, 1)
	data = io.NewSectionReader(r, off, int64(n))
	_, err = r.Seek(int64(n), 1)
	return
}

func redirStdErr(f *os.File) error {
	if runtime.GOOS == "windows" {
		return nil
	}
	return syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
}

type readerAtSeeker interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

//ID is probably the magic number, but I'n not 100% sure tbh
type ID [4]byte

func (id ID) String() string {
	return string(id[:])
}

type commonChunk struct {
	NumChans      int16
	NumSamples    int32
	BitsPerSample int16
	SampleRate    [10]byte
}
