package tmux_interface

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	loadingModal          = tview.NewModal()
	loadingModalIsLoading = false
)

func initLoadingModal() {
	loadingModal.SetTitle(" Loading ")
	loadingModal.SetBorder(true)
	loadingModal.SetBackgroundColor(tcell.ColorOrangeRed)
	loadingModal.SetTextColor(tcell.ColorWhiteSmoke)
}

// Use: go startLoadingModal("my message")
//      stopLoadingModal()
func startLoadingModal(message string) {
	// focus on modal and disable mouse
	flexBoxWrapper.ShowPage("loadingModal")
	tviewApp.SetFocus(loadingModal)
	tviewApp.EnableMouse(false)
	loadingModalIsLoading = true
	count := 0
	loader := [4]string{"|", "/", "-", "\\"}
	refreshInterval := 300 * time.Millisecond
	for {
		if !loadingModalIsLoading {
			break
		}
		time.Sleep(refreshInterval)
		tviewApp.QueueUpdateDraw(func() {
			loadingModal.SetText(loader[count] + " - " + message + " - " + loader[count])
			count++
			if count >= len(loader) {
				count = 0
			}
		})
	}
}

func stopLoadingModal() {
	tviewApp.EnableMouse(true)
	loadingModalIsLoading = false
	flexBoxWrapper.HidePage("loadingModal")
	restoreDefaultView()
}
