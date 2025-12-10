package app

import (
	"github.com/zoroqi/regex-find/internal/editor"
	"regexp"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App holds the tview application and its components.
type App struct {
	app        *tview.Application
	regexInput *tview.InputField
	editorView *tview.TextView
	editor     *editor.Editor
	flex       *tview.Flex
}

// New creates and initializes a new TUI application.
func New() *App {
	a := &App{
		app:        tview.NewApplication(),
		regexInput: tview.NewInputField(),
		editorView: tview.NewTextView(),
		editor:     editor.New(),
	}

	a.setupUI()
	a.setupEventHandlers()
	a.updateHighlight(true) // Initial draw

	return a
}

// setupUI configures the layout and appearance of the UI components.
func (a *App) setupUI() {
	// Configure Regex Input Field
	a.regexInput.SetLabel("Regex: ")
	a.regexInput.SetBorder(true)
	a.regexInput.SetTitle("Regular Expression")

	// Configure Editor View
	a.editorView.SetBorder(true)
	a.editorView.SetTitle("Text (Editable)")
	a.editorView.SetDynamicColors(true)
	a.editorView.SetScrollable(true)

	// Configure Flex Layout
	a.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.regexInput, 3, 1, true).
		AddItem(a.editorView, 0, 1, true)
}

// setupEventHandlers sets up all input handling.
func (a *App) setupEventHandlers() {
	// When regex text changes, trigger an update.
	a.regexInput.SetChangedFunc(func(text string) {
		a.updateHighlight(false)
	})

	// App-level key captures for quitting and focus cycling.
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Global quit
		if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyCtrlD {
			a.app.Stop()
			return nil
		}

		// Handle editor input if it has focus
		if a.editorView.HasFocus() {
			switch event.Key() {
			case tcell.KeyRune:
				a.editor.InsertRune(event.Rune())
			case tcell.KeyEnter:
				a.editor.InsertNewline()
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				a.editor.Backspace()
			case tcell.KeyDelete:
				a.editor.Delete()
			case tcell.KeyUp:
				a.editor.MoveCursorUp()
			case tcell.KeyDown:
				a.editor.MoveCursorDown()
			case tcell.KeyLeft:
				a.editor.MoveCursorLeft()
			case tcell.KeyRight:
				a.editor.MoveCursorRight()
			case tcell.KeyTab:
				a.cycleFocus(false)
				return nil
			case tcell.KeyBacktab:
				// Let regexInput handle back-tab if it gets it
				a.cycleFocus(true)
				return nil
			default:
				return event // Pass other keys through
			}
			a.updateHighlight(true)
			return nil // We've handled the event
		}

		// Handle focus cycling from regexInput
		if a.regexInput.HasFocus() {
			switch event.Key() {
			case tcell.KeyTab:
				a.cycleFocus(false)
				return nil
			case tcell.KeyBacktab:
				a.cycleFocus(true)
				return nil
			}
		}

		return event
	})
}

// cycleFocus switches focus between the two input widgets.
func (a *App) cycleFocus(reverse bool) {
	widgets := []tview.Primitive{a.regexInput, a.editorView}
	for i, widget := range widgets {
		if widget.HasFocus() {
			nextIndex := (i + 1) % len(widgets)
			if reverse {
				nextIndex = (i - 1 + len(widgets)) % len(widgets)
			}
			a.app.SetFocus(widgets[nextIndex])
			a.updateHighlight(true) // Redraw to show/hide cursor
			return
		}
	}
}

// updateHighlight performs the regex matching and updates the editor view.
func (a *App) updateHighlight(isEditorUpdate bool) {
	regexStr := a.regexInput.GetText()
	text := a.editor.Text()

	// 1. Handle cursor insertion
	const cursorMarker = "‹‹CURSOR››"
	textWithCursor := text
	if a.editorView.HasFocus() {
		textWithCursor = a.insertCursorMarker(text, cursorMarker)
	}

	// 2. Handle regex compilation
	var re *regexp.Regexp
	var err error
	if regexStr != "" {
		re, err = regexp.Compile(regexStr)
		if err != nil {
			a.editorView.SetText(tview.Escape(text) + "\n[red]Invalid Regular Expression")
			return
		}
	}

	// 3. Get highlighted text
	highlightedText := a.getHighlightedText(re, textWithCursor, cursorMarker)
	a.editorView.SetText(highlightedText)
	if isEditorUpdate {
		a.editorView.ScrollToHighlight()
	}
}

func (a *App) insertCursorMarker(text, marker string) string {
	cursorX, cursorY := a.editor.Cursor()
	lines := strings.Split(text, "\n")

	// Clamp cursor to be safe
	if cursorY >= len(lines) {
		cursorY = len(lines) - 1
	}
	if cursorY < 0 {
		cursorY = 0
	}
	line := lines[cursorY]
	if cursorX > len(line) {
		cursorX = len(line)
	}
	if cursorX < 0 {
		cursorX = 0
	}

	// Insert marker into the specific line
	lines[cursorY] = line[:cursorX] + marker + line[cursorX:]

	return strings.Join(lines, "\n")
}

func (a *App) getHighlightedText(re *regexp.Regexp, text string, cursorMarker string) string {
	// Final text to be displayed
	var final_text string

	// If regex is nil, just escape the text
	if re == nil {
		final_text = tview.Escape(text)
	} else {
		// Otherwise, apply highlighting
		matches := re.FindAllStringIndex(text, -1)
		if matches == nil {
			final_text = tview.Escape(text)
		} else {
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
			final_text = builder.String()
		}
	}

	// Replace the cursor marker with actual tview tags
	if cursorMarker != "" && a.editorView.HasFocus() {
		// Use a placeholder for empty space under cursor
		final_text = strings.Replace(final_text, tview.Escape(cursorMarker), "[::r] [-::]", 1)
	}

	return final_text
}

// Run starts the tview application.
func (a *App) Run() error {
	if err := a.app.SetRoot(a.flex, true).SetFocus(a.regexInput).Run(); err != nil {
		a.app.Stop()
		return err
	}
	return nil
}
