package portfolio

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Project struct {
	Name   string
	Desc   string
	Tech   []string
	URL    string
	Status string
}

var projects = []Project{
	{
		Name:   "ssh.koossaayy.tn",
		Desc:   "This very SSH portal — built with Charm's wish + bubbletea stack. Because why not.",
		Tech:   []string{"Go", "Wish", "BubbleTea", "Lipgloss"},
		URL:    "ssh ssh.koossaayy.tn",
		Status: "🟢 Live",
	},
	{
		Name:   "Project Alpha",
		Desc:   "A mysterious project shrouded in secrecy. Details classified.",
		Tech:   []string{"Rust", "WebAssembly"},
		URL:    "github.com/koossaayy/alpha",
		Status: "🔵 In Progress",
	},
	{
		Name:   "Homelab",
		Desc:   "A self-hosted everything setup powered by Coolify, Docker, and pure stubbornness.",
		Tech:   []string{"Coolify", "Docker", "Nginx", "Cloudflare"},
		URL:    "https://koossaayy.tn",
		Status: "🟢 Live",
	},
	{
		Name:   "Terminal Experiments",
		Desc:   "A graveyard of CLI tools, shell scripts, and TUI experiments.",
		Tech:   []string{"Go", "Python", "Bash"},
		URL:    "github.com/koossaayy",
		Status: "🟡 Ongoing",
	},
}

type Model struct {
	renderer *lipgloss.Renderer
	width    int
	height   int
	cursor   int
}

func New(r *lipgloss.Renderer, w, h int) Model {
	return Model{renderer: r, width: w, height: h, cursor: 0}
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
			if m.cursor < len(projects)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	r := m.renderer

	if m.width == 0 {
		m.width = 80
	}

	pink   := lipgloss.Color("#FF79C6")
	cyan   := lipgloss.Color("#8BE9FD")
	yellow := lipgloss.Color("#F1FA8C")
	green  := lipgloss.Color("#50FA7B")
	fg     := lipgloss.Color("#F8F8F2")
	subtle := lipgloss.Color("#6272A4")
	purple := lipgloss.Color("#9B72CF")

	titleStyle := r.NewStyle().Foreground(pink).Bold(true)
	footStyle  := r.NewStyle().Foreground(subtle).Italic(true)

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(titleStyle.Render("  🚀 Portfolio"))
	sb.WriteString("\n")
	sb.WriteString(r.NewStyle().Foreground(subtle).Italic(true).Render("  Things I've built, broken, and learned from."))
	sb.WriteString("\n\n")

	for i, p := range projects {
		isSelected := i == m.cursor
		borderColor := subtle
		if isSelected {
			borderColor = pink
		}

		cardStyle := r.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(0, 2).
			Width(m.width - 8)

		var tags []string
		for _, t := range p.Tech {
			tag := r.NewStyle().
				Foreground(lipgloss.Color("#282A36")).
				Background(purple).
				Padding(0, 1).
				Render(t)
			tags = append(tags, tag)
		}
		techLine := strings.Join(tags, " ")

		nameStyle := r.NewStyle().Foreground(cyan).Bold(true)
		if isSelected {
			nameStyle = r.NewStyle().Foreground(yellow).Bold(true)
		}

		content := fmt.Sprintf(
			"%s  %s\n%s\n%s  %s\n%s",
			nameStyle.Render(p.Name),
			r.NewStyle().Foreground(green).Render(p.Status),
			r.NewStyle().Foreground(fg).Render(p.Desc),
			r.NewStyle().Foreground(subtle).Render("🔗"),
			r.NewStyle().Foreground(cyan).Render(p.URL),
			techLine,
		)

		sb.WriteString(cardStyle.Render(content))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(footStyle.Render("  ↑↓ / j k to browse  •  esc to go back"))
	return sb.String()
}