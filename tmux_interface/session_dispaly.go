package tmux_interface

import (
	"strconv"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	twiz "github.com/gschnall/tmint/tmux_wizard"
)

var (
	sessionDisplay = tview.NewTreeView()
)

const (
	treeSignUpMiddle = "├── "
	treeSignUpEnding = "└─"
)

func initSessionDisplay() {
	root := tview.NewTreeNode("").
		SetSelectable(false)
	sessionDisplay.
		SetRoot(root).
		SetTopLevel(1).
		SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	// BUG - https://github.com/rivo/tview/issues/314
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
			sNode.SetColor(tview.Styles.TertiaryTextColor)
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

// _____________________
// |                   |
// | Utils and Actions |
// |                   |
// ---------------------
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

func refreshSessionDisplay() {
	result := make(chan twiz.SessionData, 1)
	go twiz.GetSessionData(sessionData.AttachedSession, sessionData.TmintSession, sessionData.IsUsingKeybindings, result)
	dataResult := <-result
	sessionData = dataResult
	close(result)
	initSessionDisplay()
}

// To Do: Deprecated - rework thsese into case statements
func toggleSessionName(session twiz.Session, node *tview.TreeNode) {
	node.SetText(getSessionDisplayName(session, node.IsExpanded()))
}
func toggleWindowName(window twiz.Window, node *tview.TreeNode) {
	node.SetText(getWindowDisplayName(window, node.IsExpanded()))
}
func togglePaneName(pane twiz.Pane, node *tview.TreeNode) {
	node.SetText(getPaneDisplayName(pane))
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

func highlightNode(index int) {
	node := sessionDisplay.GetRoot().GetChildren()[index]
	sessionDisplay.SetCurrentNode(node)
}

func highlightSessionNode(index int) {
	if index > -1 && index < len(sessionData.Sessions) {
		highlightNode(index)
	}
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

func switchTmuxToCurrentNode(node *tview.TreeNode) {
	tmux := node.GetReference()
	pathName := ""
	switch tmux.(type) {
	case twiz.Session:
		pathName = tmux.(twiz.Session).Name
	case twiz.Window:
		pathName = tmux.(twiz.Window).Path
	case twiz.Pane:
		pathName = tmux.(twiz.Pane).Path
	}

	if sessionData.HasAttachedSession {
		// for tmux -t keybinding workflow
		sessionData.AttachedSession = pathName
		tviewApp.Stop()
		twiz.SwitchToTmuxPath(pathName)
	} else {
		tviewApp.Stop()
		twiz.AttachTmuxSession(pathName)
	}
}

// To Do: Deprecated - rework thsese into case statements
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
		changeViewTo(inactiveSessionDisplay)
	} else {
		tviewApp.SetRoot(flexBoxWrapper, true)
		tviewApp.SetFocus(sessionDisplay)
	}
}

func detachSession(sessionName string) {
	tviewApp.Stop()
	twiz.DetachTmuxSession(sessionName)
}

// _______________
// |             |
// | Keybindings |
// |             |
// ---------------
func initSessionDisplayKeys() {
	sessionDisplay.SetSelectedFunc(func(node *tview.TreeNode) {
		switchTmuxToCurrentNode(node)
	})

	sessionDisplay.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			expandNode(sessionDisplay.GetCurrentNode(), false)
		case tcell.KeyLeft:
			collapseNode(sessionDisplay.GetCurrentNode(), false)
		case tcell.KeyCtrlS:
			toggleSearchBox()
		case tcell.KeyEsc:
			tviewApp.Stop()
		}

		switch event.Rune() {
		case '/':
			toggleSearchBox()
		case 'l':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				// -- Fixes vim navigation --
				// l & h don't trigger a "changed" event for some reason
				currentNode := sessionDisplay.GetCurrentNode()
				expandNode(currentNode, true)
				runCallbacksForNode(currentNode, handleChangeSession, handleChangeWindow, handleChangePane)
			}
		case 'h':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				// -- Fixes vim navigation --
				// l & h don't trigger a "changed" event for some reason
				currentNode := sessionDisplay.GetCurrentNode()
				collapseNode(sessionDisplay.GetCurrentNode(), true)
				runCallbacksForNode(currentNode, handleChangeSession, handleChangeWindow, handleChangePane)
			}
		case 'c':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				startCheatSheet()
			}
		case 'd':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				// Work in progress
				// - Detach should use the current session (not node)
				// - Needs a way to cleanly exit before detaching
				detachSession(sessionData.AttachedSession)
			}
		case 'e':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				cn := sessionDisplay.GetCurrentNode()
				if cn.IsExpanded() {
					cn.Walk(func(node *tview.TreeNode, parent *tview.TreeNode) bool {
						node.Collapse()
						runCallbacksForNode(node, toggleSessionName, toggleWindowName, togglePaneName)
						return true
					})
				} else {
					cn.Walk(func(node *tview.TreeNode, parent *tview.TreeNode) bool {
						node.Expand()
						runCallbacksForNode(node, toggleSessionName, toggleWindowName, togglePaneName)
						return true
					})
				}
			}
		case '+', '=':
			sessionDisplay.GetRoot().Walk(func(node *tview.TreeNode, parent *tview.TreeNode) bool {
				node.Expand()
				runCallbacksForNode(node, toggleSessionName, toggleWindowName, togglePaneName)
				return true
			})
		case '-':
			for _, child := range sessionDisplay.GetRoot().GetChildren() {
				child.CollapseAll()
				runCallbacksForNode(child, toggleSessionName, toggleWindowName, togglePaneName)
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
		case 'v':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				// New function that'll allow user to select where to save captured scrollback
				// What would you like to do with the scrollback history for pane ... ?
				// - (Save) to a file
				// - (View) with less
				// - (Edit) in text editor
				// startScrollbackHistoryForm()
				expandPreviewDisplay()
			}
		case 'f':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				startScrollbackHistoryForm(sessionDisplay.GetCurrentNode())
			}
		case 'a': // Temporary
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			} else {
				go startLoadingModal("Load stuff")
			}

		case 'b', 'g', 'i', 'j', 'k', 'm', 'n', 'o', 'p', 't', 'u', 'y', 'z':
			if event.Modifiers() == tcell.ModAlt {
				highlightSessionNode(int(event.Rune()) - 87)
			}
		}

		return event
	})
}
