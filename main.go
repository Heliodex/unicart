package main

import (
	"fmt"
	"os"

	"golang.org/x/image/webp"
)

// 2x2 grid
type Grid = [4]bool

var chars = []rune{' ', '▘', '▝', '▀', '▖', '▌', '▞', '▛', '▗', '▚', '▐', '▜', '▄', '▙', '▟', '█'}

func char(grid Grid) (r rune) {
	for i, b := range grid {
		if b {
			r += 1 << i
		}
	}
	return chars[r]
}

func main() {
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
}
