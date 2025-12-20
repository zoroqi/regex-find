package app

import (
	"fmt"
	"os"

	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

// updateHighlight performs the regex matching and updates the views.
func (a *App) updateHighlight() {
	regexStr := a.regexInput.GetText()
	text := a.textArea.GetText()

	// Reset match data
	a.matches = nil
	a.matchIndices = nil
	a.highlightedMatchLines = nil
	a.matchViewLines = nil
	a.currentMatchIndex = -1

	_, indices, matches, err := Search(regexStr, text)

	if err != nil {
		a.highlightedView.SetText(tview.Escape(text) + "\n[red]Invalid Regular Expression")
		a.matchView.SetText("")
		return
	}

	// If no regex, just show plain text and clear matches
	if regexStr == "" {
		a.highlightedView.SetText(tview.Escape(text))
		a.matchView.SetText("")
		return
	}

	a.matchIndices = indices
	a.updateHighlightedView(text, a.matchIndices)

	a.matches = matches
	if len(a.matches) > 0 {
		a.lastMatch = a.matches[0][0]
	} else {
		a.lastMatch = ""
	}
	a.updateMatchView(a.matches)
}

func (a *App) handleExport() {
	a.modalPages.RemovePage(ExportPage)

	// Get form data
	formatIndex, _ := a.exportForm.GetFormItemByLabel(LabelExportFormat).(*tview.DropDown).GetCurrentOption()
	destIndex, _ := a.exportForm.GetFormItemByLabel(LabelOutputTarget).(*tview.DropDown).GetCurrentOption()
	groupInput := a.exportForm.GetFormItemByLabel(LabelGroupNumbers).(*tview.InputField).GetText()
	customFormatInput := a.exportForm.GetFormItemByLabel(LabelCustomFormat).(*tview.InputField).GetText()
	filePathInput := a.exportForm.GetFormItemByLabel(LabelFilePath).(*tview.InputField).GetText()

	var outputData []byte
	var err error

	switch formatIndex {
	case 0: // JSON (all content)
		outputData, err = GenerateExportJSONAll(a.GetRegexInput(), a.matches)
	case 1: // JSON (specific groups)
		outputData, err = GenerateExportJSONGroups(a.GetRegexInput(), a.matches, groupInput)
	case 2: // Custom format
		var strData string
		strData, err = GenerateExportCustom(a.matches, customFormatInput)
		outputData = []byte(strData)
	}

	if err != nil {
		a.showResultModal(fmt.Sprintf("Error generating data: %v", err), true)
		return
	}

	switch destIndex {
	case 0: // Save to clipboard
		err = a.saveToClipboard(outputData)
	case 1: // Save to file
		err = a.saveToFile(outputData, filePathInput)
	}

	if err != nil {
		a.showResultModal(fmt.Sprintf("Error saving data: %v", err), true)
		return
	}

	a.updateHistory()
}

func (a *App) saveToClipboard(data []byte) error {
	if err := clipboard.Init(); err != nil {
		return fmt.Errorf("failed to initialize clipboard: %v", err)
	}
	clipboard.Write(clipboard.FmtText, data)
	return nil
}

func (a *App) saveToFile(data []byte, path string) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}
	return os.WriteFile(path, data, 0644)
}
