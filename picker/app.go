package picker

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dyuri/repacolor/color"
	"github.com/dyuri/repacolor/display"
)

type model struct {
	components []string
	values     []float64
	cursor     int
	color      color.RepaColor
	width      int
	height     int
	me         tea.MouseEvent
}

func initialModel(c color.RepaColor, showAlpha bool) model {
	components := []string{
		"red",
		"green",
		"blue",
	}

	if showAlpha {
		components = append(components, "alpha")
	}

	return model{
		components: components,
		values: []float64{c.R, c.G, c.B, c.A},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.MouseMsg:
		m.me = tea.MouseEvent(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			m.cursor++
			if m.cursor >= len(m.components) {
				m.cursor = 0
			}
		case "k", "up":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.components) - 1
			}
		case "h", "left":
			m.values[m.cursor] -= 0.1
			if m.values[m.cursor] < 0 {
				m.values[m.cursor] = 0
			}
		case "H", "shift+left":
			m.values[m.cursor] -= 0.01
			if m.values[m.cursor] < 0 {
				m.values[m.cursor] = 0
			}
		case "l", "right":
			m.values[m.cursor] += 0.1
			if m.values[m.cursor] > 1 {
				m.values[m.cursor] = 1
			}
		case "L", "shift+right":
			m.values[m.cursor] += 0.01
			if m.values[m.cursor] > 1 {
				m.values[m.cursor] = 1
			}
		}
	}

	m.color = color.CreateColor("rgb", m.values[0], m.values[1], m.values[2], m.values[3])

	return m, nil
}

func drawSlider(m model, i int) string {
	w := m.width - 6
	value := m.values[i]
	slider := strings.Builder{}

	for j := 0; j <= w; j++ {
		c := m.color

		// TODO rgb[a] => general
		v := float64(j) / float64(w)
		switch i {
		case 0:
			c.R = v
		case 1:
			c.G = v
		case 2:
			c.B = v
		case 3:
			c.A = v
		}

		slider.WriteString(c.AnsiBg())
		slider.WriteString(c.A11YPair().AnsiFg())
		if j == int(value * float64(w)) {
			slider.WriteString("▣")
		} else {
			slider.WriteString(" ")
		}
	}
	slider.WriteString(color.ANSI_RESET)

	return slider.String()
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	s := ""
	for i, choice := range m.components {
		value := drawSlider(m, i)
		cursor := ' '
		if i == m.cursor {
			cursor = '▸'
		}
		s += fmt.Sprintf("%c %c %s\n", cursor, choice[0], value)
	}

	ansirepr := display.RenderAnsiImage(display.GetColorAnsiImage(m.color, display.ColorAnsiImageOptions{}))
	textrepr := "\n" + display.TextColorDetails(m.color)

	s += display.MergeStringsVertically(ansirepr, textrepr)

	s += fmt.Sprintf("mouse: %d, %d, %s", m.me.X, m.me.Y, m.me)

	return s
}

func RunPicker(c color.RepaColor, showAlpha bool) {
	p := tea.NewProgram(initialModel(c, showAlpha), tea.WithMouseAllMotion(), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v", err)
		os.Exit(1)
	}

	// TODO display as it should
	fmt.Println(m.(model).color.Hex())
}

