package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
	"strings"
	"time"
)

// HistoryView represents the history display and selection modal.
type HistoryView struct {
	*tview.Box
	table         *tview.Table
	searchField   *tview.InputField
	flex          *tview.Flex
	app           *App          // Reference to the main App
	data          []HistoryItem // Original history data
	filteredData  []HistoryItem // Filtered history data
	selectedRegex string        // The regex selected by the user
	onSelect      func(regex string)
	onClose       func()
}

// NewHistoryView creates a new HistoryView.
func NewHistoryView(app *App) *HistoryView {
	hv := &HistoryView{
		Box:         tview.NewBox(),
		table:       tview.NewTable().SetSelectable(true, false).SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorDarkCyan)),
		searchField: tview.NewInputField().SetLabel("Filter: "),
		app:         app,
	}

	hv.table.SetBorder(true).SetTitle(" History ")
	hv.searchField.SetBorder(true)

	// Layout for the history view
	hv.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(hv.searchField, 3, 1, true).
		AddItem(hv.table, 0, 1, false)

	hv.searchField.SetChangedFunc(hv.filterHistory)

	// --- Input Capture for Sub-components ---

	// Capture input for the search field
	hv.searchField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter, tcell.KeyDown:
			hv.app.app.SetFocus(hv.table)
			return nil
		case tcell.KeyEsc:
			if hv.onClose != nil {
				hv.onClose()
			}
			return nil
		}
		return event
	})

	// Capture input for the table
	hv.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			row, _ := hv.table.GetSelection()
			if row > 0 && row <= len(hv.filteredData) {
				hv.selectedRegex = hv.filteredData[row-1].Regex
				if hv.onSelect != nil {
					hv.onSelect(hv.selectedRegex)
				}
			}
			return nil
		case tcell.KeyEsc:
			if hv.onClose != nil {
				hv.onClose()
			}
			return nil
		case tcell.KeyUp:
			if row, _ := hv.table.GetSelection(); row <= 1 {
				hv.app.app.SetFocus(hv.searchField)
				return nil
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case 'j':
				return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
			case 'k':
				// Need to check if we're at the top to switch focus
				if row, _ := hv.table.GetSelection(); row <= 1 {
					hv.app.app.SetFocus(hv.searchField)
					return nil
				}
				return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
			case 'q':
				if hv.onClose != nil {
					hv.onClose()
				}
				return nil
			}
		}
		return event
	})

	return hv
}

// SetHistoryData sets the data to be displayed in the history table.
func (hv *HistoryView) SetHistoryData(data []HistoryItem) {
	hv.data = data
	hv.filterHistory(hv.searchField.GetText()) // Initialize filtered data
}

// SetOnSelect sets the callback function for when a regex is selected.
func (hv *HistoryView) SetOnSelect(handler func(regex string)) {
	hv.onSelect = handler
}

// SetOnClose sets the callback function for when the view is closed without selection.
func (hv *HistoryView) SetOnClose(handler func()) {
	hv.onClose = handler
}

// filterHistory filters the history items based on the search text.
func (hv *HistoryView) filterHistory(searchText string) {
	hv.table.Clear()
	hv.filteredData = nil

	// Add table headers
	headers := []string{"Regex", "Count", "Last Used", "First Match"}
	for i, header := range headers {
		hv.table.SetCell(0, i, tview.NewTableCell(header).SetSelectable(false).SetAlign(tview.AlignCenter).SetExpansion(1).SetBackgroundColor(tcell.ColorDarkBlue))
	}

	searchText = strings.ToLower(searchText)
	rowIndex := 1
	for _, item := range hv.data {
		if searchText == "" || strings.Contains(strings.ToLower(item.Regex), searchText) || strings.Contains(strings.ToLower(item.FirstMatch), searchText) {
			hv.filteredData = append(hv.filteredData, item)
			hv.table.SetCell(rowIndex, 0, tview.NewTableCell(item.Regex).SetExpansion(10))
			hv.table.SetCell(rowIndex, 1, tview.NewTableCell(strconv.Itoa(item.Count)).SetExpansion(2).SetAlign(tview.AlignRight))
			hv.table.SetCell(rowIndex, 2, tview.NewTableCell(time.Unix(item.Timestamp, 0).Format("2006-01-02 15:04:05")).SetExpansion(2).SetAlign(tview.AlignLeft))
			hv.table.SetCell(rowIndex, 3, tview.NewTableCell(item.FirstMatch).SetExpansion(30))
			rowIndex++
		}
	}
	if len(hv.filteredData) == 0 {
		hv.table.SetCell(1, 0, tview.NewTableCell("No history items found.").SetAlign(tview.AlignCenter).SetSelectable(false))
		hv.table.SetOffset(0, 0)
	}
	hv.table.Select(1, 0)
}

// Draw implements tview.Primitive.
func (hv *HistoryView) Draw(screen tcell.Screen) {
	hv.flex.Draw(screen)
}

// GetRect implements tview.Primitive.
func (hv *HistoryView) GetRect() (int, int, int, int) {
	return hv.flex.GetRect()
}

// SetRect implements tview.Primitive.
func (hv *HistoryView) SetRect(x, y, width, height int) {
	hv.flex.SetRect(x, y, width, height)
}

// InputHandler returns the handler for this primitive.
func (hv *HistoryView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return hv.flex.InputHandler()
}

// Focus is called when this primitive receives focus.
func (hv *HistoryView) Focus(delegate func(p tview.Primitive)) {
	delegate(hv.searchField)
}

// HasFocus returns whether this primitive has focus.
func (hv *HistoryView) HasFocus() bool {
	return hv.searchField.HasFocus() || hv.table.HasFocus()
}
