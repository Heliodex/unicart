package main

import (
	"fmt"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/webp"
)

// 2x2 grid
type Grid = [4]bool

var unicodeChars = []rune{' ', '▘', '▝', '▀', '▖', '▌', '▞', '▛', '▗', '▚', '▐', '▜', '▄', '▙', '▟', '█'}

var sprigMap = map[rune][]uint32{
	'0': {0, 0, 0, 255},
	'L': {73, 80, 87, 255},
	'1': {145, 151, 156, 255},
	'2': {248, 249, 250, 255},
	'3': {235, 44, 71, 255},
	'C': {139, 65, 46, 255},
	'7': {25, 177, 248, 255},
	'5': {19, 21, 224, 255},
	'6': {254, 230, 16, 255},
	'F': {149, 140, 50, 255},
	'4': {45, 225, 62, 255},
	'D': {29, 148, 16, 255},
	'8': {245, 109, 187, 255},
	'H': {170, 58, 197, 255},
	'9': {245, 113, 23, 255},
	// '.': {0, 0, 0, 0},
}

func char(grid Grid) (r rune) {
	for i, b := range grid {
		if b {
			r += 1 << i
		}
	}
	return unicodeChars[r]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("not enough arguments")
		return
	}

	format := os.Args[1]
	switch format {
	case "unicode":
		r, err := os.Open("test.webp")
		if err != nil {
			fmt.Println(err)
			return
		}
		img, err := webp.Decode(r)
		if err != nil {
			fmt.Println(err)
			return
		}

		bounds := img.Bounds()
		x, y := bounds.Dx(), bounds.Dy()

		if x%2 != 0 || y%2 != 0 {
			fmt.Println("Image dimensions must be multiples of 2")
			return
		}

		for j := range y / 2 {
			for i := range x / 2 {
				var grid Grid
				for k := range 2 {
					for l := range 2 {
						r, g, b, _ := img.At(i*2+l, j*2+k).RGBA()
						grid[k*2+l] = r > 0x7fff || g > 0x7fff || b > 0x7fff
					}
				}
				fmt.Print(string(char(grid)))
			}
			fmt.Println()

		}
	case "sprig":
		r, err := os.Open("test.png")
		if err != nil {
			fmt.Println(err)
			return
		}
		img, err := png.Decode(r)
		if err != nil {
			fmt.Println(err)
			return
		}

		bounds := img.Bounds()
		x, y := bounds.Dx(), bounds.Dy()

		var sb strings.Builder
		sb.WriteString("\t\t\t")

		for j := range y {
			for i := range x {
				r, g, b, a := img.At(i, j).RGBA()
				r, g, b, a = r>>8, g>>8, b>>8, a>>8

				// find the closest color
				var minDist uint32 = 0xffffffff
				var minChar rune
				for char, rgba := range sprigMap {
					dist := (r-rgba[0])*(r-rgba[0]) + (g-rgba[1])*(g-rgba[1]) + (b-rgba[2])*(b-rgba[2]) + (a-rgba[3])*(a-rgba[3])
					if dist < minDist {
						minDist = dist
						minChar = char
					}
				}
				sb.WriteString(string(minChar))
			}
			sb.WriteString("\n\t\t\t")
		}

		// write to file
		f, err := os.Create("output.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		_, err = f.WriteString(sb.String())
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("invalid format")
	}
}
