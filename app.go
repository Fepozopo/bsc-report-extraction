package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// openFileWindow creates a file open dialog and calls the given callback function with the selected file.
// If the user cancels the dialog, the error argument will be set to an error with message "cancelled".
func openFileWindow(parent fyne.Window, callback func(r fyne.URIReadCloser, e error)) {
	dialog.NewFileOpen(func(r fyne.URIReadCloser, e error) {
		callback(r, e)
	}, parent).Show()
}

// selectFiles prompts the user to select a commission report and returns the path to the selected file.
func selectFiles(a fyne.App) (string, string) {
	window := a.NewWindow("Commission Report")
	window.SetContent(widget.NewLabel("Please select the commission report:"))
	window.Resize(fyne.NewSize(900, 800))

	files := make([]*widget.Entry, 3)
	buttons := make([]*widget.Button, 3)

	options := []string{"commission", "royalty"}
	list := widget.NewSelect(options, func(s string) {
	})

	for i := range files {
		files[i] = widget.NewEntry()
		buttons[i] = widget.NewButton("Browse", func(i int) func() {
			return func() {
				openFileWindow(window, func(r fyne.URIReadCloser, e error) {
					if e != nil {
						return
					}
					files[i].SetText(r.URI().Path())
				})
			}
		}(i))
	}

	var selection string
	var filePaths []string
	submitButton := widget.NewButton("Submit", func() {
		selection = list.Selected
		filePaths = make([]string, len(files))
		for i, entry := range files {
			filePaths[i] = entry.Text
		}
		window.Close()
	})

	window.SetContent(container.New(
		layout.NewVBoxLayout(),
		widget.NewLabelWithStyle("Which report type are you looking for?", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		list,
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Select the commission report:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		files[0],
		buttons[0],
		layout.NewSpacer(),
		submitButton,
	))

	window.ShowAndRun()

	window.SetCloseIntercept(func() {
		window.Close()
	})

	return selection, filePaths[0]
}
