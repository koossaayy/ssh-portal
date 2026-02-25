package ui

import (
	"fmt"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/koossaayy/ssh-portal/internal/game"
	"github.com/koossaayy/ssh-portal/internal/portfolio"
)

type view int

const (
	viewHome view = iota
	viewAbout
	viewPortfolio
	viewGame
)

const banner = `
 ██╗  ██╗ ██████╗  ██████╗ ███████╗███████╗ █████╗  █████╗ ██╗   ██╗██╗   ██╗
 ██║ ██╔╝██╔═══██╗██╔═══██╗██╔════╝██╔════╝██╔══██╗██╔══██╗╚██╗ ██╔╝╚██╗ ██╔╝
 █████╔╝ ██║   ██║██║   ██║███████╗███████╗███████║███████║ ╚████╔╝  ╚████╔╝ 
 ██╔═██╗ ██║   ██║██║   ██║╚════██║╚════██║██╔══██║██╔══██║  ╚██╔╝    ╚██╔╝  
 ██║  ██╗╚██████╔╝╚██████╔╝███████║███████║██║  ██║██║  ██║   ██║      ██║   
 ╚═╝  ╚═╝ ╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝      ╚═╝  `

var pirateQuotes = []string{
	"\"Not all treasure is silver and gold, mate.\" 🏴‍☠️",
	"\"This is the day you will always remember as the day you almost caught Captain Jack Sparrow.\" 🦜",
	"\"Why is the rum always gone? ...Oh, that's why.\" 🥃",
	"\"The problem is not the problem. The problem is your attitude about the problem.\" ☠️",
	"\"Me? I'm dishonest. And a dishonest man you can always trust to be dishonest.\" 🧭",
	"\"Nobody move! I dropped me brain.\" 💀",
	"\"I love those moments. I like to wave at them as they pass by.\" 🌊",
	"\"Did everyone see that? Because I will not be doing it again.\" 🪝",
	"\"You seem somewhat familiar. Have I threatened you before?\" ⚔️",
	"\"Wherever we want to go, we go.\" 🗺️",
	"\"“UP IS DOWN”? Well that's just maddeningly unhelpful. Why are these things never clear?\" 😕",
	"\"I've got a jar of dirt\" ⚔️",
	"\"Crazy people don't know they're crazy. I know that I'm crazy, therefore I'm not crazy. Isn't that crazy?\" 😀",
	"\"Why fight when you can negotiate?\" 🫙",
	"\"Stop blowing holes in my ship!!\" ⚓",
	"\"No! Not good! Stop! Not good! What are you doing? You burned all the food, the shade... the rum\" 🍺",

}

type menuItem struct {
	label string
	icon  string
	desc  string
	view  view
}

var menuItems = []menuItem{
	{"About & Welcome", "👋", "Who is this mysterious person?", viewAbout},
	{"Portfolio", "🚀", "Projects, work, and cool stuff", viewPortfolio},
	{"Play Snake!", "🐍", "Take a break, you deserve it", viewGame},
}

type MainModel struct {
	renderer  *lipgloss.Renderer
	width     int
	height    int
	cursor    int
	current   view
	portfolio portfolio.Model
	game      game.Model
	quote     string
}

func NewMainModel(renderer *lipgloss.Renderer, w, h int) MainModel {
	return MainModel{
		renderer:  renderer,
		width:     w,
		height:    h,
		cursor:    0,
		current:   viewHome,
		portfolio: portfolio.New(renderer, w, h),
		game:      game.New(renderer, w, h),
		quote:     pirateQuotes[rand.Intn(len(pirateQuotes))],
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.current == viewHome {
				return m, tea.Quit
			}
			m.current = viewHome
			return m, nil
		case "esc":
			if m.current != viewHome {
				m.current = viewHome
				return m, nil
			}
		}

		if m.current == viewPortfolio {
			updated, cmd := m.portfolio.Update(msg)
			m.portfolio = updated.(portfolio.Model)
			return m, cmd
		}
		if m.current == viewGame {
			updated, cmd := m.game.Update(msg)
			m.game = updated.(game.Model)
			if m.game.Quit {
				m.current = viewHome
				m.game = game.New(m.renderer, m.width, m.height)
			}
			return m, cmd
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(menuItems)-1 {
				m.cursor++
			}
		case "enter", " ":
			selected := menuItems[m.cursor]
			m.current = selected.view
			if selected.view == viewGame {
				m.game = game.New(m.renderer, m.width, m.height)
				return m, m.game.Init()
			}
		}
	}

	if m.current == viewGame {
		updated, cmd := m.game.Update(msg)
		m.game = updated.(game.Model)
		return m, cmd
	}

	return m, nil
}

func (m MainModel) View() string {
	switch m.current {
	case viewAbout:
		return m.aboutView()
	case viewPortfolio:
		return m.portfolio.View()
	case viewGame:
		return m.game.View()
	default:
		return m.homeView()
	}
}

func (m MainModel) homeView() string {
	r := m.renderer

	pink   := lipgloss.Color("#FF79C6")
	cyan   := lipgloss.Color("#8BE9FD")
	yellow := lipgloss.Color("#F1FA8C")
	fg     := lipgloss.Color("#F8F8F2")
	subtle := lipgloss.Color("#6272A4")

	bannerStyle  := r.NewStyle().Foreground(pink).Bold(true)
	taglineStyle := r.NewStyle().Foreground(cyan).Italic(true)
	selStyle     := r.NewStyle().Foreground(pink).Bold(true)
	normalStyle  := r.NewStyle().Foreground(fg)
	descStyle    := r.NewStyle().Foreground(subtle).Italic(true)
	footStyle    := r.NewStyle().Foreground(subtle).Italic(true)

	var sb strings.Builder
	sb.WriteString("\n")

	if m.width > 90 {
		sb.WriteString(bannerStyle.Render(banner))
	} else {
		sb.WriteString(bannerStyle.Render("  ✦ ssh.koossaayy.tn ✦"))
	}
	sb.WriteString("\n")
	sb.WriteString(taglineStyle.Render("  " + m.quote))
	sb.WriteString("\n\n")

	sb.WriteString(r.NewStyle().Foreground(yellow).Bold(true).Render("  Navigate"))
	sb.WriteString("\n")
	for i, item := range menuItems {
		line := fmt.Sprintf("%s  %s", item.icon, item.label)
		if i == m.cursor {
			sb.WriteString(selStyle.Render("  ▸ " + line))
			sb.WriteString("  " + descStyle.Render(item.desc))
		} else {
			sb.WriteString(normalStyle.Render("    " + line))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(footStyle.Render("  ↑↓ / j k to move  •  enter to select  •  esc / q to go back"))

	return sb.String()
}

func (m MainModel) aboutView() string {
	r := m.renderer

	pink   := lipgloss.Color("#FF79C6")
	cyan   := lipgloss.Color("#8BE9FD")
	yellow := lipgloss.Color("#F1FA8C")
	fg     := lipgloss.Color("#F8F8F2")
	subtle := lipgloss.Color("#6272A4")
	purple := lipgloss.Color("#9B72CF")

	titleStyle := r.NewStyle().Foreground(pink).Bold(true)
	footStyle  := r.NewStyle().Foreground(subtle).Italic(true)
	boxStyle   := r.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purple).
		Padding(1, 3).
		Width(m.width - 8)
	labelStyle := r.NewStyle().Foreground(yellow).Bold(true)
	valStyle   := r.NewStyle().Foreground(fg)
	hlStyle    := r.NewStyle().Foreground(cyan).Bold(true)

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(titleStyle.Render("  👋 About & Welcome"))
	sb.WriteString("\n\n")

	whoami := fmt.Sprintf(
		"%s\n\n%s\n%s\n%s\n\n%s\n%s\n\n%s\n%s %s %s %s %s %s %s",
		hlStyle.Render("  Hey, I'm Koossaayy! 👾"),
		labelStyle.Render("  What I do:"),
		valStyle.Render("  Developer, homelab nerd, terminal maximalist. I build things,"),
		valStyle.Render("  break them, learn why, and repeat. Because why not 🤷‍♂️"),
		labelStyle.Render("  Currently into:"),
		valStyle.Render("  Laravel, Serious DevSecOps, Self-hosting everything, Go, CLI aesthetics."),
		labelStyle.Render("  Stack:"),
		r.NewStyle().Background(lipgloss.Color("#F55673")).Foreground(fg).Bold(true).Padding(0, 1).Render("Laravel & PHP (I mean of course)"),
		r.NewStyle().Background(lipgloss.Color("#6272FF")).Foreground(fg).Bold(true).Padding(0, 1).Render("Finetuning Models"),
		r.NewStyle().Background(lipgloss.Color("#F1FA8C")).Foreground(lipgloss.Color("#282A36")).Bold(true).Padding(0, 1).Render("React / JS"),
		r.NewStyle().Background(lipgloss.Color("#00ADD8")).Foreground(fg).Bold(true).Padding(0, 1).Render("Go"),
		r.NewStyle().Background(lipgloss.Color("#2496ED")).Foreground(fg).Bold(true).Padding(0, 1).Render("Docker (I hate it though)"),
		r.NewStyle().Background(lipgloss.Color("#FFA500")).Foreground(lipgloss.Color("#282A36")).Bold(true).Padding(0, 1).Render("Linux but mostly Windows (MacOS soon)"),
		r.NewStyle().Background(lipgloss.Color("#7B42BC")).Foreground(fg).Bold(true).Padding(0, 1).Render("Coolify FTW"),
	)
	sb.WriteString(boxStyle.Render(whoami))
	sb.WriteString("\n\n")

	links := fmt.Sprintf(
		"%s\n%s  %s\n%s  %s\n%s  %s\n%s  %s\n%s  %s",
		labelStyle.Render("  Find me:"),
		r.NewStyle().Foreground(cyan).Render("  🌐 Web   "), valStyle.Render("https://koossaayy.tn"),
		r.NewStyle().Foreground(cyan).Render("  🐙 GitHub"), valStyle.Render("https://github.com/koossaayy"),
		r.NewStyle().Foreground(cyan).Render("  🐦 Twitter"), valStyle.Render("https://x.com/koossaayy"),
		r.NewStyle().Foreground(cyan).Render("  🔗 LinkedIn"), valStyle.Render("https://www.linkedin.com/in/koossaayy/"),
		r.NewStyle().Foreground(cyan).Render("  📡 SSH   "), valStyle.Render("ssh ssh.koossaayy.tn -p 69  ← yes, port 69. yes, on purpose. Thank you"),
	)
	sb.WriteString(boxStyle.Render(links))

	sb.WriteString("\n\n")
	sb.WriteString(footStyle.Render("  esc / q to go back"))

	return sb.String()
}