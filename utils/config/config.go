package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tidwall/gjson"
)

type Config struct {
	PlaceId int64 `json:"PlaceId"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func GetConfig(configName string) (string, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	config := gjson.ParseBytes(data)
	return config.Get(configName).String(), nil
}

func addFormItems(form *tview.Form, json gjson.Result, prefix string) {
	for key, value := range json.Map() {
		fieldName := key
		if prefix != "" {
			fieldName = prefix + "/" + key
		}
		switch value.Type {
		case gjson.String:
			form.AddInputField(fieldName, value.String(), 20, nil, nil)
		case gjson.Number:
			form.AddInputField(fieldName, value.String(), 20, nil, nil)
		case gjson.True, gjson.False:
			form.AddCheckbox(fieldName, value.Bool(), nil)
		case gjson.JSON:
			if value.IsArray() {
				form.AddInputField(fieldName, value.String(), 20, nil, nil)
			} else {
				// Recursively add form items for nested JSON structures
				addFormItems(form, value, fieldName)
			}
		}
	}
}

func EditConfig() {
	fmt.Println("Edit Configurations")
	app := tview.NewApplication()
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	defer file.Close()

	configData, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	config := gjson.ParseBytes(configData)

	form := tview.NewForm()

	form.SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite).
		SetButtonBackgroundColor(tcell.ColorGreen).
		SetButtonTextColor(tcell.ColorBlack)

	addFormItems(form, config, "")

	form.AddButton("Save", func() {
		newConfig := make(map[string]interface{})
		for i := 0; i < form.GetFormItemCount(); i++ {
			item := form.GetFormItem(i)
			label := item.GetLabel()
			switch input := item.(type) {
			case *tview.InputField:
				if num, err := strconv.ParseFloat(input.GetText(), 64); err == nil {
					newConfig[label] = num
				} else {
					newConfig[label] = input.GetText()
				}
			case *tview.Checkbox:
				newConfig[label] = input.IsChecked()
			}
		}

		nestedConfig := map[string]interface{}{}
		for k, v := range newConfig {
			nestedConfig = setNestedField(nestedConfig, k, v)
		}

		newConfigData, err := json.MarshalIndent(nestedConfig, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal new config: %v", err)
		}
		err = os.WriteFile("config.json", newConfigData, 0644)
		if err != nil {
			log.Fatalf("Failed to write new config: %v", err)
		}
		app.Stop()
		RunMenuCallback()
	})
	form.AddButton("Quit", func() {
		app.Stop()
		RunMenuCallback()
	})

	form.SetBorder(true).SetTitle("Edit Configurations (Use Tab and Shift-Tab to navigate)").SetTitleAlign(tview.AlignCenter)
	if err := app.SetRoot(form, true).Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}

func setNestedField(data map[string]interface{}, field string, value interface{}) map[string]interface{} {
	keys := strings.Split(field, "/")
	currentMap := data

	for i, key := range keys {
		if i == len(keys)-1 {
			currentMap[key] = value
		} else {
			if _, exists := currentMap[key]; !exists {
				currentMap[key] = make(map[string]interface{})
			}
			currentMap = currentMap[key].(map[string]interface{})
		}
	}
	return data
}

// This one to prevent import cycle
var RunMenuCallback func()

func SetRunMenuCallback(callback func()) {
	RunMenuCallback = callback
}
