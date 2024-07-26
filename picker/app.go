package picker

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dyuri/repacolor/color"
	"github.com/dyuri/repacolor/display"
)

type model struct {
	components []string
	values []float64
	cursor  int
	color   color.RepaColor
}

func initialModel(c color.RepaColor) model {
	return model{
		components: []string{
			"red",
			"green",
			"blue",
			"alpha",
		},
		values: []float64{c.R, c.G, c.B, c.A},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
			m.values[m.cursor] -= 0.02
			if m.values[m.cursor] < 0 {
				m.values[m.cursor] = 0
			}
		case "l", "right":
			m.values[m.cursor] += 0.02
			if m.values[m.cursor] > 1 {
				m.values[m.cursor] = 1
			}
		}
	}

	m.color = color.CreateColor("rgb", m.values[0], m.values[1], m.values[2], m.values[3])

	return m, nil
}

func (m model) View() string {
	s := ""
	for i, choice := range m.components {
		value := m.values[i]
		if i == m.cursor {
			s += fmt.Sprintf("[%s] %6s %f\n", "x", choice, value)
		} else {
			s += fmt.Sprintf("[ ] %6s %f\n", choice, value)
		}
	}

	ansirepr := display.RenderAnsiImage(display.GetColorAnsiImage(m.color, display.ColorAnsiImageOptions{}))
	textrepr := "\n" + display.TextColorDetails(m.color)

	s += display.MergeStringsVertically(ansirepr, textrepr)

	return s
}

func RunPicker(c color.RepaColor) {
	p := tea.NewProgram(initialModel(c))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v", err)
		os.Exit(1)
	}
}

