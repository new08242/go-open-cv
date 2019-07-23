package main

import (
	"net/http"
	"fmt"
	"image"
	_ "image/png"
	_ "image/jpeg"
	_ "image/gif"
	"errors"

	"go-open-cv/facedetection"
)

func GetFaceFromPictureHandler(w http.ResponseWriter, r *http.Request) {
	imgFile, _, err := r.FormFile("image")
	if err != nil {
		panic(errors.New(fmt.Sprintf("get image from request error: %s", err)))
	}
	defer imgFile.Close()

	imageDecoded, _, err := image.Decode(imgFile)
	if err != nil {
		panic(errors.New(fmt.Sprintf("png decode image error: %s", err)))
	}

	fd := facedetection.FaceDetect{}
	fp, err := fd.GetFace(imageDecoded)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{ error: "%s" }`, err.Error())))
		return
	}
	http.ServeFile(w, r, fp)
}
