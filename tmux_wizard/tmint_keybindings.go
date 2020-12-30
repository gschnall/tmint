package tmux_wizard 

import (
	"fmt"
)

func getCurrentTmuxSession() string {
	return TmuxDisplayMessage("#S")
}

// used for -t flag
func InitTmintForTmuxKeybindings(name string) {	
	currentSession := getCurrentTmuxSession()
	CreateTmuxSession(name, "~", 0)
	tmuxCommand := fmt.Sprintf("tmint -p -s \"%s\" -current-tmint-session \"%s\"", name, currentSession) 
  TmuxSendKeys(name, tmuxCommand)
	SwitchToTmuxPath(name)
}
