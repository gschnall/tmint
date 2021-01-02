package tmux_wizard 

import (
	"fmt"
)

func getCurrentTmuxSession() string {
	return TmuxDisplayMessage("#S")
}

func sessionIsDetached() bool {
	return TmuxDisplayMessage("#{session_attached}") == "0"
}

// used for -t flag
func InitTmintForTmuxKeybindings(name string) {	
	currentSession := getCurrentTmuxSession()
	tmuxCommand := fmt.Sprintf("tmint -p -s \"%s\" -current-tmint-session \"%s\"", name, currentSession) 

	CreateTmuxSession(name, "~", 0)
  TmuxSendKeys(name, tmuxCommand)

	if sessionIsDetached() {
	  AttachTmuxSession(name)  	
	} else {
		SwitchToTmuxPath(name)
	}
}
