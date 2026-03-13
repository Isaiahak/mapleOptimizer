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

type Axis int

const (
	Vertical = iota
	Horizontal
)

var xAxis Axis = Vertical
var yAxis Axis = Horizontal

var skillOffset1920x1080 = 49
var skillOffset1366x768 = 35

var blackPxl = struct {
	r uint8
	g uint8
	b uint8
	a uint8
}{
	r: 7,
	g: 15,
	b: 33,
	a: 255,
}

type Rect = struct {
	xStart int64
	xEnd   int64
	yStart int64
	yEnd   int64
	width  int64
	height int64
	size   int64
}

var skillBar1920x1080 = Rect{
	xStart: 1100,
	xEnd:   1920,
	yStart: 950,
	yEnd:   1080,
	width:  720,
	height: 130,
	size:   720 * 130,
}

var skillBar1366x768 = Rect{
	xStart: 500,
	xEnd:   920,
	yStart: 450,
	yEnd:   520,
	width:  420,
	height: 70,
	size:   420 * 70,
}

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

var skillBarMax int = 16
var skillBarMin int = 4
var numOfSkillBars = 2

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
	videoPath := "resources/testAt1920x1080.mp4"
	resolution, err := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", videoPath).Output()
	if err != nil {
		log.Fatal(err)
	}

	dim := strings.Split(strings.TrimSpace(string(resolution)), "x")

	width, err := strconv.ParseInt(dim[0], 10, 32)
	if err != nil {
		panic(err)
	}
	height, err := strconv.ParseInt(dim[1], 10, 32)
	if err != nil {
		panic(err)
	}

	fmt.Println("screen resolution : width: ", width, ", height: ", height)
	// for actual algorithm
	//cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", "fps=1", "-f", "rawvideo", "-pix_fmt", "rgb24", "-")

	//for check with raylib
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", "fps=1", "-f", "rawvideo", "-pix_fmt", "rgba", "-")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	frameSize := height * width * 4
	var buf = make([]byte, frameSize)

	var skillBar Rect
	var skillOffset int
	switch {
	case width == 1920 && height == 1080:
		skillBar = skillBar1920x1080
		skillOffset = skillOffset1920x1080
	case width == 1366 && height == 768:
		skillBar = skillBar1366x768
		skillOffset = skillOffset1366x768
	default:
		return
	}

	fmt.Println("skill bar resolution : width: ", skillBar.width, ", height: ", skillBar.height)

	var skills = make([]byte, skillBar.size)

	rl.InitWindow(int32(skillBar.width), int32(skillBar.height), "vod")
	rl.SetTargetFPS(5)
	defer rl.CloseWindow()

	img := rl.GenImageColor(int(skillBar.width), int(skillBar.height), rl.Black)
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

		skills = skills[:0]
		rowStride := width * 4
		for y := skillBar.yStart; y < skillBar.yStart+skillBar.height; y++ {
			start := y*rowStride + skillBar.xStart*4
			end := start + skillBar.width*4
			skills = append(skills, buf[start:end]...)
		}

		verticalOffset := getOffset(skills, skillBar, xAxis)
		horizontalOffset := getOffset(skills, skillBar, yAxis)

		fmt.Println("vertical offset starts at: ", verticalOffset, "\nhorizontal offsey starts at: ", horizontalOffset)

		rgba := unsafe.Slice((*color.RGBA)(unsafe.Pointer(&skills[0])), skillBar.height*skillBar.width)

		//draw vertical offset
		for i := 0; i < int(skillBar.height); i++ {
			for col := horizontalOffset; col < int(skillBar.width); col += skillOffset {
				idx := i*int(skillBar.width) + col
				rgba[idx] = color.RGBA{255, 255, 255, 255}
			}
		}

		//draw horizontal offset
		var pos int
		for j := 0; j < int(skillBar.width); j++ {
			pos = j + int(skillBar.width)*verticalOffset
			rgba[pos] = color.RGBA{255, 255, 255, 255}
		}

		rl.UpdateTexture(texture, rgba)

		/*
			takeScreenshot(texture, "resources/skillbar_snapshot_"+strconv.FormatInt(num, 10)+".png")
			num++
		*/

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawTexture(texture, 0, 0, rl.White)
		rl.EndDrawing()
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("end of video")
	}

}

func checkRowForBlack(buf []byte) bool {
	totalPixels := len(buf) / 4
	blackPixels := 0
	for i := 0; i < totalPixels; i += 4 {
		r, g, b := buf[i], buf[i+1], buf[i+2]
		if r <= blackPxl.r && g <= blackPxl.g && b <= blackPxl.b {
			blackPixels++
		}
	}
	if blackPixels != totalPixels && float32(blackPixels) > float32(totalPixels)*0.2 {
		return true
	}
	return false
}

func getOffset(buf []byte, skillBar Rect, axis Axis) int {
	histogram := checkHistogram(buf, skillBar, axis)
	correlationResult := checkCorrelation(histogram, skillBar, axis)
	minIndex := 0
	minValue := correlationResult[0]
	for i, v := range correlationResult {
		if v < minValue {
			minValue = v
			minIndex = i
		}
	}
	return minIndex

}

func checkCorrelation(histogram []float32, skillBar Rect, axis Axis) []float32 {
	var offset int
	var size int
	switch axis {
	case Vertical:
		offset = 50
		size = int(skillBar.width)

	case Horizontal:
		offset = 50
		size = int(skillBar.height)
	}

	var correlationResult = make([]float32, offset)
	i := 1
	var width = 2
	for ; i < offset; i++ {
		var correlationSum float32 = 0
		for pos := 0; pos < size-i-width; pos++ {
			colSum1 := histogram[pos] * histogram[pos+i]
			colSum2 := histogram[pos+1] * histogram[pos+i+1]
			colSum3 := histogram[pos+width] * histogram[pos+i+width]
			correlationSum += (colSum1 + colSum2 + colSum3) / 3

		}
		correlationResult[i-1] = correlationSum
	}
	return correlationResult
}

func checkHistogram(buf []byte, skillBar Rect, axis Axis) []float32 {
	var indexSize int64
	switch axis {
	case Vertical:
		indexSize = skillBar.width
	case Horizontal:
		indexSize = skillBar.height
	}
	var histogram = make([]float32, indexSize)

	var i int64 = 0
	var index int64 = 0

	switch axis {
	case Vertical:
		for ; i < indexSize; i++ {
			var j int64 = 0
			for ; j < skillBar.width; j += 4 {
				rowStride := skillBar.width*4 + i
				histogram[i] += float32(buf[rowStride+j] + buf[rowStride+j+1] + buf[rowStride+j+2])
				//histogram[i] += float32(buf[rowStride + j])
				//histogram[i] += float32(buf[rowStride+j+1])
				//histogram[i] += float32(buf[rowStride+j+2])

			}
		}

	case Horizontal:
		for ; i < indexSize; i++ {
			var j int64 = 0
			for ; j < skillBar.width*4; j += 4 {
				histogram[i] += float32(buf[j] + buf[j+1] + buf[j+2])
				//histogram[i] += float32(buf[i])
				//histogram[i] += float32(buf[i+1])
				//histogram[i] += float32(buf[i+2])

				if index == skillBar.width-1 {
					index = 0
				} else {
					index++
				}
			}
		}
	}
	return histogram
}

// texture used to create an img and filepath should be a complete path to the desired file
func takeScreenshot(texture rl.Texture2D, filePath string) {
	img := rl.LoadImageFromTexture(texture)
	rl.ExportImage(*img, filePath)

}
