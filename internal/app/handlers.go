package app

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

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
			return
		}
	}

	// If no regex, just show plain text and clear matches
	if re == nil {
		a.highlightedView.SetText(tview.Escape(text))
		a.matchView.SetText("")
		return
	}

	matchesForHighlight := re.FindAllStringIndex(text, -1)
	a.updateHighlightedView(text, matchesForHighlight)

	a.matches = re.FindAllStringSubmatch(text, -1)
	a.updateMatchView(a.matches)
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
	a.matchView.SetTitle(fmt.Sprintf("Matches(%d)", len(matches)))
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

func (a *App) handleExport() {
	a.pages.HidePage("export")

	// Get form data
	formatIndex, _ := a.exportForm.GetFormItemByLabel("Export Format").(*tview.DropDown).GetCurrentOption()
	destIndex, _ := a.exportForm.GetFormItemByLabel("Export Destination").(*tview.DropDown).GetCurrentOption()
	groupInput := a.exportForm.GetFormItemByLabel("Group Numbers (comma-separated)").(*tview.InputField).GetText()
	customFormatInput := a.exportForm.GetFormItemByLabel("Custom Format String").(*tview.InputField).GetText()
	filePathInput := a.exportForm.GetFormItemByLabel("File Path").(*tview.InputField).GetText()

	var outputData []byte
	var err error

	switch formatIndex {
	case 0: // JSON (all content)
		outputData, err = a.generateExportJSONAll()
	case 1: // JSON (specific groups)
		outputData, err = a.generateExportJSONGroups(groupInput)
	case 2: // Custom format
		var strData string
		strData, err = a.generateExportCustom(customFormatInput)
		outputData = []byte(strData)
	}

	if err != nil {
		a.showResultModal(fmt.Sprintf("Error generating data: %v", err), true)
		return
	}

	switch destIndex {
	case 0: // Save to clipboard
		err = a.saveToClipboard(outputData)
	case 1: // Save to file
		err = a.saveToFile(outputData, filePathInput)
	}

	if err != nil {
		a.showResultModal(fmt.Sprintf("Error saving data: %v", err), true)
		return
	}

	a.showResultModal("Export successful!", false)
}

func (a *App) generateExportJSONAll() ([]byte, error) {
	var resultMatches []map[string]string
	for _, match := range a.matches {
		matchMap := make(map[string]string)
		for i, group := range match {
			matchMap[strconv.Itoa(i)] = group
		}
		resultMatches = append(resultMatches, matchMap)
	}

	data := map[string]interface{}{
		"regex":   a.GetRegexInput(),
		"matches": resultMatches,
	}
	return json.MarshalIndent(data, "", "  ")
}

func (a *App) generateExportJSONGroups(groupInput string) ([]byte, error) {
	if groupInput == "" {
		return nil, fmt.Errorf("group numbers cannot be empty")
	}
	groupStrs := strings.Split(groupInput, ",")
	var groups []int
	for _, s := range groupStrs {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		g, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid group number: %s", s)
		}
		groups = append(groups, g)
	}

	var processedMatches []map[string]string
	for _, match := range a.matches {
		processedMatch := make(map[string]string)
		for _, g := range groups {
			if g >= 0 && g < len(match) {
				processedMatch[strconv.Itoa(g)] = match[g]
			}
		}
		if len(processedMatch) > 0 {
			processedMatches = append(processedMatches, processedMatch)
		}
	}

	data := map[string]interface{}{
		"regex":   a.GetRegexInput(),
		"matches": processedMatches,
	}
	return json.MarshalIndent(data, "", "  ")
}

func (a *App) generateExportCustom(format string) (string, error) {
	if format == "" {
		return "", fmt.Errorf("custom format string cannot be empty")
	}

	// Replace common escape sequences
	processedFormat := strings.NewReplacer(`\n`, "\n", `\t`, "\t", `\r`, "\r").Replace(format)

	var result strings.Builder
	for i, match := range a.matches {
		if i > 0 {
			result.WriteString("\n")
		}
		line := processedFormat
		// TODO 使用替換方案, 存在一個漏洞, 如果在一號分組文本中包含了 "$2", 會導致錯誤替換.
		for j, group := range match {
			placeholder := fmt.Sprintf("$%d", j)
			line = strings.ReplaceAll(line, placeholder, group)
		}
		result.WriteString(line)
	}
	return result.String(), nil
}

func (a *App) saveToClipboard(data []byte) error {
	if err := clipboard.Init(); err != nil {
		return fmt.Errorf("failed to initialize clipboard: %v", err)
	}
	clipboard.Write(clipboard.FmtText, data)
	return nil
}

func (a *App) saveToFile(data []byte, path string) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}
	return os.WriteFile(path, data, 0644)
}
