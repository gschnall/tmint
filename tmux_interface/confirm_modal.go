package tmux_interface

import "github.com/rivo/tview"

var (
	confirmModal = tview.NewModal()
	targetedNode *tview.TreeNode
)

func initConfirmModal() {
	confirmModal.SetTitle("ESC to Cancel")
	confirmModal.SetBorder(true)

	confirmModal.AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				killTmuxTarget(targetedNode, true)
			} else {
				flexBoxWrapper.HidePage("confirmModal")
			}
		})
	// --- Need to create custom modal to capture different keypresses
}
