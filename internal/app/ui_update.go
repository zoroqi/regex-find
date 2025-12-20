package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rivo/tview"
)

func (a *App) updateHighlightedView(text string, matches [][]int) {
	a.highlightedMatchLines = make([]int, 0, len(matches))
	colors := []string{"[white:green]", "[white:blue]"}
	var builder strings.Builder
	lastIndex := 0

	for i, match := range matches {
		start, end := match[0], match[1]
		color := colors[i%len(colors)]

		// Calculate line number for the match
		// The number of newlines before the match start + 1
		lineNumber := strings.Count(text[:start], "\n")
		a.highlightedMatchLines = append(a.highlightedMatchLines, lineNumber)

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
	a.matchView.SetTitle(fmt.Sprintf(TitleMatchesFormat, len(matches)))
	if len(matches) == 0 {
		a.matchView.SetText("(No matches)")
		return
	}

	a.matchViewLines = make([]int, 0, len(matches))
	var builder strings.Builder
	const maxLen = 80 // Max length for a match line
	lineCounter := 0

	for i, match := range matches {
		a.matchViewLines = append(a.matchViewLines, lineCounter)

		// Full match
		fullMatchText := match[0]
		fullMatchText = strconv.Quote(fullMatchText)
		fullMatchText = fullMatchText[1 : len(fullMatchText)-1] // Remove quotes

		if len(fullMatchText) > maxLen {
			fullMatchText = fullMatchText[:maxLen/2-2] + " ... " + fullMatchText[len(fullMatchText)-(maxLen/2-2):]
		}
		line := fmt.Sprintf("%d: %s\n", i, fullMatchText)
		builder.WriteString(line)
		lineCounter += strings.Count(line, "\n")

		// Capture groups
		if len(match) > 1 {
			for j, group := range match[1:] {
				groupText := strconv.Quote(group)
				groupText = groupText[1 : len(groupText)-1] // Remove quotes

				if len(groupText) > maxLen-4 { // Adjust for indentation
					groupText = groupText[:(maxLen-4)/2-2] + " ... " + groupText[len(groupText)-((maxLen-4)/2-2):]
				}
				line := fmt.Sprintf("    %d: %s\n", j+1, groupText)
				builder.WriteString(line)
				lineCounter += strings.Count(line, "\n")
			}
		}

		// Add a blank line after each match block
		builder.WriteString("\n")
		lineCounter++
	}

	a.matchView.SetText(builder.String())
}
