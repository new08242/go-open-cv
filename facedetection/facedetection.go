package facedetection

import (
	"image"
	"fmt"
	"image/color"
	"os"
	"image/png"
	"path"

	"gocv.io/x/gocv"
)

type FaceDetector interface {
	GetFace(imageData image.Image) (string, error)
}

type FaceDetect struct {}

func(fd FaceDetect) GetFace(imageData image.Image) (string, error) {
	dataPath := "/go/src/go-open-cv/resource/haarcascade_frontalface_default.xml"
	// prepare image matrix
	img, err := gocv.ImageToMatRGB(imageData)
	if err != nil {
		fmt.Errorf("gocv ImageToMatRGB error: %s", err)
		return "", err
	}
	defer img.Close()

	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(dataPath) {
		fmt.Errorf("error reading cascade file: %s", dataPath)
		return "", err
	}

	// detect faces
	rects := classifier.DetectMultiScale(img)
	fmt.Printf("found %d faces\n", len(rects))

	// draw a rectangle around each face on the original image
	for _, r := range rects {
		gocv.Rectangle(&img, r, blue, 3)
	}

	// TODO: Start cut image

	// generate image name
	//imgPath := "/go/src/go-open-cv/resource/face.jpg"
	imgName := "faceyFace.png"
	fileName := path.Join("/go/src/go-open-cv/resource/result/", imgName)

	fmt.Println("save image file at:", fileName)

	// create image
	imageResult, err := img.ToImage()
	if err != nil {
		fmt.Errorf("error mat to image: %s", err)
		return "", err
	}
	imgPath, err := CreateImageFileWithPath(imageResult, fileName)
	if err != nil {
		fmt.Errorf("CreateImageFileWithPath: %s", err)
		return "", err
	}

	return imgPath, nil
}

func CreateImageFileWithPath(img image.Image, fileName string) (string, error) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Errorf("[CreateImageFileWithPath] create file error: %s", err)
		return "", err
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		fmt.Errorf("[CreateImageFileWithPath] png encode error: %s", err)
		return "", err
	}

	return fileName, nil
}