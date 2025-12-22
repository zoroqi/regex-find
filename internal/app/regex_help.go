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
		{Title: "IPv4", Pattern: `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`},
		{Title: "IPv6", Pattern: `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`},
		{Title: "MAC Address", Pattern: `^([0-9a-fA-F][0-9a-fA-F]:){5}([0-9a-fA-F][0-9a-fA-F])$`},
		{Title: "URL", Pattern: `https://?[	\d/.-]+`},
		{Title: "Date (YYYY-MM-DD)", Pattern: `\d{4}-\d{2}-\d{2}`},
		{Title: "Username", Pattern: `^[a-zA-Z0-9_-]{3,16}$`},
		{Title: "HTML Tag", Pattern: `<([a-z]+)([^<]+)*(?:>(.*)<\/\1>|\s+\/>)`},
		{Title: "Version", Pattern: `v?(\d+\.)?(\d+\.)?(\*|\d+)`},
		{Title: "CJK", Pattern: `[\p{Han}\p{Hiragana}\p{Katakana}]`},
		{Title: "Password(Medium)", Pattern: "(?=.*\\d)(?=.*[a-z])(?=.*[A-Z]).{6,}"},
		{Title: "Password(Strong)", Pattern: ".*(?=.{6,})(?=.*\\d)(?=.*[A-Z])(?=.*[a-z])(?=.*[!@#$%^&*? ]).*"},
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
