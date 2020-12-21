package tmux_interface

import (
	"fmt"
	"time"

	twiz "github.com/gschnall/tmint/tmux_wizard"
	"github.com/rivo/tview"
)

var (
	renameForm     = tview.NewForm()
	renameFormName = "" 
)

func initRenameForm() {
	tmux := sessionDisplay.GetCurrentNode().GetReference()
	targetName := ""
	targetType := "session"
	switch tmux.(type) {
	case twiz.Session:
		targetName = tmux.(twiz.Session).Name
	case twiz.Window:
		targetType = "window"
		targetName = tmux.(twiz.Window).Name
	default:
		return
	}

	renameSession := func () {
		if renameFormName != "" {
			twiz.RenameTmuxSession(tmux.(twiz.Session).Name, renameFormName)
		}
		endRenameForm()
		time.Sleep(200 * time.Millisecond)
		refreshSessionDisplay()
	}
	renameWindow := func () {
		if renameFormName != "" {
			twiz.RenameTmuxWindow(tmux.(twiz.Window).Path, renameFormName)
		}
		endRenameForm()
		time.Sleep(200 * time.Millisecond)
		refreshSessionDisplay()
	}

	renameForm.SetTitle("Rename " + targetType + " " + targetName)
	renameForm.AddInputField("Name", targetName, 20, tmuxNameIsValid, func(n string) {
			renameFormName = n
		})

		if targetType == "window" {
			renameForm.AddButton("Save", renameWindow)
		} else {
			renameForm.AddButton("Save", renameSession)
		} 
		renameForm.AddButton("Cancel", endRenameForm)

	renameForm.SetBorder(true).SetTitle(fmt.Sprintf(" Rename %s | ESC to cancel | Ctrl-u to clear input ", targetType)).SetTitleAlign(tview.AlignLeft)

	renameForm.SetCancelFunc(func() {
		endRenameForm()
	})
}

// _____________________
// |                   |
// | Utils and Actions |
// |                   |
// ---------------------
func startRenameForm() {
	tmux := sessionDisplay.GetCurrentNode().GetReference()
	switch tmux.(type) {
	case twiz.Pane:
		return
	default:
		initRenameForm()
		changeViewTo(renameForm)
	}
}

func endRenameForm() {
	renameForm.Clear(true)
	restoreDefaultView()
}
