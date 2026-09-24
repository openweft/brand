// SPDX-License-Identifier: BSD-3-Clause

// Command check verifies the brand assets the READMEs actually depend on.
//
// ⛔ The check that matters is the LAST one: a mark that rasterises to a blank
// image. That is not hypothetical here. openweft's glyph is drawn with four
// <line> elements, gfx/svg did not implement <line>, and a shape absent from
// its element switch is not rejected -- it is silently not painted. The banner
// generated at the right size, with its wordmark, and an empty space where the
// mark belongs. Nobody published it, 56 READMEs showed a broken image for
// three months, and nothing anywhere reported a fault.
//
// Sizes and file counts would all have passed. Only looking at the pixels
// would have caught it, so that is what this does.
package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
)

func main() {
	social, err := filepath.Glob("social/*.png")
	if err != nil || len(social) == 0 {
		fmt.Println("::error::no social banners found")
		os.Exit(1)
	}
	bad := 0
	for _, p := range social {
		f, err := os.Open(p)
		if err != nil {
			fmt.Printf("::error::%s: %v\n", p, err)
			bad++
			continue
		}
		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			fmt.Printf("::error::%s: not a decodable image: %v\n", p, err)
			bad++
			continue
		}
		b := img.Bounds()
		if b.Dx() != 1280 || b.Dy() != 640 {
			fmt.Printf("::error::%s is %dx%d, not 1280x640\n", p, b.Dx(), b.Dy())
			bad++
		}
		// The glyph sits in the LEFT third. A banner whose left third carries
		// no ink is the exact shape of the <line> failure: background and
		// wordmark present, mark missing.
		if !inked(img, b.Min.X, b.Min.X+b.Dx()/3) {
			fmt.Printf("::error::%s has nothing drawn in its left third: "+
				"the mark did not render\n", p)
			bad++
		}
	}
	fmt.Printf("checked %d social banner(s)\n", len(social))
	if bad > 0 {
		os.Exit(1)
	}
}

// inked reports whether the column range holds pixels clearly lighter than the
// darkest background tone -- the white of the mark against the gradient.
func inked(img image.Image, x0, x1 int) bool {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y += 2 {
		for x := x0; x < x1; x += 2 {
			r, g, bl, _ := img.At(x, y).RGBA()
			if r > 0xE000 && g > 0xE000 && bl > 0xE000 {
				return true
			}
		}
	}
	return false
}
