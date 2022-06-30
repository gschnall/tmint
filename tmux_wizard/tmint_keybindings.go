package tmux_wizard

import (
	"fmt"
)

func getCurrentTmuxSession() string {
	return TmuxGetMessage("#S")
}

// used for -t flag
// only works if a tmux session is attached
func InitTmintForTmuxKeybindings(name string) {
	currentSession := getCurrentTmuxSession()
	tmuxCommand := fmt.Sprintf("tmint -u -p -s \"%s\" -current-tmint-session \"%s\"", name, currentSession)

	// attempt to kill session before creating
	// won't exit on error
	TryToKillTmuxSession(name)

	CreateTmuxSession(name, "~", 0)
	TmuxSendKeys(name, tmuxCommand)
	SwitchToTmuxPath(name)
}
