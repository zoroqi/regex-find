package app

// HelpItem represents a single entry in the help document.
type HelpItem struct {
	Title   string
	Pattern string
}

// HelpContent contains the different sections of the help document.
type HelpContent struct {
	Common  []HelpItem
	Escapes []HelpItem
}

// RegexHelpData holds all the predefined help information.
var RegexHelpData = HelpContent{
	Common: []HelpItem{
		{Title: "Email", Pattern: `(?i)^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$`},
		{Title: "IPv4 Address", Pattern: `\b(?:\d{1,3}\.){3}\d{1,3}\b`},
		{Title: "URL", Pattern: `https://?[	\d/.-]+`},
		{Title: "Date (YYYY-MM-DD)", Pattern: `\d{4}-\d{2}-\d{2}`},
		{Title: "Username", Pattern: `^[a-zA-Z0-9_-]{3,16}$`},
		{Title: "HTML Tag", Pattern: `<([a-z]+)([^<]+)*(?:>(.*)<\/\1>|\s+\/>)`},
	},
	Escapes: []HelpItem{
		{Title: "Digit", Pattern: `\d`},
		{Title: "Non-Digit", Pattern: `\D`},
		{Title: "Word Character", Pattern: `\w`},
		{Title: "Non-Word Character", Pattern: `\W`},
		{Title: "Whitespace", Pattern: `\s`},
		{Title: "Non-Whitespace", Pattern: `\S`},
		{Title: "Tab", Pattern: `\t`},
		{Title: "Newline", Pattern: `\n`},
		{Title: "Carriage Return", Pattern: `\r`},
		{Title: "Word Boundary", Pattern: `\b`},
		{Title: "Literal Dot", Pattern: `\.`},
		{Title: "Literal Star", Pattern: `\*`},
	},
}
