package tmux_interface

import (
	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	twiz "github.com/gschnall/tmint/tmux_wizard"
)

// Views and their state variables
var (
	wasUserPaneZoomed = false
	sessionData       twiz.SessionData
	tviewApp          = tview.NewApplication()
	flexBoxDisplay    = tview.NewFlex()
	flexBoxWrapper    = tview.NewPages()
	mainFlexBoxView   = tview.NewFlex()
)

func initNoActiveSessionInterface() {
	// Test Feature
	// - Allow users to create new session
	initNoActiveSessionDisplay()
	initDirectoryInputField()
	initConfirmModal()

	if err := tviewApp.SetRoot(inactiveSessionDisplay, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func initInterface() {
	initSessionDisplay()
	initPreviewDisplay()
	initHelpBoxDisplay()
	initCheatSheetDisplay()
	initSearchBoxDisplay()
	initDirectoryInputField()
	initConfirmModal()
	initLoadingModal()
	initScrollbackHistoryForm()

	flexBoxDisplay.
		AddItem(helpBoxDisplay, 0, 0, false). // FOR HELP MENU
		AddItem(mainFlexBoxView.SetDirection(tview.FlexRow).
			// AddItem(confirmModal, 0, 0, false).
			// AddItemgloadingModal, 0, 0, false).
			AddItem(sessionDisplay, 0, 6, true).
			AddItem(searchBoxDisplay, 0, 0, false).
			AddItem(previewDisplay, 0, 3, false), 0, 5, true)

	// Workaround - Without this, Clear functionality breaks
	tviewApp.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})

	flexBoxWrapper.
		AddPage("flexBoxDisplay", flexBoxDisplay, true, true).
		AddPage("confirmModal", confirmModal, false, false).
		AddPage("loadingModal", loadingModal, false, false)

	// flexBoxWrapper.SetBackgroundColor(tcell.ColorDefault)

	if err := tviewApp.SetRoot(flexBoxWrapper, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func setTviewColorScheme() {
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    tcell.NewRGBColor(0, 0, 0),
		ContrastBackgroundColor:     tcell.NewRGBColor(26, 28, 28),   // Dark Grey
		MoreContrastBackgroundColor: tcell.NewRGBColor(68, 215, 168), // Green
		BorderColor:                 tcell.NewRGBColor(255, 255, 255),
		TitleColor:                  tcell.NewRGBColor(255, 255, 255),
		GraphicsColor:               tcell.NewRGBColor(255, 255, 255),
		PrimaryTextColor:            tcell.NewRGBColor(255, 255, 255),
		SecondaryTextColor:          tcell.NewRGBColor(255, 222, 0),   // Yellow
		TertiaryTextColor:           tcell.NewRGBColor(152, 251, 152), // Mint Green
		InverseTextColor:            tcell.NewRGBColor(62, 164, 180),  // Mint Blue
		ContrastSecondaryTextColor:  tcell.NewRGBColor(25, 102, 255),  // Blue
	}
}

func StartResizeInterface() {
	initResizeDisplay()

	if err := tviewApp.SetRoot(resizeDisplay, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func Start(doNotZoomPane bool, currentSession string, tmintSession string, runFromKeybindings bool) {
	result := make(chan twiz.SessionData, 1)
	go twiz.GetSessionData(currentSession, tmintSession, runFromKeybindings, result)
	dataResult := <-result
	sessionData = dataResult
	close(result)

	setTviewColorScheme()

	if sessionData.HasLivingSessions == false {
		initNoActiveSessionInterface()
	} else {
		// Zoom user pane if in Tmux
		if !doNotZoomPane && sessionData.HasAttachedSession && !sessionData.HasZoomedPane {
			wasUserPaneZoomed = true
			twiz.TmuxToggleFullscreen()
		}
		// Start interface
		initInterface()
		// Unzoom pane if it wasn't zoomed before
		if wasUserPaneZoomed {
			twiz.TmuxToggleFullscreen()
		}
		// Delete session for tmux-keybinding workflow
		if currentSession != ":" {
			twiz.SwitchToTmuxPath(sessionData.AttachedSession)
			twiz.KillTmuxSession(tmintSession)
		}
	}
}
