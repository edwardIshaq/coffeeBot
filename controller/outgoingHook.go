package controller

import (
	"fmt"
	"net/http"
	"strconv"
)

/*
Handler for outgoing hooks
*/
func registerOutgoingHookRoute() {
	http.HandleFunc("/outgoingHooks", func(w http.ResponseWriter, r *http.Request) {
		text := r.PostFormValue("text")
		value, _ := strconv.ParseInt(text, 10, 32)
		message := fmt.Sprintf(`{"text" : "%d"}`, value+1)
		w.Write([]byte(message))
	})
}
