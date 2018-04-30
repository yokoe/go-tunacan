package main

import (
	"fmt"
	"image"
	"net/http"
	"os"

	"image/jpeg"
	_ "image/png"
)

const numMaxFiles = 5

func loadImage(r *http.Request, paramName string) (image.Image, error) {
	file, _, err := r.FormFile(paramName)
	if err != nil {
		if err != http.ErrMissingFile {
			return nil, err
		} else {
			return nil, nil
		}
	}
	defer file.Close()

	srcImage, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return srcImage, nil
}

func concatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("parse files")

	images := []image.Image{}
	for i := 0; i < numMaxFiles; i++ {
		paramName := fmt.Sprintf("file%d", i)
		image, err := loadImage(r, paramName)

		if image != nil {
			images = append(images, image)
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid image at %s", paramName), http.StatusInternalServerError)
			return
		}
	}

	if len(images) == 0 {
		http.Error(w, "No input images.", http.StatusInternalServerError)
		return
	}

	fmt.Printf("%d files\n", len(images))

	outputImage := concatImages(images)

	fmt.Println("write to file")

	f, err := os.Create("/tmp/tunacan.jpg")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if err := jpeg.Encode(f, outputImage, &jpeg.Options{100}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Done.")
}

func launchServer() {
	fmt.Println("Server mode")
	http.HandleFunc("/concat", concatHandler)
	http.ListenAndServe(":8080", nil)
}
