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
var yellow = struct{
	r []uint8{},
	g []uint8{},
	b []uint8{},
	a []uint8{},
}

var black = struct {
	r []uint8{},
	g []uint8{},
	b []uint8{},
	a []uint8{},
}

var white = struct {
	r []uint8{},
	g []uint8{},
	b []uint8{},
	a []uint8{},
}
*/

var cooldown = "resources/cooldown-icons/all-skills"
var active = "resources/class-skills"



func LoadAllImages(option int) (images []image.Image, err error) {
	var dir string
	if option == 1 {
		dir = cooldown
	} else if option == 2 {
		dir = active 
	} else {
		return images, err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir(){
			pngs := getImagesFromSubDir(file)
			for png := range pngs {
				imageDir := dir + "/" + png.Name()
				image := getImage(imageDir, png.Name())
				checkForCooldown(image, png)
			}
		} else {
			imageDir := dir + "/" + file.Name()
			image := getImage(imageDir, file.Name())
			checkForCooldown(image, file.Name())
		}
	}

	return images, err
}


func checkForCooldown(png image.Image, fileName string){
	center_x_min, center_x_max, center_y_min, center_y_max := 13, 19, 13, 19
		bounds := png.Bounds()
		min := bounds.Min
		max := bounds.Max

		fmt.Println("\n", file)

		corner_y_min, corner_x_max, corner_y_min, corner_y_max := 27, 30, 4, 12
		for x := corner_x_min; x <= corner_x_max; x++ {
			fmt.Println("x: ", x)
			for y := corner_y_min; y <= corner_x_max; y++ {
				if corner_x_min > min.X && corner_x_max < max.X && corner_y_min > min.Y && corner_y_max < max.Y {
					color := png.At(x, y)
					r, g, b, a := color.RGBA()
					fmt.Println(byte(r), byte(g), byte(b), byte(a))
				} else {
					fmt.Printf(" (%v, %v) not in image\n", x, y)
				}
			}
		}

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

func getImagesFromSubDir(file fs.DirEntry) []fs.DirEntry{
	images,err := os.ReadDir(file)
	if err != nil {
		panic(err)
	}
	return images

}

func getImage(file fs.DirEntry, name string) image.Image{
	data, err := os.ReadFile(dir + "/" + file.Name())
	if err != nil {
		panic(err)
	}
	image, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	checkForCooldown(image, file)
}
