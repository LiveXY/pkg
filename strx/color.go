package strx

import (
	"fmt"
	"image/color"
	"strings"
)

func ToColor(scol string, def color.Color) color.Color {
	if len(scol) == 0 {
		return color.Black
	}
	if strings.Index(scol, "#") != 0 {
		scol = "#" + scol
	}
	format := "#%02x%02x%02x"
	var r, g, b uint8
	n, err := fmt.Sscanf(scol, format, &r, &g, &b)
	if err != nil {
		return color.Black
	}
	if n != 3 {
		return color.Black
	}
	col := color.RGBA{uint8(r), uint8(g), uint8(b), 0xff}
	return col
}