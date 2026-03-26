package observability

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static
var staticFiles embed.FS

type WebUI struct {
	handler http.Handler
}

func NewWebUI() (*WebUI, error) {
	subFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, err
	}

	return &WebUI{
		handler: http.FileServer(http.FS(subFS)),
	}, nil
}

func (ui *WebUI) HandleUI(w http.ResponseWriter, r *http.Request) {
	ui.handler.ServeHTTP(w, r)
}
