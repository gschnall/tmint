package tmux_wizard 

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var (
	tmintWindowName = "__| Tmint | a Tmux session manager |__"
)

// Command line Utils
func HandleExecError(err error, f string) {
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
	// fmt.Println("Result: " + out.String())
}
func clear() {
	cmd := exec.Command("clear")
	stdout, err := cmd.Output()
	HandleExecError(err, "clear")
	fmt.Printf(string(stdout))
}
func StartTmux() {
	cmd := exec.Command("tmux")
	// We need to connect tmux to our terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	// -- Tmux isn't running -- 
	HandleExecError(err, "startTmux")
}
func getTmuxLs() (string, bool) {
	format := "#{session_attached} :: #{session_name} :: #{session_activity} :: #{window_zoomed_flag}"
	cmd := exec.Command("tmux", "list-sessions", "-F", format)
	stdout, err := cmd.Output()
	// -- Tmux isn't running -- 
	if err != nil {
		return "", false
	}
	// HandleExecError(err, "getTmuxLs")
	return string(stdout), true
}
func getTmuxListWindows(sessionString string) string {
	format := "#{window_active} :: #{window_index} :: #{window_name} :: #{window_activity}"
	cmd := exec.Command("tmux", "list-window", "-F", format, "-t", sessionString)
	stdout, err := cmd.Output()
	HandleExecError(err, "getTmuxListWindows")
	return string(stdout)
}
func getTmuxListPanes(windowPath string) string {
	format := "#{pane_active} :: #{pane_index} :: #{pane_current_path} :: #{pane_current_command}"
	cmd := exec.Command("tmux", "list-panes", "-F", format, "-t", windowPath)
	stdout, err := cmd.Output()
	HandleExecError(err, windowPath)
	return string(stdout)
}
func getTmuxCapturePane(panePath string) string {
	cmd := exec.Command("tmux", "capture-pane", "-pe", "-t", panePath)
	stdout, err := cmd.Output()
	HandleExecError(err, "getTmuxCapturePane")
	return string(stdout)
}
func SwitchToTmuxPath(path string) {
	cmd := exec.Command("tmux", "switch-client", "-t", path)
	_, err := cmd.Output()
	HandleExecError(err, "SwitchToTmuxPath")
}
func AttachTmuxSession(sessionName string) {
	cmd := exec.Command("tmux", "attach-session", "-t", sessionName)
	// We need to connect tmux to our terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	HandleExecError(err, "AttachTmuxSession")
}
func RenameTmuxSession(sessionName string, newSessionName string) {
	cmd := exec.Command("tmux", "rename-session", "-t", sessionName, newSessionName)
	_, err := cmd.Output()
	HandleExecError(err, "RenameTmuxSession")
}
func RenameTmuxWindow(windowPath string, newWindowName string) {
	cmd := exec.Command("tmux", "rename-window", "-t", windowPath, newWindowName)
	_, err := cmd.Output()
	HandleExecError(err, "RenameTmuxWindow")
}
func KillTmuxSession(sessionName string) {
	cmd := exec.Command("tmux", "kill-session", "-t", sessionName)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	HandleExecError(err, "KillTmuxSession")
}
func KillTmuxWindow(windowPath string) {
	cmd := exec.Command("tmux", "kill-window", "-t", windowPath)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	HandleExecError(err, "KillTmuxWindow")
}
func KillTmuxPane(panePath string) {
	cmd := exec.Command("tmux", "kill-pane", "-t", panePath)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	HandleExecError(err, "KillTmuxPane")
}
func DetachTmuxSession(sessionName string) {
	cmd := exec.Command("tmux", "detach", "-s", sessionName)
	_, err := cmd.Output()
	HandleExecError(err, "DetachTmuxSession")
}

func tmuxSplitPaneV(panePath string) {
	cmd := exec.Command("tmux", "split-window", "-d", "-v", "-c", "#{pane_current_path}", "-t", panePath)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	HandleExecError(err, "tmuxSplitPaneV")
}
func tmuxSplitPaneH(panePath string) {
	cmd := exec.Command("tmux", "split-window", "-d", "-h", "-c", "#{pane_current_path}", "-t", panePath)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	HandleExecError(err, "tmuxSplitPaneH")
}
func CreateTmuxSession(name string, dir string, numberOfPanes int) {
	cmd := exec.Command("tmux", "new", "-d", "-s", strings.TrimSpace(name), "-c", dir)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	HandleExecError(err, "CreateTmuxSession")
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

func CreateTmuxWindow(name string, dir string, numberOfPanes int, targetSession string) {
	// BUG - Creating tmux window within a target session with an int as name
	// Need to target the next available int "name:nextAvailableInt" 4:1
	// - tmux new-window -d -n "cool cat" -c ~/Documents -t 4:1
	cmd := exec.Command("tmux", "new-window", "-d", "-n", strings.TrimSpace(name), "-c", dir, "-t", targetSession + ":")
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	HandleExecError(err, "CreateTmuxWindow")
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

func TmuxToggleFullscreen() {
	cmd := exec.Command("tmux", "resize-pane", "-Z")
	_, err := cmd.Output()
	HandleExecError(err, "TmuxToggleFullscreen")
}

func tmuxSendKeys(panePath string, command string) {
	cmd := exec.Command("tmux", "send-keys", "-t", panePath, command, "Enter")
	_, err := cmd.Output()
	HandleExecError(err, "tmuxSendKeys")
}

func getTmuxSessionList() ([]string, bool) {
	tls, tmuxIsRunning := getTmuxLs()
	if tmuxIsRunning == false {
		return make([]string, 0), false
	}
	return strings.Split(strings.TrimSpace(tls), "\n"), true
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
	HandleExecError(err, "parseSessionString")
	return numberAttachedTo > 0, a[1], a[2], zoomedFlag > 0
}
func parseWindowString(windowString string) (bool, string, string, string) {
	a := strings.Split(windowString, " :: ")
	isActive, err := strconv.Atoi(a[0])
	HandleExecError(err, "parseWindowString")
	return isActive > 0, a[1], a[2], a[3]
}
func parsePaneString(paneString string) (bool, string, string, string) {
	a := strings.Split(paneString, " :: ")
	isActive, err := strconv.Atoi(a[0])
	HandleExecError(err, "parsePaneString")
	return isActive > 0, a[1], a[2], a[3]
}

func getPanePreview(panePath string) string {
	return getTmuxCapturePane(panePath)
}

// Tmux Session Data 
type Session struct {
	Name              string
	CreationDate      string
	Preview           string
	Id                int
	ActiveWindowIndex int
	NumberOfWindows   int
	IsAttached        bool
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
	Session         string
	IsActive        bool
	IsExpanded      bool
	ActivePaneIndex int
	Panes           []Pane
}

type Pane struct {
	Name      string
	Path      string
	Index     string
	Command   string
	Directory string
	Session   string
	IsActive  bool
	Preview   string
}

type SessionData struct {
	HasAttachedSession   bool
	HasLivingSessions    bool
	HasZoomedPane        bool
	TmintWindowIsActive  bool
	MaxSessionNameLength int
	AttachedSession      string
	Sessions             []Session
}

type Tmux interface {
	kill()
}

func (s Session) kill() {
	KillTmuxSession(s.Name)
}
func (w Window) kill() {
	KillTmuxWindow(w.Path)
}
func (p Pane) kill() {
	KillTmuxPane(p.Path)
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

func getWindowPanes(windowPath string, session string) ([]Pane, int) {
	tmuxPanes := getTmuxPaneList(windowPath)
	sliceOfPanes := make([]Pane, len(tmuxPanes))
	activeIndex := 0

	for ind, paneString := range tmuxPanes {
		isActive, paneIndex, dir, command := parsePaneString(paneString)

		pane := Pane{IsActive: isActive, Index: paneIndex, Directory: dir, Command: command}
		pane.Session = session
		pane.Name = paneIndex
		pane.Path = windowPath + "." + pane.Name
		pane.Preview = getPanePreview(pane.Path)
		sliceOfPanes[ind] = pane

		if pane.IsActive {
			activeIndex = ind
		}
	}
	return sliceOfPanes, activeIndex
}

func getSessionWindows(sessionName string) ([]Window, int, bool) {
	tmuxWindows := getTmuxWindowList(sessionName)
	// Get last active window time (DATE ONLY for now)
	//time.Unix(1588109472, 0)
	sliceOfWindows := make([]Window, len(tmuxWindows))
	activeIndex := 0
	tmuxWindowIsActive := false

	for ind, windowString := range tmuxWindows {
		isActive, windowInd, name, _ := parseWindowString(windowString)

		window := Window{IsActive: isActive, Name: name, Index: windowInd}
		window.Session = sessionName
		window.Path = sessionName + ":" + window.Index
		window.Panes, window.ActivePaneIndex = getWindowPanes(window.Path, window.Session)
		window.Preview = window.Panes[window.ActivePaneIndex].Preview
		sliceOfWindows[ind] = window

		if window.IsActive {
			activeIndex = ind
		}

		if name == tmintWindowName {
			tmuxWindowIsActive = true	
		}
	}
	return sliceOfWindows, activeIndex, tmuxWindowIsActive
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

func GetSessionData() SessionData {
	sessionNameLimiter := 100
	sessionData := SessionData{HasAttachedSession: false, MaxSessionNameLength: 0}
	tmuxLsList, tmuxIsRunning := getTmuxSessionList()

	if tmuxIsRunning == false {
		return sessionData 
	}

	sliceOfSessions := make([]Session, len(tmuxLsList))

	for ind, sessionString := range tmuxLsList {
		isAttached, name, _, isZoomed := parseSessionString(sessionString)
		session := Session{IsAttached: isAttached, Name: name, Id: ind, IsZoomed: isZoomed}
		if isAttached {
			sessionData.HasAttachedSession = true
			sessionData.AttachedSession = name
			sessionData.HasZoomedPane = isZoomed
		}
		session.Windows, session.ActiveWindowIndex, sessionData.TmintWindowIsActive = getSessionWindows(session.Name)
		session.Preview = session.Windows[session.ActiveWindowIndex].Preview
		sliceOfSessions[ind] = session
		sessionData.MaxSessionNameLength = getMaxInt(sessionData.MaxSessionNameLength, len(session.Name))
	}

	sessionData.Sessions = sliceOfSessions
	sessionData.HasLivingSessions = len(sliceOfSessions) != 0 && sliceOfSessions[0].Name != ""
	sessionData.MaxSessionNameLength = getMinInt(sessionNameLimiter, sessionData.MaxSessionNameLength)
	return sessionData
}
