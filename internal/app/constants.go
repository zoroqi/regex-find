package app

// Page Names
const (
	MainPage            = "main"
	RegexHelpPage       = "regex_help"
	KeybindingsHelpPage = "keybindings_help"
	HistoryPage         = "history_page"
	ExportPage          = "export"
	ResultPage          = "result"
)

// Widget Titles
const (
	TitleRegex         = "Regular Expression"
	TitleText          = "Text Input"
	TitleHighlighted   = "Highlighted"
	TitleMatches       = "Matches"
	TitleHelp          = "Help"
	TitleExportOptions = "Export Matches"
	TitleSuccess       = "Success"
	TitleError         = "Error"
	TitleMatchesFormat = "Matches (%d)"
)

// Form Labels & Button Text
const (
	LabelRegex        = "Regex: "
	LabelExportFormat = "Export Format"
	LabelCustomFormat = "Custom Format String"
	LabelGroupNumbers = "Group Numbers (comma-separated)"
	LabelOutputTarget = "Export Destination"
	LabelFilePath     = "File Path"
	ButtonExport      = "Export"
	ButtonCancel      = "Cancel"
	ButtonOK          = "OK"
)

// Export Options
const (
	OptJsonAll    = "JSON (all content)"
	OptJsonGroups = "JSON (specific groups)"
	OptCustom     = "Custom format"
)

// Output Targets
const (
	TargetClipboard = "Save to clipboard"
	TargetFile      = "Save to file"
)

// Help text blocks
const (
	HelpKeybindings = `[yellow]KEYBINDINGS:

[green]F1[white]:           Show this help modal
[green]F2[white]:           Show regex pattern help
[green]Ctrl+E[white]:       Show export options
[green]Tab / Shift+Tab[white]: Cycle focus between windows
[green]Ctrl+C / Ctrl+D[white]: Quit the application
[green]ESC[white]:          Close help or modals`

	HelpScrolling = `[yellow]SCROLLING (in 'Highlighted' and 'Matches' windows):

- [green]Arrow Keys[white]: Scroll up, down, left, right
- [green]h, j, k, l[white]:  Vim-style scrolling (left, down, up, right)`

	HintHelp = "F1 Helps | F2 Regex Help | F3 History | Ctrl+E Export | Ctrl+C Quit"
)
