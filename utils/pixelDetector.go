package utils

import (
	"bytes"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unsafe"
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

func LoadVideo() {
	resolution, err := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", "resources/test.mp4").Output()
	if err != nil {
		log.Fatal(err)
	}

	dim := strings.Split(strings.TrimSpace(string(resolution)), "x")

	height, err := strconv.ParseInt(dim[0], 10, 32)
	if err != nil {
		panic(err)
	}
	width, err := strconv.ParseInt(dim[1], 10, 32)
	if err != nil {
		panic(err)
	}

	fmt.Println(height, width)
	// for actual algorithm
	//cmd := exec.Command("ffmpeg", "-i", "resources/test.mp4", "-vf", "fps=1", "-f", "rawvideo", "-pix_fmt", "rgb24", "-")
	//for check with raylib
	cmd := exec.Command("ffmpeg", "-i", "resources/test.mp4", "-vf", "fps=60", "-f", "rawvideo", "-pix_fmt", "rgba", "-")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	frameSize := height * width * 4
	var buf = make([]byte, frameSize)
	rl.InitWindow(int32(width), int32(height), "screen shot")
	rl.SetTargetFPS(60)
	defer rl.CloseWindow()
	img := rl.GenImageColor(int(width), int(height), rl.Black)
	texture := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	for !rl.WindowShouldClose() {
		_, err = io.ReadFull(stdout, buf)
		if err != nil {
			if err == io.EOF {
				log.Fatal(err)
			} else {
				log.Fatal(err)
			}
		}

		rgba := unsafe.Slice((*color.RGBA)(unsafe.Pointer(&buf[0])), width*height)
		rl.UpdateTexture(texture, rgba)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawTexture(texture, 0, 0, rl.White)
		rl.EndDrawing()
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

}
