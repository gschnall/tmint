package tmux_interface

import (
	"fmt"
	"strconv"

	twiz "github.com/gschnall/tmint/tmux_wizard"
	"github.com/rivo/tview"
)

var (
	scrollbackHistoryForm        = tview.NewForm()
	scrollback_history_file_name = "tmint_scrollback.sh"
	scrollbackPanePath           = ""
	saveToTmintDir               = true
	scrollbackHistoryLimit       = 2000
)

func setScrollbackHistoryFormColorTheme() {
	scrollbackHistoryForm.SetBackgroundColor(tview.Styles.ContrastBackgroundColor)
	scrollbackHistoryForm.SetFieldTextColor(tview.Styles.PrimaryTextColor)
	scrollbackHistoryForm.SetFieldBackgroundColor(tview.Styles.ContrastSecondaryTextColor)
	scrollbackHistoryForm.SetLabelColor(tview.Styles.SecondaryTextColor)
	scrollbackHistoryForm.SetTitleColor(tview.Styles.PrimaryTextColor)
	scrollbackHistoryForm.SetButtonTextColor(tview.Styles.PrimitiveBackgroundColor)
	scrollbackHistoryForm.SetButtonBackgroundColor(tview.Styles.InverseTextColor)
}

func initScrollbackHistoryForm() {
	scrollbackHistoryLimit = sessionData.HistoryLimit

	scrollbackHistoryForm.SetBorder(true).SetTitle(" Save Scrollback for pane | ESC to cancel").SetTitleAlign(tview.AlignLeft)
	// scrollbackActions := []string{"View", "Save", "Save & Edit", "Save & View"}
	// scrollbackHistoryForm.AddDropDown("Action", scrollbackActions, 0, func(option string, optionIndex int) {
	// 	// Case statement for each action
	// 	switch option {
	// 	case "Save":
	// 		saveScrollbackHistory(panePathForScrollbackHistory)
	// 	}
	// })
	scrollbackHistoryForm.AddCheckbox("Save to ~/.tmux_interface", true, func(checked bool) {
		saveToTmintDir = checked
	})

	scrollbackHistoryForm.AddInputField("Number of Lines", strconv.Itoa(scrollbackHistoryLimit), 30,
		func(textToCheck string, lastChar rune) bool {
			_, err := strconv.Atoi(textToCheck)
			if err != nil || len(textToCheck) > 7 {
				return false
			}
			return true
		},
		func(text string) {
			newHistoryLimit, err := strconv.Atoi(text)
			if err != nil {
				fmt.Println(err.Error())
			}
			scrollbackHistoryLimit = newHistoryLimit
		})
	scrollbackHistoryForm.GetFormItemByLabel("Number of Lines")

	scrollbackHistoryForm.AddInputField("File Name", "tmint_scrollback.sh", 30, nil, func(text string) {
		scrollback_history_file_name = text
	})

	scrollbackHistoryForm.AddButton("Save", saveScrollbackHistory)
	scrollbackHistoryForm.AddButton("Cancel", endScrollbackHistoryForm)

	scrollbackHistoryForm.SetCancelFunc(endScrollbackHistoryForm)

	setScrollbackHistoryFormColorTheme()
}

// _____________________
// |                   |
// | Utils and Actions |
// |                   |
// ---------------------
func startScrollbackHistoryForm(node *tview.TreeNode) {
	tmux := node.GetReference()
	switch tmux.(type) {
	case twiz.Session:
		scrollbackPanePath = tmux.(twiz.Session).PanePath
	case twiz.Window:
		scrollbackPanePath = tmux.(twiz.Window).PanePath
	case twiz.Pane:
		scrollbackPanePath = tmux.(twiz.Pane).Path
	}

	scrollbackHistoryForm.SetFocus(1)

	changeViewTo(scrollbackHistoryForm)
}
func endScrollbackHistoryForm() {
	restoreDefaultView()
}
func saveScrollbackHistory() {
	twiz.SaveTmuxScrollback(scrollbackPanePath, scrollback_history_file_name, saveToTmintDir, scrollbackHistoryLimit)
	endScrollbackHistoryForm()
}
