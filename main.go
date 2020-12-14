package main
import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

  twiz "github.com/gschnall/tmint/tmux_wizard"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/urfave/cli"
)

// Views and their state variables
var (
	isUserPaneZoomed       = false
	sessionData            = twiz.GetSessionData()
	tviewApp               = tview.NewApplication()
	previewDisplay         = tview.NewTextView()
	sessionDisplay         = tview.NewTreeView()
	targetedNode           *tview.TreeNode
	renameForm             = tview.NewForm()
	renameFormName         = "" 
	creationForm           = tview.NewForm()
	creationFormName       = ""
	creationFormSession    = ""
	creationFormPaneCount  = 0
	helpBoxDisplay         = tview.NewTable()
	helpBoxDisplayWidth    = 0
	searchBoxDisplay       = tview.NewInputField()
	searchBoxDisplayHeight = 0
	confirmModal           = tview.NewModal()
	directoryInputField    = tview.NewInputField()
	noActiveSessionDisplay = tview.NewTextView()
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

func getHeader() string {
	h := "\n\n   :  :  : :   ::  ::  ::  :: :: :\n"
	h += "   |___-_|-__-|-__|--_|--__--__--_--_-|\n"
	h += "   -----| Tmint - a Tumux interface"
	h += " > > >\n"
	h += "   |___-_|-__-|-__|--_|--__--__--_--_-|\n"
	return h
}

func getFormattedSessionName(name string, isAttached bool) string {
	if isAttached {
		return name + " (attached)"
	}
	return name
}

func formatSessionId(id int) string {
	if id < 10 {
		return "(" + strconv.Itoa(id) + ")  "
	} else if id < 36 {
		letter := twiz.ToCharStr(id - 9)
		return "(M-" + letter + ")"
	}
	return ""
}

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

func expandAllChildNodes(node *tview.TreeNode) {
	node.ExpandAll()
}

func expandNode(node *tview.TreeNode, moveNodes bool) {
	if !node.IsExpanded() {
		node.SetExpanded(true)
	}
	if moveNodes {
		sessionDisplay.SetCurrentNode(getNextNodeInTree(node))
	}
	runCallbacksForNode(node, toggleSessionName, toggleWindowName, togglePaneName)
}
func collapseNode(node *tview.TreeNode, moveNodes bool) {
	node.Collapse()

	if moveNodes {
		sessionDisplay.SetCurrentNode(getPreviousNodeInTree(node))
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
func startCreateForm(creation string) {
	initCreationForm(creation)
	changeViewTo(creationForm)
}
func endCreateSession() {
	creationForm.Clear(true)
	restoreDefaultView()
}
func endRenameForm() {
	renameForm.Clear(true)
	restoreDefaultView()
}

func highlightNode(index int) {
	node := sessionDisplay.GetRoot().GetChildren()[index]
	sessionDisplay.SetCurrentNode(node)
}

func highlightSessionNode(index int) {
	if index > -1 && index < len(sessionData.Sessions) {
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
func initRenameFormKeys() {
	renameForm.SetCancelFunc(func() {
		endRenameForm()
	})
}

func initNoActiveSessionDisplayKeys() {
	noActiveSessionDisplay.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			tviewApp.Stop()
			twiz.StartTmux()
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

func initSessionDisplayKeys() {
	sessionDisplay.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			expandNode(sessionDisplay.GetCurrentNode(), false)
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
				expandNode(sessionDisplay.GetCurrentNode(), true)
			}
		case 'd':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				detachSession(sessionDisplay.GetCurrentNode())
			}
		case 'e':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				cn := sessionDisplay.GetCurrentNode()
				if cn.IsExpanded() {
					cn.CollapseAll()
				} else {
					cn.ExpandAll()
				}
			}
		case '+', '=':
			for _, child := range sessionDisplay.GetRoot().GetChildren() {
				child.ExpandAll()
			}
		case '-':
			for _, child := range sessionDisplay.GetRoot().GetChildren() {
				child.CollapseAll()
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
		case 'w':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				startCreateForm("window")
			}
		case 'x':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				killTmuxTarget(sessionDisplay.GetCurrentNode(), false)
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			highlightSessionNode(int(event.Rune() - '0'))
		case 's':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				startCreateForm("session")
			}
		case 'r':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				startRenameForm()
			}

		case 'a', 'b', 'c', 'f', 'g', 'i', 'j', 'k', 'm', 'n', 'o', 'p', 't', 'u', 'v', 'y', 'z':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			}
		}

		return event
	})
}

func getSessionDisplayName(session twiz.Session, isExpanded bool) string {
	mark := "+ "
	if isExpanded {
		mark = "- "
	}
	return formatSessionId(session.Id) + mark + getFormattedSessionName(session.Name, session.IsAttached)
}

func getWindowDisplayName(window twiz.Window, isExpanded bool) string {
	if isExpanded {
		return " - " + window.Index + ": " + window.Name + " (" + strconv.Itoa(len(window.Panes)) + " panes)"
	}
	return " + " + window.Index + ": " + window.Name + " (" + strconv.Itoa(len(window.Panes)) + " panes)"
}
func getPaneDisplayName(pane twiz.Pane) string {
	return pane.Name + ": " + pane.Command
}

func toggleSessionName(session twiz.Session, node *tview.TreeNode) {
	node.SetText(getSessionDisplayName(session, node.IsExpanded()))
}
func toggleWindowName(window twiz.Window, node *tview.TreeNode) {
	node.SetText(getWindowDisplayName(window, node.IsExpanded()))
}
func togglePaneName(pane twiz.Pane, node *tview.TreeNode) {
	node.SetText(getPaneDisplayName(pane))
}

func getSessionFromRef(ids []int) twiz.Session {
	return sessionData.Sessions[ids[0]]
}
func getWindowFromRef(ids []int) twiz.Window {
	return sessionData.Sessions[ids[0]].Windows[ids[1]]
}
func getPaneFromRef(ids []int) twiz.Pane {
	return sessionData.Sessions[ids[0]].Windows[ids[1]].Panes[ids[2]]
}

type SessionFunc func(twiz.Session, *tview.TreeNode)
type WindowFunc func(twiz.Window, *tview.TreeNode)
type PaneFunc func(twiz.Pane, *tview.TreeNode)

func runCallbacksForNode(node *tview.TreeNode, sfunc SessionFunc, wfunc WindowFunc, pfunc PaneFunc) {
	tmux := node.GetReference()
	switch tmux.(type) {
	case twiz.Session:
		sfunc(tmux.(twiz.Session), node)
	case twiz.Window:
		wfunc(tmux.(twiz.Window), node)
	case twiz.Pane:
		pfunc(tmux.(twiz.Pane), node)
	}
}

func initSessionDisplay() {
	root := tview.NewTreeNode("").
		SetSelectable(false)
	sessionDisplay.
		SetRoot(root).
		SetTopLevel(1)

	// - BUG - https://github.com/rivo/tview/issues/314
	// SetBackgroundColor(tcell.ColorDefault)
	// --> Can't hide modals and highlighted colors seem strange

	for sInd, session := range sessionData.Sessions {
		sNode := tview.NewTreeNode(getSessionDisplayName(session, false))
		sNode.SetSelectable(true).SetExpanded(false)
		sNode.SetReference(session)
		for _, window := range session.Windows {
			wNode := tview.NewTreeNode(getWindowDisplayName(window, false))
			wNode.SetReference(window)
			wNode.SetExpanded(false)
			wNode.SetSelectable(true)
			wNode.SetIndent(6)
			for _, pane := range window.Panes {
				pNode := tview.NewTreeNode(getPaneDisplayName(pane))
				pNode.SetReference(pane)
				pNode.SetIndent(3)
				wNode.AddChild(pNode)
			}
			sNode.AddChild(wNode)
		}
		if sInd == 0 {
			sessionDisplay.SetCurrentNode(sNode)
		}
		if session.IsAttached {
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
}
func refreshSessionDisplay() {
	sessionData = twiz.GetSessionData()
	initSessionDisplay()
}

func initPreviewDisplay() {
	previewDisplay.SetDynamicColors(true)
	previewDisplay.SetBackgroundColor(tcell.ColorDefault)
	previewDisplay.SetBorder(true).SetTitle(" " + sessionData.Sessions[0].Name).SetTitleAlign(0)
	previewDisplay.SetText(tview.TranslateANSI(sessionData.Sessions[0].Preview))
	previewDisplay.ScrollToBeginning()
}

func initHelpBoxDisplay() {
	helpBoxDisplay.SetBorder(true)
	helpBoxDisplay.SetBorderColor(tcell.ColorGreen)
	helpBoxDisplay.SetTitle("KEYS")
	helpBoxDisplay.SetTitleAlign(tview.AlignLeft)

	helpTextList := [][]string{
		{"", "arrow keys to navigate"},
		{"", "or hjkl to navigate"},
		{"", ""},
		{"/", "search | Ctrl-s"},
		{"s", "create session"},
		{"w", "create window"},
		{"x", "delete target"},
		{"r", "rename target"},
		{"e", "toggle expand"},
		{"?", "toggle help"},
		{"+", "toggle expand all"},
		{"-", "toggle collapse all"},
	}

	for ind, rowList := range helpTextList {
		helpBoxDisplay.SetCell(ind, 0, tview.NewTableCell(rowList[0]).SetTextColor(tcell.ColorLightGreen))
		helpBoxDisplay.SetCell(ind, 1, tview.NewTableCell(rowList[1]))
	}
}

func getPreviousNodeInTree(node *tview.TreeNode) *tview.TreeNode {
	tmux := node.GetReference()

	tmuxType := "session"	
	switch tmux.(type) {
	case twiz.Window:
		tmuxType = "window" 
	case twiz.Pane:
		tmuxType = "pane" 
	}

	parent := getParentOfNode(node)
	children := parent.GetChildren()

	for index, child := range children {
		if child == node {
			if tmuxType == "session" {
				if index == 0 {
					return node
				}
				return getLastVisibleChildNodeInSession(children[index-1])
			}

			if tmuxType == "window" && index > 0 {
				previousSib := getPreviousSibling(node)				
				previousSibChildren := previousSib.GetChildren()
				return previousSibChildren[len(previousSibChildren)-1]
			}

			if tmuxType == "pane" && index > 0 {
				return getPreviousSibling(node)
			}

			return getParentOfNode(node) 
		}
	}
	return node
}
func getLastVisibleChildNodeInSession(node *tview.TreeNode) *tview.TreeNode {
	sessionChildren := node.GetChildren()
	lastWindow := sessionChildren[len(sessionChildren)-1]
	windowChildren := lastWindow.GetChildren()
	lastPane := windowChildren[len(windowChildren)-1]

	if lastWindow.IsExpanded() {
		return lastPane
	} else if node.IsExpanded() {
		return lastWindow
	} 
	return node 
}

func getNextNodeInTree(node *tview.TreeNode) *tview.TreeNode {
	tmux := node.GetReference()

	tmuxType := "session"	
	switch tmux.(type) {
	case twiz.Window:
		tmuxType = "window" 
	case twiz.Pane:
		tmuxType = "pane" 
	}

	parent := getParentOfNode(node)
	children := parent.GetChildren()

	for index, child := range children {
		if child == node {
			if index == len(children)-1 && tmuxType == "pane" {
				// need to check if parent window has a next sibling
				nextWindowSibling := getNextSibling(parent)
				if nextWindowSibling != nil {
					return nextWindowSibling
				}
				nextSessionSibling := getNextSibling(getSessionFromNode(parent))
				if nextSessionSibling != nil {
					return nextSessionSibling
				}
				return node
			}

			if tmuxType == "pane" {
				return children[index+1]
			} else if tmuxType == "window" || tmuxType == "session" {
				return node.GetChildren()[0]
			} 
		}
	}
	return node
}

func getNextSibling(node *tview.TreeNode) *tview.TreeNode {
	parent := getParentOfNode(node)	
	children := parent.GetChildren()
	for index, child := range children {
		if child == node && index < len(children)-1 {
			return children[index + 1]
		}
	}
	return nil
}
func getPreviousSibling(node *tview.TreeNode) *tview.TreeNode {
	parent := getParentOfNode(node)	
	children := parent.GetChildren()
	for index, child := range children {
		if child == node {
			return children[index - 1]
		}
	}
	return nil
}

func getSessionFromNode(node *tview.TreeNode) *tview.TreeNode {
	tmux := node.GetReference()
	switch tmux.(type) {
	case twiz.Session:
		return node
	case twiz.Window:
		return getParentOfNode(node)
	case twiz.Pane:
		windowNode := getParentOfNode(node)
		return getParentOfNode(windowNode)
	}
	return nil
}

func getParentOfNode(node *tview.TreeNode) *tview.TreeNode {
	root  := sessionDisplay.GetRoot()
	queue := make([]*tview.TreeNode, 0)
	queue = append(queue, root)

	for len(queue) > 0 {
		parent := queue[0]	
		queue = queue[1:]


		nextUpChildren := parent.GetChildren()
		if len(nextUpChildren) > 0 {
			for _, child := range nextUpChildren {
				if child == node {
					return parent
				}
				queue = append(queue, child)
			}
		}
	}

	return nil
}

func breadthFirstSearch(root *tview.TreeNode, search string) *tview.TreeNode {
	if search == "" {
		return sessionDisplay.GetCurrentNode() 
	}
	search = strings.ToLower(search)

	queue := make([]*tview.TreeNode, 0)
	queue = append(queue, root)

	matchContains := []*tview.TreeNode{}

	for len(queue) > 0 {
		nextUp := queue[0]
		queue = queue[1:]

		nodeName := ""

		tmux := nextUp.GetReference()
		switch tmux.(type) {
		case twiz.Session:
			nodeName = strings.ToLower(getSessionDisplayName(tmux.(twiz.Session), nextUp.IsExpanded()))
		case twiz.Window:
			nodeName = strings.ToLower(getWindowDisplayName(tmux.(twiz.Window), nextUp.IsExpanded()))
		case twiz.Pane:
			// nodeName = strings.ToLower(tmux.(twiz.Pane).Name)
			nodeName = strings.ToLower(getPaneDisplayName(tmux.(twiz.Pane)))
		}

		matchStart, _ := regexp.MatchString("(?i)^"+search, nodeName)
		if matchStart {
			return nextUp
		} else if strings.Contains(nodeName, search) {
			matchContains = append(matchContains, nextUp)	
		}

		nextUpChildren := nextUp.GetChildren()
		if len(nextUpChildren) > 0 {
			for _, child := range nextUpChildren {
				queue = append(queue, child)
			}
		}
	}

	if len(matchContains) > 0 {
		return matchContains[0]
	}
	return nil
}

func initSearchBoxDisplay() {
	searchBoxDisplay.
		SetLabel("Search: ").
		SetFieldWidth(32)
	searchBoxDisplay.SetBorder(true)
	searchBoxDisplay.Box.SetBorderPadding(1, 1, 1, 1)
	searchBoxDisplay.SetChangedFunc(func(text string) {
		searchedForNode := breadthFirstSearch(sessionDisplay.GetRoot(), text)
		if searchedForNode != nil {
			sessionDisplay.SetCurrentNode(searchedForNode)
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
			// -- Tab key could not be overriden - Attempted June 1st 2020 --
			// case tcell.KeyTab:
			// 	text := directoryInputField.GetText()
			// 	if text[len(text)-1:] != "/" {
			// 		directoryInputField.SetText(text + "/")
			// 	}
			// -- --------------------------- --
		}
		return event
	})

	directoryInputField.SetAutocompleteFunc(func(text string) (entries []string) {
		if len(text) == 0 || directoryInputField.HasFocus() == false {
			return
		} else if len(text) > 1 && text[0:2] == "~/" {
			home, homeDirErr := os.UserHomeDir()
			if homeDirErr != nil {
				twiz.HandleExecError(homeDirErr, "initDirectoryInputField")
			}
			text = home + "/" + text[2:]
		}
		dirParts := strings.Split(text, "/")
		directory := dirParts[0]
		searchTerm := ""
		if dirParts[0] == "" && len(dirParts) <= 2 { //Fix for directories with leading "/" char
			searchTerm = dirParts[len(dirParts)-1]
			directory = "/" + strings.Join(dirParts[0:len(dirParts)-1], "/")
		} else if len(dirParts) > 1 {
			searchTerm = dirParts[len(dirParts)-1]
			directory = strings.Join(dirParts[0:len(dirParts)-1], "/")
		}
		files, _ := ioutil.ReadDir(directory)

		if directory == "/" { //Fix for directories with leading "/" char
			return getSliceOfDirectories(searchTerm, directory, files)
		}
		return getSliceOfDirectories(searchTerm, directory+"/", files)
	})
}

func tmuxNameIsValid(name string, lastChar rune) bool {
	return !(strings.Contains(name, ".") || strings.Contains(name, ":"))
}

// Add regex to prevent invalid characters
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

// Add regex to prevent invalid characters
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
func switchTmuxToCurrentNode() {
	tviewApp.Stop()

	tmux := sessionDisplay.GetCurrentNode().GetReference()
	if sessionData.HasAttachedSession {
		switch tmux.(type) {
		case twiz.Session:
			twiz.SwitchToTmuxPath(tmux.(twiz.Session).Name)
		case twiz.Window:
			twiz.SwitchToTmuxPath(tmux.(twiz.Window).Path)
		case twiz.Pane:
			twiz.SwitchToTmuxPath(tmux.(twiz.Pane).Path)
		}
	} else {
		switch tmux.(type) {
		case twiz.Session:
			twiz.AttachTmuxSession(tmux.(twiz.Session).Name)
		case twiz.Window:
			twiz.AttachTmuxSession(tmux.(twiz.Window).Session)
		case twiz.Pane:
			twiz.AttachTmuxSession(tmux.(twiz.Pane).Session)
		}
	}
}
func handleKillSession(session twiz.Session, node *tview.TreeNode) {
	twiz.KillTmuxSession(session.Name)
}
func handleKillWindow(window twiz.Window, node *tview.TreeNode) {
	twiz.KillTmuxWindow(window.Path)
}
func handleKillPane(pane twiz.Pane, node *tview.TreeNode) {
	twiz.KillTmuxPane(pane.Path)
}
func confirmKillSession(session twiz.Session, node *tview.TreeNode) {
	confirmKillTarget("Session", session.Name)
}
func confirmKillWindow(window twiz.Window, node *tview.TreeNode) {
	confirmKillTarget("Window", window.Index)
}
func confirmKillPane(pane twiz.Pane, node *tview.TreeNode) {
	confirmKillTarget("Pane", pane.Index)
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
	if sessionData.HasLivingSessions == false {
		changeViewTo(noActiveSessionDisplay)
	} else {
		tviewApp.SetRoot(flexBoxWrapper, true)
		tviewApp.SetFocus(sessionDisplay)
	}
}

func detachSession(node *tview.TreeNode) {
	tmux := node.GetReference()
	switch tmux.(type) {
	case twiz.Session:
		if tmux.(twiz.Session).IsAttached {
			tviewApp.Stop()
			time.Sleep(10 * time.Millisecond)
			twiz.DetachTmuxSession(tmux.(twiz.Session).Name)
		}
	}
}

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

	initRenameFormKeys()
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
		// -- Possibly a future feature: let users select # of panes --
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
	creationForm.SetBorder(true).SetTitle(fmt.Sprintf(" Create New %s | ESC to cancel | Ctrl-u to clear input ", creation)).SetTitleAlign(tview.AlignLeft)

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
	// --- Need to create custom modal to capture different keypresses
}

func initNoActiveSessionDisplay() {
	noActiveSessionDisplay.SetTitle("Tmux has not been started")
	noActiveSessionDisplay.SetBorder(true)
  // Displayed Text
	headerText := getHeader() + "\n\n"
	info   := "   No active sessions available\n\n\n"
	start  := "   Press ENTER or 't' to start Tmux\n\n"
	create := "   ----- 's' to create a session\n\n"
	quit   := "   ----- 'q' to quit" 
	noActiveSessionDisplay.SetText(headerText + info + start + create + quit)
	// Keys
	initNoActiveSessionDisplayKeys()
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

	// Workaround - Without this, Clear functionality breaks
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

func initNoActiveSessionInterface() {
  // - Allow users to create new session	
	initConfirmModal()
	initNoActiveSessionDisplay()

	// -------- This is problaby inactive code - look to delete ---------
	flexBoxDisplay.
		AddItem(helpBoxDisplay, 0, 0, false). // FOR HELP MENU
		AddItem(mainFlexBoxView.SetDirection(tview.FlexRow).
			AddItem(confirmModal, 0, 0, false).
			AddItem(noActiveSessionDisplay, 0, 3, false), 0, 5, true)

	tviewApp.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})

	flexBoxWrapper.
		AddPage("flexBoxDisplay", flexBoxDisplay, true, true).
		AddPage("confirmModal", confirmModal, false, false)

	flexBoxWrapper.SetBackgroundColor(tcell.ColorDefault)
	// -------- ------------------------------- -------------------------

	if err := tviewApp.SetRoot(noActiveSessionDisplay, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func startApp() {
	if sessionData.HasLivingSessions == false {
		initNoActiveSessionInterface()
	} else {
		initInterface()
	}
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
			if sessionData.HasAttachedSession && !sessionData.HasZoomedPane {
				twiz.TmuxToggleFullscreen()
			}
			isUserPaneZoomed = sessionData.HasZoomedPane
			startApp()
		}
		return nil
	}
}

func testOutputInHelpBox(s string) {
	helpBoxDisplay.SetTitle(s)
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
		twiz.TmuxToggleFullscreen()
	}
}
