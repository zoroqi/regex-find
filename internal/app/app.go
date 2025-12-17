package app

import (
	"github.com/rivo/tview"
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
	focusables            []tview.Primitive
	matches               [][]string // Store matches for export
	matchIndices          [][]int    // Store match indices for navigation
	highlightedMatchLines []int      // Store line numbers for each match in highlighted view
	matchViewLines        []int      // Store line numbers for each match in match view
	currentMatchIndex     int        // For navigating between matches
	exportForm            *tview.Form
}

// New creates and initializes a new TUI application.
func New(initialText string) *App {
	a := &App{
		app:               tview.NewApplication(),
		regexInput:        tview.NewInputField(),
		textArea:          tview.NewTextArea(),
		highlightedView:   tview.NewTextView(),
		matchView:         tview.NewTextView(),
		helpHintView:      tview.NewTextView(),
		pages:             tview.NewPages(),
		modalPages:        tview.NewPages(),
		currentMatchIndex: -1, // No match selected initially
	}

	a.textArea.SetText(initialText, false)
	a.setupUI()
	a.setupEventHandlers()
	a.focusables = []tview.Primitive{a.regexInput, a.textArea, a.highlightedView, a.matchView}
	a.updateHighlight()

	return a
}

// Run starts the tview application.
func (a *App) Run() error {
	if err := a.app.SetRoot(a.modalPages, true).SetFocus(a.regexInput).Run(); err != nil {
		a.app.Stop()
		return err
	}
	return nil
}

// GetRegexInput returns the current text in the regex input field.
func (a *App) GetRegexInput() string {
	return a.regexInput.GetText()
}