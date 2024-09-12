// send current screen cap and system information from EmuWatcher/pkg/sys/system.go
//

package webhook

import (
	"EmuWatcher/pkg/sys"
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/vova616/screenshot"
)

type Webhook struct {
	WebhookURL string
	Client     *http.Client
}

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
	Timestamp   string `json:"timestamp"`
}

// NewWebhook creates a new Webhook struct with a custom HTTP client
func NewWebhook(webhookURL string) *Webhook {
	return &Webhook{
		WebhookURL: webhookURL,
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// SendWebhook sends system information and a screenshot to Discord
func (w *Webhook) SendWebhook() error {
	systemInfo, err := sys.GetSystemInfo()
	if err != nil {
		return fmt.Errorf("failed to get system info: %w", err)
	}

	screen, err := captureScreen()
	if err != nil {
		return fmt.Errorf("failed to capture screen: %w", err)
	}

	embed := Embed{
		Title: "System Information",
		Description: fmt.Sprintf("**Core Count:** %d\n"+
			"**EmuWatcher Running Thread Count:** %d\n"+
			"**CPU Utilization:** %.2f%%\n"+
			"**Total RAM:** %d GB\n"+
			"**Available RAM:** %d GB\n"+
			"**Total Disk:** %d GB\n"+
			"**Free Disk:** %d GB\n",
			systemInfo.CoreCount, systemInfo.ThreadCount, systemInfo.CPUUtilization,
			systemInfo.TotalRAM, systemInfo.AvailableRAM, systemInfo.TotalDisk, systemInfo.FreeDisk,
		),
		Color:     0x00ff00, // Green color
		Timestamp: time.Now().Format(time.RFC3339),
	}

	embedData, err := json.Marshal(struct {
		Embeds []Embed `json:"embeds"`
	}{
		Embeds: []Embed{embed},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal embed data: %w", err)
	}

	var formData bytes.Buffer
	writer := multipart.NewWriter(&formData)

	if err := addFormField(writer, "payload_json", embedData); err != nil {
		return err
	}

	if err := addFormFile(writer, "file", "screenshot.png", screen); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", w.WebhookURL, &formData)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := w.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook request failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func addFormField(writer *multipart.Writer, fieldName string, fieldData []byte) error {
	part, err := writer.CreateFormField(fieldName)
	if err != nil {
		return fmt.Errorf("failed to create form field %s: %w", fieldName, err)
	}

	if _, err := part.Write(fieldData); err != nil {
		return fmt.Errorf("failed to write form field %s: %w", fieldName, err)
	}

	return nil
}

func addFormFile(writer *multipart.Writer, fieldName, fileName string, fileData *bytes.Buffer) error {
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return fmt.Errorf("failed to create form file %s: %w", fieldName, err)
	}

	if _, err := part.Write(fileData.Bytes()); err != nil {
		return fmt.Errorf("failed to write form file %s: %w", fieldName, err)
	}

	return nil
}

func captureScreen() (*bytes.Buffer, error) {
	screen, err := screenshot.CaptureScreen()
	if err != nil {
		return nil, fmt.Errorf("failed to capture screen: %w", err)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, screen); err != nil {
		return nil, fmt.Errorf("failed to encode screenshot: %w", err)
	}

	return &buf, nil
}
