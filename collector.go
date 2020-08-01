package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Command line funcs
func handleExecError(err error, f string) {
	if err != nil {
		fmt.Println("Exec Error: " + f)
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
func debugExecError(cmd *exec.Cmd) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("Result: " + out.String())
}
func clear() {
	cmd := exec.Command("clear")
	stdout, err := cmd.Output()
	handleExecError(err, "clear")
	fmt.Printf(string(stdout))
}
func getTmuxLs() string {
	format := "#{session_attached} :: #{session_name} :: #{session_activity} :: #{window_zoomed_flag}"
	cmd := exec.Command("tmux", "list-sessions", "-F", format)
	stdout, err := cmd.Output()
	handleExecError(err, "getTmuxLs")
	return string(stdout)
}
func getTmuxListWindows(sessionString string) string {
	format := "#{window_active} :: #{window_index} :: #{window_name} :: #{window_activity}"
	cmd := exec.Command("tmux", "list-window", "-F", format, "-t", sessionString)
	stdout, err := cmd.Output()
	handleExecError(err, "getTmuxListWindows")
	return string(stdout)
}
func getTmuxListPanes(windowPath string) string {
	format := "#{pane_active} :: #{pane_index} :: #{pane_current_path} :: #{pane_current_command}"
	cmd := exec.Command("tmux", "list-panes", "-F", format, "-t", windowPath)
	stdout, err := cmd.Output()
	handleExecError(err, windowPath)
	return string(stdout)
}
func getTmuxCapturePane(panePath string) string {
	cmd := exec.Command("tmux", "capture-pane", "-pe", "-t", panePath)
	stdout, err := cmd.Output()
	handleExecError(err, "getTmuxCapturePane")
	return string(stdout)
}
func switchToTmuxPath(path string) {
	cmd := exec.Command("tmux", "switch-client", "-t", path)
	_, err := cmd.Output()
	handleExecError(err, "switchToTmuxPath")
}
func attachTmuxSession(sessionName string) {
	cmd := exec.Command("tmux", "attach-session", "-t", sessionName)
	// We need to connect tmux to our terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	handleExecError(err, "attachTmuxSession")
}
func killTmuxSession(sessionName string) {
	cmd := exec.Command("tmux", "kill-session", "-t", sessionName)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	handleExecError(err, "killTmuxSession")
}
func killTmuxWindow(windowPath string) {
	cmd := exec.Command("tmux", "kill-window", "-t", windowPath)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	handleExecError(err, "killTmuxWindow")
}
func killTmuxPane(panePath string) {
	cmd := exec.Command("tmux", "kill-pane", "-t", panePath)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	handleExecError(err, "killTmuxPane")
}
func detachTmuxSession(sessionName string) {
	cmd := exec.Command("tmux", "detach", "-s", sessionName)
	_, err := cmd.Output()
	handleExecError(err, "detachTmuxSession")
}

func tmuxSplitPaneV(panePath string) {
	cmd := exec.Command("tmux", "split-window", "-d", "-v", "-c", "#{pane_current_path}", "-t", panePath)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	handleExecError(err, "tmuxSplitPaneV")
}
func tmuxSplitPaneH(panePath string) {
	cmd := exec.Command("tmux", "split-window", "-d", "-h", "-c", "#{pane_current_path}", "-t", panePath)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	handleExecError(err, "tmuxSplitPaneH")
}
func createTmuxSession(name string, dir string, numberOfPanes int) {
	cmd := exec.Command("tmux", "new", "-d", "-s", strings.TrimSpace(name), "-c", dir)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	handleExecError(err, "createTmuxSession")
	if numberOfPanes == 2 {
		tmuxSplitPaneV(name + ":0")
	} else if numberOfPanes == 3 {
		tmuxSplitPaneV(name + ":0")
		tmuxSplitPaneH(name + ":0.1")
	} else if numberOfPanes == 4 {
		tmuxSplitPaneV(name + ":0")
		tmuxSplitPaneH(name + ":0.0")
		tmuxSplitPaneH(name + ":0.2")
	}
}

func createTmuxWindow(name string, dir string, numberOfPanes int) {
	cmd := exec.Command("tmux", "new-window", "-d", "-n", strings.TrimSpace(name), "-c", dir)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	handleExecError(err, "createTmuxWindow")
	if numberOfPanes == 2 {
		tmuxSplitPaneV(name + ":0")
	} else if numberOfPanes == 3 {
		tmuxSplitPaneV(name + ":0")
		tmuxSplitPaneH(name + ":0.1")
	} else if numberOfPanes == 4 {
		tmuxSplitPaneV(name + ":0")
		tmuxSplitPaneH(name + ":0.0")
		tmuxSplitPaneH(name + ":0.2")
	}
}

func tmuxToggleFullscreen() {
	cmd := exec.Command("tmux", "resize-pane", "-Z")
	_, err := cmd.Output()
	handleExecError(err, "tmuxToggleFullscreen")
}

func tmuxSendKeys(panePath string, command string) {
	cmd := exec.Command("tmux", "send-keys", "-t", panePath, command, "Enter")
	_, err := cmd.Output()
	handleExecError(err, "tmuxSendKeys")
}

func getTmuxSessionList() []string {
	tls := getTmuxLs()
	return strings.Split(strings.TrimSpace(tls), "\n")
}
func getTmuxWindowList(sessionName string) []string {
	wls := getTmuxListWindows(sessionName)
	return strings.Split(strings.TrimSpace(wls), "\n")
}
func getTmuxPaneList(windowPath string) []string {
	pls := getTmuxListPanes(windowPath)
	return strings.Split(strings.TrimSpace(pls), "\n")
}

func parseSessionString(sessionString string) (bool, string, string, bool) {
	a := strings.Split(sessionString, " :: ")
	numberAttachedTo, err := strconv.Atoi(a[0])
	zoomedFlag, err := strconv.Atoi(a[3])
	handleExecError(err, "parseSessionString")
	return numberAttachedTo > 0, a[1], a[2], zoomedFlag > 0
}
func parseWindowString(windowString string) (bool, string, string, string) {
	a := strings.Split(windowString, " :: ")
	isActive, err := strconv.Atoi(a[0])
	handleExecError(err, "parseWindowString")
	return isActive > 0, a[1], a[2], a[3]
}
func parsePaneString(paneString string) (bool, string, string, string) {
	a := strings.Split(paneString, " :: ")
	isActive, err := strconv.Atoi(a[0])
	handleExecError(err, "parsePaneString")
	return isActive > 0, a[1], a[2], a[3]
}

func getPanePreview(panePath string) string {
	return getTmuxCapturePane(panePath)
}

type Session struct {
	name              string
	creationDate      string
	preview           string
	id                int
	activeWindowIndex int
	numberOfWindows   int
	isAttached        bool
	isZoomed          bool
	isExpanded        bool
	windows           []Window
}

type Window struct {
	name            string
	index           string
	activeDate      string
	path            string
	preview         string
	session         string
	isActive        bool
	isExpanded      bool
	activePaneIndex int
	panes           []Pane
}

type Pane struct {
	name      string
	path      string
	index     string
	command   string
	directory string
	session   string
	isActive  bool
	preview   string
}

type SessionData struct {
	hasAttachedSession   bool
	hasLivingSessions    bool
	hasZoomedPane        bool
	maxSessionNameLength int
	attachedSession      string
	sessions             []Session
}

type Tmux interface {
	kill()
}

func (s Session) kill() {
	killTmuxSession(s.name)
}
func (w Window) kill() {
	killTmuxWindow(w.path)
}
func (p Pane) kill() {
	killTmuxPane(p.path)
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

func sessionIsAttached(sessionString string) bool {
	return strings.HasSuffix(sessionString, "] (attached)")
}
func nodeIsActive(nodeString string) bool {
	return strings.HasSuffix(nodeString, " (active)")
}
func getNameFromSessionString(sessionString string) string {
	return strings.Split(sessionString, ":")[0]
}

func toCharStr(i int) string {
	return string('a' - 1 + i)
}

func getWindowPanes(windowPath string, session string) ([]Pane, int) {
	tmuxPanes := getTmuxPaneList(windowPath)
	sliceOfPanes := make([]Pane, len(tmuxPanes))
	activeIndex := 0

	for ind, paneString := range tmuxPanes {
		isActive, paneIndex, dir, command := parsePaneString(paneString)

		pane := Pane{isActive: isActive, index: paneIndex, directory: dir, command: command}
		pane.session = session
		pane.name = paneIndex
		pane.path = windowPath + "." + pane.name
		pane.preview = getPanePreview(pane.path)
		sliceOfPanes[ind] = pane

		if pane.isActive {
			activeIndex = ind
		}
	}
	return sliceOfPanes, activeIndex
}

func getSessionWindows(sessionName string) ([]Window, int) {
	tmuxWindows := getTmuxWindowList(sessionName)
	// Get last active window time (DATE ONLY for now)
	//time.Unix(1588109472, 0)
	sliceOfWindows := make([]Window, len(tmuxWindows))
	activeIndex := 0

	for ind, windowString := range tmuxWindows {
		isActive, windowInd, name, _ := parseWindowString(windowString)

		window := Window{isActive: isActive, name: name, index: windowInd}
		window.session = sessionName
		window.path = sessionName + ":" + window.index
		window.panes, window.activePaneIndex = getWindowPanes(window.path, window.session)
		window.preview = window.panes[window.activePaneIndex].preview
		sliceOfWindows[ind] = window

		if window.isActive {
			activeIndex = ind
		}
	}
	return sliceOfWindows, activeIndex
}

func getNumberOfWindows(sessionString string) int {
	re := regexp.MustCompile(`: (\d+) windows \(`)
	match := re.FindString(sessionString)
	i, err := strconv.Atoi(match)
	if err != nil {
		handleExecError(err, "getNumberOfWindows")
	}
	return i
}

func getNumberOfPanes(windowString string) int {
	re := regexp.MustCompile(`\((\d+) panes\) \[`)
	match := re.FindString(windowString)
	i, err := strconv.Atoi(match)
	if err != nil {
		handleExecError(err, "getNumberOfPanes")
	}
	return i
}

func getSessionData() SessionData {
	sessionNameLimiter := 100
	sessionData := SessionData{hasAttachedSession: false, maxSessionNameLength: 0}
	tmuxLsList := getTmuxSessionList()

	sliceOfSessions := make([]Session, len(tmuxLsList))

	for ind, sessionString := range tmuxLsList {
		isAttached, name, _, isZoomed := parseSessionString(sessionString)
		session := Session{isAttached: isAttached, name: name, id: ind, isZoomed: isZoomed}
		if isAttached {
			sessionData.hasAttachedSession = true
			sessionData.attachedSession = name
			sessionData.hasZoomedPane = isZoomed
		}
		session.windows, session.activeWindowIndex = getSessionWindows(session.name)
		session.preview = session.windows[session.activeWindowIndex].preview
		sliceOfSessions[ind] = session
		sessionData.maxSessionNameLength = getMaxInt(sessionData.maxSessionNameLength, len(session.name))
	}

	sessionData.sessions = sliceOfSessions
	sessionData.hasLivingSessions = len(sliceOfSessions) != 0 && sliceOfSessions[0].name != ""
	sessionData.maxSessionNameLength = getMinInt(sessionNameLimiter, sessionData.maxSessionNameLength)
	return sessionData
}
