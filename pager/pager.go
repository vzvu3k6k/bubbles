package pager

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/pkg/errors"
)

/* TODO:
- add rendered markdown
- show status (similar to paginator?)
- add search functionality - similar to neovim
*/

const useHighPerformanceRenderer = false

type Model struct {
    content string
    ready bool 
    viewport viewport.Model
    errors []error
}

func New(content string) Model {
    return Model{content: content}
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var (
        cmd tea.Cmd
        cmds []tea.Cmd
    )
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        if !m.ready {
            m.viewport = viewport.New(msg.Width, msg.Height)
            m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
            renderedContent, err := m.renderContent(msg.Width)
            if err != nil {
                m.errors = append(m.errors, err)
            }
            m.viewport.SetContent(renderedContent)
            m.ready = true
        } else {
            m.viewport.Width = msg.Width
            m.viewport.Height = msg.Height
        }
        
        if useHighPerformanceRenderer {
            cmds = append(cmds, viewport.Sync(m.viewport))
        }
    }
    m.viewport, cmd = m.viewport.Update(msg)
    cmds = append(cmds, cmd)
    return m, tea.Batch(cmds...)
    // TODO: scrolling
    // TODO: filtering
}

func (m Model) renderContent(width int) (string, error) {
    r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
        return "", errors.Wrap(err, "could not init glamour renderer")
	}
    rendered, err := r.Render(m.content)
	if err != nil {
        return "", errors.Wrap(err, "could not render content")
	}
    return rendered, nil
}

func (m Model) View() string {
    return m.viewport.View()
}
