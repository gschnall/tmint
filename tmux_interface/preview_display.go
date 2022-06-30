package tmux_interface

import (
	tcell "github.com/gdamore/tcell/v2"
	twiz "github.com/gschnall/tmint/tmux_wizard"
	"github.com/rivo/tview"
)

var (
	previewDisplay         = tview.NewTextView()
	currentPreviewPanePath = ""
)

func initPreviewDisplay() {
	previewDisplay.SetDynamicColors(true)
	previewDisplay.SetBackgroundColor(tcell.ColorDefault)
	previewDisplay.SetBorder(true).SetTitle(" " + sessionData.Sessions[0].Name).SetTitleAlign(0)
	previewDisplay.SetText(tview.TranslateANSI(sessionData.Sessions[0].Preview))
	previewDisplay.SetScrollable(true)
	previewDisplay.ScrollToEnd()
	initPreviewDisplayKeys()
}

// _____________________
// |                   |
// | Utils and Actions |
// |                   |
// --------------------
func handleChangeSession(session twiz.Session, node *tview.TreeNode) {
	currentPreviewPanePath = session.PanePath
	previewDisplay.ScrollToEnd()
	previewDisplay.Clear()
	previewDisplay.SetText(tview.TranslateANSI(session.Preview))
	previewDisplay.SetTitle(" " + session.Name + " ")
}
func handleChangeWindow(window twiz.Window, node *tview.TreeNode) {
	currentPreviewPanePath = window.PanePath
	previewDisplay.ScrollToEnd()
	previewDisplay.Clear()
	previewDisplay.SetText(tview.TranslateANSI(window.Preview))
	previewDisplay.SetTitle(" " + window.Index + " (" + window.Name + ")" + " ")
}
func handleChangePane(pane twiz.Pane, node *tview.TreeNode) {
	currentPreviewPanePath = pane.Path
	previewDisplay.ScrollToEnd()
	previewDisplay.Clear()
	previewDisplay.SetText(tview.TranslateANSI(pane.Preview))
	previewDisplay.SetTitle(" " + pane.Name + " - " + pane.Directory + " ")
}

// _____________________
// |                   |
// | Utils and Actions |
// |                   |
// ---------------------
func expandPreviewDisplay() {
	scrollback := twiz.GetColoredPaneScrollback(currentPreviewPanePath, sessionData.HistoryLimit)
	previewDisplay.SetText(tview.TranslateANSI(scrollback))
	previewDisplay.SetScrollable(true)
	previewDisplay.ScrollToBeginning()
	changeViewTo(previewDisplay)
}
func collapsePreviewDisplay() {
	previewDisplay.SetScrollable(false)
	previewDisplay.ScrollToEnd()
	restoreDefaultView()
}

// _______________
// |             |
// | Keybindings |
// |             |
// ---------------
func initPreviewDisplayKeys() {
	previewDisplay.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			collapsePreviewDisplay()
		case tcell.KeyEsc:
			collapsePreviewDisplay()
		}

		switch event.Rune() {
		case 'q':
			collapsePreviewDisplay()
		case 'f':
			startScrollbackHistoryForm(sessionDisplay.GetCurrentNode())
		}

		return event
	})
}
