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
	textView   *tview.TextView
	flex       *tview.Flex
	rawText    *strings.Builder
}

// New creates and initializes a new TUI application.
func New() *App {
	a := &App{
		app:        tview.NewApplication(),
		regexInput: tview.NewInputField(),
		textView:   tview.NewTextView(),
		rawText:    &strings.Builder{},
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

	// Configure Combined Text View
	a.textView.SetBorder(true)
	a.textView.SetTitle("Text")
	a.textView.SetDynamicColors(true)

	// Configure Flex Layout. Both items are focusable.
	a.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.regexInput, 3, 1, true).
		AddItem(a.textView, 0, 10, true)
}

// setupEventHandlers sets up all input handling.
func (a *App) setupEventHandlers() {
	// A single, centralized input handler for the entire application.
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Route exit keys first.
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

		// If textView has focus, handle its input manually.
		if a.textView.HasFocus() {
			switch event.Key() {
			case tcell.KeyRune:
				a.rawText.WriteRune(event.Rune())
				a.updateHighlight()
				return nil // We've handled the event.
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if a.rawText.Len() > 0 {
					runes := []rune(a.rawText.String())
					if len(runes) > 0 {
						runes = runes[:len(runes)-1]
						a.rawText.Reset()
						a.rawText.WriteString(string(runes))
						a.updateHighlight()
					}
				}
				return nil // We've handled the event.
			}
		}

		// For all other cases, return the event to be handled by the default
		// handler of the focused widget. This is crucial for the InputField.
		return event
	})

	// When regex text changes, trigger an update.
	a.regexInput.SetChangedFunc(func(text string) {
		a.updateHighlight()
	})
}

// cycleFocus switches focus between the two input widgets.
func (a *App) cycleFocus(reverse bool) {
	widgets := []tview.Primitive{a.regexInput, a.textView}
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

// updateHighlight performs the regex matching and updates the text view's content.
func (a *App) updateHighlight() {
	regexStr := a.regexInput.GetText()
	text := a.rawText.String()

	escapedText := tview.Escape(text)

	// If there's no regex, just show the plain text.
	if regexStr == "" {
		a.textView.SetText(escapedText)
		return
	}

	re, err := regexp.Compile(regexStr)
	if err != nil {
		// Show an indicator for invalid regex.
		a.textView.SetText(escapedText + " [red](Invalid Regex)[-]")
		return
	}

	matches := re.FindAllStringIndex(text, -1)
	if matches == nil {
		// If no matches, show plain text.
		a.textView.SetText(escapedText)
		return
	}

	// Build the string with highlight tags.
	var builder strings.Builder
	lastIndex := 0
	for _, match := range matches {
		start, end := match[0], match[1]
		builder.WriteString(tview.Escape(text[lastIndex:start]))
		builder.WriteString("[yellow]")
		builder.WriteString(tview.Escape(text[start:end]))
		builder.WriteString("[-]")
		lastIndex = end
	}
	builder.WriteString(tview.Escape(text[lastIndex:]))

	a.textView.SetText(builder.String())
}

// Run starts the tview application.
func (a *App) Run() error {
	if err := a.app.SetRoot(a.flex, true).SetFocus(a.regexInput).Run(); err != nil {
		a.app.Stop()
		return err
	}
	return nil
}