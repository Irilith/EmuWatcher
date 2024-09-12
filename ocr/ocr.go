package ocr

import (
	"EmuWatcher/adb"
	"EmuWatcher/utils/imghandler"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

const (
	tesseract string = "./tools/ocr/tesseract"
	tessdata  string = "./tools/ocr/tessdata"
)

func CheckCrash(emuName string) (bool, error) {
	imageData, err := adb.CaptureScreen(emuName)
	if err != nil {
		fmt.Println("Error capturing screen:", err)
		return false, fmt.Errorf("Error capturing screen: %v", err)
	}

	var outputBuffer bytes.Buffer
	var errBuffer bytes.Buffer

	args := []string{
		"stdin",
		"stdout",
		"-l",
		"vie",
		"--tessdata-dir",
		"./tools/ocr/tessdata",
	}
	cmd := exec.Command(tesseract, args...)

	imgGray, err := imghandler.ToGrayScale(bytes.NewBuffer(imageData))
	if err != nil {
		fmt.Println("Error converting image to grayscale:", err)
		return false, fmt.Errorf("Error converting image to grayscale: %v", err)
	}

	cmd.Stdin = imgGray
	// Could use CombinedOutput()
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &errBuffer

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running Tesseract: %s", err)
		return false, fmt.Errorf("Error running Tesseract: %v", err)
	}

	text := outputBuffer.String()
	errString := []string{"Disconnected", "Mất kết nối", ": 278)", ": 277)", ": 264)", ": 524)", "Lỗi Khi Gia Nhập", ": 529"}

	return containsAny(text, errString), nil
}

func containsAny(text string, errString []string) bool {
	for _, substring := range errString {
		pattern := regexp.MustCompile(regexp.QuoteMeta(substring))
		if pattern.MatchString(text) {
			return true
		}
	}
	return false
}
