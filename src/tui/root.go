package tui

import (
	_ "embed"
	"encoding/json"
	"fzwds/src/game"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page int
const (
	splash_page page = iota
	about_page
	settings_page
    stats_page
	game_page
	game_over_page
)

type size int
const (
	undersized size = iota
	small
	medium
	large
)

type footerCmd struct {
	key		string
	value	string
}

type gameState struct {
	restrict_input		bool
	validation_msg		string
}

type state struct {
	press_play			pressPlayState
	// TODO: change to game_input state and move GameState here?
	game				gameState
	settings			settingsState
}

type gameTimer struct {
	remaining_time      time.Duration
	done                bool
}

type model struct {
	debug 				bool
    ready               bool
	game_active			bool
	game_over			bool
	switched			bool
    has_scroll          bool

	game_over_msg		string

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
	footer_cmds			[]footerCmd

	text_input			textinput.Model
	default_prompt_style	lipgloss.Style

	state				state
	game_state			game.GameState
	game_settings		*game.Settings
	game_settings_copy	game.Settings
	settings_menu_json	[]game.Config // TODO rename
	settings_path		string

	game_start_time		time.Time
    game_timer          gameTimer
}

//go:embed settings_info.json
var settings_info_json []byte

func NewModel() tea.Model {
	cfg_dir, err := os.UserConfigDir()
	if err != nil {
		slog.Error("Config dir not found, using tmp dir to save settings instead", "error", err)
		cfg_dir = os.TempDir()
	}
	fzwds_cfg_path := filepath.Join(cfg_dir, "fuzzwords")
	os.MkdirAll(fzwds_cfg_path, os.ModePerm)

	var game_settings game.Settings
	settings_file_path := filepath.Join(fzwds_cfg_path, "settings.json")
	contents, err := os.ReadFile(settings_file_path)
    if err != nil {
		game_settings = game.InitializeSettings()
	}

    if err := json.Unmarshal(contents, &game_settings); err != nil {
		slog.Error("Error parsing settings.json - restoring default settings", "error", err)
		game_settings = game.InitializeSettings()
	} else {
		game_settings.ValidateSettings()
	}

	marshaled_settings, err := json.MarshalIndent(game_settings, "", "    ")
	if err != nil {
		slog.Error("Error marshaling validated settings JSON", "error", err)
	}

	if err := os.WriteFile(settings_file_path, marshaled_settings, 0644); err != nil {
		slog.Error("Error writing settings.json", "error", err)
	}

	var settings_info_parsed []game.Config
	if err := json.Unmarshal(settings_info_json, &settings_info_parsed); err != nil {
		slog.Error("Error parsing settings_info.json", "error", err)
		os.Exit(1)
	}

	text := textinput.New()
	text.Placeholder = "Answer"
	text.Focus()
	text.Prompt = " > "
	text.CharLimit = 40
	text.Width = 40

	renderer := lipgloss.DefaultRenderer()

	return model{
		// debug: true,
		game_active: false,

		renderer: renderer,
		theme: BasicTheme(renderer),

		text_input: text,
		default_prompt_style: text.PromptStyle,

		footer_cmds: []footerCmd{
			{key: "q", value: "quit"},
		},

		game_settings: &game_settings,
		game_settings_copy: game_settings,
		settings_menu_json: settings_info_parsed,
		settings_path: settings_file_path,

		state: state{
			press_play: pressPlayState{ visible: true },
			settings: settingsState{ selected: 0 },
		},
	}
}

func (m model) Init() tea.Cmd {
	// TODO: batch async cmds - I/O, db loading, settings json etc
	return m.PressPlayInit()
}

type EnableInputMsg time.Time
func (m *model) debounceInputCmd(duration_ms int) tea.Cmd {
    m.state.game.restrict_input = true

    return tea.Tick(time.Millisecond * time.Duration(duration_ms), func(t time.Time) tea.Msg {
		return EnableInputMsg(t)
	})
}

type TurnTimerTickMsg struct{}
func setTurnTickerCmd() tea.Cmd {
	return tea.Tick(time.Millisecond * 100, func(t time.Time) tea.Msg {
		return TurnTimerTickMsg{}
	})
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
		m = m.updateViewport()
	case PressPlayTickMsg:
		m, cmd := m.PressPlayUpdate(msg)
        // m.viewport.SetContent(m.getContent())
		return m, cmd
	case EnableInputMsg:
		m.state.game.restrict_input = false
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
	case game_over_page:
		m, cmd = m.GameOverUpdate(msg)
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
		m = m.updateViewport()
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
		var header string
		if m.page == game_page {
			header = m.GameHudView()
		} else {
			header = m.HeaderView()
		}

		game_input := m.GameInputView()
		footer := m.FooterView()

		height := m.height_container
		height -= lipgloss.Height(header)
		height -= lipgloss.Height(game_input)
		height -= lipgloss.Height(footer)

		content_style := m.theme.Base().Height(height).Padding(0, 1)
		if m.page == splash_page || m.page == game_page || m.page == game_over_page {
			// center all content on screen
			content_style = content_style.AlignVertical(lipgloss.Center)
		} else {
			content_style = content_style.Width(m.width_container)
		}
        content := m.viewport.View()

        var view string
		if m.has_scroll {
			view = lipgloss.JoinHorizontal(
				lipgloss.Top,
				content,
				m.theme.Base().Width(1).Render(), // space between content and scrollbar
				m.getScrollbar(),
			)
		} else {
			view = m.getContent()
		}

		child := lipgloss.JoinVertical(
			lipgloss.Center,
			header,
            content_style.Render(view),
			game_input,
			footer,
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
	case stats_page:
		page = m.StatsView()
	case game_page:
		// TODO: possible to return game hud, prompt, and input as single page?
		page = m.GamePromptView()
	case game_over_page:
		page = m.GameOverView()
	}

	return page
}

func (m model) updateViewport() model {
    header_height := lipgloss.Height(m.HeaderView())
    footer_height := lipgloss.Height(m.FooterView())
    vertical_margin_height := header_height + footer_height

    width := m.width_container - 4
	m.height_content = m.height_container - vertical_margin_height
    m.width_content = m.width_container - 4

    if !m.ready {
        m.viewport = viewport.New(width, m.height_content)
        m.viewport.YPosition = header_height
        m.viewport.HighPerformanceRendering = false
        m.ready = true

        // m.viewport.YPosition = headerHeight + 1
        m.viewport.YPosition = header_height
    } else {
        m.viewport.Width = width
		m.viewport.Height = m.height_content
		m.viewport.GotoTop()
    }

	m.has_scroll = m.viewport.VisibleLineCount() < m.viewport.TotalLineCount()

	// if m.has_scroll {
	// 	m.width_content = m.width_container - 4
	// } else {
	// 	// m.width_content = m.width_container - 2
	// 	m.width_content = m.width_container - 4
	// }

    return m
}

func (m model) getScrollbar() string {
	y := m.viewport.YOffset
	viewport_height := m.viewport.Height
	content_height := lipgloss.Height(m.getContent())
	if viewport_height >= content_height {
		return ""
	}

	scrollbar_height := (viewport_height * viewport_height) / content_height
	max_scroll := content_height - viewport_height
	scrollbar_pos := 1.0 - (float64(y) / float64(max_scroll))
	if scrollbar_pos <= 0 {
		scrollbar_pos = 1
	} else if scrollbar_pos >= 1 {
		scrollbar_pos = 0
	}

    bar := m.theme.Base().
		Height(scrollbar_height).
		Width(1).
		Background(m.theme.Accent()).
		Render()

	style := m.theme.Base().Width(1).Height(viewport_height)
    return style.Render(
        lipgloss.PlaceVertical(
            viewport_height,
            lipgloss.Position(scrollbar_pos),
            bar,
        ),
    )
}

