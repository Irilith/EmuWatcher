package ui

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// In = Inside the form
func InShowLabelModal(app *tview.Application, form *tview.Form, message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		// If the user presses Enter, hide the modal
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(form, true).SetFocus(form)
		})

	app.SetRoot(modal, false).SetFocus(modal)
}

func OutShowLabelModal(message string) {
	app := tview.NewApplication()
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.Stop()
		})

	if err := app.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
		os.Exit(1)
	}
}

func OutShowFatalErrorModal(message string, callback func()) {
	app := tview.NewApplication()
	fatalErrorColor := tcell.ColorRed
	backgroundColor := tcell.ColorBlack

	modal := tview.NewModal().
		SetText("[::b]Error\n" + message).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Ok" {
				app.Stop()
				callback()
			}
		})
	modal.SetBackgroundColor(backgroundColor)
	modal.SetTextColor(fatalErrorColor)

	if err := app.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
		os.Exit(1)
	}
}

func OutShowWarningModal(message string, callback func()) {
	app := tview.NewApplication()
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Confirm", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Confirm" {
				app.Stop()
				callback()
			} else {
				app.Stop()
			}
		})

	if err := app.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
		os.Exit(1)
	}
}

func OutShowTextboxModal(message string, callback func()) {
	app := tview.NewApplication()
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Confirm", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Confirm" {
				app.Stop()
				callback()
			} else {
				app.Stop()
			}
		})

	if err := app.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
		os.Exit(1)
	}
}

// This is used to create a TextBox struct with the given values
// And simplify the creation of the form
type SimpleTextBox struct {
	Label   string
	Require bool
}

// This is used to create a form with the given values
// And more customization to the input field such as you can add placeholder, field width, etc.
type TextBox struct {
	Label   string
	Input   *tview.InputField
	Require bool
}

type FormItem interface{}

func createTextBox(simpleBox []SimpleTextBox) []TextBox {
	var textBox []TextBox
	for _, tb := range simpleBox {
		textBox = append(textBox, TextBox{
			Label:   tb.Label,
			Input:   tview.NewInputField().SetLabel(tb.Label + ": ").SetFieldWidth(20),
			Require: tb.Require,
		})
	}
	return textBox
}

// NOTE: I mean, i just create this but i think it had no use anyway, but i will keep it here so i can use it in future
// or you can use it in your code if you want
func ShowForm(items []FormItem, callback func(values map[string]string)) {
	var textBox []TextBox

	// Check if the item is a SimpleTextBox or TextBox
	for _, i := range items {
		switch tb := i.(type) {
		case SimpleTextBox:
			textBox = append(textBox, createTextBox([]SimpleTextBox{tb})...)
		case TextBox:
			textBox = append(textBox, tb)
		}
	}

	app := tview.NewApplication()
	form := tview.NewForm()

	for _, tb := range textBox {
		form.AddFormItem(tb.Input)
	}

	form.AddButton("Submit", func() {
		values := make(map[string]string)
		valid := true
		// Check if the required field is empty
		for _, tb := range textBox {
			if tb.Require && tb.Input.GetText() == "" {
				tb.Input.SetFieldBackgroundColor(tcell.ColorRed)
				tb.Input.SetFieldTextColor(tcell.ColorWhite)
				valid = false
			} else {
				values[tb.Label] = tb.Input.GetText()
			}
		}
		if valid {
			app.Stop()
			callback(values)
		} else {
			InShowLabelModal(app, form, "Please fill all required fields")
		}
	})

	form.AddButton("Cancel", func() {
		app.Stop()
	})

	if err := app.SetRoot(form, true).SetFocus(form).Run(); err != nil {
		os.Exit(1)
	}
}
