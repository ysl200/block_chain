package web

import (
	"embed"
	"net/http"
)

//go:embed static
var content embed.FS

func StartWebServer() {
	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.FS(content))))
}
