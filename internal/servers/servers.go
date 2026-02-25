package servers

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Server struct {
	Name string
	Host string
	Desc string
	Icon string
	Tag  string
}

// ── Add your servers here! ──────────────────────────────────────────────────
var serverList = []Server{
	{
		Name: "This Portal",
		Host: "ssh ssh.koossaayy.tn -p 2222",
		Desc: "You are here. Very meta.",
		Icon: "🌀",
		Tag:  "portal",
	},
	{
		Name: "Main Server",
		Host: "ssh koossaayy.tn",
		Desc: "The homelab overlord. Runs everything.",
		Icon: "🖥️",
		Tag:  "homelab",
	},
	{
		Name: "Dev Box",
		Host: "ssh dev.koossaayy.tn",
		Desc: "Where code goes to be born (and sometimes die).",
		Icon: "💻",
		Tag:  "dev",
	},
	{
		Name: "Staging",
		Host: "ssh staging.koossaayy.tn",
		Desc: "It works on staging, I swear.",
		Icon: "🧪",
		Tag:  "staging",
	},
}

// ── Model ────────────────────────────────────────────────────────────────────

type Model struct {
	renderer *lipgloss.Renderer
	width    int
	height   int
	cursor   int
}

func New(r *lipgloss.Renderer, w, h int) Model {
	return Model{renderer: r, width: w, height: h}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(serverList)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	r := m.renderer

	pink   := lipgloss.Color("#FF79C6")
	cyan   := lipgloss.Color("#8BE9FD")
	yellow := lipgloss.Color("#F1FA8C")
	fg     := lipgloss.Color("#F8F8F2")
	subtle := lipgloss.Color("#6272A4")
	green  := lipgloss.Color("#50FA7B")
	purple := lipgloss.Color("#9B72CF")

	titleStyle := r.NewStyle().Foreground(pink).Bold(true)
	footStyle  := r.NewStyle().Foreground(subtle).Italic(true)

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(titleStyle.Render("  🖧  Server Directory"))
	sb.WriteString("\n")
	sb.WriteString(r.NewStyle().Foreground(subtle).Italic(true).Render("  SSH into the machines of the realm.\n\n"))

	// Header row
	headerStyle := r.NewStyle().Foreground(subtle).Bold(true)
	sb.WriteString(fmt.Sprintf("  %-3s  %-18s  %-32s  %s\n", "", headerStyle.Render("NAME"), headerStyle.Render("COMMAND"), headerStyle.Render("DESCRIPTION")))
	sb.WriteString(r.NewStyle().Foreground(subtle).Render("  " + strings.Repeat("─", m.width-6) + "\n"))

	for i, s := range serverList {
		isSelected := i == m.cursor

		nameStyle := r.NewStyle().Foreground(fg)
		cmdStyle  := r.NewStyle().Foreground(cyan)
		rowPrefix := "  "

		if isSelected {
			nameStyle = r.NewStyle().Foreground(yellow).Bold(true)
			cmdStyle  = r.NewStyle().Foreground(green).Bold(true)
			rowPrefix = r.NewStyle().Foreground(pink).Bold(true).Render("▸ ")
		}

		tagStyle := r.NewStyle().
			Foreground(lipgloss.Color("#282A36")).
			Background(purple).
			Padding(0, 1)

		line := fmt.Sprintf("%s%s  %-18s  %-32s  %s  %s",
			rowPrefix,
			s.Icon,
			nameStyle.Render(s.Name),
			cmdStyle.Render(s.Host),
			r.NewStyle().Foreground(subtle).Italic(true).Render(s.Desc),
			tagStyle.Render(s.Tag),
		)
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(r.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purple).
		Padding(0, 2).
		Render(r.NewStyle().Foreground(yellow).Render("💡 Tip: ") + r.NewStyle().Foreground(fg).Render("Copy the command and run it in a new terminal to connect!")))
	sb.WriteString("\n\n")
	sb.WriteString(footStyle.Render("  ↑↓ / j k to browse  •  esc to go back"))

	return sb.String()
}
