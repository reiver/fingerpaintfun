package ui

import (
	"bytes"
	"embed"
	"log"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
)

//go:embed sounds
var soundFS embed.FS

var (
	soundBrush   *beep.Buffer
	soundPop     *beep.Buffer
	soundWhoosh  *beep.Buffer
	soundSplash  *beep.Buffer
	soundMuted   bool
	soundReady   bool
	soundFormat  beep.Format
)

// initSound loads all sound effects and initializes the speaker.
func initSound() {
	format := beep.Format{
		SampleRate:  beep.SampleRate(44100),
		NumChannels: 2,
		Precision:   2,
	}
	soundFormat = format

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Millisecond*50)); err != nil {
		log.Printf("sound: failed to init speaker: %v", err)
		return
	}

	soundBrush = loadWAV("sounds/brush.wav")
	soundPop = loadWAV("sounds/pop.wav")
	soundWhoosh = loadWAV("sounds/whoosh.wav")
	soundSplash = loadWAV("sounds/splash.wav")
	soundReady = true
}

func loadWAV(name string) *beep.Buffer {
	data, err := soundFS.ReadFile(name)
	if err != nil {
		log.Printf("sound: %s not found, skipping", name)
		return nil
	}
	streamer, format, err := wav.Decode(bytes.NewReader(data))
	if err != nil {
		log.Printf("sound: failed to decode %s: %v", name, err)
		return nil
	}
	defer streamer.Close()

	buf := beep.NewBuffer(format)
	buf.Append(streamer)
	return buf
}

func playBuffer(buf *beep.Buffer) {
	if !soundReady || soundMuted || buf == nil {
		return
	}
	speaker.Play(buf.Streamer(0, buf.Len()))
}

// PlayBrush plays the drawing sound.
func PlayBrush()   { playBuffer(soundBrush) }

// PlayPop plays the color/tool switch sound.
func PlayPop()     { playBuffer(soundPop) }

// PlayWhoosh plays the undo sound.
func PlayWhoosh()  { playBuffer(soundWhoosh) }

// PlaySplash plays the clear sound.
func PlaySplash()  { playBuffer(soundSplash) }

// SetSoundMuted toggles sound on/off.
func SetSoundMuted(muted bool) { soundMuted = muted }

// IsSoundMuted returns whether sound is muted.
func IsSoundMuted() bool { return soundMuted }
