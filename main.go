package main

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"strings"
)

var (
	IMG_DIR_PATH   string = os.Getenv("IMG_DIR_PATH")
	THUMB_DIR_PATH string = os.Getenv("THUMB_DIR_PATH")
)

const (
	THUMB_WIDTH   = 300
	THUMB_HEIGHT  = 300
	THUMB_QUALITY = 95
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func main() {
	if IMG_DIR_PATH == "" || THUMB_DIR_PATH == "" {
		ErrorHandler(errors.New("img or thumb path is empty"))
	}

	err := filepath.Walk(IMG_DIR_PATH, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is an image
		if !info.IsDir() && isImageFile(info.Name()) {
			err := printImageSize(path)
			if err != nil {
				fmt.Printf("Error processing %s: %v\n", path, err)
			}
		}

		return nil
	})

	ErrorHandler(err)
}

// Check if a file is an image based on its extension
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	}
	return false
}

func printImageSize(filename string) error {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("could not decode image: %v", err)
	}

	// Get the size of the image
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	centralPoint := image.Point{(0 + width) / 2, (0 + height) / 2}

	cropSize := image.Rect(0, 0, THUMB_WIDTH, THUMB_HEIGHT).Add(centralPoint)
	cropImage := img.(SubImager).SubImage(cropSize)

	newCropFilePath := THUMB_DIR_PATH + "/thumb_" + filepath.Base(filename)
	croppedImageFile, err := os.Create(newCropFilePath)
	if err != nil {
		panic(err)
	}

	defer croppedImageFile.Close()

	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg":
		if err := jpeg.Encode(croppedImageFile, cropImage, &jpeg.Options{
			Quality: THUMB_QUALITY,
		}); err != nil {
			panic(err)
		}
	case ".jpeg":
		if err := jpeg.Encode(croppedImageFile, cropImage, &jpeg.Options{
			Quality: THUMB_QUALITY,
		}); err != nil {
			panic(err)
		}
	}

	return nil
}

func ErrorHandler(err error) {
	if err != nil {
		panic(err)
	}
}
