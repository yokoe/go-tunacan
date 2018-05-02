package tunacan

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"golang.org/x/image/draw"

	_ "image/png"
)

func concat(sourceFilenames []string, outputFilename string) error {
	images, err := loadImages(sourceFilenames)

	if err != nil {
		return err
	}

	if len(images) == 0 {
		return errors.New("No valid input images.")
	}

	outputImage := concatImages(images)

	file, _ := os.Create(outputFilename)
	defer file.Close()

	if err := jpeg.Encode(file, outputImage, &jpeg.Options{100}); err != nil {
		return err
	}

	return nil
}

func concatImages(images []image.Image) image.Image {
	canvasWidth := 0
	canvasHeight := 0

	minHeight := 0
	if len(images) > 0 {
		minHeight = images[0].Bounds().Size().Y
	}
	for _, srcImg := range images {
		height := srcImg.Bounds().Size().Y
		if height < minHeight {
			minHeight = height
		}
	}

	fmt.Println("Min height: ", minHeight)

	canvasHeight = minHeight

	for _, srcImg := range images {
		scale := float64(minHeight) / float64(srcImg.Bounds().Size().Y)
		canvasWidth += int(float64(srcImg.Bounds().Size().X) * scale)

		fmt.Println(srcImg.Bounds().Size())
	}

	fmt.Println("Canvas size: ", canvasWidth, canvasHeight)

	outputImage := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))

	x := 0
	for _, srcImg := range images {
		scale := float64(minHeight) / float64(srcImg.Bounds().Size().Y)
		scaledWidth := int(float64(srcImg.Bounds().Size().X) * scale)

		targetRect := image.Rect(x, 0, scaledWidth+x, canvasHeight)
		draw.BiLinear.Scale(outputImage, targetRect, srcImg, srcImg.Bounds(), draw.Over, nil)

		x += scaledWidth
	}
	return outputImage
}

func loadImages(filenames []string) ([]image.Image, error) {
	images := []image.Image{}
	for _, filename := range filenames {
		fmt.Println(filename)

		src, _ := os.Open(filename)
		defer src.Close()

		srcImg, _, err := image.Decode(src)
		if err != nil {
			return nil, err
		}

		images = append(images, srcImg)
	}
	return images, nil
}
