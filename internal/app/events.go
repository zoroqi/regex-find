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
		a.updateHighlight()
	})

	a.regexInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
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
		// If a modal page is currently displayed, don't allow main page shortcuts.
		// The modals have their own input handling (or it's handled globally here).
		if a.modalPages.HasPage(ExportPage) || a.modalPages.HasPage(HistoryPage) || a.modalPages.HasPage(RegexHelpPage) || a.modalPages.HasPage(KeybindingsHelpPage) {
			// Check for modal-closing keys
			switch event.Key() {
			case tcell.KeyEsc:
				if a.modalPages.HasPage(ExportPage) {
					a.modalPages.RemovePage(ExportPage)
				} else if a.modalPages.HasPage(HistoryPage) {
					a.modalPages.RemovePage(HistoryPage)
				} else if a.modalPages.HasPage(RegexHelpPage) {
					a.modalPages.RemovePage(RegexHelpPage)
				} else if a.modalPages.HasPage(KeybindingsHelpPage) {
					a.modalPages.RemovePage(KeybindingsHelpPage)
				}
				a.app.SetFocus(a.regexInput)
				return nil
			case tcell.KeyF1:
				if a.modalPages.HasPage(KeybindingsHelpPage) {
					a.modalPages.RemovePage(KeybindingsHelpPage)
					a.app.SetFocus(a.regexInput)
					return nil
				}
			case tcell.KeyF2:
				if a.modalPages.HasPage(RegexHelpPage) {
					a.modalPages.RemovePage(RegexHelpPage)
					a.app.SetFocus(a.regexInput)
					return nil
				}
			case tcell.KeyF3:
				if a.modalPages.HasPage(HistoryPage) {
					a.modalPages.RemovePage(HistoryPage)
					a.app.SetFocus(a.regexInput)
					return nil
				}
			}
			// If not a closing key, let the modal handle it
			return event
		}

		// If no modal page is active, handle global application shortcuts.
		switch event.Key() {
		case tcell.KeyCtrlC, tcell.KeyCtrlD:
			a.app.Stop()
			return nil
		case tcell.KeyF1: // Show Keybindings Help
			a.modalPages.AddPage(KeybindingsHelpPage, a.keybindingsModal, true, true)
			a.app.SetFocus(a.keybindingsModal)
			return nil
		case tcell.KeyF2: // Show Regex Help
			a.modalPages.AddPage(RegexHelpPage, a.helpView, true, true)
			a.app.SetFocus(a.helpView)
			return nil
		case tcell.KeyF3: // Show History Page
			a.modalPages.AddPage(HistoryPage, a.historyPageFlex, true, true)
			a.app.SetFocus(a.historyView) // Set focus to the history view
			return nil
		case tcell.KeyCtrlE: // Show Export Options
			a.modalPages.AddPage(ExportPage, a.exportPage, true, true)
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
