package tmux_wizard

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"time"
)

var (
	tmintDirName = ".tmux_interface"
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

func TmuxDisplayMessage(message string) string {
	cmd := exec.Command("tmux", "display-message", "-F", message)
	stdout, err := cmd.Output()
	HandleExecError(err, "tmuxDisplayMessage")
	return strings.TrimSpace(string(stdout))
}
func tmuxDisplayLongMessage(message string, count int) {
	for {
		if count <= 0 {
			break
		}
		TmuxDisplayMessage(message)
		time.Sleep(300 * time.Millisecond)
		count -= 1
	}
}

func TmuxGetMessage(message string) string {
	cmd := exec.Command("tmux", "display-message", "-p", message)
	stdout, err := cmd.Output()
	HandleExecError(err, "tmuxDisplayMessage")
	return strings.TrimSpace(string(stdout))
}

func getTmuxHistoryLimit() int {
	historyLimitStr := TmuxGetMessage("#{history_limit}")
	historyLimit, err := strconv.Atoi(historyLimitStr)
	HandleExecError(err, "TmuxGetHistoryLimit")
	return historyLimit
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
	err := cmd.Start()
	HandleExecError(err, "KillTmuxSession")
	err = cmd.Wait()
	HandleExecError(err, "KillTmuxSession")
}

func TryToKillTmuxSession(sessionName string) {
	cmd := exec.Command("tmux", "kill-session", "-t", sessionName)
	cmd.Output()
}

func KillTmuxWindow(windowPath string) {
	cmd := exec.Command("tmux", "kill-window", "-t", windowPath)
	_, err := cmd.Output()
	HandleExecError(err, "KillTmuxWindow")
}

func KillTmuxPane(panePath string) {
	cmd := exec.Command("tmux", "kill-pane", "-t", panePath)
	err := cmd.Start()
	HandleExecError(err, "KillTmuxPane")
	err = cmd.Wait()
	HandleExecError(err, "KillTmuxPane")
}

func DetachTmuxSession(sessionName string) {
	cmd := exec.Command("tmux", "detach", "-s", sessionName)
	_, err := cmd.Output()
	HandleExecError(err, "DetachTmuxSession")
}

func tmuxSplitPaneV(panePath string) {
	cmd := exec.Command("tmux", "split-window", "-d", "-v", "-c", "#{pane_current_path}", "-t", panePath)
	_, err := cmd.Output()
	HandleExecError(err, "tmuxSplitPaneV")
}

func tmuxSplitPaneH(panePath string) {
	cmd := exec.Command("tmux", "split-window", "-d", "-h", "-c", "#{pane_current_path}", "-t", panePath)
	_, err := cmd.Output()
	HandleExecError(err, "tmuxSplitPaneH")
}

func TmuxResizePaneUp() {
	cmd := exec.Command("tmux", "resize-pane", "-U")
	_, err := cmd.Output()
	HandleExecError(err, "tmuxResizePaneUp")
}
func TmuxResizePaneDown() {
	cmd := exec.Command("tmux", "resize-pane", "-D")
	_, err := cmd.Output()
	HandleExecError(err, "tmuxResizePaneDown")
}
func TmuxResizePaneRight() {
	cmd := exec.Command("tmux", "resize-pane", "-R")
	_, err := cmd.Output()
	HandleExecError(err, "tmuxResizePaneRight")
}
func TmuxResizePaneLeft() {
	cmd := exec.Command("tmux", "resize-pane", "-L")
	_, err := cmd.Output()
	HandleExecError(err, "tmuxResizePaneLeft")
}

func CreateTmuxSession(name string, dir string, numberOfPanes int) {
	cmd := exec.Command("tmux", "new", "-d", "-s", strings.TrimSpace(name), "-c", dir)
	err := cmd.Start()
	HandleExecError(err, "CreateTmuxSession")
	err = cmd.Wait()
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
	cmd := exec.Command("tmux", "new-window", "-d", "-n", strings.TrimSpace(name), "-c", dir, "-t", targetSession+":")
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

func TmuxSendKeys(panePath string, command string) {
	cmd := exec.Command("tmux", "send-keys", "-t", panePath, command, "Enter")
	_, err := cmd.Output()
	HandleExecError(err, "TmuxSendKeys")
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
	// tmux capture-pane -pS -${1} -t "T mint:0.0" > ~/Documents/captured-stuff.sh
	return strings.TrimSpace(getTmuxCapturePane(panePath))
}

func GetColoredPaneScrollback(panePath string, numberOfScrollbackLines int) string {
	numberOfScrollbackLinesArg := "-" + strconv.Itoa(numberOfScrollbackLines)

	cmd := exec.Command("tmux", "capture-pane", "-pe", "-S", numberOfScrollbackLinesArg, "-t", panePath)
	stdout, err := cmd.Output()
	HandleExecError(err, "GetColoredPaneScrollback")
	return strings.TrimSpace(string(stdout))
}

func getPaneScrollback(panePath string, numberOfScrollbackLines int) string {
	numberOfScrollbackLinesArg := "-" + strconv.Itoa(numberOfScrollbackLines)

	cmd := exec.Command("tmux", "capture-pane", "-p", "-S", numberOfScrollbackLinesArg, "-t", panePath)
	stdout, err := cmd.Output()
	HandleExecError(err, "getPaneScrollback")
	return strings.TrimSpace(string(stdout))
}

func generateTmintDirectory() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	homeDirectory := user.HomeDir

	tmintDir := homeDirectory + "/" + tmintDirName
	if _, err := os.Stat(tmintDir); os.IsNotExist(err) {
		os.Mkdir(tmintDir, 0755)
	}
	return tmintDir
}

func SaveTmuxScrollback(panePath string, fileName string, saveToTmintDir bool, numberOfLinesToSave int) {
	tmintDir := ""
	if saveToTmintDir {
		tmintDir = generateTmintDirectory() + "/"
	}
	// TODO: let user choose number of scrollback lines
	paneScrollback := getPaneScrollback(panePath, numberOfLinesToSave)

	filePath := tmintDir + fileName
	f, err := os.Create(filePath)
	if err != nil {
		HandleExecError(err, "saveTmuxScrollback")
		fmt.Println(err)
	}
	// close the file with defer
	defer f.Close()

	f.WriteString(paneScrollback)
	tmuxDisplayLongMessage("File Saved: "+filePath, 4)
}
