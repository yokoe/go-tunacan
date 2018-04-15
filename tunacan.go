package main

import (
	"flag"
	"fmt"
	"image"
	"os"
)

import _ "image/jpeg"

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

	for i := range sourceFilenames {
		fmt.Println("Source: ", sourceFilenames[i])
	}
	fmt.Println("Output: ", outputFilename)

	concat(sourceFilenames)
}

func concat(sourceFilenames []string) {
	for i := range sourceFilenames {
		filename := sourceFilenames[i]
		fmt.Println(filename)

		src, _ := os.Open(filename)
		defer src.Close()

		srcImg, _, err := image.Decode(src)
		if err != nil {
			panic(err)
		}

		fmt.Println(srcImg.Bounds().Size())
	}
}
