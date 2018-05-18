package controller

import (
	"fmt"
	"net/http"
	"strings"
)

func registerHelloRoute() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		message := r.URL.Path
		message = strings.TrimPrefix(message, "/")
		message = "Hello " + message
		message += " " + r.Method
		message += "\nBody:" + fmt.Sprintf("%s", r.Body)
		w.Write([]byte(message))
	})

}
