package main

import (
	"fmt"
	"image/color"
	"image"
	"os"
	"errors"

	"gocv.io/x/gocv"
)

//FIXME: refactor this pls

func main() {
	////////////FIXME: get image from request/////////////////
	imagePath := "/resource/face.jpg"
	// open image
	fileImage, err := os.Open(imagePath)
	if err != nil {
		panic(errors.New(fmt.Sprintf("open file image error: %s", err)))
	}
	defer fileImage.Close()

	imageDecode, _, err := image.Decode(fileImage)
	if err != nil {
		panic(errors.New(fmt.Sprintf("decode image error: %s", err)))
	}
	defer fileImage.Close()

	// prepare image matrix
	img, err := gocv.ImageToMatRGB(imageDecode)
	if err != nil {
		panic(errors.New(fmt.Sprintf("gocv ImageToMatRGB error: %s", err)))
	}
	defer img.Close()
	////////////FIXME: get image from request/////////////////

	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load("data/haarcascade_frontalface_default.xml") {
		fmt.Println("Error reading cascade file: data/haarcascade_frontalface_default.xml")
		return
	}

	// Start cut image
	// detect faces
	rects := classifier.DetectMultiScale(img)
	fmt.Printf("found %d faces\n", len(rects))

	// draw a rectangle around each face on the original image
	for _, r := range rects {
		gocv.Rectangle(&img, r, blue, 3)
	}

	// send img to some where

}