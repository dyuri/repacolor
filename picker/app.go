package picker

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	"github.com/dyuri/repacolor/color"
	"github.com/dyuri/repacolor/display"
)

const SLIDER_LGAP = 4
const SLIDER_RGAP = 1

type model struct {
	components []string
	values     []float64
	cursor     int
	color      color.RepaColor
	width      int
	height     int
	step	   float64
}

func getSliderWidth(width int) int {
	return width - SLIDER_LGAP - SLIDER_RGAP - 1
}

func mousePick(x, y int, m model) model {
	if y >= len(m.components) || x < SLIDER_LGAP || x > SLIDER_LGAP + getSliderWidth(m.width) {
		return m
	}

	v := float64(x - SLIDER_LGAP) / float64(getSliderWidth(m.width))
	m.values[y] = v

	return m
}

func mouseWheel(y int, change float64, m model) model {
	if y >= len(m.components) {
		return m
	}

	m.values[y] += change

	return m
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
		m.step = 1.0 / float64(getSliderWidth(m.width))
	case tea.MouseMsg:
		e := tea.MouseEvent(msg)
		switch e.Button {
		case tea.MouseButtonLeft:
			m = mousePick(e.X, e.Y, m)
		case tea.MouseButtonWheelUp:
			m = mouseWheel(e.Y, m.step, m)
		case tea.MouseButtonWheelDown:
			m = mouseWheel(e.Y, -m.step, m)
		}
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
			m.values[m.cursor] -= m.step
		case "H", "shift+left":
			m.values[m.cursor] -= m.step / 10
		case "l", "right":
			m.values[m.cursor] += m.step
		case "L", "shift+right":
			m.values[m.cursor] += m.step / 10
		}
	}

	for i := 0; i < len(m.values); i++ {
		if m.values[i] < 0 {
			m.values[i] = 0
		} else if m.values[i] > 1 {
			m.values[i] = 1
		}
	}
	m.color = color.CreateColor("rgb", m.values[0], m.values[1], m.values[2], m.values[3])

	return m, nil
}

func drawSlider(m model, i int) string {
	w := getSliderWidth(m.width)
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

	if m.height >= 16 {
		ansirepr := display.RenderAnsiImage(display.GetColorAnsiImage(m.color, display.ColorAnsiImageOptions{}))
		textrepr := "\n" + display.TextColorDetails(m.color)

		s += display.MergeStringsVertically(ansirepr, textrepr, 24)
	} else if m.height >= 5 {
		s += "\n" + m.color.AnsiBg() + m.color.A11YPair().AnsiFg() + m.color.Hex() + color.ANSI_RESET + "\n"
	}

	return s
}

func RunPicker(c color.RepaColor, showAlpha bool) {
	p := tea.NewProgram(initialModel(c, showAlpha), tea.WithMouseAllMotion(), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v", err)
		os.Exit(1)
	}

	// TODO display it as requested (hex by default)
	fmt.Println(m.(model).color.Hex())
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	args := s.Command()
	c := color.WHITE
	if len(args) > 0 {
		pc, err := color.ParseColor(args[0], true)
		if err != nil {
			log.Error("Could not parse color", "error", err)
		} else {
			c = pc
		}
	}

	return initialModel(c, false), []tea.ProgramOption{tea.WithMouseAllMotion(), tea.WithAltScreen()}
}

func ServePicker(port string) {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort("0.0.0.0", port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting server", "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Server error", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not gracefully shutdown server", "error", err)
	}
}

