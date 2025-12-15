package app

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

// App holds the tview application and its components.
type App struct {
	app             *tview.Application
	regexInput      *tview.InputField
	textArea        *tview.TextArea
	highlightedView *tview.TextView
	matchView       *tview.TextView
	helpHintView    *tview.TextView
	flex            *tview.Flex
	pages           *tview.Pages
	modalPages      *tview.Pages
	focusables      []tview.Primitive
	matches         [][]string // Store matches for export
	exportForm      *tview.Form
}

// New creates and initializes a new TUI application.
func New(initialText string) *App {
	a := &App{
		app:             tview.NewApplication(),
		regexInput:      tview.NewInputField(),
		textArea:        tview.NewTextArea(),
		highlightedView: tview.NewTextView(),
		matchView:       tview.NewTextView(),
		helpHintView:    tview.NewTextView(),
		pages:           tview.NewPages(),
		modalPages:      tview.NewPages(),
	}

	a.textArea.SetText(initialText, false)
	a.setupUI()
	a.setupEventHandlers()
	a.focusables = []tview.Primitive{a.regexInput, a.textArea, a.highlightedView, a.matchView}
	a.updateHighlight()

	return a
}

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

// setupEventHandlers sets up all input handling.
func (a *App) setupEventHandlers() {
	a.regexInput.SetChangedFunc(func(text string) {
		a.updateHighlight()
	})

	a.textArea.SetChangedFunc(func() {
		a.updateHighlight()
	})

	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if a.exportForm.HasFocus() {
			return event
		}

		// Global quit
		if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyCtrlD {
			a.app.Stop()
			return nil
		}

		// Show help modal
		if event.Key() == tcell.KeyF1 {
			a.pages.ShowPage("help")
			return nil
		}

		// Show export modal
		if event.Key() == tcell.KeyCtrlE {
			a.pages.ShowPage("export")
			a.app.SetFocus(a.exportForm)
			return nil
		}

		// Handle Tab and Backtab for focus cycling
		if event.Key() == tcell.KeyTab {
			a.cycleFocus(false)
			return nil
		}
		if event.Key() == tcell.KeyBacktab {
			a.cycleFocus(true)
			return nil
		}

		// Scrolling for Highlighted View
		if a.highlightedView.HasFocus() {
			return a.handleScrolling(a.highlightedView, event)
		}

		// Scrolling for Match View
		if a.matchView.HasFocus() {
			return a.handleScrolling(a.matchView, event)
		}

		return event
	})
}

// handleScrolling provides vim-like and arrow key scrolling for a TextView.
func (a *App) handleScrolling(tv *tview.TextView, event *tcell.EventKey) *tcell.EventKey {
	row, col := tv.GetScrollOffset()
	switch event.Key() {
	case tcell.KeyUp:
		tv.ScrollTo(row-1, col)
		return nil
	case tcell.KeyDown:
		tv.ScrollTo(row+1, col)
		return nil
	case tcell.KeyLeft:
		tv.ScrollTo(row, col-1)
		return nil
	case tcell.KeyRight:
		tv.ScrollTo(row, col+1)
		return nil
	}

	switch event.Rune() {
	case 'k':
		tv.ScrollTo(row-1, col)
		return nil
	case 'j':
		tv.ScrollTo(row+1, col)
		return nil
	case 'h':
		tv.ScrollTo(row, col-1)
		return nil
	case 'l':
		tv.ScrollTo(row, col+1)
		return nil
	}

	return event
}

// cycleFocus switches focus between the input widgets.
func (a *App) cycleFocus(reverse bool) {
	for i, widget := range a.focusables {
		if widget.HasFocus() {
			nextIndex := (i + 1) % len(a.focusables)
			if reverse {
				nextIndex = (i - 1 + len(a.focusables)) % len(a.focusables)
			}
			a.app.SetFocus(a.focusables[nextIndex])
			return
		}
	}
}

// updateHighlight performs the regex matching and updates the views.
func (a *App) updateHighlight() {
	regexStr := a.regexInput.GetText()
	text := a.textArea.GetText()

	// Compile regex
	var re *regexp.Regexp
	var err error
	if regexStr != "" {
		re, err = regexp.Compile(regexStr)
		if err != nil {
			a.highlightedView.SetText(tview.Escape(text) + "\n[red]Invalid Regular Expression")
			a.matchView.SetText("")
			return
		}
	}

	// If no regex, just show plain text and clear matches
	if re == nil {
		a.highlightedView.SetText(tview.Escape(text))
		a.matchView.SetText("")
		return
	}

	matchesForHighlight := re.FindAllStringIndex(text, -1)
	a.updateHighlightedView(text, matchesForHighlight)

	a.matches = re.FindAllStringSubmatch(text, -1)
	a.updateMatchView(a.matches)
}

func (a *App) updateHighlightedView(text string, matches [][]int) {
	colors := []string{"[white:green]", "[white:blue]"}
	var builder strings.Builder
	lastIndex := 0

	for i, match := range matches {
		start, end := match[0], match[1]
		color := colors[i%len(colors)]

		builder.WriteString(tview.Escape(text[lastIndex:start]))
		builder.WriteString(color)
		builder.WriteString(tview.Escape(text[start:end]))
		builder.WriteString("[:-]")

		lastIndex = end
	}
	builder.WriteString(tview.Escape(text[lastIndex:]))
	a.highlightedView.SetText(builder.String())
}

func (a *App) updateMatchView(matches [][]string) {
	a.matchView.SetTitle(fmt.Sprintf("Matches(%d)", len(matches)))
	if len(matches) == 0 {
		a.matchView.SetText("(No matches)")
		return
	}

	var builder strings.Builder
	const maxLen = 80 // Max length for a match line

	for i, match := range matches {
		// Full match
		fullMatchText := match[0]
		fullMatchText = strconv.Quote(fullMatchText)
		fullMatchText = fullMatchText[1 : len(fullMatchText)-1] // Remove quotes

		if len(fullMatchText) > maxLen {
			fullMatchText = fullMatchText[:maxLen/2-2] + " ... " + fullMatchText[len(fullMatchText)-(maxLen/2-2):]
		}
		builder.WriteString(fmt.Sprintf("%d: %s\n", i, fullMatchText))

		// Capture groups
		if len(match) > 1 {
			for j, group := range match[1:] {
				groupText := strconv.Quote(group)
				groupText = groupText[1 : len(groupText)-1] // Remove quotes

				if len(groupText) > maxLen-4 { // Adjust for indentation
					groupText = groupText[:(maxLen-4)/2-2] + " ... " + groupText[len(groupText)-((maxLen-4)/2-2):]
				}
				builder.WriteString(fmt.Sprintf("    %d: %s\n", j+1, groupText))
			}
		}

		// Add a blank line after each match block
		builder.WriteString("\n")
	}

	a.matchView.SetText(builder.String())
}

// Run starts the tview application.
func (a *App) Run() error {
	if err := a.app.SetRoot(a.modalPages, true).SetFocus(a.regexInput).Run(); err != nil {
		a.app.Stop()
		return err
	}
	return nil
}

func (a *App) createHelpModal() *tview.Modal {
	helpText := `[yellow]Scrolling Help:

- [green]Arrow Keys[white]: Scroll up, down, left, right.
- [green]h, j, k, l[white]: Vim-style scrolling (left, down, up, right).

This applies to both the 'Highlighted' and 'Matches' views when they are in focus.

Press ESC to dismiss.`

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

// GetRegexInput returns the current text in the regex input field.
func (a *App) GetRegexInput() string {
	return a.regexInput.GetText()
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

func (a *App) handleExport() {
	a.pages.HidePage("export")

	// Get form data
	formatIndex, _ := a.exportForm.GetFormItemByLabel("Export Format").(*tview.DropDown).GetCurrentOption()
	destIndex, _ := a.exportForm.GetFormItemByLabel("Export Destination").(*tview.DropDown).GetCurrentOption()
	groupInput := a.exportForm.GetFormItemByLabel("Group Numbers (comma-separated)").(*tview.InputField).GetText()
	customFormatInput := a.exportForm.GetFormItemByLabel("Custom Format String").(*tview.InputField).GetText()
	filePathInput := a.exportForm.GetFormItemByLabel("File Path").(*tview.InputField).GetText()

	var outputData []byte
	var err error

	switch formatIndex {
	case 0: // JSON (all content)
		outputData, err = a.generateExportJSONAll()
	case 1: // JSON (specific groups)
		outputData, err = a.generateExportJSONGroups(groupInput)
	case 2: // Custom format
		var strData string
		strData, err = a.generateExportCustom(customFormatInput)
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

	a.showResultModal("Export successful!", false)
}

func (a *App) generateExportJSONAll() ([]byte, error) {
	var resultMatches []map[string]string
	for _, match := range a.matches {
		matchMap := make(map[string]string)
		for i, group := range match {
			matchMap[strconv.Itoa(i)] = group
		}
		resultMatches = append(resultMatches, matchMap)
	}

	data := map[string]interface{}{
		"regex":   a.GetRegexInput(),
		"matches": resultMatches,
	}
	return json.MarshalIndent(data, "", "  ")
}

func (a *App) generateExportJSONGroups(groupInput string) ([]byte, error) {
	if groupInput == "" {
		return nil, fmt.Errorf("group numbers cannot be empty")
	}
	groupStrs := strings.Split(groupInput, ",")
	var groups []int
	for _, s := range groupStrs {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		g, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid group number: %s", s)
		}
		groups = append(groups, g)
	}

	var processedMatches []map[string]string
	for _, match := range a.matches {
		processedMatch := make(map[string]string)
		for _, g := range groups {
			if g >= 0 && g < len(match) {
				processedMatch[strconv.Itoa(g)] = match[g]
			}
		}
		if len(processedMatch) > 0 {
			processedMatches = append(processedMatches, processedMatch)
		}
	}

	data := map[string]interface{}{
		"regex":   a.GetRegexInput(),
		"matches": processedMatches,
	}
	return json.MarshalIndent(data, "", "  ")
}

func (a *App) generateExportCustom(format string) (string, error) {
	if format == "" {
		return "", fmt.Errorf("custom format string cannot be empty")
	}

	var result strings.Builder
	for i, match := range a.matches {
		if i > 0 {
			result.WriteString("\n")
		}
		line := format
		for j, group := range match {
			placeholder := fmt.Sprintf("$%d", j)
			line = strings.ReplaceAll(line, placeholder, group)
		}
		result.WriteString(line)
	}
	return result.String(), nil
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
