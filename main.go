package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/jedib0t/go-pretty/text"
	"github.com/urfave/cli"
)

// List of views
var (
	isUserPaneZoomed       = false
	viewArr                = []string{"header", "sessions", "preview"}
	sessionData            = getSessionData()
	active                 = 0
	tviewApp               = tview.NewApplication()
	previewDisplay         = tview.NewTextView()
	sessionDisplay         = tview.NewTreeView()
	sessionDisplaySess     *tview.TreeNode
	sessionDisplayWind     *tview.TreeNode
	targetedNode           *tview.TreeNode
	creationForm           = tview.NewForm()
	creationFormDir        = ""
	creationFormName       = ""
	creationFormPaneCount  = 0
	helpBoxDisplay         = tview.NewTable()
	helpBoxDisplayWidth    = 0
	searchBoxDisplay       = tview.NewInputField()
	searchBoxDisplayHeight = 0
	confirmModal           = tview.NewModal()
	directoryInputField    = tview.NewInputField()
	firstTabCompletionTerm = ""
	flexBoxDisplay         = tview.NewFlex()
	flexBoxWrapper         = tview.NewPages()
	mainFlexBoxView        = tview.NewFlex()
)

const (
	treeSignDash     = "─"
	treeSignVertical = "│─"
	treeSignUpMiddle = "├── "
	treeSignUpMid    = "├─ "
	treeSignUpEnding = "└─"
)

func printHeader() string {
	h := text.FgGreen.Sprint("    :  :  : :   ::  ::  ::  :: :: :\n")
	h += text.FgBlue.Sprint("|___-_|-__-|-__|--_|--__--__--_--_-|\n")
	h += text.Bold.Sprint("--------| Tmux-Interface")
	h += text.FgRed.Sprint(" > > > > > > >\n")
	h += text.FgBlue.Sprint("|___-_|-__-|-__|--_|--__--__--_--_-|\n")
	return h
	// fmt.Println(h)
}

func getColoredSessionName(name string, isAttached bool) string {
	if isAttached {
		return name + " (attached)"
	}
	return name
}

func getColoredSessionNumber(index int, isAttached bool) string {
	if isAttached {
		return text.FgGreen.Sprint(getSessionNumberWithPadding(index))
	}
	return text.FgBlue.Sprint(getSessionNumberWithPadding(index))
}

func getSessionNumberWithPadding(index int) string {
	sn := index + 1
	if sn <= 9 {
		return " " + strconv.Itoa(sn)
	}
	return strconv.Itoa(sn)
}

func getColoredSwitchText(index int, isAttached bool, firstSessionIsAttached bool) string {
	if isAttached {
		return text.FgYellow.Sprint(" d") + ":DETACH"
	}
	if index == 0 || (firstSessionIsAttached && index == 1) {
		return text.FgGreen.Sprint(getSessionNumberWithPadding(index+1)) + ":SWITCH"
	}
	return text.FgGreen.Sprint(getSessionNumberWithPadding(index + 1))
}
func getColoredRenameText(index int) string {
	if index == 0 {
		return text.FgYellow.Sprint("r"+strconv.Itoa(index+1)) + ":RENAME"
	}
	return text.FgYellow.Sprint("r" + strconv.Itoa(index+1))
}
func getColoredKillText(index int) string {
	if index == 0 {
		return text.FgRed.Sprint("k"+strconv.Itoa(index+1)) + ":KILL"
	}
	return text.FgRed.Sprint("k" + strconv.Itoa(index+1))
}

func formatSessionId(id int) string {
	if id < 10 {
		return "(" + strconv.Itoa(id) + ")  "
	} else if id < 36 {
		letter := toCharStr(id - 9)
		return "(M-" + letter + ")"
	}
	return ""
}

func handleChangeSession(session Session, node *tview.TreeNode) {
	sessionDisplaySess = node
	previewDisplay.Clear()
	previewDisplay.SetText(tview.TranslateANSI(session.preview))
	previewDisplay.SetTitle(" " + session.name)
	previewDisplay.ScrollToBeginning()
}
func handleChangeWindow(window Window, node *tview.TreeNode) {
	sessionDisplayWind = node
	previewDisplay.Clear()
	previewDisplay.SetText(tview.TranslateANSI(window.preview))
	previewDisplay.SetTitle(" " + window.index + " (" + window.name + ")")
	previewDisplay.ScrollToBeginning()
}
func handleChangePane(pane Pane, node *tview.TreeNode) {
	previewDisplay.Clear()
	previewDisplay.SetText(tview.TranslateANSI(pane.preview))
	previewDisplay.SetTitle(" " + pane.name + " - " + pane.directory)
	previewDisplay.ScrollToBeginning()
}

func expandNode(node *tview.TreeNode) {
	if !node.IsExpanded() {
		node.SetExpanded(true)
	}
	runCallbacksForNode(node, toggleSessionName, toggleWindowName, togglePaneName)
}
func collapseNode(node *tview.TreeNode, moveNodes bool) {
	node.Collapse()

	if moveNodes {
		tmux := node.GetReference()
		switch tmux.(type) {
		case Window:
			sessionDisplay.SetCurrentNode(sessionDisplaySess)
		case Pane:
			sessionDisplay.SetCurrentNode(sessionDisplayWind)
		}
	}
	runCallbacksForNode(node, toggleSessionName, toggleWindowName, togglePaneName)
}
func toggleHelpBox() {
	width := helpBoxDisplayWidth
	if width > 0 {
		width = 0
	} else {
		width = 1
	}
	helpBoxDisplayWidth = width
	flexBoxDisplay.ResizeItem(helpBoxDisplay, 0, width)
}
func toggleSearchBox() {
	height := searchBoxDisplayHeight
	if height > 0 {
		height = 0
		tviewApp.SetFocus(sessionDisplay)
		searchBoxDisplay.SetText("")
	} else {
		height = 1
		tviewApp.SetFocus(searchBoxDisplay)
	}
	searchBoxDisplayHeight = height
	mainFlexBoxView.ResizeItem(searchBoxDisplay, 0, height)
}
func startCreateSession() {
	initCreationForm()
	changeViewTo(creationForm)
}
func endCreateSession() {
	creationForm.Clear(true)
	restoreDefaultView()
}

func highlightNode(index int) {
	node := sessionDisplay.GetRoot().GetChildren()[index]
	sessionDisplay.SetCurrentNode(node)
}
func highlightNodeDotPath(path string) {
	parts := strings.Split(path, ".")
	currentNode := sessionDisplay.GetCurrentNode()
	sessionChildren := currentNode.GetChildren()

	if len(parts) > 1 && len(parts[1]) > 0 {
		//---------------------
		windowSearch := parts[1]
		// -- Session should already be selected by this point
		// -- so parts[0] isn't necessary
		// - if windowSearch is int, convert it and hightlight that index
		windowIndex, err := strconv.Atoi(windowSearch)
		if err == nil {
			currentNode.Expand()
			node := sessionChildren[windowIndex]
			// searchBoxDisplay.SetLabel(node.GetText())
			sessionDisplay.SetCurrentNode(node)
		} else {
			node := sessionChildren[windowIndex]
			sessionDisplay.SetCurrentNode(node)
			// -- for loop through currentNodeChildren
			//			find the matching window name using algo used in parent function
			//      & highlight it
		}
		//---------------------
	} else {
		currentNode.Expand()
	}

	if len(parts) > 2 {
		//---------------------
		// -- Current node will be the window found in previous if statement
		// - currentNode = selected/highlighted node
		// - currentNodeChildren = currentNode.GetChildren()
		// - paneSearch = parts[2]
		// - if paneSearch cannot be converted to int, perform no-op
		// - if paneSearch is int, convert it and hightlight that index from currentNodeChildren
		//---------------------
	}
}
func highlightSessionNode(index int) {
	if index > -1 && index < len(sessionData.sessions) {
		highlightNode(index)
	}
}

func initSearchDisplayKeys() {
	searchBoxDisplay.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter, tcell.KeyCtrlS, tcell.KeyEsc:
			toggleSearchBox()
		}
		return event
	})
}

func initCreationFormKeys() {
	creationForm.SetCancelFunc(func() {
		endCreateSession()
	})
}

func initSessionDisplayKeys() {
	sessionDisplay.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			expandNode(sessionDisplay.GetCurrentNode())
		case tcell.KeyLeft:
			collapseNode(sessionDisplay.GetCurrentNode(), false)
		case tcell.KeyCtrlS:
			toggleSearchBox()
		case tcell.KeyEnter:
			switchTmuxToCurrentNode()
		}

		switch event.Rune() {
		case '/':
			toggleSearchBox()
		case 'l':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				expandNode(sessionDisplay.GetCurrentNode())
			}
		case 'd':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				detachSession(sessionDisplay.GetCurrentNode())
			}
		case 'h':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				collapseNode(sessionDisplay.GetCurrentNode(), true)
			}
		case 'q':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				tviewApp.Stop()
			}
		case '?':
			toggleHelpBox()
		case 's':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			}
		case 'x':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				killTmuxTarget(sessionDisplay.GetCurrentNode(), false)
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			highlightSessionNode(int(event.Rune() - '0'))
		case 'c':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				startCreateSession()
			}

		case 'a', 'b', 'e', 'f', 'g', 'i', 'j', 'k', 'm', 'n', 'o', 'p', 'r', 't', 'u', 'v', 'w', 'y', 'z':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			}
		}

		return event
	})
}

func getSessionDisplayName(session Session, isExpanded bool) string {
	mark := "+ "
	if isExpanded {
		mark = "- "
	}
	return formatSessionId(session.id) + mark + getColoredSessionName(session.name, session.isAttached)
}

func getWindowDisplayName(window Window, isExpanded bool) string {
	if isExpanded {
		return " - " + window.index + ": " + window.name + " (" + strconv.Itoa(len(window.panes)) + " panes)"
	}
	return " + " + window.index + ": " + window.name + " (" + strconv.Itoa(len(window.panes)) + " panes)"
}
func getPaneDisplayName(pane Pane, isExpanded bool) string {
	if isExpanded {
		return pane.name + ": " + pane.command
	}
	return pane.name + ": " + pane.command
}

func toggleSessionName(session Session, node *tview.TreeNode) {
	node.SetText(getSessionDisplayName(session, node.IsExpanded()))
}
func toggleWindowName(window Window, node *tview.TreeNode) {
	node.SetText(getWindowDisplayName(window, node.IsExpanded()))
}
func togglePaneName(pane Pane, node *tview.TreeNode) {
	node.SetText(getPaneDisplayName(pane, node.IsExpanded()))
}

func getSessionFromRef(ids []int) Session {
	return sessionData.sessions[ids[0]]
}
func getWindowFromRef(ids []int) Window {
	return sessionData.sessions[ids[0]].windows[ids[1]]
}
func getPaneFromRef(ids []int) Pane {
	return sessionData.sessions[ids[0]].windows[ids[1]].panes[ids[2]]
}

type SessionFunc func(Session, *tview.TreeNode)
type WindowFunc func(Window, *tview.TreeNode)
type PaneFunc func(Pane, *tview.TreeNode)

func runCallbacksForNode(node *tview.TreeNode, sfunc SessionFunc, wfunc WindowFunc, pfunc PaneFunc) {
	// idArr := node.GetReference().([]int)
	tmux := node.GetReference()
	switch tmux.(type) {
	case Session:
		sfunc(tmux.(Session), node)
	case Window:
		wfunc(tmux.(Window), node)
	case Pane:
		pfunc(tmux.(Pane), node)
	}

	// if len(idArr) == 1 {
	// 	sfunc(getSessionFromRef(idArr), node)
	// } else if len(idArr) == 2 {
	// 	wfunc(getWindowFromRef(idArr), node)
	// } else if len(idArr) == 3 {
	// 	pfunc(getPaneFromRef(idArr), node)
	// }
}

func initSessionDisplay() {
	// box := tview.NewBox().SetBorder(true)
	// box.HasFocus(true)
	// box.
	root := tview.NewTreeNode("").
		SetSelectable(false)
	sessionDisplay.
		SetRoot(root).
		SetTopLevel(1)
		// - BUG - Breaks right now - https://github.com/rivo/tview/issues/314
		// SetBackgroundColor(tcell.ColorDefault)

	for sInd, session := range sessionData.sessions {
		sNode := tview.NewTreeNode(getSessionDisplayName(session, false))
		sNode.SetSelectable(true).SetExpanded(false)
		// sNode.SetReference([]int{session.id})
		sNode.SetReference(session)
		for _, window := range session.windows {
			wNode := tview.NewTreeNode(getWindowDisplayName(window, false))
			// wNode.SetReference([]int{session.id, wInd})
			wNode.SetReference(window)
			wNode.SetExpanded(false)
			wNode.SetSelectable(true)
			wNode.SetIndent(6)
			for _, pane := range window.panes {
				pNode := tview.NewTreeNode(getPaneDisplayName(pane, false))
				// pNode.SetReference([]int{session.id, wInd, pInd})
				pNode.SetReference(pane)
				pNode.SetIndent(3)
				wNode.AddChild(pNode)
			}
			sNode.AddChild(wNode)
		}
		if sInd == 0 {
			sessionDisplay.SetCurrentNode(sNode)
		}
		if session.isAttached {
			sNode.SetColor(tcell.ColorLimeGreen)
		}
		root.AddChild(sNode)
	}

	sessionDisplay.
		SetGraphics(false).
		SetPrefixes([]string{" ", treeSignUpEnding, treeSignUpMiddle})

	sessionDisplay.Box.SetBorder(false)

	initSessionDisplayKeys()

	// Invoked on hover
	sessionDisplay.SetChangedFunc(func(node *tview.TreeNode) {
		runCallbacksForNode(node, handleChangeSession, handleChangeWindow, handleChangePane)
	})
	// sessionDisplay.SetBackgroundColor(tcell.ColorDefault)

	// Invoked on selected
	// sessionDisplay.SetSelectedFunc(func(node *tview.TreeNode) {
	// node.SetExpanded(!node.IsExpanded())
	// node.SetColor(tcell.ColorRed)
	// })
}
func refreshSessionDisplay() {
	sessionData = getSessionData()
	initSessionDisplay()
}

func initPreviewDisplay() {
	previewDisplay.
		SetDynamicColors(true)
	// SetChangedFunc(func() {
	// 	tviewApp.Draw()
	// })
	previewDisplay.SetBackgroundColor(tcell.ColorDefault)
	previewDisplay.SetBorder(true).SetTitle(" " + sessionData.sessions[0].name).SetTitleAlign(0)
	previewDisplay.SetText(tview.TranslateANSI(sessionData.sessions[0].preview))
	previewDisplay.ScrollToBeginning()
}

func initHelpBoxDisplay() {
	helpBoxDisplay.SetBorder(true)
	helpBoxDisplay.SetBorderColor(tcell.ColorGreen)
	helpBoxDisplay.SetTitle("KEYS")
	helpBoxDisplay.SetTitleAlign(tview.AlignLeft)
	// Create full guide for users to view from app
	helpTextList := [][]string{
		{"", "arrow keys to navigate"},
		{"", "hjkl to navigate"},
		{"", ""},
		{"c", "create session"},
		{"x", "delete target"},
		{"?", "toggle help"},
		{"/", "search | Ctrl-s"},
	}

	for ind, rowList := range helpTextList {
		helpBoxDisplay.SetCell(ind, 0, tview.NewTableCell(rowList[0]).SetTextColor(tcell.ColorLightGreen))
		helpBoxDisplay.SetCell(ind, 1, tview.NewTableCell(rowList[1]))
	}
}

func initSearchBoxDisplay() {
	searchBoxDisplay.
		SetLabel("Search: ").
		SetFieldWidth(32)
	searchBoxDisplay.SetBorder(true)
	searchBoxDisplay.Box.SetBorderPadding(1, 1, 1, 1)
	// TODO: DOT . navigation to children of highlighted nodes
	// ----- SearchTerm.0.10 -> should accept any number child
	searchBoxDisplay.SetChangedFunc(func(text string) {

		// sessionDisplay.GetRoot().CollapseAll()

		// parts := strings.Split(text, ".")
		// currentNode := sessionDisplay.GetCurrentNode()
		// sessionChildren := currentNode.GetChildren()

		if strings.Contains(text, ".") {
			highlightNodeDotPath(text)
		} else {
			sNodes := sessionDisplay.GetRoot().GetChildren()
			matchedPrefix := -1
			matchedContain := -2

			for i := len(sNodes) - 1; i >= 0; i-- {
				sNode := sNodes[i]
				// sNodeIDArr := sNode.GetReference().([]int)
				session := sNode.GetReference().(Session)
				// session := getSessionFromRef(sNodeIDArr)
				if text == "" {
					return
				}

				matchStart, _ := regexp.MatchString("(?i)^"+text, session.name)
				matchContain := strings.Contains(session.name, text)
				if matchStart {
					matchedPrefix = session.id
				} else if matchContain {
					matchedContain = session.id
				}
			}

			if matchedPrefix != -1 {
				highlightSessionNode(matchedPrefix)
			} else if matchedContain != -1 {
				highlightSessionNode(matchedContain)
			}
		}
	})

	initSearchDisplayKeys()
}

func getSliceOfDirectories(search string, currentPath string, files []os.FileInfo) []string {
	dirs := []string{}
	if search == "" {
		for _, file := range files {
			if file.IsDir() && len(file.Name()) > 0 {
				dirs = append(dirs, currentPath+file.Name()+"/")
			}
		}
		return dirs
	} else {
		matchedPrefixes := []string{}
		matchedGeneral := []string{}

		for _, file := range files {
			if file.IsDir() && len(file.Name()) > 0 {
				fName := strings.ToLower(file.Name())
				sTerm := strings.ToLower(search)
				if strings.HasPrefix(fName, sTerm) {
					matchedPrefixes = append(matchedPrefixes, currentPath+file.Name())
				} else if strings.Contains(fName, sTerm) {
					matchedGeneral = append(matchedGeneral, currentPath+file.Name())
				}
			}
		}

		if len(matchedPrefixes) > 0 {
			firstTabCompletionTerm = matchedPrefixes[0] + "/"
		} else if len(matchedGeneral) > 0 {
			firstTabCompletionTerm = matchedGeneral[0] + "/"
		}
		return append(matchedPrefixes, matchedGeneral...)
	}
}

func initDirectoryInputField() {
	directoryInputField.SetText("~/")
	directoryInputField.SetLabel("Directory")
	directoryInputField.SetFieldWidth(30)

	// Enables selection of first autocomplte term on ENTER
	directoryInputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			directoryInputField.SetText(firstTabCompletionTerm)
			// -- tab key cannot be overriden --
			// case tcell.KeyTab:
			// 	text := directoryInputField.GetText()
			// 	if text[len(text)-1:] != "/" {
			// 		directoryInputField.SetText(text + "/")
			// 	}
		}
		return event
	})

	// TODO: curenctly fails on /App
	directoryInputField.SetAutocompleteFunc(func(text string) (entries []string) {
		if len(text) == 0 || directoryInputField.HasFocus() == false {
			return
			// } // if only 1 "/" found, then run this block
			// else if len(text) < 2 && text[0:1] == "/" {
			// 	files, _ := ioutil.ReadDir(text)
			// 	return getSliceOfDirectories("", text, files)
		} else if len(text) > 1 && text[0:2] == "~/" {
			home, homeDirErr := os.UserHomeDir()
			if homeDirErr != nil {
				handleExecError(homeDirErr, "initDirectoryInputField")
			}
			text = home + "/" + text[2:]
		}
		dirParts := strings.Split(text, "/")
		directory := dirParts[0]
		searchTerm := ""
		if len(dirParts) > 1 {
			searchTerm = dirParts[len(dirParts)-1]
			directory = strings.Join(dirParts[0:len(dirParts)-1], "/")
		}
		files, _ := ioutil.ReadDir(directory)
		return getSliceOfDirectories(searchTerm, directory+"/", files)
	})
}

func createNewSession() {
	sName := creationFormName
	pCount := creationFormPaneCount
	sDir := directoryInputField.GetText()
	if sDir == "" {
		homeDir, homeDirErr := os.UserHomeDir()
		handleExecError(homeDirErr, "createNewSession")
		sDir = homeDir
	}
	if pCount == 1 {
		previewDisplay.SetTitle(sName + " " + sDir)
	}
	createTmuxSession(creationFormName, directoryInputField.GetText(), creationFormPaneCount)
	endCreateSession()
	time.Sleep(200 * time.Millisecond)
	refreshSessionDisplay()
}
func switchTmuxToCurrentNode() {
	tviewApp.Stop()

	tmux := sessionDisplay.GetCurrentNode().GetReference()
	if sessionData.hasAttachedSession {
		switch tmux.(type) {
		case Session:
			switchToTmuxPath(tmux.(Session).name)
		case Window:
			switchToTmuxPath(tmux.(Window).path)
		case Pane:
			switchToTmuxPath(tmux.(Pane).path)
		}
	} else {
		switch tmux.(type) {
		case Session:
			attachTmuxSession(tmux.(Session).name)
		case Window:
			attachTmuxSession(tmux.(Window).session)
		case Pane:
			attachTmuxSession(tmux.(Pane).session)
		}
	}
}
func handleKillSession(session Session, node *tview.TreeNode) {
	killTmuxSession(session.name)
}
func handleKillWindow(window Window, node *tview.TreeNode) {
	killTmuxWindow(window.path)
}
func handleKillPane(pane Pane, node *tview.TreeNode) {
	killTmuxPane(pane.path)
}
func confirmKillSession(session Session, node *tview.TreeNode) {
	confirmKillTarget("Session", session.name)
}
func confirmKillWindow(window Window, node *tview.TreeNode) {
	confirmKillTarget("Window", window.index)
}
func confirmKillPane(pane Pane, node *tview.TreeNode) {
	confirmKillTarget("Pane", pane.index)
}
func killTmuxTarget(node *tview.TreeNode, killTarget bool) {
	if killTarget {
		runCallbacksForNode(node, handleKillSession, handleKillWindow, handleKillPane)
		time.Sleep(10 * time.Millisecond)
		refreshSessionDisplay()
		flexBoxWrapper.HidePage("confirmModal")
	} else {
		targetedNode = node
		runCallbacksForNode(node, confirmKillSession, confirmKillWindow, confirmKillPane)
	}
}
func confirmKillTarget(targetType string, targetName string) {
	confirmModal.SetText("Kill " + targetType + " " + targetName + "?")
	confirmModal.SetFocus(0)
	flexBoxWrapper.ShowPage("confirmModal")
}

func changeViewTo(view tview.Primitive) {
	tviewApp.SetRoot(view, true)
	tviewApp.SetFocus(view)
}

func restoreDefaultView() {
	tviewApp.SetRoot(flexBoxWrapper, true)
	tviewApp.SetFocus(sessionDisplay)
}

func detachSession(node *tview.TreeNode) {
	tmux := node.GetReference()
	switch tmux.(type) {
	case Session:
		if tmux.(Session).isAttached {
			tviewApp.Stop()
			time.Sleep(10 * time.Millisecond)
			detachTmuxSession(tmux.(Session).name)
		}
	}
}

func initCreationForm() {
	creationForm.
		AddInputField("Name", "", 20, nil, func(n string) {
			creationFormName = n
		}).
		AddFormItem(directoryInputField).
		// AddDropDown("Panes", []string{"1", "2", "3", "4"}, 0, func(option string, ind int) {
		// 	creationFormPaneCount = ind + 1
		// }).
		AddButton("Save", createNewSession).
		AddButton("Cancel", endCreateSession)
	creationForm.SetBorder(true).SetTitle(" Create New Session | ESC to cancel ").SetTitleAlign(tview.AlignLeft)

	initCreationFormKeys()
}

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
}

func initInterface() {
	initSessionDisplay()
	initPreviewDisplay()
	initHelpBoxDisplay()
	initSearchBoxDisplay()
	initDirectoryInputField()
	initConfirmModal()

	flexBoxDisplay.
		AddItem(helpBoxDisplay, 0, 0, false). // FOR HELP MENU
		AddItem(mainFlexBoxView.SetDirection(tview.FlexRow).
			AddItem(confirmModal, 0, 0, false).
			AddItem(sessionDisplay, 0, 6, true).
			AddItem(searchBoxDisplay, 0, 0, false).
			AddItem(previewDisplay, 0, 3, false), 0, 5, true)

	// IMPORTANT - Workaround - Without this Clear functionality breaks
	tviewApp.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})

	flexBoxWrapper.
		AddPage("flexBoxDisplay", flexBoxDisplay, true, true).
		AddPage("confirmModal", confirmModal, false, false)

	flexBoxWrapper.SetBackgroundColor(tcell.ColorDefault)

	if err := tviewApp.SetRoot(flexBoxWrapper, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func startApp() {
	initInterface()
}

var app = cli.App{
	Action: func(c *cli.Context) error {
		return nil
	},
}

func setupCliApp() {
	app.Name = "Tmux Session Interface"
	app.Usage = "A feature packed tmux tree for managing sessions, windows, and panes."
	app.Author = "Gabe Schnall"

	app.Commands = []cli.Command{}
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:     "i",
			Usage:    "Initiate interface in current terminal. Without this, tmint will attempt to fullscreen the app in a new window.",
			Required: false,
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.Bool("i") {
			fmt.Println(sessionData)
		} else {
			if sessionData.hasAttachedSession && !sessionData.hasZoomedPane {
				tmuxToggleFullscreen()
			}
			isUserPaneZoomed = sessionData.hasZoomedPane
			startApp()
		}
		return nil
	}
}

func main() {
	setupCliApp()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	// On data refresh this will always be true - have to avoid setting it again
	if !isUserPaneZoomed {
		tmuxToggleFullscreen()
	}
}
