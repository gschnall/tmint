package tmux_interface

import (
	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	helpBoxDisplay      = tview.NewTable()
	helpBoxDisplayWidth = 0
)

func initHelpBoxDisplay() {
	helpBoxDisplay.SetBorder(true)
	helpBoxDisplay.SetBorderColor(tcell.ColorGreen)
	helpBoxDisplay.SetTitle("KEYS")
	helpBoxDisplay.SetTitleColor(tview.Styles.PrimaryTextColor)
	helpBoxDisplay.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	helpBoxDisplay.SetTitleAlign(tview.AlignLeft)

	helpTextList := [][]string{
		{"", "arrow keys to navigate"},
		{"", "or hjkl to navigate"},
		{"", ""},
		{"", "ENTER to attach target"},
		{"", ""},
		{"/", "search | Ctrl-s"},
		{"s", "create session"},
		{"w", "create window"},
		{"d", "detach session"},
		{"x", "delete target"},
		{"r", "rename target"},
		{"c", "tmux cheat sheet"},
		{"v", "view scrollback history"},
		{"f", "save scrollback to file"},
		{"e", "toggle expand"},
		{"+", "expand all"},
		{"-", "collapse all"},
		{"?", "toggle help"},
	}

	for ind, rowList := range helpTextList {
		helpBoxDisplay.SetCell(ind, 0, tview.NewTableCell(rowList[0]).SetTextColor(tview.Styles.TertiaryTextColor))
		helpBoxDisplay.SetCell(ind, 1, tview.NewTableCell(rowList[1]))
	}
}

// _____________________
// |                   |
// | Utils and Actions |
// |                   |
// ---------------------
func toggleHelpBox() {
	width := helpBoxDisplayWidth
	if width > 0 {
		width = 0
	} else {
		width = 1
	}
	helpBoxDisplayWidth = width
	flexBoxDisplay.ResizeItem(helpBoxDisplay, width*27, width)
}

// - Only used for debugging
func testOutputInHelpBox(s string) {
	helpBoxDisplay.SetTitle(s)
}
