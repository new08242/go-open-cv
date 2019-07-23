package facedetection

import (
	"image"
	"fmt"
	"errors"
	"image/color"

	"gocv.io/x/gocv"
)

type FaceDetector interface {
	GetFace(imageData image.Image) string
}

type FaceDetect struct {}

func(fd FaceDetect) GetFace(imageData image.Image) string {
	dataPath := "./vendor/data/haarcascade_frontalface_default.xml"
	// prepare image matrix
	img, err := gocv.ImageToMatRGB(imageData)
	if err != nil {
		panic(errors.New(fmt.Sprintf("gocv ImageToMatRGB error: %s", err)))
	}
	defer img.Close()

	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(dataPath) {
		fmt.Println("Error reading cascade file:", dataPath)
		return ""
	}

	// Start cut image
	// detect faces
	rects := classifier.DetectMultiScale(img)
	fmt.Printf("found %d faces\n", len(rects))

	// draw a rectangle around each face on the original image
	for _, r := range rects {
		gocv.Rectangle(&img, r, blue, 3)
	}

	//FIXME: save image in resource
	imgPath := "/resource/face.jpg"

	return imgPath
}