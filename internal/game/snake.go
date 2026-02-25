package game

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type point struct{ x, y int }

type direction int

const (
	dirUp direction = iota
	dirDown
	dirLeft
	dirRight
)

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type gameState int

const (
	statePlaying gameState = iota
	stateGameOver
)

type Model struct {
	renderer  *lipgloss.Renderer
	width     int
	height    int
	boardW    int
	boardH    int
	snake     []point
	dir       direction
	nextDir   direction
	food      point
	score     int
	highScore int
	state     gameState
	Quit      bool
}

func New(r *lipgloss.Renderer, w, h int) Model {
	m := Model{
		renderer: r,
		width:    w,
		height:   h,
		boardW:   60,
		boardH:   25,
		dir:      dirRight,
		nextDir:  dirRight,
		state:    statePlaying,
	}
	m.reset()
	return m
}

func (m *Model) reset() {
	cx, cy := m.boardW/2, m.boardH/2
	m.snake = []point{
		{cx, cy},
		{cx - 1, cy},
		{cx - 2, cy},
	}
	m.dir = dirRight
	m.nextDir = dirRight
	m.score = 0
	m.state = statePlaying
	m.spawnFood()
}

func (m *Model) spawnFood() {
	occupied := map[point]bool{}
	for _, s := range m.snake {
		occupied[s] = true
	}
	for {
		p := point{rand.Intn(m.boardW), rand.Intn(m.boardH)}
		if !occupied[p] {
			m.food = p
			return
		}
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.Quit = true
			return m, nil
		case "up", "k", "w":
			if m.dir != dirDown {
				m.nextDir = dirUp
			}
		case "down", "j", "s":
			if m.dir != dirUp {
				m.nextDir = dirDown
			}
		case "left", "h", "a":
			if m.dir != dirRight {
				m.nextDir = dirLeft
			}
		case "right", "l", "d":
			if m.dir != dirLeft {
				m.nextDir = dirRight
			}
		case "enter", " ":
			if m.state == stateGameOver {
				m.reset()
				return m, tick()
			}
		}

	case tickMsg:
		if m.state == stateGameOver {
			return m, nil
		}
		m.dir = m.nextDir

		head := m.snake[0]
		var newHead point
		switch m.dir {
		case dirUp:
			newHead = point{head.x, head.y - 1}
		case dirDown:
			newHead = point{head.x, head.y + 1}
		case dirLeft:
			newHead = point{head.x - 1, head.y}
		case dirRight:
			newHead = point{head.x + 1, head.y}
		}

		if newHead.x < 0 || newHead.x >= m.boardW || newHead.y < 0 || newHead.y >= m.boardH {
			m.state = stateGameOver
			if m.score > m.highScore {
				m.highScore = m.score
			}
			return m, nil
		}

		for _, s := range m.snake {
			if s == newHead {
				m.state = stateGameOver
				if m.score > m.highScore {
					m.highScore = m.score
				}
				return m, nil
			}
		}

		ate := newHead == m.food
		m.snake = append([]point{newHead}, m.snake...)
		if ate {
			m.score++
			m.spawnFood()
		} else {
			m.snake = m.snake[:len(m.snake)-1]
		}

		return m, tick()
	}

	return m, nil
}

func (m Model) View() string {
	r := m.renderer

	pink   := lipgloss.Color("#FF79C6")
	cyan   := lipgloss.Color("#8BE9FD")
	yellow := lipgloss.Color("#F1FA8C")
	green  := lipgloss.Color("#50FA7B")
	red    := lipgloss.Color("#FF5555")
	fg     := lipgloss.Color("#F8F8F2")
	subtle := lipgloss.Color("#6272A4")
	purple := lipgloss.Color("#9B72CF")

	grid := make([][]rune, m.boardH)
	for y := range grid {
		grid[y] = make([]rune, m.boardW)
		for x := range grid[y] {
			grid[y][x] = ' '
		}
	}

	for i, s := range m.snake {
		if s.x >= 0 && s.x < m.boardW && s.y >= 0 && s.y < m.boardH {
			if i == 0 {
				grid[s.y][s.x] = '●'
			} else {
				grid[s.y][s.x] = '○'
			}
		}
	}

	grid[m.food.y][m.food.x] = '❤'

boardStyle := r.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(purple).
    Padding(1, 1).  // ← change Padding(0, 1) to Padding(1, 1)
    MarginTop(2)    // ← add this

	var boardSb strings.Builder
	for y, row := range grid {
		for x, cell := range row {
			pt := point{x, y}
			switch {
			case pt == m.snake[0]:
				boardSb.WriteString(r.NewStyle().Foreground(green).Bold(true).Render(string(cell)))
			case m.isSnakeBody(pt):
				boardSb.WriteString(r.NewStyle().Foreground(cyan).Render(string(cell)))
			case pt == m.food:
				boardSb.WriteString(r.NewStyle().Foreground(red).Render(string(cell)))
			default:
				boardSb.WriteString(string(cell))
			}
		}
		if y < m.boardH-1 {
			boardSb.WriteString("\n")
		}
	}

	board := boardStyle.Render(boardSb.String())

	statsStyle := r.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cyan).
		Padding(1, 2).
	    MarginTop(2).   // ← add this
		Width(18)

	stats := fmt.Sprintf(
		"%s\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s",
		r.NewStyle().Foreground(yellow).Bold(true).Render("🐍 SNAKE"),
		r.NewStyle().Foreground(subtle).Render(""),
		r.NewStyle().Foreground(subtle).Render("SCORE"),
		r.NewStyle().Foreground(pink).Bold(true).Render(fmt.Sprintf(" %d", m.score)),
		r.NewStyle().Foreground(subtle).Render("HIGH SCORE"),
		r.NewStyle().Foreground(yellow).Bold(true).Render(fmt.Sprintf(" %d", m.highScore)),
		r.NewStyle().Foreground(subtle).Render("w a s d\n↑ ↓ ← →\nh j k l"),
	)

	statsPanel := statsStyle.Render(stats)
	gameArea := lipgloss.JoinHorizontal(lipgloss.Top, board, "  ", statsPanel)

	var sb strings.Builder
	sb.WriteString("\n\n\n")
sb.WriteString(r.NewStyle().Foreground(pink).Bold(true).Render("  🎮 Snake — take a break!"))
sb.WriteString("\n\n\n")  // ← was \n\n, add one more \n here
sb.WriteString("  ")
sb.WriteString(gameArea)
	sb.WriteString("\n")

	if m.state == stateGameOver {
		overlay := r.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(red).
			Padding(1, 4).
		    MarginTop(2).   // ← add this
			Render(fmt.Sprintf(
				"%s\n%s\n%s",
				r.NewStyle().Foreground(red).Bold(true).Render("  💀 GAME OVER  "),
				r.NewStyle().Foreground(fg).Render(fmt.Sprintf("  Final Score: %s", r.NewStyle().Foreground(yellow).Bold(true).Render(fmt.Sprintf("%d", m.score)))),
				r.NewStyle().Foreground(subtle).Render("  enter to restart • esc to go back"),
			))
		sb.WriteString("\n  ")
		sb.WriteString(overlay)
	} else {
		sb.WriteString(r.NewStyle().Foreground(subtle).Italic(true).Render("  esc to go back to menu"))
	}

	return sb.String()
}

func (m Model) isSnakeBody(p point) bool {
	for _, s := range m.snake[1:] {
		if s == p {
			return true
		}
	}
	return false
}