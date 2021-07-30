package room

import (
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/ysmood/gson"
)

type team int
type color colorful.Color

const (
	TeamSpectator team = 0
	TeamRed       team = 1
	TeamBlue      team = 2
)

func (c color) Byte() string {
	return strings.Replace(colorful.Color(c).Hex(), "#", "0x", 1)
}

func ColorHex(hex string) (color, error) {
	clr, err := colorful.Hex(hex)
	return color(clr), err
}

func ColorRgb(r, g, b uint8) color {
	return color(colorful.FastLinearRgb(float64(r/255), float64(g/255), float64(b/255)))
}

type scores struct {
	red        int
	blue       int
	time       float32
	scoreLimit int
	timeLimit  float32
}

func newScores(data map[string]gson.JSON) *scores {
	return &scores{
		red:        data["red"].Int(),
		blue:       data["blue"].Int(),
		time:       float32(data["time"].Num()),
		scoreLimit: data["scoreLimit"].Int(),
		timeLimit:  float32(data["time"].Num()),
	}
}

func (s *scores) Red() int {
	return s.red
}

func (s *scores) Blue() int {
	return s.blue
}

func (s *scores) Time() float32 {
	return s.time
}

func (s *scores) ScoreLimit() int {
	return s.scoreLimit
}

func (s *scores) TimeLimit() float32 {
	return s.timeLimit
}
