package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/fs"
	"os"
)

/*
for center
var yellow = struct{
	r []uint8{},
	g []uint8{},
	b []uint8{},
	a []uint8{},
darkest 41,41,5
brightest 245, 231, 11
}

var black = struct {
	r []uint8{},
	g []uint8{},
	b []uint8{},
	a []uint8{},
darkest 4 4 5
brightest 7, 15, 33
}


for corner
var black = struct {
	r []uint8{},
	g []uint8{},
	b []uint8{},
	a []uint8{},
darkest 0, 0, 0
}
var white = struct {
	r []uint8{},
	g []uint8{},
	b []uint8{},
	a []uint8{},
}
*/

var cooldown = "resources/cooldown-icons/all-skills"
var active = "resources/skill-icons"

func LoadAllImages(option int) (images []image.Image, err error) {
	var dir string

	switch option {
	case 1:
		dir = cooldown
	case 2:
		dir = active
	default:
		return images, err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			imageDir := dir + "/" + file.Name()
			pngs := getImagesFromSubDir(imageDir)
			for _, png := range pngs {
				imageDir = dir + "/" + file.Name() + "/" + png.Name()
				image := getImage(imageDir)
				checkForCooldown(image, png.Name())
			}
		} else {
			imageDir := dir + "/" + file.Name()
			image := getImage(imageDir)
			checkForCooldown(image, file.Name())
		}
	}

	return images, err
}

func checkForCooldown(png image.Image, fileName string) {
	bounds := png.Bounds()
	min := bounds.Min
	max := bounds.Max

	fmt.Println("\n", fileName)

	corner_x_min, corner_x_max, corner_y_min, corner_y_max := 27, 30, 4, 14
	for x := corner_x_min; x <= corner_x_max; x++ {
		for y := corner_y_min; y <= corner_y_max; y++ {
			if corner_x_min > min.X && corner_x_max < max.X && corner_y_min > min.Y && corner_y_max < max.Y {
				color := png.At(x, y)
				r, g, b, a := color.RGBA()
				fmt.Println("x: ", x, ", y: ", y, ", rgba: ", byte(r), byte(g), byte(b), byte(a))
			} else {
				fmt.Printf(" (%v, %v) not in image\n", x, y)
			}
		}
	}

	center_x_min, center_x_max, center_y_min, center_y_max := 14, 22, 11, 21
	for x := center_x_min; x < center_x_max; x++ {
		for y := center_y_min; y < center_y_max; y++ {
			if center_x_min > min.X && center_x_max < max.X && center_y_min > min.Y && center_y_max < max.Y {
				color := png.At(x, y)
				r, g, b, a := color.RGBA()
				fmt.Println("x: ", x, ", y: ", y, ", rgba: ", byte(r), byte(g), byte(b), byte(a))
			} else {
				fmt.Printf(" (%v, %v) not in image\n", x, y)
			}
		}
	}

}

func getImagesFromSubDir(file string) []fs.DirEntry {
	images, err := os.ReadDir(file)
	if err != nil {
		panic(err)
	}
	return images

}

func getImage(name string) image.Image {
	data, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	image, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	return image
}
