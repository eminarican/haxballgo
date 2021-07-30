package room

import (
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

type color colorful.Color

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
