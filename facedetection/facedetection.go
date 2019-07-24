package facedetection

import (
	"image"
	"fmt"
	"os"
	"image/png"
	"path"
	"math"
	"errors"

	"gocv.io/x/gocv"
	"github.com/google/uuid"
	"github.com/oliamb/cutter"
	"github.com/nfnt/resize"
)

type FaceDetector interface {
	GetFace(imageData image.Image) (string, error)
}

type FaceDetect struct {}

func(fd FaceDetect) GetFace(imageData image.Image) (string, error) {
	dataModelPath := "/go/src/go-open-cv/resource/haarcascade_frontalface_default.xml"

	// process image before detect face
	imgResize, err := ProcessResizeImage(imageData, 80000, 40000)
	if err != nil {
		fmt.Errorf("ProcessResizeImage error: %s", err)
		return "", err
	}

	// prepare image matrix
	img, err := gocv.ImageToMatRGB(*imgResize)
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

	if !classifier.Load(dataModelPath) {
		fmt.Errorf("error reading cascade file: %s", dataModelPath)
		return "", err
	}

	// detect faces
	rects := classifier.DetectMultiScale(img)
	fmt.Printf("found %d faces\n", len(rects))

	rect := image.Rectangle{}
	// draw a rectangle around each face on the original image
	//for _, r := range rects {
		//gocv.Rectangle(&img, r, blue, 3)
	//}

	if len(rects) != 1 {
		return "", errors.New(fmt.Sprintf("not found face or found more than one face, len found: %d", len(rects)))
	}
	rect = rects[0]

	// TODO: Start cut image
	//imageResult, err := img.ToImage()
	//if err != nil {
	//	fmt.Errorf("error mat to image: %s", err)
	//	return "", err
	//}
	croppedImg, err := cutter.Crop(*imgResize, cutter.Config{
		Width: rect.Dx() + 40,
		Height: rect.Dy() + 40,
		Anchor: image.Point{
			X: rect.Min.X - 20,
			Y: rect.Min.Y - 20,
		},
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

func ProcessResizeImage(img image.Image, maxRes, minRes int) (*image.Image, error) {
	rect := img.Bounds()
	width := rect.Dx()
	height := rect.Dy()
	imgScale := float64(rect.Dx()) / float64(rect.Dy())
	resolution := width * height

	fmt.Println("image resolution:", resolution, "image width to height scale:", imgScale)

	if resolution <= maxRes && resolution > minRes {
		// ok resolution
		return &img, nil
	}

	matSrc, err := gocv.ImageToMatRGB(img)
	if err != nil {
		fmt.Errorf("gocv ImageToMatRGB error: %s", err)
		return nil, err
	}
	defer matSrc.Close()

	// not ok resolution calculate new size to resize
	widthReSize, heightReSize := width, height
	targetRes := 0
	if resolution > maxRes {
		targetRes = maxRes

	} else if resolution <= minRes {
		targetRes = minRes
	}

	heightReSize = int(math.Round(math.Sqrt(float64(targetRes)/ imgScale)))
	widthReSize = int(math.Round(imgScale * float64(heightReSize)))

	imgResize := resize.Resize(uint(widthReSize), uint(heightReSize), img, resize.Lanczos3)

	fmt.Println("new image height:", heightReSize, "width:", widthReSize)

	//gocv.Resize(matSrc, &matSrc, image.Point{}, scaleX, scaleY, gocv.InterpolationArea)

	//} else if resolution <= minRes { // not ok resolution too small
	//	gocv.Resize(matSrc, &matSrc, image.Point{}, 0, 0, gocv.InterpolationArea)
	//}

	//imgResize, err := matSrc.ToImage()
	//if err != nil {
	//	fmt.Errorf("mat to image error: %s", err)
	//	return nil, err
	//}

	resolution = imgResize.Bounds().Dx() * imgResize.Bounds().Dy()

	fmt.Println("image resolution after resize:", resolution)

	return &imgResize, nil
}