package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// setupUI configures the layout and appearance of the UI components.
func (a *App) setupUI() {
	// Configure Regex Input Field
	a.regexInput.SetLabel(LabelRegex)
	a.regexInput.SetBorder(true)
	a.regexInput.SetTitle(TitleRegex)
	a.regexInput.SetFieldBackgroundColor(tcell.ColorDefault) // Remove background color

	// Configure Text Area
	a.textArea.SetBorder(true)
	a.textArea.SetTitle(TitleText)

	// Configure Highlighted View
	a.highlightedView.SetBorder(true)
	a.highlightedView.SetTitle(TitleHighlighted)
	a.highlightedView.SetDynamicColors(true)
	a.highlightedView.SetScrollable(true)

	// Configure Match View
	a.matchView.SetBorder(true)
	a.matchView.SetTitle(TitleMatches)
	a.matchView.SetScrollable(true)

	// Configure Status Bar components
	a.helpHintView.SetText(HintHelp)

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

	a.pages.AddPage(MainPage, a.flex, true, true)
	a.pages.AddPage(HelpPage, helpModal, true, false)
	a.pages.AddPage(ExportPage, exportPage, true, false)

	a.modalPages.AddPage(MainPage, a.pages, true, true)
}

func (a *App) showResultModal(message string, isError bool) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{ButtonOK}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.modalPages.RemovePage(ResultPage)
		})

	if isError {
		modal.SetTitle(TitleError)
		modal.SetBackgroundColor(tcell.ColorRed)
	} else {
		modal.SetTitle(TitleSuccess)
	}

	a.modalPages.AddPage(ResultPage, modal, true, true)
}

func (a *App) createHelpModal() *tview.Modal {
	helpText := HelpKeybindings + "\n\n" + HelpScrolling

	modal := tview.NewModal().SetText(helpText)
	modal.SetBorder(true).SetTitle(TitleHelp)
	modal.SetBackgroundColor(tcell.ColorDefault)
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			a.pages.HidePage(HelpPage)
			return nil
		}
		return event
	})
	return modal
}

func (a *App) createExportForm() *tview.Form {
	form := tview.NewForm().
		AddDropDown(LabelExportFormat, []string{OptJsonAll, OptJsonGroups, OptCustom}, 2, nil).
		AddInputField(LabelCustomFormat, "$1", 40, nil, nil).
		AddInputField(LabelGroupNumbers, "", 40, nil, nil).
		AddDropDown(LabelOutputTarget, []string{TargetClipboard, TargetFile}, 0, nil).
		AddInputField(LabelFilePath, "", 40, nil, nil).
		AddButton(ButtonExport, a.handleExport).
		AddButton(ButtonCancel, func() {
			a.pages.HidePage(ExportPage)
		})

	form.SetBorder(true).SetTitle(TitleExportOptions).SetTitleAlign(tview.AlignLeft)
	return form
}
