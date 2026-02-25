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
	Emoji  string
}

var projects = []Project{
	{
		Name:   "SSH Portal",
		Desc:   "This very SSH portal — built with Charm's wish + bubbletea stack. Because why not.",
		Tech:   []string{"Go", "Wish", "BubbleTea", "Lipgloss"},
		URL:    "ssh ssh.koossaayy.tn -p 69",
		Status: "Live",
		Emoji:  "🟢",
	},
	{
		Name:   "Laralingo",
		Desc:   "Manage your localization & translation process as code, and never miss anything.",
		Tech:   []string{"Laravel", "React", "Inertia", "GitHub", "GitLab", "AI & Translation APIs"},
		URL:    "laralingo.app",
		Status: "Closed Preview",
		Emoji:  "🔵",
	},
	{
		Name:   "Personal Blog",
		Desc:   "Well, I got to write something somewhere right?",
		Tech:   []string{"Laravel", "Statamic"},
		URL:    "https://koossaayy.tn",
		Status: "Live",
		Emoji:  "🟢",
	},
	{
		Name:   "Devs.tn",
		Desc:   "Linktree but for Tunisian devs. Because we deserve our own corner of the internet.",
		Tech:   []string{"Laravel", "React", "Inertia"},
		URL:    "devs.tn",
		Status: "Ongoing",
		Emoji:  "🟡",
	},
	{
		Name:   "SUPER DUPER SECRET PROJECT",
		Desc:   "SUPER DUPER SECRET DESCRIPTION. But lawyers gonna love it. 🤫",
		Tech:   []string{"Laravel", "React", "Inertia"},
		URL:    "¯\\_(ツ)_/¯",
		Status: "Ongoing",
		Emoji:  "🔴",
	},
}

// Tech badge colors — each tech gets its own vibe
var techColors = map[string]struct{ bg, fg string }{
	"Go":         {"#00ADD8", "#FFFFFF"},
	"Wish":       {"#BD93F9", "#FFFFFF"},
	"BubbleTea":  {"#FF79C6", "#282A36"},
	"Lipgloss":   {"#FF5555", "#FFFFFF"},
	"Laravel":    {"#F55673", "#FFFFFF"},
	"React":      {"#61DAFB", "#282A36"},
	"Inertia":    {"#9553E9", "#FFFFFF"},
	"GitHub":     {"#F8F8F2", "#282A36"},
	"GitLab":     {"#FC6D26", "#FFFFFF"},
	"AI APIs":    {"#F1FA8C", "#282A36"},
	"Statamic":   {"#0069FF", "#FFFFFF"},
	"Docker":     {"#2496ED", "#FFFFFF"},
	"Nginx":      {"#009639", "#FFFFFF"},
	"Cloudflare": {"#F48120", "#FFFFFF"},
	"Rust":       {"#CE422B", "#FFFFFF"},
}

var statusColors = map[string]string{
	"Live":        "#50FA7B",
	"In Progress": "#8BE9FD",
	"Ongoing":     "#F1FA8C",
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
	fg     := lipgloss.Color("#F8F8F2")
	subtle := lipgloss.Color("#6272A4")

	titleStyle := r.NewStyle().Foreground(pink).Bold(true)
	footStyle  := r.NewStyle().Foreground(subtle).Italic(true)
	countStyle := r.NewStyle().Foreground(subtle)

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(titleStyle.Render("  🚀 Portfolio"))
	sb.WriteString("  ")
	sb.WriteString(countStyle.Render(fmt.Sprintf("(%d/%d)", m.cursor+1, len(projects))))
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

		// Status badge
		statusColor := lipgloss.Color("#6272A4")
		if c, ok := statusColors[p.Status]; ok {
			statusColor = lipgloss.Color(c)
		}
		statusBadge := r.NewStyle().
			Foreground(lipgloss.Color("#282A36")).
			Background(statusColor).
			Bold(true).
			Padding(0, 1).
			Render(p.Emoji + " " + p.Status)

		// Name style
		nameStyle := r.NewStyle().Foreground(cyan).Bold(true)
		if isSelected {
			nameStyle = r.NewStyle().Foreground(yellow).Bold(true)
		}

		// Tech badges with per-tech colors
		var tags []string
		for _, t := range p.Tech {
			bgHex := "#6272A4"
			fgHex := "#F8F8F2"
			if colors, ok := techColors[t]; ok {
				bgHex = colors.bg
				fgHex = colors.fg
			}
			tag := r.NewStyle().
				Foreground(lipgloss.Color(fgHex)).
				Background(lipgloss.Color(bgHex)).
				Bold(true).
				Padding(0, 1).
				Render(t)
			tags = append(tags, tag)
		}
		techLine := strings.Join(tags, " ")

		content := fmt.Sprintf(
			"%s  %s\n\n%s\n\n%s  %s\n%s",
			nameStyle.Render(p.Name),
			statusBadge,
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