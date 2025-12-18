package app

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
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
	a.helpHintView.SetText(HintHelp) // Updated hint text

	// Create a status bar
	statusBar := tview.NewFlex().
		AddItem(a.helpHintView, 100, 1, false) // Adjusted width for new hint

	// Configure Flex Layout for the main page
	bottomPane := tview.NewFlex().
		AddItem(a.highlightedView, 0, 1, false).
		AddItem(a.matchView, 0, 1, false)

	a.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.regexInput, 3, 1, true).
		AddItem(a.textArea, 0, 3, true).
		AddItem(bottomPane, 0, 2, false).
		AddItem(statusBar, 1, 0, false)

	// --- Create Pages ---

	// Main application page
	a.pages.AddPage(MainPage, a.flex, true, true)

	// F1 Help Page (Keybindings)
	keybindingsModal := a.createHelpModal()
	a.pages.AddPage(KeybindingsHelpPage, keybindingsModal, true, false)

	// F2 Help Page (Regex)
	a.setupHelpPage()
	a.pages.AddPage(RegexHelpPage, a.helpView, true, false)

	// Export Page
	a.exportForm = a.createExportForm()
	exportPage := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(a.exportForm, 80, 0, true).
			AddItem(nil, 0, 1, false), 0, 1, true).
		AddItem(nil, 0, 1, false)
	a.pages.AddPage(ExportPage, exportPage, true, false)

	// Modal pages holder (for popups over everything)
	a.modalPages.AddPage(MainPage, a.pages, true, true)
}

func (a *App) setupHelpPage() {
	var builder strings.Builder
	builder.WriteString("[yellow]Common Patterns\n")
	builder.WriteString("[yellow]---------------\n")
	for _, item := range RegexHelpData.Common {
		builder.WriteString(fmt.Sprintf("[green]%-20s [white]%s\n", item.Title, tview.Escape(item.Pattern)))
	}
	builder.WriteString("\n[yellow]Basic Escape Characters\n")
	builder.WriteString("[yellow]-----------------------\n")
	for _, item := range RegexHelpData.Escapes {
		builder.WriteString(fmt.Sprintf("[green]%-20s [white]%s\n", item.Title, tview.Escape(item.Pattern)))
	}

	a.helpView.SetText(builder.String())
	a.helpView.SetTitle("Regex Patterns Help (Press F2 or Esc to close)")
	a.helpView.SetBorder(true)
	a.helpView.SetDynamicColors(true)
	a.helpView.SetScrollable(true)
}

func (a *App) createHelpModal() *tview.Modal {
	helpText := HelpKeybindings + "\n\n" + HelpScrolling

	modal := tview.NewModal().SetText(helpText)
	modal.SetBorder(true).SetTitle(TitleHelp)
	modal.SetBackgroundColor(tcell.ColorDefault)
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyF1 {
			a.pages.HidePage(KeybindingsHelpPage)
			return nil
		}
		return event
	})
	return modal
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
