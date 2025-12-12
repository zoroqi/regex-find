package app

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App holds the tview application and its components.
type App struct {
	app             *tview.Application
	regexInput      *tview.InputField
	textArea        *tview.TextArea
	highlightedView *tview.TextView
	matchView       *tview.TextView
	statsView       *tview.TextView
	helpHintView    *tview.TextView
	flex            *tview.Flex
	pages           *tview.Pages
	focusables      []tview.Primitive
}

// New creates and initializes a new TUI application.
func New(initialText string) *App {
	a := &App{
		app:             tview.NewApplication(),
		regexInput:      tview.NewInputField(),
		textArea:        tview.NewTextArea(),
		highlightedView: tview.NewTextView(),
		matchView:       tview.NewTextView(),
		statsView:       tview.NewTextView(),
		helpHintView:    tview.NewTextView(),
		pages:           tview.NewPages(),
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
	a.statsView.SetTextAlign(tview.AlignRight)
	a.helpHintView.SetText("F1 Help")

	// Create a status bar
	statusBar := tview.NewFlex().
		AddItem(a.helpHintView, 10, 1, false).
		AddItem(tview.NewBox(), 0, 1, false). // Spacer
		AddItem(a.statsView, 20, 1, false)

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
	a.pages.AddPage("main", a.flex, true, true)
	a.pages.AddPage("help", helpModal, true, false)
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
			a.statsView.SetText("Matches: 0")
			return
		}
	}

	// If no regex, just show plain text and clear matches
	if re == nil {
		a.highlightedView.SetText(tview.Escape(text))
		a.matchView.SetText("")
		a.statsView.SetText("Matches: 0")
		return
	}

	matchesForHighlight := re.FindAllStringIndex(text, -1)
	a.updateHighlightedView(text, matchesForHighlight)

	matchesForView := re.FindAllStringSubmatch(text, -1)
	a.updateMatchView(matchesForView)
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
	a.statsView.SetText(fmt.Sprintf("Matches: %d", len(matches)))
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
	if err := a.app.SetRoot(a.pages, true).SetFocus(a.regexInput).Run(); err != nil {
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
