package tunacan

import (
	"encoding/json"
	"flag"
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

type ServerCommand struct {
}

func (c *ServerCommand) Synopsis() string {
	return "Launch HTTP server."
}

func (c *ServerCommand) Help() string {
	return "Usage: tunacan server"
}

func (c *ServerCommand) Run(args []string) int {
	port := "8080"
	flags := flag.NewFlagSet("server", flag.ContinueOnError)
	flags.StringVar(&port, "p", "8080", "Port number to listen.")
	flags.Parse(args)

	launchServer(port)
	return 0
}

type Response struct {
	Status string `json:"status"`
	Url    string `json:"url"`
}

var tmpDir string

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

	outputImage := tunacan.ConcatImages(images)

	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	outputFilename := timestamp + ".jpg"
	outputFilepath := filepath.Join(tmpDir, outputFilename)

	fmt.Println("write to file: " + outputFilepath)

	f, err := os.Create(outputFilepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if err := jpeg.Encode(f, outputImage, &jpeg.Options{100}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imageUrl := "http://" + r.Host + "/" + imagesDirName + "/" + outputFilename

	json.NewEncoder(w).Encode(Response{Status: "ok", Url: imageUrl})
}

func launchServer(port string) {
	fmt.Println("Server listening at " + port)
	var err error

	tmpDir, err = ioutil.TempDir("", "tunacan")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Temp directory: " + tmpDir)

	http.HandleFunc("/concat", concatHandler)
	http.Handle("/"+imagesDirName+"/", http.StripPrefix("/"+imagesDirName+"/", http.FileServer(http.Dir(tmpDir))))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
