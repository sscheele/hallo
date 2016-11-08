package sphinx

import (
	"errors"
	"log"
	"unsafe"

	"github.com/xlab/closer"
	"github.com/xlab/pocketsphinx-go/sphinx"
	"github.com/xlab/portaudio-go/portaudio"
)

const (
	samplesPerChannel = 512
	sampleRate        = 16000
	channels          = 1
	sampleFormat      = portaudio.PaInt16
)

var (
	hmm  = "/usr/local/share/pocketsphinx/model/en-us/en-us"              //Sets directory containing acoustic model files for sphinx.
	dict = "/usr/local/share/pocketsphinx/model/en-us/cmudict-en-us.dict" //Sets main pronunciation dictionary (lexicon) input file..
	lm   = "/usr/local/share/pocketsphinx/model/en-us/en-us.lm.bin"       //Sets word trigram language model input file.
)

//AppRun is essentially main
func AppRun() error {
	defer closer.Close()
	closer.Bind(func() {
		log.Println("Bye!")
	})
	if err := portaudio.Initialize(); paError(err) {
		return err
	}
	closer.Bind(func() {
		if err := portaudio.Terminate(); paError(err) {
			return err
		}
	})

	// Init CMUSphinx
	cfg := sphinx.NewConfig(
		sphinx.HMMDirOption(hmm),
		sphinx.DictFileOption(dict),
		sphinx.LMFileOption(lm),
		sphinx.SampleRateOption(sampleRate),
	)

	dec, err := sphinx.NewDecoder(cfg)
	if err != nil {
		return err
	}
	closer.Bind(func() {
		dec.Destroy()
	})
	l := &Listener{
		dec: dec,
	}

	var stream *portaudio.Stream
	if err := portaudio.OpenDefaultStream(&stream, channels, 0, sampleFormat, sampleRate, samplesPerChannel, l.paCallback, nil); paError(err) {
		return err
	}
	closer.Bind(func() {
		if err := portaudio.CloseStream(stream); paError(err) {
			return err
		}
	})

	if err := portaudio.StartStream(stream); paError(err) {
		return err
	}
	closer.Bind(func() {
		if err := portaudio.StopStream(stream); paError(err) {
			return err
		}
	})

	if !dec.StartUtt() {
		return errors.New("[ERR] Sphinx failed to start utterance")
	}
	log.Println(banner)
	closer.Hold()
}

type Listener struct {
	inSpeech   bool
	uttStarted bool
	dec        *sphinx.Decoder
}

// paCallback: for simplicity reasons we process raw audio with sphinx in the this stream callback,
// never do that for any serious applications, use a buffered channel instead.
func (l *Listener) paCallback(input unsafe.Pointer, _ unsafe.Pointer, sampleCount uint,
	_ *portaudio.StreamCallbackTimeInfo, _ portaudio.StreamCallbackFlags, _ unsafe.Pointer) int32 {

	const (
		statusContinue = int32(portaudio.PaContinue)
		statusAbort    = int32(portaudio.PaAbort)
	)

	in := (*(*[1 << 24]int16)(input))[:int(sampleCount)*channels]
	// ProcessRaw with disabled search because callback needs to be relatime
	_, ok := l.dec.ProcessRaw(in, true, false)
	// log.Printf("processed: %d frames, ok: %v", frames, ok)
	if !ok {
		return statusAbort
	}
	if l.dec.IsInSpeech() {
		l.inSpeech = true
		if !l.uttStarted {
			l.uttStarted = true
			log.Println("Listening..")
		}
	} else if l.uttStarted {
		// speech -> silence transition, time to start new utterance
		l.dec.EndUtt()
		l.uttStarted = false
		l.report() // report results
		if !l.dec.StartUtt() {
			closer.Fatalln("[ERR] Sphinx failed to start utterance")
		}
	}
	return statusContinue
}

func (l *Listener) report() {
	hyp, _ := l.dec.Hypothesis()
	if len(hyp) > 0 {
		//log.Printf("    > hypothesis: %s", hyp)
		//TODO: HANDLE USER INPUT HERE!!!
	}
}

func paError(err portaudio.Error) bool {
	return portaudio.ErrorCode(err) != portaudio.PaNoError
}

func paErrorText(err portaudio.Error) string {
	return portaudio.GetErrorText(err)
}

const banner = `
 |_| _ | | _
 | |(_|| |(_)
`
