package ui

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/vertexai/genai"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/tehbooom/gobot/internal/vertex"
)

type Model struct {
	ctx         context.Context
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	aiModel     *genai.ChatSession
	senderStyle lipgloss.Style
	vertexStyle lipgloss.Style
	err         error
	renderer    *glamour.TermRenderer
}

func InitialModel(project_id, region, model_name string) Model {
	ctx := context.Background()
	c := vertex.Client(project_id, region, ctx)
	geminiModel := c.GenerativeModel(model_name)
	ta := textarea.New()
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(50, 5)
	vp.SetContent(`Start chatting!`)
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return Model{
		textarea:    ta,
		ctx:         ctx,
		messages:    []string{},
		viewport:    vp,
		aiModel:     geminiModel.StartChat(),
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		vertexStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
		err:         nil,
	}
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{tea.EnterAltScreen}
	return tea.Batch(cmds...)
}

type (
	errMsg error
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:

			return m, tea.Quit
		case tea.KeyEnter:
			// m.messages = append(m.messages, m.RenderConversation(m.viewport.Width, m.textarea.Value()))
			m.viewport.SetContent(m.RenderConversation(m.viewport.Width, m.textarea.Value()))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func (m Model) RenderConversation(maxWidth int, question string) string {
	var sb strings.Builder
	resp := vertex.Send(m.ctx, m.aiModel, m.textarea.Value())

	sb.WriteString(m.senderStyle.Render("You: "))
	youContent := wordwrap.String(question, maxWidth-5)
	sb.WriteString(youContent)

	content := vertex.Response(resp)
	sb.WriteString(m.vertexStyle.Render("Vertex: "))
	content = wordwrap.String(content, maxWidth-5)
	sb.WriteString(content)
	return sb.String()
}
