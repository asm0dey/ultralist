package ultralist

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Webapp is the main struct of this file.
type Webapp struct {
	Router *httprouter.Router
	server *http.Server
}

// Run is starting the ultralist webapp.
func (w *Webapp) Run() {
	w.server = &http.Server{Addr: ":9976"}

	http.HandleFunc("/", w.handleAuthResponse)
	http.HandleFunc("/favicon.ico", w.handleFavicon)

	w.server.ListenAndServe()
}

func (w *Webapp) handleFavicon(writer http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(writer, "")
}

func (w *Webapp) handleAuthResponse(writer http.ResponseWriter, r *http.Request) {
	cliTokens, ok := r.URL.Query()["cli_token"]
	if ok == false {
		fmt.Println("Something went wrong... I did not get a CLI token back.")
		os.Exit(0)
	}

	backend := NewBackend()
	backend.WriteCreds(cliTokens[0])
	fmt.Println("Authorization successful!")

	webTokens, ok := r.URL.Query()["web_token"]
	if ok == false {
		fmt.Println("Something went wrong... I did not get a web token back.")
		os.Exit(0)
	}

	signup, _ := r.URL.Query()["signup"]

	http.Redirect(writer, r, w.frontendUrl(webTokens[0], signup[0]), http.StatusSeeOther)

	// sleep 1 second before shutting server down, so we can display msg on web.
	go func() {
		time.Sleep(1 * time.Second)
		w.server.Shutdown(nil)
	}()
}

func (w *Webapp) frontendUrl(token string, signup string) string {
	envFrontendURL := os.Getenv("ULTRALIST_FRONTEND_URL")

	if envFrontendURL != "" {
		return envFrontendURL + "/auth?cli_auth=true&signup=" + signup + "&token=" + token
	}

	return "https://app.ultralist.io/auth?cli_auth=true&token=" + token
}
