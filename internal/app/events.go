package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
