package tui

import (
	_ "embed"
	"encoding/json"
	"fzwds/src/game"
	"fzwds/src/utils"
	"fzwds/src/constants"
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

type footer_keymaps struct {
	key		string
	value	string
}

type GameUIState struct {
	start_time			time.Time
	timer 				time.Duration

	game_active			bool
	game_over_msg		string

	damage_anim_padding	int
	player_damaged		bool

	input_restricted	bool
	validation_msg		string
}

type SplashScreenState struct {
	logo_anim_active	bool
	logo_anim_complete	bool
	logo_anim_idx		int
	logo_hidden			bool
}

type State struct {
	game				game.GameState
	game_ui				GameUIState
	press_play			PressPlayState
	settings			SettingsState
	title				SplashScreenState
}

type model struct {
	debug 				bool
	debug_map			map[string]string

    ready               bool
	switched			bool
    has_scroll          bool

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
	footer_keymaps		[]footer_keymaps

	text_input			textinput.Model

	state				State
	game_settings		*game.Settings
	game_settings_copy	game.Settings
	settings_schema		game.SettingsSchema
	settings_path		string

	anim_fps			int
	enable_animations	bool
}

//go:embed game_settings_schema.json
var game_settings_schema_json []byte

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

	var game_settings_schema_parsed game.SettingsSchema
	if err := json.Unmarshal(game_settings_schema_json, &game_settings_schema_parsed); err != nil {
		slog.Error("Error parsing game_settings_schema.json", "error", err)
		os.Exit(1)
	}

    if err := json.Unmarshal(contents, &game_settings); err != nil {
		slog.Error("Error parsing settings.json - restoring default settings", "error", err)
		game_settings = game.InitializeSettings()
	} else {
		game_settings.ValidateSettings(game_settings_schema_parsed)
	}

	marshaled_settings, err := json.MarshalIndent(game_settings, "", "    ")
	if err != nil {
		slog.Error("Error marshaling validated settings JSON", "error", err)
	}

	if err := os.WriteFile(settings_file_path, marshaled_settings, 0644); err != nil {
		slog.Error("Error writing settings.json", "error", err)
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
		debug_map: make(map[string]string),

		renderer: renderer,
		theme: BasicTheme(renderer),

		text_input: text,

		footer_keymaps: []footer_keymaps{
			{key: "q", value: "quit"},
		},

		game_settings: &game_settings,
		game_settings_copy: game_settings,
		settings_schema: game_settings_schema_parsed,
		settings_path: settings_file_path,

		state: State{
			press_play: PressPlayState{ visible: true },
			settings: SettingsState{ selected: 0 },
			game_ui: GameUIState{
				game_active: false,
				game_over_msg: "",

				damage_anim_padding: 0,
				player_damaged: false,

				input_restricted: false,
				validation_msg: "",
			},
			title: SplashScreenState {
				logo_anim_active: 	false,
				logo_anim_complete: false,
				logo_anim_idx: 		0,
				logo_hidden: 		false,
			},
		},

		anim_fps: 30,
		enable_animations: true,
	}
}

func (m model) Init() tea.Cmd {
	// TODO: batch async cmds - I/O, db loading, settings json etc
	return tea.Batch(
		m.MainMenuInit(),
		m.PressPlayInit(),
	)
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
		m.state.game_ui.input_restricted = false
	case TogglePlayerDamagedMsg:
		m.state.game_ui.player_damaged = false
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case LogoInitMsg:
		m.state.title.logo_anim_active = true
		return m, m.mainMenuLogoUpdateCmd()
	case LogoTickMsg:
		if m.state.title.logo_anim_idx < len(constants.GAME_TITLE) {
			m.state.title.logo_anim_idx++
			return m, m.mainMenuLogoUpdateCmd()
		}
	case LogoCompleteMsg:
		m.state.title.logo_anim_complete = true
		return m, tea.Tick(10 * time.Second, func(t time.Time) tea.Msg {
			return LogoRestartMsg{}
		})
	case LogoRestartMsg:
		m.state.title = SplashScreenState {
			logo_anim_active: 	false,
			logo_anim_complete: false,
			logo_anim_idx: 		0,
			logo_hidden:		true,
		}
		return m, tea.Sequence(
			tea.Tick(750 * time.Millisecond, func(t time.Time) tea.Msg {
				return LogoUnhideMsg{}
			}),
			m.initMainMenuLogoAnimCmd(),
		)
	case LogoUnhideMsg:
		m.state.title.logo_hidden = false
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

// Apply spacing to string to produce a shaking animation.
// The first return value is the padded string.
// The second return value is the number of padding spaces applied.
func (m model) applyDamageShakeAnimation(str string) (string, int) {
	if m.state.game_ui.damage_anim_padding <= 0 {
		return str, 0
	}

	result := str

	pad := m.state.game_ui.damage_anim_padding / 2
	var padding_spaces int
	if pad % 2 == 0 {
		padding_spaces = pad
		result = utils.RightPad(result, padding_spaces)
	} else {
		padding_spaces = m.state.game_ui.damage_anim_padding
		result = utils.LeftPad(result, padding_spaces)
	}

	return result, padding_spaces
}
