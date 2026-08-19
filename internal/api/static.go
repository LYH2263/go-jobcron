package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func staticDir(root string) http.Handler {
	fs := http.FileServer(http.Dir(root))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(root, "index.html"))
			return
		}
		p := strings.TrimPrefix(r.URL.Path, "/")
		if _, err := os.Stat(filepath.Join(root, p)); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(root, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})
}
