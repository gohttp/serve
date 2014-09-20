package serve

import "path/filepath"
import "net/http"

// New file serving middleware.
func New(root string) func(http.Handler) http.Handler {
	fs := http.Dir(root)

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// open
			f, err := fs.Open(r.URL.Path)
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}
			defer f.Close()

			// stat
			stat, err := f.Stat()
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}

			// dir
			if stat.IsDir() {
				h.ServeHTTP(w, r)
				return
			}

			// file
			path := filepath.Join(root, stat.Name())
			http.ServeFile(w, r, path)
		})
	}
}
