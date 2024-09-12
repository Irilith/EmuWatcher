package imghandler

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
)

// since the image package does not have a built-in function
// to convert an image to grayscale i have to do it manually
// by looping through each pixel in the image and calculating
// the luminance of the pixel and setting the pixel to the new gray color.
// *
// This is used to convert the image to grayscale before passing it to the OCR so its can focus on the text (i think)
func ToGrayScale(imgBuf *bytes.Buffer) (*bytes.Buffer, error) {
	img, err := png.Decode(imgBuf)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			oldPixel := img.At(x, y)
			r, g, b, _ := oldPixel.RGBA()

			luminance := 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
			grayColor := uint8(luminance / 256)
			gray.Set(x, y, color.Gray{Y: grayColor})
		}
	}

	var grayBuf bytes.Buffer
	err = png.Encode(&grayBuf, gray)
	if err != nil {
		return nil, err
	}
	return &grayBuf, nil
}

func DetectColorRange(imgBuf *bytes.Buffer) (float64, error) {
	img, err := png.Decode(imgBuf)
	if err != nil {
		return 0, fmt.Errorf("failed to decode image: %w", err)
	}

	// 30, 31, 33, 1 | HOME BACKGROUND COLOR
	// 55 59 60 1 | LOADING SCREEN COLOR
	targetColor := color.RGBA{R: 55, G: 59, B: 60, A: 1}
	// Set the tolerance range for the color matching
	torlerance := 30

	colorCount := 0
	totalPixels := img.Bounds()
	for x := 0; x < totalPixels.Dx(); x++ {
		for y := 0; y < totalPixels.Dy(); y++ {
			pixel := img.At(x, y)
			r, g, b, _ := pixel.RGBA()
			if withinTolerance(r, g, b, targetColor, torlerance) {
				colorCount++
			}
		}
	}
	// Calculate the percentage of pixels that match the target color
	// colorPercentage = (colorCount / (totalPixelsWidth * totalPixelsHeight)) * 100
	colorPercentage := float64(colorCount) / float64(totalPixels.Dx()*totalPixels.Dy()) * 100
	return colorPercentage, nil
}

func withinTolerance(r, g, b uint32, targetColor color.RGBA, torlerance int) bool {
	return abs(int(r>>8)-int(targetColor.R)) <= torlerance &&
		abs(int(g>>8)-int(targetColor.G)) <= torlerance &&
		abs(int(b>>8)-int(targetColor.B)) <= torlerance
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
