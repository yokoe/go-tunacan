package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

import _ "image/png"

// https://stackoverflow.com/questions/28322997/how-to-get-a-list-of-values-into-a-flag-in-golang?utm_medium=organic&utm_source=google_rich_qa&utm_campaign=google_rich_qa
type arrayFlags []string

func (i *arrayFlags) String() string {
	return "arrayFlags"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var outputFilename string
	var sourceFilenames arrayFlags

	flag.Var(&sourceFilenames, "s", "Source image filepaths.")
	flag.StringVar(&outputFilename, "o", "", "Output filepath.")
	flag.Parse()

	concat(sourceFilenames, outputFilename)
}

func concat(sourceFilenames []string, outputFilename string) {
	canvasWidth := 0
	canvasHeight := 0

	images := LoadImages(sourceFilenames)

	for i := range images {
		srcImg := images[i]
		canvasWidth += srcImg.Bounds().Size().X
		canvasHeight = Max(canvasHeight, srcImg.Bounds().Size().Y)

		fmt.Println(srcImg.Bounds().Size())
	}

	fmt.Println("Canvas size: ", canvasWidth, canvasHeight)

	outputImage := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))

	x := 0
	for i := range images {
		srcImage := images[i]

		srcRect := image.Rect(x, 0, srcImage.Bounds().Size().X+x, srcImage.Bounds().Size().Y)
		draw.Draw(outputImage, srcRect, srcImage, image.Pt(0, 0), draw.Src)

		x += srcImage.Bounds().Size().X
	}

	file, _ := os.Create(outputFilename)
	defer file.Close()

	if err := jpeg.Encode(file, outputImage, &jpeg.Options{100}); err != nil {
		panic(err)
	}
}

func LoadImages(filenames []string) []image.Image {
	images := []image.Image{}
	for i := range filenames {
		filename := filenames[i]
		fmt.Println(filename)

		src, _ := os.Open(filename)
		defer src.Close()

		srcImg, _, err := image.Decode(src)
		if err != nil {
			panic(err)
		}

		images = append(images, srcImg)
	}
	return images
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
