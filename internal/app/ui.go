package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// setupUI configures the layout and appearance of the UI components.
func (a *App) setupUI() {
	// Configure Regex Input Field
	a.regexInput.SetLabel("Regex: ")
	a.regexInput.SetBorder(true)
	a.regexInput.SetTitle("Regular Expression")
	a.regexInput.SetFieldBackgroundColor(tcell.ColorDefault) // Remove background color

	// Configure Text Area
	a.textArea.SetBorder(true)
	a.textArea.SetTitle("Text Input")

	// Configure Highlighted View
	a.highlightedView.SetBorder(true)
	a.highlightedView.SetTitle("Highlighted")
	a.highlightedView.SetDynamicColors(true)
	a.highlightedView.SetScrollable(true)

	// Configure Match View
	a.matchView.SetBorder(true)
	a.matchView.SetTitle("Matches")
	a.matchView.SetScrollable(true)

	// Configure Status Bar components
	a.helpHintView.SetText("F1 Help")

	// Create a status bar
	statusBar := tview.NewFlex().
		AddItem(a.helpHintView, 10, 1, false).
		AddItem(tview.NewBox(), 0, 1, false) // Spacer

	// Configure Flex Layout
	bottomPane := tview.NewFlex().
		AddItem(a.highlightedView, 0, 1, false).
		AddItem(a.matchView, 0, 1, false)

	a.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.regexInput, 3, 1, true).
		AddItem(a.textArea, 0, 3, true).
		AddItem(bottomPane, 0, 2, false).
		AddItem(statusBar, 1, 0, false)

	// Create and add pages
	helpModal := a.createHelpModal()
	a.exportForm = a.createExportForm()
	exportPage := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(a.exportForm, 80, 0, true).
			AddItem(nil, 0, 1, false), 0, 1, true).
		AddItem(nil, 0, 1, false)

	a.pages.AddPage("main", a.flex, true, true)
	a.pages.AddPage("help", helpModal, true, false)
	a.pages.AddPage("export", exportPage, true, false)

	a.modalPages.AddPage("main", a.pages, true, true)
}

func (a *App) showResultModal(message string, isError bool) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.modalPages.RemovePage("result")
		})

	if isError {
		modal.SetTitle("Error")
		modal.SetBackgroundColor(tcell.ColorRed)
	} else {
		modal.SetTitle("Success")
	}

	a.modalPages.AddPage("result", modal, true, true)
}

func (a *App) createHelpModal() *tview.Modal {
	helpText := `[yellow]KEYBINDINGS:

[green]F1[white]:           Show this help modal
[green]Ctrl+E[white]:       Show export options
[green]Tab / Shift+Tab[white]: Cycle focus between windows
[green]Ctrl+C / Ctrl+D[white]: Quit the application
[green]ESC[white]:          Close this modal

[yellow]SCROLLING (in 'Highlighted' and 'Matches' windows):

- [green]Arrow Keys[white]: Scroll up, down, left, right
- [green]h, j, k, l[white]:  Vim-style scrolling (left, down, up, right)`

	modal := tview.NewModal().SetText(helpText)
	modal.SetBorder(true).SetTitle("Help")
	modal.SetBackgroundColor(tcell.ColorDefault)
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			a.pages.HidePage("help")
			return nil
		}
		return event
	})
	return modal
}

func (a *App) createExportForm() *tview.Form {
	form := tview.NewForm().
		AddDropDown("Export Format", []string{"JSON (all content)", "JSON (specific groups)", "Custom format"}, 2, nil).
		AddInputField("Custom Format String", "$1", 40, nil, nil).
		AddInputField("Group Numbers (comma-separated)", "", 40, nil, nil).
		AddDropDown("Export Destination", []string{"Save to clipboard", "Save to file"}, 0, nil).
		AddInputField("File Path", "", 40, nil, nil).
		AddButton("Export", a.handleExport).
		AddButton("Cancel", func() {
			a.pages.HidePage("export")
		})

	form.SetBorder(true).SetTitle("Export Matches").SetTitleAlign(tview.AlignLeft)
	return form
}
