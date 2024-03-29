package tmux_interface

import (
	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	twiz "github.com/gschnall/tmint/tmux_wizard"
)

var (
	inactiveSessionDisplay = tview.NewTextView()
)

func initNoActiveSessionDisplay() {
	inactiveSessionDisplay.SetTitle("Tmux has not been started")
	inactiveSessionDisplay.SetBorder(true)
	// Displayed Text
	headerText := getTmintHeader() + "\n\n"
	info := "   No active sessions available\n\n\n"
	start := "   Press ENTER or 't' to start Tmux\n\n"
	create := "   ----- 's' to create a session\n\n"
	quit := "   ----- 'q' to quit"
	inactiveSessionDisplay.SetText(headerText + info + start + create + quit)

	initNoActiveSessionDisplayKeys()
}

// _____________________
// |                   |
// | Utils and Actions |
// |                   |
// ---------------------
func getTmintHeader() string {
	h := "\n\n:  :  : :   ::  ::  ::  :: :: : :  :  : \n"
	h += "   -----| tmint - a Tmux interface\n"
	h += ":  :  : :   ::  ::  ::  :: :: : :  :  : \n"
	return h
}

// _______________
// |             |
// | Keybindings |
// |             |
// ---------------
func initNoActiveSessionDisplayKeys() {
	inactiveSessionDisplay.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			tviewApp.Stop()
			twiz.StartTmux()
		case tcell.KeyEsc:
			tviewApp.Stop()
		}

		switch event.Rune() {
		case 'q':
			tviewApp.Stop()
		case 's':
			startCreateForm("session")
		case 't':
			tviewApp.Stop()
			twiz.StartTmux()
		}

		return event
	})
}
