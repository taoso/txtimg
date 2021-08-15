package main

import (
	"bufio"
	"flag"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	dpi      = flag.Float64("dpi", 240, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 14, "font size in points")
	width    = flag.Int("width", 940, "image width in points")
	padding  = flag.Int("padding", 10, "text left and right padding")
	height   = flag.Int("height", 400, "image height in points")
	chars    = flag.Int("chars", 20, "chars displayed per line")
	spacing  = flag.Float64("spacing", 1.0, "line spacing")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
)

func main() {
	parseFlags()

	f := readFont()

	fg, bg := image.Black, image.White
	if *wonb {
		fg, bg = bg, fg
	}

	rgba := image.NewRGBA(image.Rect(0, 0, *width, *height))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)

	// Freetype context
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetSrc(fg)
	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	opts := truetype.Options{}
	opts.Size = *size
	opts.DPI = *dpi
	face := truetype.NewFace(f, &opts)

	// Calculate the widths and print to image
	pt := freetype.Pt(*padding, c.PointToFixed(*size).Round())
	newline := func() {
		pt.X = fixed.Int26_6(*padding) << 6
		pt.Y += c.PointToFixed(*size * *spacing)
	}

	var err error
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		for _, x := range []rune(scanner.Text()) {
			w, _ := face.GlyphAdvance(x)
			if x == '\t' {
				x = ' '
			} else if f.Index(x) == 0 {
				continue
			} else if pt.X.Round()+w.Round() > *width-*padding {
				newline()
			}

			pt, err = c.DrawString(string(x), pt)
			if err != nil {
				log.Fatal(err)
			}
		}
		newline()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	saveImage(rgba)
}

func parseFlags() {
	flag.Parse()

	if *chars > 0 {
		*size = float64(*width-*padding*2) / float64(*chars) * 72 / *dpi
	}
}

func readFont() *truetype.Font {
	b, err := os.ReadFile(*fontfile)
	if err != nil {
		log.Panic(err)
	}
	f, err := truetype.Parse(b)
	if err != nil {
		log.Panic(err)
	}

	return f
}

func saveImage(rgba *image.RGBA) {
	bf := bufio.NewWriter(os.Stdout)
	if err := png.Encode(bf, rgba); err != nil {
		log.Panic(err)
	}
	if err := bf.Flush(); err != nil {
		log.Panic(err)
	}
}
