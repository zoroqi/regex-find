package app

import (
	"fmt"
	"github.com/rivo/tview"
)

// App holds the tview application and its components.
type App struct {
	app                   *tview.Application
	regexInput            *tview.InputField
	textArea              *tview.TextArea
	highlightedView       *tview.TextView
	matchView             *tview.TextView
	helpHintView          *tview.TextView
	flex                  *tview.Flex
	pages                 *tview.Pages
	modalPages            *tview.Pages
	helpView              *tview.TextView // For the help screen
	focusables            []tview.Primitive
	matches               [][]string // Store matches for export
	matchIndices          [][]int    // Store match indices for navigation
	highlightedMatchLines []int      // Store line numbers for each match in highlighted view
	matchViewLines        []int      // Store line numbers for each match in match view
	currentMatchIndex     int        // For navigating between matches

	// UI components for modal pages
	exportForm       *tview.Form
	historyView      *HistoryView // Instance of the history view
	keybindingsModal *tview.Modal
	exportPage       *tview.Flex
	historyPageFlex  *tview.Flex

	// History and Help state
	historyFilePath string
}

// New creates and initializes a new TUI application.
func New(initialText string, historyPath string) (*App, error) {
	history, err := LoadHistory(historyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load history: %w", err)
	}

	a := &App{
		app:               tview.NewApplication(),
		regexInput:        tview.NewInputField(),
		textArea:          tview.NewTextArea(),
		highlightedView:   tview.NewTextView(),
		matchView:         tview.NewTextView(),
		helpHintView:      tview.NewTextView(),
		pages:             tview.NewPages(),
		modalPages:        tview.NewPages(),
		helpView:          tview.NewTextView(),
		currentMatchIndex: -1, // No match selected initially
		historyFilePath:   historyPath,
	}

	a.textArea.SetText(initialText, false)
	a.setupUI()
	a.historyView.InitData(history.Patterns)
	a.setupEventHandlers()
	a.focusables = []tview.Primitive{a.regexInput, a.textArea, a.highlightedView, a.matchView}
	a.updateHighlight()

	return a, nil
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

// SaveHistory persists the current history to the file.
func (a *App) SaveHistory() error {
	return SaveHistory(a.historyFilePath, History{Patterns: a.historyView.GetItems()})
}

// UpdateHistoryWithCurrentRegex adds the current regex from the input field to the history.
func (a *App) UpdateHistoryWithCurrentRegex() {
	a.updateHistory()
}

// updateHistory adds the current regex to the history if it's new.
func (a *App) updateHistory() {
	pattern := a.GetRegexInput()
	if pattern == "" {
		return
	}
	a.historyView.AddItem(pattern, a.lastMatch())
}
