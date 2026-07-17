package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"

	"github.com/rdrseraphim/blog/internal/frontmatter"
	"github.com/rdrseraphim/blog/internal/render"
)

type htmlPreviewServer struct {
	path string

	mu      sync.Mutex
	version int
	page    string
	err     error
}

func (s *htmlPreviewServer) reload() {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		s.err = err
		s.version++
		return
	}
	fm, body, hasFM := frontmatter.Split(data)
	title := filepath.Base(s.path)
	if hasFM {
		if t := extractYAMLString(fm, "title"); t != "" {
			title = t
		}
	}

	fragment, err := render.ToHTML(body)
	if err != nil {
		s.err = err
		s.version++
		return
	}
	s.err = nil
	s.page = render.Page(title, fragment)
	s.version++
}

func (s *htmlPreviewServer) snapshot() (page string, version int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.page, s.version, s.err
}

const liveReloadScript = `
<script>
(function() {
  var current = %d;
  setInterval(function() {
    fetch('/__version').then(function(r) { return r.text(); }).then(function(v) {
      if (parseInt(v, 10) !== current) location.reload();
    }).catch(function() {});
  }, 1000);
})();
</script>
</body>`

func (s *htmlPreviewServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	page, version, err := s.snapshot()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page = strings.Replace(page, "</body>", fmt.Sprintf(liveReloadScript, version), 1)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, page)
}

func (s *htmlPreviewServer) handleVersion(w http.ResponseWriter, r *http.Request) {
	_, version, _ := s.snapshot()
	fmt.Fprint(w, strconv.Itoa(version))
}

func previewHTML(path string, noOpen bool) error {
	s := &htmlPreviewServer{path: path}
	s.reload()
	if s.err != nil {
		return s.err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	dir := filepath.Dir(path)
	if err := watcher.Add(dir); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case ev, ok := <-watcher.Events:
				if !ok {
					return
				}
				if ev.Name == path && (ev.Op&(fsnotify.Write|fsnotify.Create) != 0) {
					s.reload()
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/__version", s.handleVersion)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s/", ln.Addr().String())
	fmt.Fprintf(os.Stderr, "serving live HTML preview of %s at %s (auto-reloads on save, ctrl+c to stop)\n", path, url)

	if !noOpen {
		openBrowser(url)
	}

	return http.Serve(ln, mux)
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
