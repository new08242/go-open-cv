package facedetection

import (
	"image"
	"fmt"
	"os"
	"image/png"
	"path"

	"gocv.io/x/gocv"
	"github.com/google/uuid"
	"github.com/oliamb/cutter"
)

type FaceDetector interface {
	GetFace(imageData image.Image) (string, error)
}

type FaceDetect struct {}

func(fd FaceDetect) GetFace(imageData image.Image) (string, error) {
	dataPath := "/go/src/go-open-cv/resource/haarcascade_frontalface_default.xml"

	// process image before detect face
	imgResize, err := ProcessResizeImage(imageData, 50000, 40000)
	if err != nil {
		fmt.Errorf("ProcessResizeImage error: %s", err)
		return "", err
	}

	// prepare image matrix
	img, err := gocv.ImageToMatRGB(imgResize)
	if err != nil {
		fmt.Errorf("gocv ImageToMatRGB error: %s", err)
		return "", err
	}
	defer img.Close()

	// color for the rect when faces detected
	//blue := color.RGBA{0, 0, 255, 0}

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

	rect := image.Rectangle{}
	// draw a rectangle around each face on the original image
	for _, r := range rects {
		//gocv.Rectangle(&img, r, blue, 3)
		rect = r
	}

	// TODO: Start cut image
	//imageResult, err := img.ToImage()
	//if err != nil {
	//	fmt.Errorf("error mat to image: %s", err)
	//	return "", err
	//}
	croppedImg, err := cutter.Crop(imgResize, cutter.Config{
		Width: rect.Dx(),
		Height: rect.Dy(),
		Anchor: image.Point{rect.Min.X, rect.Min.Y},
		Mode:   cutter.TopLeft,
	})
	if err != nil {
		fmt.Errorf("crop image error: %s", err)
		return "", err
	}

	// generate image name
	uid := uuid.New()
	imgName := fmt.Sprintf("%s.png", uid)
	fileName := path.Join("/go/src/go-open-cv/resource/result/", imgName)

	fmt.Println("save image file at:", fileName)

	// create image
	imgPath, err := CreateImageFileWithPath(croppedImg, fileName)
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

func ProcessResizeImage(img image.Image, maxRes, minRes int) (image.Image, error) {
	rect := img.Bounds()
	width := rect.Dx()
	height := rect.Dy()
	resolution := width * height

	fmt.Println("image resolution:", resolution)

	if resolution <= maxRes && resolution > minRes {
		// ok resolution
		return img, nil
	}

	matSrc, err := gocv.ImageToMatRGB(img)
	if err != nil {
		fmt.Errorf("gocv ImageToMatRGB error: %s", err)
		return nil, err
	}
	defer matSrc.Close()

	// not ok resolution too big
	//if resolution > maxRes {
	gocv.Resize(matSrc, &matSrc, image.Point{}, 300, 150, gocv.InterpolationArea)

	//} else if resolution <= minRes { // not ok resolution too small
	//	gocv.Resize(matSrc, &matSrc, image.Point{}, 0, 0, gocv.InterpolationArea)
	//}

	imgResize, err := matSrc.ToImage()
	if err != nil {
		fmt.Errorf("mat to image error: %s", err)
		return nil, err
	}

	resolution = imgResize.Bounds().Dx() * imgResize.Bounds().Dy()

	fmt.Println("image resolution after resize:", resolution)

	return imgResize, nil
}