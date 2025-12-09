package app

import (
	"regexp"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App holds the tview application and its components.
type App struct {
	app        *tview.Application
	regexInput *tview.InputField
	textArea   *tview.TextArea
	resultView *tview.TextView
	flex       *tview.Flex
}

// New creates and initializes a new TUI application.
func New() *App {
	a := &App{
		app:        tview.NewApplication(),
		regexInput: tview.NewInputField(),
		textArea:   tview.NewTextArea(),
		resultView: tview.NewTextView(),
	}

	a.setupUI()
	a.setupEventHandlers()

	return a
}

// setupUI configures the layout and appearance of the UI components.
func (a *App) setupUI() {
	// Configure Regex Input Field
	a.regexInput.SetLabel("Regex: ")
	a.regexInput.SetBorder(true)
	a.regexInput.SetTitle("Regular Expression")

	// Configure Text Area for input
	a.textArea.SetBorder(true)
	a.textArea.SetTitle("Text Input")

	// Configure Text View for output
	a.resultView.SetBorder(true)
	a.resultView.SetTitle("Result (Highlighted)")
	a.resultView.SetDynamicColors(true)

	// Configure Flex Layout
	a.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.regexInput, 3, 1, true).
		AddItem(a.textArea, 0, 10, true).
		AddItem(a.resultView, 0, 10, false)
}

// setupEventHandlers sets up all input handling.
func (a *App) setupEventHandlers() {
	// When regex text changes, trigger an update.
	a.regexInput.SetChangedFunc(func(text string) {
		a.updateHighlight()
	})

	// When the text area text changes, trigger an update.
	a.textArea.SetChangedFunc(func() {
		a.updateHighlight()
	})

	// App-level key captures for quitting and focus cycling.
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC, tcell.KeyCtrlD:
			a.app.Stop()
			return nil
		case tcell.KeyTab:
			a.cycleFocus(false)
			return nil
		case tcell.KeyBacktab:
			a.cycleFocus(true)
			return nil
		}
		return event
	})
}

// cycleFocus switches focus between the two input widgets.
func (a *App) cycleFocus(reverse bool) {
	widgets := []tview.Primitive{a.regexInput, a.textArea}
	for i, widget := range widgets {
		if widget.HasFocus() {
			nextIndex := (i + 1) % len(widgets)
			if reverse {
				nextIndex = (i - 1 + len(widgets)) % len(widgets)
			}
			a.app.SetFocus(widgets[nextIndex])
			return
		}
	}
}

// updateHighlight performs the regex matching and updates the result view's content.
func (a *App) updateHighlight() {
	regexStr := a.regexInput.GetText()
	text := a.textArea.GetText()

	if text == "" {
		a.resultView.SetText("")
		return
	}

	var re *regexp.Regexp
	var err error
	if regexStr != "" {
		re, err = regexp.Compile(regexStr)
		if err != nil {
			a.resultView.SetText(tview.Escape(text) + "\n[red]Invalid Regular Expression")
			return
		}
	}

	if re == nil {
		a.resultView.SetText(tview.Escape(text))
		return
	}

	matches := re.FindAllStringIndex(text, -1)
	if matches == nil {
		a.resultView.SetText(tview.Escape(text))
		return
	}

	colors := []string{"[white:green]", "[white:blue]"}
	var builder strings.Builder
	lastIndex := 0

	for i, match := range matches {
		start, end := match[0], match[1]
		color := colors[i%len(colors)]

		// Append text before the match
		builder.WriteString(tview.Escape(text[lastIndex:start]))
		// Append the match with color
		builder.WriteString(color)
		builder.WriteString(tview.Escape(text[start:end]))
		builder.WriteString("[:-]")

		lastIndex = end
	}

	// Append any remaining text after the last match
	builder.WriteString(tview.Escape(text[lastIndex:]))

	a.resultView.SetText(builder.String())
}

// Run starts the tview application.
func (a *App) Run() error {
	if err := a.app.SetRoot(a.flex, true).SetFocus(a.regexInput).Run(); err != nil {
		a.app.Stop()
		return err
	}
	return nil
}
