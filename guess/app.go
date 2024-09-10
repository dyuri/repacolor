package guess

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	"github.com/dyuri/repacolor/color"
)

type model struct {
	color      color.RepaColor
	numChoices int
	choices	   []color.RepaColor
	width      int
	height     int
	points	   int
	rounds	   int
}

func getRandomColor() color.RepaColor {
	r := rand.Float64()
	g := rand.Float64()
	b := rand.Float64()

	c := color.CreateColor(color.CS_RGB, r, g, b, 1)
	l, a, b := c.Lab()
	
	if l < .5 {
		c = color.CreateColor(color.CS_LAB, 1 - l, a, b, 1)
	}

	return c
}

func getChoices(c color.RepaColor, numChoices int) []color.RepaColor {
	choices := make([]color.RepaColor, numChoices)
	choices[0] = c

	for i := 1; i < numChoices; i++ {
		choices[i] = getRandomColor()
	}

	// shuffle
	for i := range choices {
		j := rand.Intn(i + 1)
		choices[i], choices[j] = choices[j], choices[i]
	}
	
	return choices
}

func initialModel() model {
	numChoices := 2
	rounds := 10
	c := getRandomColor()
	choices := getChoices(c, numChoices)

	return model{
		color: c,
		choices: choices,
		numChoices: numChoices,
		rounds: rounds,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO mouse?
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			if m.choices[0] == m.color {
				m.points++
			}
			m.rounds--
			m.color = getRandomColor()
			m.choices = getChoices(m.color, m.numChoices)
		case "2":
			if m.choices[1] == m.color {
				m.points++
			}
			m.rounds--
			m.color = getRandomColor()
			m.choices = getChoices(m.color, m.numChoices)
		}
	}

	if m.rounds == 0 {
		return m, tea.Quit
	}

	return m, nil
}


func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	colorarea := ""

	for i := 0; i < 5; i++ {
		colorarea += " " + m.color.AnsiBg()
		for j := 0; j < m.width - 2; j++ {
			colorarea += " "
		}
		colorarea += color.ANSI_RESET + "\n"
	}

	s := fmt.Sprintf("%s\n1: %s - 2: %s\nPoints: %d [%d left]\n", colorarea, m.choices[0], m.choices[1], m.points, m.rounds)

	return s
}

func RunGuess() {
	p := tea.NewProgram(initialModel(), tea.WithMouseAllMotion(), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v", err)
		os.Exit(1)
	}

	// TODO display it as requested (hex by default)
	fmt.Printf("Finished with %d points.\n", m.(model).points)
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return initialModel(), []tea.ProgramOption{tea.WithMouseAllMotion(), tea.WithAltScreen()}
}

func ServeGuess(port string) {
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

