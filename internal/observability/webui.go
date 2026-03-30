package observability

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static
var staticFiles embed.FS

type WebUI struct {
	handler   http.Handler
	password  string
	lanAccess bool
}

func NewWebUI(password string, lanAccess bool) (*WebUI, error) {
	subFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, err
	}

	return &WebUI{
		handler:   http.FileServer(http.FS(subFS)),
		password:  password,
		lanAccess: lanAccess,
	}, nil
}

func (ui *WebUI) HandleUI(w http.ResponseWriter, r *http.Request) {
	if ui.password != "" {
		if !ui.isAuthorized(r) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Admin"`)
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}
	}
	ui.handler.ServeHTTP(w, r)
}

func (ui *WebUI) isAuthorized(r *http.Request) bool {
	if ui.password == "" {
		return true
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		return false
	}
	return username == "admin" && password == ui.password
}
