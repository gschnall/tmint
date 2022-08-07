package tmux_interface

import (
	"github.com/rivo/tview"

	twiz "github.com/gschnall/tmint/tmux_wizard"
)

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
				if !previousSib.IsExpanded() {
					return previousSib
				}
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

	if !node.IsExpanded() {
		return node
	} else if lastWindow.IsExpanded() {
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
			return children[index+1]
		}
	}
	return nil
}
func getPreviousSibling(node *tview.TreeNode) *tview.TreeNode {
	parent := getParentOfNode(node)
	children := parent.GetChildren()
	for index, child := range children {
		if child == node && index > 0 {
			return children[index-1]
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
	root := sessionDisplay.GetRoot()
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
