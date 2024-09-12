package roblox

import (
	"EmuWatcher/utils/ui"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var cookieFiles = "./data/cookies.txt"

type UserInfo struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type ErrorResponse struct {
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func GetCookies() ([]string, error) {
	var cookies []string
	file, err := os.Open(cookieFiles)
	if err != nil {
		return cookies, fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" { // Skip empty lines
			cookies = append(cookies, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return cookies, fmt.Errorf("Error reading file: %v", err)
	}

	if len(cookies) == 0 {
		return cookies, fmt.Errorf("No cookies found")
	}

	return cookies, nil
}

func ValidCookies(cookie string) (bool, UserInfo, error) {
	apiURL := "https://users.roblox.com/v1/users/authenticated"

	url, err := url.Parse(apiURL)
	if err != nil {
		return false, UserInfo{}, fmt.Errorf("Error parsing URL: %v", err)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return false, UserInfo{}, fmt.Errorf("Error creating cookie jar: %v", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return false, UserInfo{}, fmt.Errorf("Error creating request: %v", err)
	}
	cookie = strings.Trim(cookie, "\n")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,vi;q=0.8,ja;q=0.7")
	req.Header.Set("cookie", fmt.Sprintf(".ROBLOSECURITY=%s;", cookie))
	req.Header.Set("Sec-Ch-Ua", "\"Not)A;Brand\";v=\"99\", \"Microsoft Edge\";v=\"127\", \"Chromium\";v=\"127\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36 Edg/127.0.0.0")

	resp, err := client.Do(req)
	if err != nil {
		return false, UserInfo{}, fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, UserInfo{}, fmt.Errorf("Error reading response body: %v", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var userInfo UserInfo
		if err := json.Unmarshal(body, &userInfo); err != nil {
			return false, UserInfo{}, fmt.Errorf("Error unmarshalling response body: %v", err)
		}
		return true, userInfo, nil
	case http.StatusUnauthorized:
		var errorResponse ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return false, UserInfo{}, fmt.Errorf("Error unmarshalling error response body: %v", err)
		}
		if len(errorResponse.Errors) > 0 {
			return false, UserInfo{}, fmt.Errorf(errorResponse.Errors[0].Message)
		}
		return false, UserInfo{}, fmt.Errorf("Unauthorized access without error message")
	default:
		return false, UserInfo{}, fmt.Errorf("Unexpected response status code: %d", resp.StatusCode)
	}
}

// Promt a text editor using tview to add cookie to cookies.txt and save
func AddCookies() {
	app := tview.NewApplication()
	form := tview.NewForm()
	form.SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite).
		SetButtonBackgroundColor(tcell.ColorGreen).
		SetButtonTextColor(tcell.ColorBlack)
	textBox := tview.NewTextArea()
	textBox.SetLabel("Cookie")
	form.AddFormItem(textBox)
	form.AddButton("Save", func() {
		cookie := textBox.GetText()
		if cookie == "" {
			ui.InShowLabelModal(app, form, "Cookie cannot be empty!")
			return
		}
		// Save the cookie to the file cookies.txt
		file, err := os.OpenFile(cookieFiles, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		defer file.Close()
		file.WriteString(cookie + "\n")
		ui.InShowLabelModal(app, form, "Cookie saved successfully!")
		app.Stop()
		RunMenuCallback()
	})
	form.AddButton("Cancel", func() {
		app.Stop()
		RunMenuCallback()
	})
	if err := app.SetRoot(form, true).SetFocus(form).Run(); err != nil {
		panic(err)
	}
}

var RunMenuCallback func()

func SetRunMenuCallback(callback func()) {
	RunMenuCallback = callback
}
