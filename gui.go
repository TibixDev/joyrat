package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateGui(cfg *Settings) (fyne.App, fyne.Window) {
	a := app.New()
	w := a.NewWindow("Joyrat")

	tempCfg := *cfg

	// Create input fields and populate them using reflection
	formItems := []struct {
		Label string
		Value interface{}
	}{
		{"Mouse Speed", &tempCfg.MOUSE_SPEED},
		{"Mouse Speed (Low)", &tempCfg.MOUSE_SPEED_LOW},
		{"Mouse Speed (High)", &tempCfg.MOUSE_SPEED_HIGH},
		{"Scroll Speed", &tempCfg.SCROLL_SPEED},
		{"Joystick Deadzone", &tempCfg.JOYSTICK_DEADZONE},
		{"Axis - Left Trigger", &tempCfg.AXIS_LT},
		{"Axis - Right Trigger", &tempCfg.AXIS_RT},
		{"Axis - Left Stick Horizontal", &tempCfg.AXIS_LS_X},
		{"Axis - Left Stick Vertical", &tempCfg.AXIS_LS_Y},
		{"Axis - Right Stick Horizontal", &tempCfg.AXIS_RS_X},
		{"Axis - Right Stick Vertical", &tempCfg.AXIS_RS_Y},
	}

	var formWidgets []*widget.FormItem
	for _, item := range formItems {
		entry := widget.NewEntry()
		switch v := item.Value.(type) {
		case *int:
			entry.SetText(strconv.Itoa(*v))
		case *int16:
			entry.SetText(strconv.Itoa(int(*v)))
		case *uint8:
			entry.SetText(strconv.Itoa(int(*v)))
		}
		itemCopy := item
		entry.OnChanged = func(text string) {
			switch v := itemCopy.Value.(type) {
			case *int:
				*v, _ = strconv.Atoi(text)
			case *int16:
				value, _ := strconv.Atoi(text)
				*v = int16(value)
			case *uint8:
				value, _ := strconv.Atoi(text)
				*v = uint8(value)
			}
		}
		formWidgets = append(formWidgets, widget.NewFormItem(item.Label, entry))
	}

	// Create Save button
	saveButton := widget.NewButton("Save", func() {
		CopyConfig(&tempCfg, cfg)
		speed = cfg.MOUSE_SPEED
		SaveCfg(cfg)
		fmt.Println("Settings saved")
	})

	// Create Reset to default button
	resetButton := widget.NewButton("Reset to Default", func() {
		CopyConfig(&DefaultCfg, cfg)
		CopyConfig(&DefaultCfg, &tempCfg)
		SaveCfg(cfg)
		speed = cfg.MOUSE_SPEED
		for i, item := range formItems {
			entry := formWidgets[i].Widget.(*widget.Entry)
			switch v := item.Value.(type) {
			case *int:
				entry.SetText(strconv.Itoa(*v))
			case *int16:
				entry.SetText(strconv.Itoa(int(*v)))
			case *uint8:
				entry.SetText(strconv.Itoa(int(*v)))
			}
		}
		fmt.Println("Settings reset to default")
	})

	// Header
	header := widget.NewRichTextFromMarkdown("# Joyrat")
	settingsLabel := widget.NewRichTextFromMarkdown("## Settings")

	// Layout
	form := container.NewVBox(
		container.NewCenter(header),
		widget.NewSeparator(),
		settingsLabel,
		widget.NewForm(formWidgets...),
		widget.NewSeparator(),
		saveButton,
		resetButton,
	)

	w.SetContent(form)
	w.Resize(fyne.NewSize(400, 400))

	return a, w
}
