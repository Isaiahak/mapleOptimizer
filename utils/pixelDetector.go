package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
)

var yellow = [][]uint8{
	[]uint8{},
	[]uint8{},
	[]uint8{},
	[]uint8{},
}

func LoadAllImages() (images []image.Image, err error) {
	dir := "resources/cooldown-icons/all-skills"
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		data, err := os.ReadFile(dir + "/" + file.Name())
		if err != nil {
			panic(err)
		}
		png, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
		center_x_min, center_x_max, center_y_min, center_y_max := 13, 19, 13, 19
		bounds := png.Bounds()
		min := bounds.Min
		max := bounds.Max

		fmt.Println("\n", file.Name())

		for x := center_x_min; x <= center_x_max; x++ {
			fmt.Println("x: ", x)
			for y := center_y_min; y <= center_x_max; y++ {
				if center_x_min > min.X && center_x_max < max.X && center_y_min > min.Y && center_y_max < max.Y {
					color := png.At(x, y)
					r, g, b, a := color.RGBA()
					fmt.Println(byte(r), byte(g), byte(b), byte(a))
				} else {
					fmt.Printf(" (%v, %v) not in image\n", x, y)
				}
			}
		}
	}

	return images, err
}
