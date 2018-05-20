package tunacan

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"image/jpeg"
	_ "image/png"

	"github.com/yokoe/tunacan"
)

type Response struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

var tmpDir string
var bucketName string

const numMaxFiles = 5
const imagesDirName = "images"

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

func concatImagesToTempFile(images []image.Image, filename string) (*string, error) {
	outputImage := tunacan.ConcatImages(images)

	outputFilepath := filepath.Join(tmpDir, filename)

	fmt.Println("write to file: " + outputFilepath)

	f, err := os.Create(outputFilepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := jpeg.Encode(f, outputImage, &jpeg.Options{Quality: 100}); err != nil {
		return nil, err
	}

	return &outputFilepath, nil
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

	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	filename := timestamp + ".jpg"
	tempFilepath, err := concatImagesToTempFile(images, filename)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if bucketName != "" {
		uploadedURL, err := uploadToCloudStorage(*tempFilepath, bucketName, filename, "image/jpeg")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(Response{Status: "ok", URL: *uploadedURL})
	} else {
		imageURL := "http://" + r.Host + "/" + imagesDirName + "/" + filename
		json.NewEncoder(w).Encode(Response{Status: "ok", URL: imageURL})
	}
}

func launchServer(port string, optBucketName string) {
	fmt.Println("Server listening at " + port)
	var err error

	tmpDir, err = ioutil.TempDir("", "tunacan")
	if err != nil {
		log.Fatal(err)
	}

	bucketName = optBucketName

	fmt.Println("Temp directory: " + tmpDir)

	http.HandleFunc("/concat", concatHandler)
	http.Handle("/"+imagesDirName+"/", http.StripPrefix("/"+imagesDirName+"/", http.FileServer(http.Dir(tmpDir))))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
