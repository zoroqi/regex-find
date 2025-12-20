package app

import "regexp"

// Search performs the regex matching on the provided text.
// It returns the compiled regex, match indices, submatches, and any error encountered.
func Search(regexStr, text string) (*regexp.Regexp, [][]int, [][]string, error) {
	if regexStr == "" {
		return nil, nil, nil, nil
	}

	re, err := regexp.Compile(regexStr)
	if err != nil {
		return nil, nil, nil, err
	}

	indices := re.FindAllStringIndex(text, -1)
	matches := re.FindAllStringSubmatch(text, -1)

	return re, indices, matches, nil
}
