package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// setupEventHandlers sets up all input handling.
func (a *App) setupEventHandlers() {
	// --- Regex Input Field specific handlers ---
	a.regexInput.SetChangedFunc(func(text string) {
		// Reset history navigation on manual input
		a.historyIndex = -1
		a.updateHighlight()
	})

	a.regexInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			if len(a.history) > 0 {
				if a.historyIndex < len(a.history)-1 {
					a.historyIndex++
					a.regexInput.SetText(a.history[a.historyIndex].Regex)
				}
			}
			return nil
		case tcell.KeyDown:
			if a.historyIndex > 0 {
				a.historyIndex--
				a.regexInput.SetText(a.history[a.historyIndex].Regex)
			} else if a.historyIndex == 0 {
				a.historyIndex = -1
				a.regexInput.SetText("")
			}
			return nil
		case tcell.KeyEnter:
			a.updateHighlight()
			return nil
		}
		return event
	})

	a.textArea.SetChangedFunc(func() {
		a.updateHighlight()
	})

	// Set input capture for views that have special navigation
	a.highlightedView.SetInputCapture(a.handleViewNavigation)
	a.matchView.SetInputCapture(a.handleViewNavigation)

	// Set global input capture for app-wide events
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// If RegexHelpPage is visible, it gets priority for some keys
		if a.showHelp {
			if event.Key() == tcell.KeyF2 || event.Key() == tcell.KeyEsc {
				a.pages.HidePage(RegexHelpPage)
				a.showHelp = false
				return nil
			}
			// Let the help view's own capture handle scrolling
			a.helpView.InputHandler()(event, func(p tview.Primitive) {
				a.app.SetFocus(p)
			})
			return nil
		}

		// If a form or modal is active, let it handle its own events.
		// The F1 modal has its own input capture. The export form also captures input.
		if a.exportForm.HasFocus() {
			return event
		}

		switch event.Key() {
		case tcell.KeyCtrlC, tcell.KeyCtrlD:
			a.app.Stop()
			return nil
		case tcell.KeyF1:
			a.pages.ShowPage(KeybindingsHelpPage)
			return nil
		case tcell.KeyF2:
			a.pages.ShowPage(RegexHelpPage)
			a.showHelp = true
			return nil
		case tcell.KeyCtrlE:
			a.pages.ShowPage(ExportPage)
			a.app.SetFocus(a.exportForm)
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

// handleViewNavigation provides advanced navigation for TextViews.
func (a *App) handleViewNavigation(event *tcell.EventKey) *tcell.EventKey {
	var view *tview.TextView
	if a.highlightedView.HasFocus() {
		view = a.highlightedView
	} else if a.matchView.HasFocus() {
		view = a.matchView
	} else {
		return event // Should not happen if capture is set correctly
	}

	// Handle standard scrolling
	row, col := view.GetScrollOffset()
	switch event.Key() {
	case tcell.KeyUp:
		view.ScrollTo(row-1, col)
		return nil
	case tcell.KeyDown:
		view.ScrollTo(row+1, col)
		return nil
	case tcell.KeyLeft:
		view.ScrollTo(row, col-1)
		return nil
	case tcell.KeyRight:
		view.ScrollTo(row, col+1)
		return nil
	case tcell.KeyHome:
		view.ScrollToBeginning()
		return nil
	case tcell.KeyEnd:
		view.ScrollToEnd()
		return nil
	case tcell.KeyPgUp, tcell.KeyCtrlB:
		_, _, _, height := view.GetInnerRect()
		view.ScrollTo(row-height, col)
		return nil
	case tcell.KeyPgDn, tcell.KeyCtrlF:
		_, _, _, height := view.GetInnerRect()
		view.ScrollTo(row+height, col)
		return nil
	}

	// Handle custom navigation
	switch event.Rune() {
	case 'k':
		view.ScrollTo(row-1, col)
		return nil
	case 'j':
		view.ScrollTo(row+1, col)
		return nil
	case 'h':
		view.ScrollTo(row, col-1)
		return nil
	case 'l':
		view.ScrollTo(row, col+1)
		return nil
	case 'g':
		view.ScrollToBeginning()
		return nil
	case 'G':
		view.ScrollToEnd()
		return nil
	case 'n':
		a.navigateToMatch(1) // Next
		return nil
	case 'N':
		a.navigateToMatch(-1) // Previous
		return nil
	}

	return event
}

// navigateToMatch jumps to the next or previous match in the focused view.
func (a *App) navigateToMatch(direction int) {
	if len(a.matches) == 0 {
		return // No matches to navigate
	}

	// Calculate next index
	a.currentMatchIndex += direction
	if a.currentMatchIndex < 0 {
		a.currentMatchIndex = len(a.matches) - 1
	} else if a.currentMatchIndex >= len(a.matches) {
		a.currentMatchIndex = 0
	}

	// Get the correct line number based on focused view
	var line int
	var view *tview.TextView
	if a.highlightedView.HasFocus() {
		if a.currentMatchIndex < len(a.highlightedMatchLines) {
			line = a.highlightedMatchLines[a.currentMatchIndex]
		}
		view = a.highlightedView
	} else if a.matchView.HasFocus() {
		if a.currentMatchIndex < len(a.matchViewLines) {
			line = a.matchViewLines[a.currentMatchIndex]
		}
		view = a.matchView
	}

	if view != nil {
		view.ScrollTo(line, 0)
	}
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
