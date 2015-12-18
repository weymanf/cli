package yolo

import (
	"encoding/binary"
	"io"
	"os"
	"os/signal"
	"path/filepath"

	"code.google.com/p/portaudio-go/portaudio"

	"github.com/mkb218/gosndfile/sndfile"
)

const CF_YOLO = "CF_YOLO"

type Yolo struct{}

func (y Yolo) Activate() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	var c commonChunk
	var audio io.Reader
	var i sndfile.Info
	cwd, err := os.Getwd()
	var yoloDir string = filepath.Join(cwd, "cf/yolo")
	var sandstormMp3 = filepath.Join(yoloDir, "sandstorm.ogg")
	file, err := sndfile.Open(sandstormMp3, sndfile.Read, &i)

	id, data, err := readChunk(file)
	_, err = data.Read(id[:])

	out := make([]int32, 8192)
	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, len(out), &out)

	defer stream.Close()

	stream.Start()

	for remaining := int(c.NumSamples); remaining > 0; remaining -= len(out) {
		if len(out) > remaining {
			out = out[:remaining]
		}
		err := binary.Read(audio, binary.BigEndian, out)
		if err == io.EOF {
			break
		}
		chk(err)
		chk(stream.Write())
		select {
		case <-sig:
			return
		default:
		}
	}
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

type readerAtSeeker interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

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

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
