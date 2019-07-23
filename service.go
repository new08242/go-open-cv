package main

import (
	"net/http"
	"os"
	"fmt"
	"image"
	"errors"

	"go-open-cv/facedetection"
)

func GetFaceFromPictureHandler(w http.ResponseWriter, r *http.Request) {
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

	fd := facedetection.FaceDetect{}
	fp := fd.GetFace(imageDecode)
	http.ServeFile(w, r, fp)
}
