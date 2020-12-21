package tmux_interface

import (
	"github.com/rivo/tview"
	tcell "github.com/gdamore/tcell/v2"

	twiz "github.com/gschnall/tmint/tmux_wizard"
)

var (
	previewDisplay = tview.NewTextView()
)

func initPreviewDisplay() {
	previewDisplay.SetDynamicColors(true)
	previewDisplay.SetBackgroundColor(tcell.ColorDefault)
	previewDisplay.SetBorder(true).SetTitle(" " + sessionData.Sessions[0].Name).SetTitleAlign(0)
	previewDisplay.SetText(tview.TranslateANSI(sessionData.Sessions[0].Preview))
	previewDisplay.ScrollToBeginning()
}

// _____________________
// |                   |
// | Utils and Actions |
// |                   |
// ---------------------
func handleChangeSession(session twiz.Session, node *tview.TreeNode) {
	previewDisplay.Clear()
	previewDisplay.SetText(tview.TranslateANSI(session.Preview))
	previewDisplay.SetTitle(" " + session.Name)
	previewDisplay.ScrollToBeginning()
}
func handleChangeWindow(window twiz.Window, node *tview.TreeNode) {
	previewDisplay.Clear()
	previewDisplay.SetText(tview.TranslateANSI(window.Preview))
	previewDisplay.SetTitle(" " + window.Index + " (" + window.Name + ")")
	previewDisplay.ScrollToBeginning()
}
func handleChangePane(pane twiz.Pane, node *tview.TreeNode) {
	previewDisplay.Clear()
	previewDisplay.SetText(tview.TranslateANSI(pane.Preview))
	previewDisplay.SetTitle(" " + pane.Name + " - " + pane.Directory)
	previewDisplay.ScrollToBeginning()
}
