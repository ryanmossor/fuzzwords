package tui

import (
	"fzw/src/game"
	"math"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// type page = int
type page int
const (
	splash_page page = iota
	about_page
	settings_page
	game_page
	game_over_page // stats, keymaps for play again, return to menu, exit game, etc
)

// type size = int
type size int
const (
	undersized size = iota
	small
	medium
	large
)

type model struct {
	switched			bool
	page				page
	viewport			viewport.Model
	viewport_width   	int
	viewport_height  	int
	width_container  	int
	height_container 	int
	width_content    	int
	height_content   	int
	renderer        	*lipgloss.Renderer
	theme 				theme
	size				size

	settings			game.Settings
}

func NewModel(renderer *lipgloss.Renderer) tea.Model {
	return model{
		renderer: renderer,
		theme: BasicTheme(renderer, nil),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport_width = msg.Width
		m.viewport_height = msg.Height

		switch {
		case m.viewport_width < 20 || m.viewport_height < 10:
			m.size = undersized
			m.width_container = m.viewport_width
			m.height_container = m.viewport_height
		case m.viewport_width < 40:
			m.size = small
			m.width_container = m.viewport_width
			m.height_container = m.viewport_height
		case m.viewport_width < 70:
			m.size = medium
			m.width_container = 40
			m.height_container = int(math.Min(float64(msg.Height), 30))
		default:
			m.size = large
			m.width_container = 70
			m.height_container = int(math.Min(float64(msg.Height), 30))
		}

		m.width_content = m.width_container - 4
		// m = m.updateViewport() 
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd

	switch m.page {
	case splash_page:
		m, cmd = m.MainMenuUpdate(msg)
	case about_page:
		m, cmd = m.AboutUpdate(msg)
	case settings_page:
		m, cmd = m.SettingsUpdate(msg)
	case game_page:
		m, cmd = m.GameUpdate(msg)
	// case gameOverPage:
		// fmt.Println("updating game over page")
	}

	var header_cmd tea.Cmd
	m, header_cmd = m.HeaderUpdate(msg)
	cmds := []tea.Cmd{header_cmd}

	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	m.viewport.SetContent(m.getContent()) 
	m.viewport, cmd = m.viewport.Update(msg)
	if m.switched {
		// m = m.updateViewport()
		m.switched = false
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.page {
	// case menuPage:
		// return m.MenuView()
		// return ""
	default:
		header := m.HeaderView()
		debug := m.DebugView()
		view := m.getContent()

		height := m.height_container
		height -= lipgloss.Height(debug)
		height -= lipgloss.Height(header)

		var v string
		if m.page == splash_page {
			v = m.theme.Base().
				// Width(m.widthContainer). // commenting out centers "Press p to play"
				// Align(lipgloss.Center).
				AlignVertical(lipgloss.Center).
				AlignHorizontal(lipgloss.Center).
				Height(height).
				Padding(0, 1).
				Render(view)
		} else {
			v = m.theme.Base().
				Width(m.width_container). // commenting out centers "Press p to play"
				// Align(lipgloss.Center).
				// AlignVertical(lipgloss.Center).
				// AlignHorizontal(lipgloss.Center).
				Height(height).
				Padding(0, 1).
				Render(view)
		}

		child := lipgloss.JoinVertical(
			lipgloss.Center,
			debug,
			header,
			v,
			// m.theme.Base().
			// 	Width(m.widthContainer). // commenting out centers "Press p to play"
			// 	// Align(lipgloss.Center).
			// 	AlignVertical(lipgloss.Center).
			// 	AlignHorizontal(lipgloss.Center).
			// 	Height(height).
			// 	Padding(0, 1).
			// 	Render(view),
		) 

		return m.renderer.Place(
			m.viewport_width,
			m.viewport_height,
			lipgloss.Center,
			lipgloss.Center,
			m.theme.Base().
				MaxWidth(m.width_container).
				MaxHeight(m.height_container).
				Render(child),
		) 
	}
}

func (m model) SwitchPage(page page) model {
	m.page = page
	m.switched = true
	return m
}

func (m model) getContent() string {
	page := "unknown"

	switch m.page {
	case splash_page:
		page = m.MainMenuView()
	case about_page:
		page = m.AboutView()
	case settings_page:
		page = m.SettingsView()
	case game_page:
		page = m.GameView()
	}

	return page
}