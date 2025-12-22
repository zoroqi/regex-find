package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type exportJson struct {
	Regex   string              `json:"regex"`
	Matches []map[string]string `json:"matches"`
}

// GenerateExportJSONAll generates a JSON byte slice containing the regex and all matches.
func GenerateExportJSONAll(regexStr string, matches [][]string) ([]byte, error) {
	var resultMatches []map[string]string
	for _, match := range matches {
		matchMap := make(map[string]string)
		for i, group := range match {
			matchMap[strconv.Itoa(i)] = group
		}
		resultMatches = append(resultMatches, matchMap)
	}

	data := exportJson{
		Regex:   regexStr,
		Matches: resultMatches,
	}
	return json.MarshalIndent(data, "", "  ")
}

// GenerateExportJSONGroups generates a JSON byte slice containing the regex and specific capture groups.
func GenerateExportJSONGroups(regexStr string, matches [][]string, groupInput string) ([]byte, error) {
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
	for _, match := range matches {
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

	data := exportJson{
		Regex:   regexStr,
		Matches: processedMatches,
	}
	return json.MarshalIndent(data, "", "  ")
}

// GenerateExportCustom generates a formatted string based on a custom format string (e.g., "$1 - $2").
func GenerateExportCustom(matches [][]string, format string) ([]byte, error) {
	if format == "" {
		return nil, fmt.Errorf("custom format string cannot be empty")
	}

	regexGroupSeq := regexp.MustCompile(`\$(\d+)`)
	regexGroup := regexGroupSeq.FindAllStringSubmatch(format, -1)
	groupSeq := make([]int, len(regexGroup))
	for i, g := range regexGroup {
		num, err := strconv.Atoi(g[1])
		if err != nil {
			return nil, fmt.Errorf("invalid group number in format: %s", g[0])
		}
		groupSeq[i] = num
	}

	// First, escape any literal '%' characters so Sprintf doesn't interpret them.
	tempFormat := strings.ReplaceAll(format, "%", "%%")
	// Next, handle common escape sequences.
	tempFormat = strings.NewReplacer(`\n`, "\n", `\t`, "\t", `\r`, "\r").Replace(tempFormat)
	// Finally, replace the $N placeholders with %s. This is safe now.
	processedFormat := regexGroupSeq.ReplaceAllString(tempFormat, "%s")

	var result bytes.Buffer
	for i, match := range matches {
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
	return result.Bytes(), nil
}
