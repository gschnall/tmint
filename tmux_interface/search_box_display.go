package tmux_interface

import (
	"regexp"
	"strings"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	twiz "github.com/gschnall/tmint/tmux_wizard"
)

var (
	searchBoxDisplay       = tview.NewInputField()
	searchBoxDisplayHeight = 0
)

func setSearchBoxColorTheme() {
	searchBoxDisplay.SetBackgroundColor(tview.Styles.ContrastBackgroundColor)
	searchBoxDisplay.SetFieldTextColor(tview.Styles.PrimaryTextColor)
	searchBoxDisplay.SetFieldBackgroundColor(tview.Styles.ContrastSecondaryTextColor)
	searchBoxDisplay.SetLabelColor(tview.Styles.SecondaryTextColor)
	searchBoxDisplay.SetTitleColor(tview.Styles.PrimaryTextColor)
}

func initSearchBoxDisplay() {
	searchBoxDisplay.
		SetLabel("Search: ").
		SetFieldWidth(32)
	searchBoxDisplay.SetBorder(true)
	searchBoxDisplay.SetChangedFunc(func(text string) {
		searchedForNode := breadthFirstSearch(sessionDisplay.GetRoot(), text)
		if searchedForNode != nil {
			sessionDisplay.SetCurrentNode(searchedForNode)
		}
	})

	initSearchDisplayKeys()
	setSearchBoxColorTheme()
}

// _____________________
// |                   |
// | Utils and Actions |
// |                   |
// ---------------------
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
		}
		return event
	})
}
