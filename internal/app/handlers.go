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

	// Reset match data
	a.matches = nil
	a.matchIndices = nil
	a.highlightedMatchLines = nil
	a.matchViewLines = nil
	a.currentMatchIndex = -1

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

	a.matchIndices = re.FindAllStringIndex(text, -1)
	a.updateHighlightedView(text, a.matchIndices)

	a.matches = re.FindAllStringSubmatch(text, -1)
	a.updateMatchView(a.matches)
}

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

func (a *App) handleExport() {
	a.pages.HidePage("export")

	// Get form data
	formatIndex, _ := a.exportForm.GetFormItemByLabel(LabelExportFormat).(*tview.DropDown).GetCurrentOption()
	destIndex, _ := a.exportForm.GetFormItemByLabel(LabelOutputTarget).(*tview.DropDown).GetCurrentOption()
	groupInput := a.exportForm.GetFormItemByLabel(LabelGroupNumbers).(*tview.InputField).GetText()
	customFormatInput := a.exportForm.GetFormItemByLabel(LabelCustomFormat).(*tview.InputField).GetText()
	filePathInput := a.exportForm.GetFormItemByLabel(LabelFilePath).(*tview.InputField).GetText()

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

	regexGroupSeq := regexp.MustCompile(`\$(\d+)`)
	regexGroup := regexGroupSeq.FindAllStringSubmatch(format, -1)
	groupSeq := make([]int, len(regexGroup))
	for i, g := range regexGroup {
		num, err := strconv.Atoi(g[1])
		if err != nil {
			return "", fmt.Errorf("invalid group number in format: %s", g[0])
		}
		groupSeq[i] = num
	}

	// Replace common escape sequences
	processedFormat := regexGroupSeq.ReplaceAllString(strings.NewReplacer(`\n`, "\n", `\t`, "\t", `\r`, "\r").Replace(format), "%s")

	var result strings.Builder
	for i, match := range a.matches {
		if i > 0 {
			result.WriteString("\n")
		}
		args := make([]any, len(groupSeq))
		for j, g := range groupSeq {
			if g >= 0 && g < len(match) {
				args[j] = match[g]
			} else {
				// TODO 暫時沒有想好是使用空字符串還是原樣輸出 `$n`, 當前選擇後者.
				// Sublime Text 是使用空字符串.
				// 原樣輸出可以快速發現 format 寫錯了, 或正則寫錯了.
				args[j] = fmt.Sprintf("$%d", g)
			}
		}
		result.WriteString(fmt.Sprintf(processedFormat, args...))
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
