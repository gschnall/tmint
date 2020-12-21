package tmux_interface

import (
	"io/ioutil"
	"os"
	"strings"


	twiz "github.com/gschnall/tmint/tmux_wizard"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	directoryInputField    = tview.NewInputField()
	firstTabCompletionTerm = ""
)


func initDirectoryInputField() {
	directoryInputField.SetText("~/")
	directoryInputField.SetLabel("Directory")
	directoryInputField.SetFieldWidth(30)

	// Enables selection of first autocomplte term on ENTER
	directoryInputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			directoryInputField.SetText(firstTabCompletionTerm)
			// -- Tab key could not be overriden - Attempted June 1st 2020 --
			// case tcell.KeyTab:
			// 	text := directoryInputField.GetText()
			// 	if text[len(text)-1:] != "/" {
			// 		directoryInputField.SetText(text + "/")
			// 	}
			// -- --------------------------- --
		}
		return event
	})

	directoryInputField.SetAutocompleteFunc(func(text string) (entries []string) {
		if len(text) == 0 || directoryInputField.HasFocus() == false {
			return
		} else if len(text) > 1 && text[0:2] == "~/" {
			home, homeDirErr := os.UserHomeDir()
			if homeDirErr != nil {
				twiz.HandleExecError(homeDirErr, "initDirectoryInputField")
			}
			text = home + "/" + text[2:]
		}
		dirParts := strings.Split(text, "/")
		directory := dirParts[0]
		searchTerm := ""
		if dirParts[0] == "" && len(dirParts) <= 2 { //Fix for directories with leading "/" char
			searchTerm = dirParts[len(dirParts)-1]
			directory = "/" + strings.Join(dirParts[0:len(dirParts)-1], "/")
		} else if len(dirParts) > 1 {
			searchTerm = dirParts[len(dirParts)-1]
			directory = strings.Join(dirParts[0:len(dirParts)-1], "/")
		}
		files, _ := ioutil.ReadDir(directory)

		if directory == "/" { //Fix for directories with leading "/" char
			return getSliceOfDirectories(searchTerm, directory, files)
		}
		return getSliceOfDirectories(searchTerm, directory+"/", files)
	})
}

// _____________________
// |                   |
// | Utils and Actions |
// |                   |
// ---------------------
func getSliceOfDirectories(search string, currentPath string, files []os.FileInfo) []string {
	dirs := []string{}
	if search == "" {
		for _, file := range files {
			if file.IsDir() && len(file.Name()) > 0 {
				dirs = append(dirs, currentPath+file.Name()+"/")
			}
		}
		return dirs
	} else {
		matchedPrefixes := []string{}
		matchedGeneral := []string{}

		for _, file := range files {
			if file.IsDir() && len(file.Name()) > 0 {
				fName := strings.ToLower(file.Name())
				sTerm := strings.ToLower(search)
				if strings.HasPrefix(fName, sTerm) {
					matchedPrefixes = append(matchedPrefixes, currentPath+file.Name())
				} else if strings.Contains(fName, sTerm) {
					matchedGeneral = append(matchedGeneral, currentPath+file.Name())
				}
			}
		}

		if len(matchedPrefixes) > 0 {
			firstTabCompletionTerm = matchedPrefixes[0] + "/"
		} else if len(matchedGeneral) > 0 {
			firstTabCompletionTerm = matchedGeneral[0] + "/"
		}
		return append(matchedPrefixes, matchedGeneral...)
	}
}
