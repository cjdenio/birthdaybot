package commands

import "net/http"

// HandleHelpCommand handles /birthday help
func HandleHelpCommand(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("I'm not sure what you're trying to do :thinking_face:"))
}
