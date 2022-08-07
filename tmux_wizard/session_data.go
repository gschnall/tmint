package tmux_wizard

import (
	"regexp"
	"strconv"
	"strings"
)

type Session struct {
	Name              string
	CreationDate      string
	Preview           string
	PanePath          string
	ActiveWindowName  string
	Id                int
	ActiveWindowIndex int
	NumberOfWindows   int
	IsAttached        bool
	IsCurrentNode     bool
	IsZoomed          bool
	IsExpanded        bool
	Windows           []Window
}

type Window struct {
	Name            string
	Index           string
	ActiveDate      string
	Path            string
	Preview         string
	PanePath        string
	Session         string
	IsActive        bool
	IsCurrentNode   bool
	IsExpanded      bool
	ActivePaneIndex int
	Panes           []Pane
}

type Pane struct {
	Name          string
	Path          string
	Index         string
	Command       string
	Directory     string
	Session       string
	Window        string
	IsActive      bool
	IsCurrentNode bool
	Preview       string
}

type SessionData struct {
	HasAttachedSession   bool
	HasLivingSessions    bool
	HasZoomedPane        bool
	IsUsingKeybindings   bool
	MaxSessionNameLength int
	HistoryLimit         int
	AttachedSession      string
	TmintSession         string
	Sessions             []Session
}

func getMaxInt(x int, y int) int {
	if x >= y {
		return x
	}
	return y
}

func getMinInt(x int, y int) int {
	if x <= y {
		return x
	}
	return y
}

func ToCharStr(i int) string {
	return string('a' - 1 + i)
}

func getWindowPanes(windowPath string, session string, windowName string, selectedPanePath string) ([]Pane, int, bool) {
	tmuxPanes := getTmuxPaneList(windowPath)
	sliceOfPanes := make([]Pane, len(tmuxPanes))
	hasSelectedCurrentPane := false
	activeIndex := 0

	for ind, paneString := range tmuxPanes {
		isActive, paneIndex, dir, command := parsePaneString(paneString)

		pane := Pane{IsActive: isActive, Index: paneIndex, Directory: dir, Command: command}
		pane.Session = session
		pane.Window = windowName
		pane.Name = paneIndex
		pane.Path = windowPath + "." + pane.Name
		pane.Preview = getPanePreview(pane.Path)

		if pane.Path == selectedPanePath {
			hasSelectedCurrentPane = true
			pane.IsCurrentNode = true
		}

		if pane.IsActive {
			activeIndex = ind
		}

		sliceOfPanes[ind] = pane
	}
	return sliceOfPanes, activeIndex, hasSelectedCurrentPane
}

func getSessionWindows(sessionName string, selectedWindowPath string) ([]Window, int, string, bool) {
	tmuxWindows := getTmuxWindowList(sessionName)
	// Get last active window time (DATE ONLY for now)
	//time.Unix(1588109472, 0)
	sliceOfWindows := make([]Window, len(tmuxWindows))
	hasSelectedCurrentNode := false
	hasSelectedCurrentPane := false
	activeIndex := 0
	activeName := ""
	// - Not currently being used
	// look to delete

	for ind, windowString := range tmuxWindows {
		isActive, windowInd, name, _ := parseWindowString(windowString)

		window := Window{IsActive: isActive, Name: name, Index: windowInd}
		window.Session = sessionName
		window.Path = sessionName + ":" + window.Index
		window.Panes, window.ActivePaneIndex, hasSelectedCurrentPane = getWindowPanes(window.Path, window.Session, window.Name, selectedWindowPath)
		window.Preview = window.Panes[window.ActivePaneIndex].Preview
		window.PanePath = window.Panes[window.ActivePaneIndex].Path

		pathArr := strings.Split(selectedWindowPath, ":")
		if window.Path == selectedWindowPath || (pathArr[0] == sessionName && pathArr[1] == window.Name) {
			hasSelectedCurrentNode = true
			window.IsCurrentNode = true
		} else if hasSelectedCurrentPane {
			window.IsExpanded = true
			hasSelectedCurrentNode = true
		}

		if window.IsActive {
			activeIndex = ind
			activeName = window.Name
		}

		sliceOfWindows[ind] = window
	}
	return sliceOfWindows, activeIndex, activeName, hasSelectedCurrentNode
}

func getNumberOfWindows(sessionString string) int {
	re := regexp.MustCompile(`: (\d+) windows \(`)
	match := re.FindString(sessionString)
	i, err := strconv.Atoi(match)
	if err != nil {
		HandleExecError(err, "getNumberOfWindows")
	}
	return i
}

func getNumberOfPanes(windowString string) int {
	re := regexp.MustCompile(`\((\d+) panes\) \[`)
	match := re.FindString(windowString)
	i, err := strconv.Atoi(match)
	if err != nil {
		HandleExecError(err, "getNumberOfPanes")
	}
	return i
}

func GetSessionData(currentSession string, tmintSession string, runFromKeybindings bool, result chan SessionData, targetName string, targetType string, targetPath string) {
	sessionNameLimiter := 100
	sessionData := SessionData{HasAttachedSession: false, MaxSessionNameLength: 0, TmintSession: tmintSession, IsUsingKeybindings: runFromKeybindings}
	tmuxLsList, tmuxIsRunning := getTmuxSessionList()
	// No tmux Commands before this line, EXCEPTION: getTmuxSessionList
	if tmuxIsRunning == false {
		result <- sessionData
		return
	}
	// History limit is number of scrollback lines
	sessionData.HistoryLimit = getTmuxHistoryLimit()

	if currentSession != ":" && runFromKeybindings {
		tmuxLsList = tmuxLsList[:len(tmuxLsList)-1]
	}
	sliceOfSessions := make([]Session, len(tmuxLsList))

	// On refresh we want to set the selected node
	selectedWindowPath := ""
	selectedSessionName := ""
	hasSelectedCurrentNode := false
	if targetType == "session" {
		selectedSessionName = targetName
	} else {
		// We need to make sure the session gets expanded
		selectedWindowPath = targetPath
	}

	for ind, sessionString := range tmuxLsList {
		isAttached, name, _, isZoomed := parseSessionString(sessionString)
		session := Session{IsAttached: isAttached, Name: name, Id: ind, IsZoomed: isZoomed}
		if isAttached || currentSession == name {
			sessionData.HasZoomedPane = isZoomed
			sessionData.AttachedSession = name
			sessionData.HasAttachedSession = true
			// Needed for tmux-keybindings workflow
			session.IsAttached = true
		}
		session.Windows, session.ActiveWindowIndex, session.ActiveWindowName, hasSelectedCurrentNode = getSessionWindows(session.Name, selectedWindowPath)
		session.Preview = session.Windows[session.ActiveWindowIndex].Preview
		session.PanePath = session.Windows[session.ActiveWindowIndex].PanePath

		if session.Name == selectedSessionName {
			session.IsCurrentNode = true
		} else if hasSelectedCurrentNode {
			session.IsExpanded = true
		}
		sliceOfSessions[ind] = session

		sessionData.MaxSessionNameLength = getMaxInt(sessionData.MaxSessionNameLength, len(session.Name))
	}

	sessionData.Sessions = sliceOfSessions
	sessionData.HasLivingSessions = len(sliceOfSessions) != 0 && sliceOfSessions[0].Name != ""
	sessionData.MaxSessionNameLength = getMinInt(sessionNameLimiter, sessionData.MaxSessionNameLength)

	result <- sessionData
	return
}
