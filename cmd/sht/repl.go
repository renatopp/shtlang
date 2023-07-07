package repl

import (
	"fmt"
	"os"
	"regexp"
	"sht/lang"
	"sht/lang/runtime"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	accentColor       = lipgloss.AdaptiveColor{Light: "#FFA41B", Dark: "#FFA41B"}
	promptChar        = "→ "
	promptPlaceholder = "Enter command"
	// cmdChar = "⟜⥁🜂▵Ⲇ "
	// cmdChar     = "[sht] "
	cmdChar     = "▵ "
	spinnerType = spinner.Points

	promptStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			BorderLeft(false).
			BorderRight(false).
			BorderBottom(false).
			BorderTop(true).
			Padding(0, 1)
	}()
	contentStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().
			Padding(0, 1)
	}()
	contentCmdStyle = func() lipgloss.Style {
		return lipgloss.NewStyle()
	}()
	contentResultStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().
			Foreground(accentColor)
	}()
	contentAnnounceStyle = func() lipgloss.Style {
		return lipgloss.NewStyle()
	}()
	contentErrorStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555"))
	}()
	loadingStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.HiddenBorder()).
			Padding(0, 1).
			Foreground(accentColor)
	}()
)

type model struct {
	ready  bool
	width  int
	height int

	prompt   textarea.Model
	viewport viewport.Model
	spinner  spinner.Model

	history []string
	content string
	builder strings.Builder

	cursor int
	buffer string

	runtime *runtime.Runtime
}

func Start(r *runtime.Runtime) {
	fmt.Print("\033[H\033[2J") // Clear screen
	prompt := textarea.New()
	prompt.Placeholder = promptPlaceholder
	prompt.Prompt = promptChar
	prompt.SetHeight(1)
	prompt.ShowLineNumbers = false
	prompt.Focus()

	viewport := viewport.Model{}

	spin := spinner.New()
	spin.Spinner = spinnerType

	m := model{
		prompt:   prompt,
		viewport: viewport,
		spinner:  spin,
		history:  []string{},
		content:  "",
		builder:  strings.Builder{},
		runtime:  r,
	}

	if _, err := tea.NewProgram(&m).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (m *model) Init() tea.Cmd {
	if m.ready {
		return textinput.Blink
	}

	return m.spinner.Tick
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, m.onKey(msg)
	case tea.WindowSizeMsg:
		return m, m.onResize()
	case spinner.TickMsg:
		return m, m.onTick(msg)
	}

	return m, nil
}

func (m *model) View() string {
	if !m.ready {
		return loadingStyle.Width(m.width).Height(m.height).Render(
			lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "Loading "+m.spinner.View()),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		contentStyle.Width(m.width).Render(m.viewport.View()),
		promptStyle.Width(m.width).Render(m.prompt.View()),
	)
}

func (m *model) onResize() tea.Cmd {
	var cmd tea.Cmd

	w, h, _ := term.GetSize(0)
	m.width = w - 1
	m.height = h

	if !m.ready {
		m.ready = true
	}

	rh := lipgloss.Height(promptStyle.Render(m.prompt.View()))

	m.viewport.YPosition = 0
	m.viewport.Width = m.width
	m.viewport.Height = m.height - rh
	// m.prompt.Width = m.width
	m.prompt.SetWidth(m.width - 5)

	return cmd
}

func (m *model) onTick(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return cmd
}

var doubleNewLines, _ = regexp.Compile(`\n\n`)

func (m *model) onKey(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyEnter:

		val := m.prompt.Value()
		if areBracketsBalanced(val) || doubleNewLines.MatchString(val) {
			return m.onCommand(val)
		}

	case tea.KeyCtrlC:
		return tea.Quit

	case tea.KeyCtrlUp:
		return m.onHistory(1)

	case tea.KeyCtrlDown:
		return m.onHistory(1)

	case tea.KeyShiftDown:
		return m.onHistory(-1)

	case tea.KeyShiftUp:
		m.viewport, cmd = m.viewport.Update(tea.KeyMsg{Type: tea.KeyUp})
		return cmd

	case tea.KeyEsc:
		if m.cursor != -1 {
			m.cursor = 0
			return m.onHistory(-1)
		}

	}

	m.cursor = -1
	add := 1
	if msg.Type == tea.KeyEnter {
		add += 1
	} else if msg.Type == tea.KeyBackspace {
		add -= 1
	}
	m.prompt.SetHeight(strings.Count(m.prompt.Value(), "\n") + add)
	m.prompt, cmd = m.prompt.Update(msg)
	m.onResize()

	return cmd
}

func (m *model) onHistory(dir int) tea.Cmd {
	if m.cursor == -1 {
		m.buffer = m.prompt.Value()
	}

	m.cursor += dir
	size := len(m.history)

	if m.cursor < -1 {
		m.cursor = -1
	} else if m.cursor >= size {
		m.cursor = size - 1
	}

	if m.cursor == -1 {
		m.prompt.SetValue(m.buffer)
		m.prompt.CursorEnd()
		return nil
	}

	m.prompt.SetValue(m.history[size-m.cursor-1])
	m.prompt.CursorEnd()

	return nil
}

func (m *model) onCommand(cmd string) tea.Cmd {
	if cmd == "" {
		return nil
	}

	m.prompt.SetValue("")
	m.buffer = ""
	m.cursor = -1
	m.prompt.SetHeight(1)
	m.onResize()
	m.prompt, _ = m.prompt.Update(tea.KeyLeft)

	if cmd == "exit" {
		return tea.Quit
	}

	if cmd == "clear" {
		m.history = []string{}
		m.content = ""
		m.builder.Reset()
		m.viewport.SetContent(m.content)
		m.prompt.SetValue("")
		return nil
	}

	if cmd == "help" {
		m.appendAnnounce("help!")
		return nil
	}

	m.history = append(m.history, cmd)
	m.appendCommand(cmd)

	tree, err := lang.Parse([]byte(cmd))
	if err != nil {
		m.appendError(err.Error())
		return nil
	}

	res, err := m.runtime.Run(tree)
	if err != nil {
		m.appendResult(res)
	}

	return nil
}

func (m *model) appendCommand(x string) {
	cmd := contentCmdStyle.Render(cmdChar + x + "\n")

	m.builder.WriteString(strings.TrimRight(cmd, " "))
	m.content = m.builder.String()
	m.viewport.SetContent(m.content)
	m.viewport.GotoBottom()
}

func (m *model) appendResult(x string) {
	x = strings.ReplaceAll(x, "\n", "\n")
	cmd := contentResultStyle.Render("" + x + "\n")

	m.builder.WriteString(strings.TrimRight(cmd, " "))
	m.content = m.builder.String()
	m.viewport.SetContent(m.content)
	m.viewport.GotoBottom()
}

func (m *model) appendError(x string) {
	x = strings.ReplaceAll(x, "\n", "\n")
	cmd := contentErrorStyle.Render("" + x + "\n")

	m.builder.WriteString(strings.TrimRight(cmd, " "))
	m.content = m.builder.String()
	m.viewport.SetContent(m.content)
	m.viewport.GotoBottom()
}

func (m *model) appendAnnounce(x string) {
	cmd := contentAnnounceStyle.Render(x + "\n")
	m.builder.WriteString(strings.TrimRight(cmd, " "))
	m.content = m.builder.String()
	m.viewport.SetContent(m.content)
	m.viewport.GotoBottom()
}
