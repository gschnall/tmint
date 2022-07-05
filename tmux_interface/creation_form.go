package tmux_interface

import (
	"fmt"
	"os"
	"strings"
	"time"

	twiz "github.com/gschnall/tmint/tmux_wizard"

	"github.com/rivo/tview"
)

var (
	creationForm          = tview.NewForm()
	creationFormName      = ""
	creationFormSession   = ""
	creationFormPaneCount = 0
)

func setCreationFormColorTheme() {
	creationForm.SetBackgroundColor(tview.Styles.ContrastBackgroundColor)
	creationForm.SetFieldTextColor(tview.Styles.PrimaryTextColor)
	creationForm.SetFieldBackgroundColor(tview.Styles.ContrastSecondaryTextColor)
	creationForm.SetLabelColor(tview.Styles.SecondaryTextColor)
	creationForm.SetTitleColor(tview.Styles.PrimaryTextColor)
	creationForm.SetButtonTextColor(tview.Styles.PrimitiveBackgroundColor)
	creationForm.SetButtonBackgroundColor(tview.Styles.InverseTextColor)
}

func initCreationForm(creation string) {
	creationForm.SetTitle("Create new " + creation)

	if creation == "window" {
		currentSessionName := getSessionFromNode(sessionDisplay.GetCurrentNode()).GetReference().(twiz.Session).Name
		creationFormSession = currentSessionName
		sessionNames := []string{}
		initialOption := 0
		for index, session := range sessionData.Sessions {
			sessionNames = append(sessionNames, session.Name)
			if session.Name == currentSessionName {
				initialOption = index
			}
		}

		creationForm.AddDropDown("Session", sessionNames, initialOption, func(s string, option int) {
			creationFormSession = s
		})
	}

	creationForm.
		AddInputField("Name", "", 20, tmuxNameIsValid, func(n string) {
			creationFormName = n
		}).
		AddFormItem(directoryInputField)
	// -- Thought: feature - let users select # of panes --
	// AddDropDown("Panes", []string{"1", "2", "3", "4"}, 0, func(option string, ind int) {
	// 	creationFormPaneCount = ind + 1
	// }).
	// -- -------------------------- --
	if creation == "window" {
		creationForm.AddButton("Save", createNewWindow)
		creationForm.SetFocus(1)
	} else {
		creationForm.AddButton("Save", createNewSession)
	}
	creationForm.AddButton("Cancel", endCreateSession)
	creationForm.SetBorder(true).SetTitle(fmt.Sprintf(" Create New %s | ESC to cancel ", creation)).SetTitleAlign(tview.AlignLeft)

	creationForm.SetCancelFunc(func() {
		endCreateSession()
	})

	setCreationFormColorTheme()
}

// _____________________
// |                   |
// | Utils and Actions |
// |                   |
// ---------------------
func startCreateForm(creation string) {
	initCreationForm(creation)
	changeViewTo(creationForm)
}
func endCreateSession() {
	creationForm.Clear(true)
	restoreDefaultView()
}

func createNewWindow() {
	sName := creationFormName
	pCount := creationFormPaneCount
	sDir := directoryInputField.GetText()
	if sDir == "" {
		homeDir, homeDirErr := os.UserHomeDir()
		twiz.HandleExecError(homeDirErr, "createNewWindow")
		sDir = homeDir
	}
	if pCount == 1 {
		previewDisplay.SetTitle(sName + " " + sDir)
	}
	twiz.CreateTmuxWindow(creationFormName, directoryInputField.GetText(), creationFormPaneCount, creationFormSession)
	endCreateSession()
	time.Sleep(200 * time.Millisecond)
	refreshSessionDisplay()
}

func createNewSession() {
	sName := creationFormName
	pCount := creationFormPaneCount
	sDir := directoryInputField.GetText()
	if sDir == "" {
		homeDir, homeDirErr := os.UserHomeDir()
		twiz.HandleExecError(homeDirErr, "createNewSession")
		sDir = homeDir
	}
	if pCount == 1 {
		previewDisplay.SetTitle(sName + " " + sDir)
	}
	twiz.CreateTmuxSession(creationFormName, directoryInputField.GetText(), creationFormPaneCount)
	endCreateSession()
	time.Sleep(200 * time.Millisecond)
	refreshSessionDisplay()
}

func tmuxNameIsValid(name string, lastChar rune) bool {
	return !(strings.Contains(name, ".") || strings.Contains(name, ":"))
}
