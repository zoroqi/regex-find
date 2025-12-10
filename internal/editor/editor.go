package editor

import "strings"

// Editor manages the text content and cursor position for a text view.
type Editor struct {
	lines    [][]rune
	cursorX  int
	cursorY  int
	scrolled int // Top visible line
}

// New creates a new Editor.
func New() *Editor {
	return &Editor{
		lines: [][]rune{
			{},
		},
	}
}

// Text returns the entire content of the editor as a single string.
func (e *Editor) Text() string {
	var builder strings.Builder
	for i, line := range e.lines {
		builder.WriteString(string(line))
		if i < len(e.lines)-1 {
			builder.WriteRune('\n')
		}
	}
	return builder.String()
}

// Cursor returns the current cursor position.
func (e *Editor) Cursor() (int, int) {
	return e.cursorX, e.cursorY
}

// clampCursor ensures the cursor is within the valid bounds of the text.
func (e *Editor) clampCursor() {
	if e.cursorY < 0 {
		e.cursorY = 0
	}
	if e.cursorY >= len(e.lines) {
		e.cursorY = len(e.lines) - 1
	}
	if e.cursorX < 0 {
		e.cursorX = 0
	}
	if e.cursorX > len(e.lines[e.cursorY]) {
		e.cursorX = len(e.lines[e.cursorY])
	}
}

// MoveCursorUp moves the cursor up.
func (e *Editor) MoveCursorUp() {
	e.cursorY--
	e.clampCursor()
}

// MoveCursorDown moves the cursor down.
func (e *Editor) MoveCursorDown() {
	e.cursorY++
	e.clampCursor()
}

// MoveCursorLeft moves the cursor left.
func (e *Editor) MoveCursorLeft() {
	e.cursorX--
	if e.cursorX < 0 {
		if e.cursorY > 0 {
			e.cursorY--
			e.cursorX = len(e.lines[e.cursorY])
		} else {
			e.cursorX = 0
		}
	}
	e.clampCursor()
}

// MoveCursorRight moves the cursor right.
func (e *Editor) MoveCursorRight() {
	if e.cursorX < len(e.lines[e.cursorY]) {
		e.cursorX++
	} else if e.cursorY < len(e.lines)-1 {
		e.cursorY++
		e.cursorX = 0
	}
	e.clampCursor()
}

// InsertRune inserts a character at the cursor position.
func (e *Editor) InsertRune(r rune) {
	currentLine := e.lines[e.cursorY]
	// Grow the slice if necessary
	currentLine = append(currentLine, 0)
	// Shift characters to the right
	copy(currentLine[e.cursorX+1:], currentLine[e.cursorX:])
	// Insert the new rune
	currentLine[e.cursorX] = r
	e.lines[e.cursorY] = currentLine
	e.cursorX++
}

// InsertNewline inserts a newline at the cursor position.
func (e *Editor) InsertNewline() {
	currentLine := e.lines[e.cursorY]
	rest := make([]rune, len(currentLine[e.cursorX:]))
	copy(rest, currentLine[e.cursorX:])

	// Truncate the current line
	e.lines[e.cursorY] = currentLine[:e.cursorX]

	// Insert the new line
	e.lines = append(e.lines, nil) // Grow lines slice
	copy(e.lines[e.cursorY+2:], e.lines[e.cursorY+1:])
	e.lines[e.cursorY+1] = rest

	// Move cursor
	e.cursorY++
	e.cursorX = 0
}

// Backspace deletes the character before the cursor.
func (e *Editor) Backspace() {
	if e.cursorX == 0 && e.cursorY == 0 {
		return // Nothing to delete
	}

	if e.cursorX > 0 {
		// Delete character in the same line
		currentLine := e.lines[e.cursorY]
		copy(currentLine[e.cursorX-1:], currentLine[e.cursorX:])
		e.lines[e.cursorY] = currentLine[:len(currentLine)-1]
		e.cursorX--
	} else {
		// Join with the previous line
		prevLineLen := len(e.lines[e.cursorY-1])
		e.lines[e.cursorY-1] = append(e.lines[e.cursorY-1], e.lines[e.cursorY]...)
		copy(e.lines[e.cursorY:], e.lines[e.cursorY+1:])
		e.lines = e.lines[:len(e.lines)-1]
		e.cursorY--
		e.cursorX = prevLineLen
	}
}

// Delete deletes the character at the cursor.
func (e *Editor) Delete() {
	if e.cursorX == len(e.lines[e.cursorY]) && e.cursorY == len(e.lines)-1 {
		return // At the end of the text
	}

	if e.cursorX < len(e.lines[e.cursorY]) {
		// Delete character in the same line
		currentLine := e.lines[e.cursorY]
		copy(currentLine[e.cursorX:], currentLine[e.cursorX+1:])
		e.lines[e.cursorY] = currentLine[:len(currentLine)-1]
	} else {
		// Join with the next line
		e.lines[e.cursorY] = append(e.lines[e.cursorY], e.lines[e.cursorY+1]...)
		copy(e.lines[e.cursorY+1:], e.lines[e.cursorY+2:])
		e.lines = e.lines[:len(e.lines)-1]
	}
}
