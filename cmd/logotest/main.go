package main

import (
	"fmt"
	"image/color"

	"github.com/JyotirmoyDas05/openpilot/internal/tui/components/logo"
)

func main() {
	opts := logo.Opts{
		FieldColor:   color.RGBA{R: 150, G: 150, B: 150, A: 255},
		TitleColorA:  color.RGBA{R: 200, G: 100, B: 200, A: 255},
		TitleColorB:  color.RGBA{R: 100, G: 200, B: 200, A: 255},
		BrandColor:   color.RGBA{R: 255, G: 100, B: 200, A: 255},
		VersionColor: color.RGBA{R: 200, G: 200, B: 200, A: 255},
		Width:        80,
	}
	fmt.Println(logo.Render("v0.0.0", false, opts))
}
