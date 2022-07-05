package tmux_interface

import (
	tcell "github.com/gdamore/tcell/v2"
	twiz "github.com/gschnall/tmint/tmux_wizard"
	"github.com/rivo/tview"
)

var (
	resizeDisplay = tview.NewTextView()
)

func initResizeDisplay() {
	resizeDisplay.SetDynamicColors(true)
	resizeDisplay.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	resizeDisplay.SetTextColor(tview.Styles.PrimaryTextColor)
	resizeDisplay.SetBorder(true).SetTitle(" Tmint | Resize Pane ").SetTitleAlign(0)
	setupResizeTextDisplay()

	initResizeDisplayKeys()
}

func setupResizeTextDisplay() {
	resizeText := `
  |----------------------|
  |    Resize Pane       |
  |                      |
  |  arrow keys or hjkl  |
  |----------------------|
  |   quit: q or esc     |
  ------------------------`
	resizeDisplay.SetText(resizeText)
}

// _____________________
// |                   |
// | Utils and Actions |
// |                   |
// ---------------------

// _______________
// |             |
// | Keybindings |
// |             |
// ---------------
func initResizeDisplayKeys() {
	resizeDisplay.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			tviewApp.Stop()
		case tcell.KeyEsc:
			tviewApp.Stop()
		case tcell.KeyUp:
			twiz.TmuxResizePaneUp()
		case tcell.KeyDown:
			twiz.TmuxResizePaneDown()
		case tcell.KeyRight:
			twiz.TmuxResizePaneRight()
		case tcell.KeyLeft:
			twiz.TmuxResizePaneLeft()
		}

		switch event.Rune() {
		case 'q':
			tviewApp.Stop()
		case 'j':
			twiz.TmuxResizePaneDown()
		case 'k':
			twiz.TmuxResizePaneUp()
		case 'h':
			twiz.TmuxResizePaneLeft()
		case 'l':
			twiz.TmuxResizePaneRight()
		}

		return event
	})
}
