package audio

import (
	"errors"
	"io"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

func GetAudioPlayer(r io.Reader) (*oto.Player, error) {
	decodedMp3, err := mp3.NewDecoder(r)
	if err != nil {
		return nil, errors.New("Failed to decode audio: " + err.Error())
	}

	op := &oto.NewContextOptions{}
	op.SampleRate = 44100
	op.ChannelCount = 2
	op.Format = oto.FormatSignedInt16LE
	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		return nil, errors.New("Failed to create oto context: " + err.Error())
	}
	<-readyChan

	player := otoCtx.NewPlayer(decodedMp3)
	return player, nil
}
