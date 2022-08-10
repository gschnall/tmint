package tmux_interface

import (
	"strconv"
	"strings"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	twiz "github.com/gschnall/tmint/tmux_wizard"
)

var (
	searchBoxDisplay                         = tview.NewInputField()
	searchBoxDisplayHeight                   = 0
	searchBoxListOfResults                   = make([]*tview.TreeNode, 0)
	searchBoxLiistOfResultsIndex             = 0
	searchBoxWindowNodeShouldNotBeCollapsed  = true
	searchBoxSessionNodeShouldNotBeCollapsed = true
	searchBoxExpandedHistory                 = make([]*tview.TreeNode, 2) // session, window
)

func setSearchBoxColorTheme() {
	searchBoxDisplay.SetBackgroundColor(tview.Styles.ContrastBackgroundColor)
	searchBoxDisplay.SetFieldTextColor(tview.Styles.PrimaryTextColor)
	searchBoxDisplay.SetFieldBackgroundColor(tview.Styles.ContrastSecondaryTextColor)
	searchBoxDisplay.SetLabelColor(tview.Styles.SecondaryTextColor)
	searchBoxDisplay.SetTitleColor(tview.Styles.TertiaryTextColor)
}

func initSearchBoxDisplay() {
	searchBoxDisplay.
		SetLabel("Search: ").
		SetFieldWidth(32)
	searchBoxDisplay.SetBorder(true)
	searchBoxDisplay.SetChangedFunc(searchAndSelectNodeByText)
	searchBoxDisplay.SetTitle("")
	searchBoxDisplay.SetTitleAlign(0)

	initSearchDisplayKeys()
	setSearchBoxColorTheme()
}

func setSearchBoxResultsTitle() {
	if len(searchBoxListOfResults) > 1 {
		searchBoxDisplay.SetTitle(" Results: " + strconv.Itoa(searchBoxLiistOfResultsIndex+1) + " of " + strconv.Itoa(len(searchBoxListOfResults)) + " | CTRL-n next | CTRL-b back")
	} else {
		searchBoxDisplay.SetTitle(" Results: " + strconv.Itoa(searchBoxLiistOfResultsIndex+1) + " of " + strconv.Itoa(len(searchBoxListOfResults)) + " ")
	}
}

func saveSearchBoxExpandedHistory(node *tview.TreeNode) {
	tmux := node.GetReference()
	switch tmux.(type) {
	case twiz.Window:
		sessionNode := getParentOfNode(node)
		searchBoxSessionNodeShouldNotBeCollapsed = sessionNode.IsExpanded()
		searchBoxWindowNodeShouldNotBeCollapsed = true
		if sessionNode.IsExpanded() == false {
			sessionNode.Expand()
		}
		searchBoxExpandedHistory[0] = sessionNode
		searchBoxExpandedHistory[1] = node
	case twiz.Pane:
		windowNode := getParentOfNode(node)
		sessionNode := getParentOfNode(windowNode)
		searchBoxSessionNodeShouldNotBeCollapsed = sessionNode.IsExpanded()
		searchBoxWindowNodeShouldNotBeCollapsed = windowNode.IsExpanded()
		if sessionNode != nil && sessionNode.IsExpanded() == false {
			sessionNode.Expand()
		}
		if windowNode != nil && windowNode.IsExpanded() == false {
			windowNode.Expand()
		}
		searchBoxExpandedHistory[0] = sessionNode
		searchBoxExpandedHistory[1] = windowNode
	}
}

func undoSearchBoxExpandedHistory() {
	if searchBoxSessionNodeShouldNotBeCollapsed == false {
		searchBoxExpandedHistory[0].Collapse()
	}
	if searchBoxWindowNodeShouldNotBeCollapsed == false {
		searchBoxExpandedHistory[1].Collapse()
	}
}

func searchAndSelectNodeByText(text string) {
	if text == "" {
		searchBoxDisplay.SetTitle("")
		return
	}

	searchBoxListOfResults, searchBoxLiistOfResultsIndex = breadthFirstTextSearch(sessionDisplay.GetRoot(), sessionDisplay.GetCurrentNode(), text)
	if len(searchBoxListOfResults) > 0 {
		if len(searchBoxListOfResults) < searchBoxLiistOfResultsIndex+1 {
			searchBoxLiistOfResultsIndex = 0
		}
		currentSearchNode := searchBoxListOfResults[searchBoxLiistOfResultsIndex]

		undoSearchBoxExpandedHistory()
		saveSearchBoxExpandedHistory(currentSearchNode)
		// expandAllNodeParents(currentSearchNode)

		sessionDisplay.SetCurrentNode(currentSearchNode)
		setSearchBoxResultsTitle()
	} else {
		searchBoxDisplay.SetTitle(" Results: 0 ")
	}
}

func moveSearchSelectionForward() {
	if len(searchBoxListOfResults) > 0 {
		searchBoxLiistOfResultsIndex += 1
		if searchBoxLiistOfResultsIndex >= len(searchBoxListOfResults) {
			searchBoxLiistOfResultsIndex = 0
		}
		currentSearchNode := searchBoxListOfResults[searchBoxLiistOfResultsIndex]

		undoSearchBoxExpandedHistory()
		saveSearchBoxExpandedHistory(currentSearchNode)
		// expandAllNodeParents(currentSearchNode)

		sessionDisplay.SetCurrentNode(currentSearchNode)
		setSearchBoxResultsTitle()
	}
}
func moveSearchSelectionBack() {
	if len(searchBoxListOfResults) > 0 {
		searchBoxLiistOfResultsIndex -= 1
		if searchBoxLiistOfResultsIndex < 0 {
			searchBoxLiistOfResultsIndex = len(searchBoxListOfResults) - 1
		}
		currentSearchNode := searchBoxListOfResults[searchBoxLiistOfResultsIndex]

		undoSearchBoxExpandedHistory()
		saveSearchBoxExpandedHistory(currentSearchNode)
		// expandAllNodeParents(currentSearchNode)

		sessionDisplay.SetCurrentNode(currentSearchNode)
		setSearchBoxResultsTitle()
	}
	// TODO: set some text about results somewhere here
}

// _____________________
// |                   |
// | Utils and Actions |
// |                   |
// ---------------------
func toggleSearchBox() {
	// Reset these values
	searchBoxWindowNodeShouldNotBeCollapsed = true
	searchBoxSessionNodeShouldNotBeCollapsed = true

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

func breadthFirstTextSearch(root *tview.TreeNode, currentSelectedNode *tview.TreeNode, search string) ([]*tview.TreeNode, int) {
	// Maybe only need contains...
	// Starts with might be overkill
	// TODO: investigate if depth first makes more sense to user
	if search == "" {
		return make([]*tview.TreeNode, 0), 0
	}
	search = strings.ToLower(search)

	queue := make([]*tview.TreeNode, 0)
	queue = append(queue, root)

	// matchStarts := []*tview.TreeNode{}
	currentSelectedNodeIndex := 0
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
			nodeName = strings.ToLower(getPaneDisplayName(tmux.(twiz.Pane)))
		}

		// matchStart, _ := regexp.MatchString("(?i)^"+search, nodeName)
		// if matchStart {
		// 	matchContains = append(matchStarts, nextUp)
		if strings.Contains(nodeName, search) {
			matchContains = append(matchContains, nextUp)
		}

		if nextUp == currentSelectedNode {
			currentSelectedNodeIndex = len(matchContains)
		}

		nextUpChildren := nextUp.GetChildren()
		if len(nextUpChildren) > 0 {
			for _, child := range nextUpChildren {
				queue = append(queue, child)
			}
		}
	}

	if len(matchContains) > 0 {
		return matchContains, currentSelectedNodeIndex
	}
	return make([]*tview.TreeNode, 0), 0
}

// _______________
// |             |
// | Keybindings |
// |             |
// ---------------
func initSearchDisplayKeys() {
	searchBoxDisplay.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter, tcell.KeyCtrlS, tcell.KeyEsc:
			toggleSearchBox()
		case tcell.KeyCtrlN:
			moveSearchSelectionForward()
		case tcell.KeyCtrlB:
			moveSearchSelectionBack()
		}
		return event
	})
}
