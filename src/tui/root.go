package tui

import (
	_ "embed"
	"encoding/json"
	"fzwds/src/game"
	anim "fzwds/src/tui/animations"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page int
const (
	about_page page = iota
	game_over_page
	game_page
	pokemon_gen_selector
	settings_page
	splash_page
    stats_page
)

type size int
const (
	undersized size = iota
	small
	medium
	large
)

type FooterKeymap struct {
	key		string
	value	string
}

type GameUIState struct {
	prev_answer				string
	validation_msg			string
	game_over_msg			string
	player_damaged			bool
	input_restricted		bool
}

type FooterState struct {
	footer_msg		string
}

type State struct {
	game					game.GameState
	game_ui					GameUIState
	press_play				PressPlayState
	settings				SettingsState
	pokemon_gen_selector	PokemonGenSelectorState
	footer					FooterState
}

// TODO: refactor root model to have context prop that is passed to views
// Root: orchestrator, delegates updates/view renders to pages
// Page: branches of root, interface
//	- Update(...)
//	- View(...) string
// 	- state struct
// Components: leaves of pages, move more complicated rendering here
//	- eg review summary rows/detail tables, scrollable menu items (eg settings)?
type UIContext struct {
	Size				size

	ContainerWidth		int
	ContainerHeight		int

	ContentWidth		int
	ContentHeight		int

	viewportWidth		int
	viewportHeight		int

	Theme				theme

	// ?
	// anim?
}

type model struct {
	debug 				bool
	debug_map			map[string]string

	ctrl_c_pressed		bool

    ready               bool
	switched			bool
    has_scroll          bool

	goto_top			bool
	goto_bottom			bool

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
	footer_keymaps		[]FooterKeymap

	text_input			textinput.Model

	state				State
	app_settings		*game.Settings
	app_settings_copy	game.Settings
	app_settings_schema	game.SettingsSchema
	app_settings_path	string

	FPS					int
	anim_mgr			anim.AnimationManager
}

//go:embed game_settings_schema.json
var game_settings_schema_json []byte

func NewModel(renderer *lipgloss.Renderer, debug bool) tea.Model {
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
		game_settings = game.GetDefaultSettings()
	}

	var game_settings_schema_parsed game.SettingsSchema
	if err := json.Unmarshal(game_settings_schema_json, &game_settings_schema_parsed); err != nil {
		slog.Error("Error parsing game_settings_schema.json", "error", err)
		os.Exit(1)
	}

    if err := json.Unmarshal(contents, &game_settings); err != nil {
		slog.Error("Error parsing settings.json - restoring default settings", "error", err)
		game_settings = game.GetDefaultSettings()
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

	theme := BasicTheme(renderer)

	title_logo_anim := anim.NewTitleScreenLogoAnim(theme.GetRainbowColors())
	extra_life_anim := anim.NewRainbowScrollAnim(anim.ExtraLife, 30, false, theme.GetRainbowColors())
	validation_msg_dmg_anim := anim.NewDamageShakeAnim(anim.ValidationMessage, 10)
	strike_dmg_anim := anim.NewDamageShakeAnim(anim.StrikeCounter, 8)
	win_anim := anim.NewRainbowScrollAnim(anim.GameOverWin, 0, true, theme.GetRainbowColors())

	mgr := anim.NewAnimationManager(game_settings.Prefs.AnimationsEnabled)
	mgr.Register(
		title_logo_anim,
		extra_life_anim,
		validation_msg_dmg_anim,
		strike_dmg_anim,
		win_anim,
	)

	return model {
		debug: debug,
		debug_map: make(map[string]string),

		renderer: renderer,
		theme: theme,

		footer_keymaps: []FooterKeymap {
			{key: "ctrl+p", value: "preferences"},
			{key: "q", value: "quit"},
		},

		app_settings: &game_settings,
		app_settings_copy: game_settings,
		app_settings_schema: game_settings_schema_parsed,
		app_settings_path: settings_file_path,

		page: splash_page,

		state: State {
			press_play: PressPlayState {
				visible: true,
			},

			settings: SettingsState {
				selected: 0,
			},

			pokemon_gen_selector: PokemonGenSelectorState {
				gen_list: []int{},
				gen_state: initSelectedPokemonGens(&game_settings),
				selected: 1,
			},

			game_ui: GameUIState {
				prev_answer: 		"",
				validation_msg: 	"",
				game_over_msg: 		"",
				player_damaged: 	false,
				input_restricted: 	false,
			},
		},

		FPS: 30,
		anim_mgr: mgr,
	}
}

type TickMsg struct {
	Time	time.Time
}

// Global tick timer
func (m model) tickCmd() tea.Cmd {
	return tea.Tick(time.Second / time.Duration(m.FPS), func(t time.Time) tea.Msg {
		return TickMsg{t}
	})
}

func (m model) Init() tea.Cmd {
	// TODO: batch async cmds - I/O, db loading, settings json etc
	return tea.Batch(
		m.MainMenuInit(),
		m.PressPlayInit(),
		m.tickCmd(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		var cmds []tea.Cmd
		now := msg.Time

		if m.state.game.GameActive && m.state.game.TimeRemaining() <= 0 {
			cmds = append(cmds, m.turnTimerExpiredCmd())
		}

		if m.app_settings.Prefs.AnimationsEnabled {
			m.anim_mgr.Update(now)
		}

		cmds = append(cmds, m.tickCmd())
		return m, tea.Batch(cmds...)

	case tea.WindowSizeMsg:
		m.viewport_width = msg.Width
		m.viewport_height = msg.Height

		switch {
		case m.viewport_width < 20 || m.viewport_height < 10:
			m.size = undersized
			m.width_container = m.viewport_width
			m.height_container = m.viewport_height

		case m.viewport_width < 50:
			m.size = small
			m.width_container = m.viewport_width
			m.height_container = m.viewport_height

		case m.viewport_width < 80:
			m.size = medium
			m.width_container = 50
			m.height_container = min(msg.Height, 30)

		default:
			m.size = large
			m.width_container = 80
			m.height_container = min(msg.Height, 30)
		}

		m.width_content = m.width_container - 4
		m = m.updateViewport()

	case PressPlayTickMsg:
		m, cmd := m.PressPlayUpdate(msg)
		return m, cmd

	case EnableInputMsg:
		m.state.game_ui.input_restricted = false

	case TogglePlayerDamagedMsg:
		m.state.game_ui.player_damaged = false

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.ctrl_c_pressed {
				return m, tea.Quit
			} else {
				m.ctrl_c_pressed = true
				m.state.footer.footer_msg = m.theme.TextRed().Bold(true).Render("Press ctrl+c again to quit")
				return m, nil
			}
		default:
			if m.ctrl_c_pressed {
				m.ctrl_c_pressed = false
				m.state.footer.footer_msg = ""
			}
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
	case pokemon_gen_selector:
		m, cmd = m.PokemonGenSelectorUpdate(msg)
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

	if m.goto_top {
		m.viewport.GotoTop()
		m.goto_top = false
	} else if m.goto_bottom {
		m.viewport.GotoBottom()
		m.goto_bottom = false
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	header := m.HeaderView()
	footer := m.FooterView()

	height := m.height_container - lipgloss.Height(header) - lipgloss.Height(footer)
	content_style := m.theme.Base().
		Height(height).
		Padding(0, 1).
		AlignVertical(lipgloss.Center) // center all content on screen

	if m.page == about_page {
		content_style = content_style.
			Width(m.width_container).
			AlignVertical(lipgloss.Top).
			PaddingTop(1)
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
		m.DebugView(),
		header,
		content_style.Render(view),
		footer,
	) 

	v := m.renderer.Place(
		m.viewport_width,
		m.viewport_height,
		lipgloss.Center,
		lipgloss.Center,
		m.theme.Base().
			MaxWidth(m.viewport_width).
			MaxHeight(m.viewport_height).
			Render(child),
		)

	if m.debug {
		m.debug_map["viewSize"] = strconv.Itoa(len(v))
		m.debug_map["runeCount"] = strconv.Itoa(utf8.RuneCountInString(v))
	}

	return v
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
	case pokemon_gen_selector:
		page = m.PokemonGenSelectorView()
	case stats_page:
		page = m.StatsView()
	case game_page:
		page = m.GameView()
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

	if m.page == game_page {
		m.has_scroll = false
	} else {
		m.has_scroll = m.viewport.VisibleLineCount() < m.viewport.TotalLineCount()
	}

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
