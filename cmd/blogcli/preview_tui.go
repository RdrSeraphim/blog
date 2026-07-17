package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/fsnotify/fsnotify"

	"github.com/rdrseraphim/blog/internal/frontmatter"
)

var headerStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("212")).
	Padding(0, 1)

var footerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240")).
	Padding(0, 1)

type fileChangedMsg struct{}
type reloadErrMsg struct{ err error }

type tuiModel struct {
	path     string
	viewport viewport.Model
	ready    bool
	title    string
	err      error
	events   chan tea.Msg
}

func newTUIModel(path string) *tuiModel {
	return &tuiModel{path: path, events: make(chan tea.Msg, 8)}
}

func (m *tuiModel) Init() tea.Cmd {
	return tea.Batch(m.waitForEvent, m.watchFile)
}

// waitForEvent turns the next value on m.events into a tea.Msg, and
// re-arms itself so the watcher keeps feeding the program.
func (m *tuiModel) waitForEvent() tea.Msg {
	return <-m.events
}

// watchFile runs an fsnotify watcher on the file's directory (editors
// commonly replace-and-rename on save, which drops a direct file watch)
// and forwards write/create events for our file onto m.events.
func (m *tuiModel) watchFile() tea.Msg {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return reloadErrMsg{err}
	}
	dir := "."
	if idx := strings.LastIndexByte(m.path, '/'); idx != -1 {
		dir = m.path[:idx]
	}
	if err := watcher.Add(dir); err != nil {
		return reloadErrMsg{err}
	}
	go func() {
		for {
			select {
			case ev, ok := <-watcher.Events:
				if !ok {
					return
				}
				if ev.Name == m.path && (ev.Op&(fsnotify.Write|fsnotify.Create) != 0) {
					m.events <- fileChangedMsg{}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				m.events <- reloadErrMsg{err}
			}
		}
	}()
	return nil
}

func (m *tuiModel) render(width int) (string, error) {
	data, err := os.ReadFile(m.path)
	if err != nil {
		return "", err
	}
	fm, body, hasFM := frontmatter.Split(data)
	m.title = m.path
	if hasFM {
		if t := extractYAMLString(fm, "title"); t != "" {
			m.title = t
		}
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", err
	}
	return r.Render(body)
}

func (m *tuiModel) reload() {
	out, err := m.render(m.viewport.Width)
	if err != nil {
		m.err = err
		return
	}
	m.err = nil
	atTop := m.viewport.AtTop()
	m.viewport.SetContent(out)
	if atTop {
		m.viewport.GotoTop()
	}
}

func (m *tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "r":
			m.reload()
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMargin := headerHeight + footerHeight
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMargin)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMargin
		}
		m.reload()
	case fileChangedMsg:
		m.reload()
		cmds = append(cmds, m.waitForEvent)
	case reloadErrMsg:
		m.err = msg.err
		cmds = append(cmds, m.waitForEvent)
	}

	vp, cmd := m.viewport.Update(msg)
	m.viewport = vp
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *tuiModel) headerView() string {
	return headerStyle.Render(m.title)
}

func (m *tuiModel) footerView() string {
	status := fmt.Sprintf("%3.f%%  q: quit  r: reload  auto-reloads on save", m.viewport.ScrollPercent()*100)
	if m.err != nil {
		status = "error: " + m.err.Error()
	}
	return footerStyle.Render(status)
}

func (m *tuiModel) View() string {
	if !m.ready {
		return "loading..."
	}
	return m.headerView() + "\n" + m.viewport.View() + "\n" + m.footerView()
}

func previewTUI(path string) error {
	m := newTUIModel(path)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// extractYAMLString does a minimal best-effort lookup of a top-level
// "key: value" scalar in a YAML front matter block, stripping simple
// quoting. It's intentionally not a full YAML parser: it's only used to
// pull a display title for the TUI header.
func extractYAMLString(fm, key string) string {
	for _, line := range strings.Split(fm, "\n") {
		prefix := key + ":"
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		v := strings.TrimSpace(line[len(prefix):])
		v = strings.Trim(v, `"'`)
		return v
	}
	return ""
}
