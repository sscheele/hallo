package audio

import (
	"testing"
)

func TestPlay(t *testing.T) {
	if err := PlayFile("bell.aiff"); err != nil {
		t.Error(err)
	}
}
