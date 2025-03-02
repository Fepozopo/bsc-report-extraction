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

// selectFiles creates a GUI window to select the product line to update and the paths to the hotsheet, stock report, and sales report files.
// It then returns the selection and the paths as strings.
func selectFiles(a fyne.App) string {
	window := a.NewWindow("Commission Report")
	window.SetContent(widget.NewLabel("Please select the commission report:"))
	window.Resize(fyne.NewSize(900, 800))

	files := make([]*widget.Entry, 3)
	buttons := make([]*widget.Button, 3)

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

	var filePaths []string
	submitButton := widget.NewButton("Submit", func() {
		filePaths = make([]string, len(files))
		for i, entry := range files {
			filePaths[i] = entry.Text
		}
		window.Close()
	})

	window.SetContent(container.New(
		layout.NewVBoxLayout(),
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

	return filePaths[0]
}
