package tui

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"fzwds/pkg/game"
	anim "fzwds/pkg/tui/animations"
	"fzwds/pkg/tui/pages"
	"fzwds/pkg/tui/styles"
	"fzwds/pkg/tui/theme"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type size int
const (
	undersized size = iota
	small
	medium
	large
)

type footerKeymap struct {
	key		string
	value	string
}

type state struct {
	game			gameUIState
	gameReview		gameReviewState
	gameOver		gameOverState
	pressPlay		pressPlayState
	settings		settingsState
	pokemonMenu		pokemonMenuState
	footer			footerState
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

	// ?
	// anim?
	// debug map? to allow pages/components to write to it
	// footer msg?
}

type model struct {
	ready               bool
	debug 				bool
	debug_map			map[string]string

	switched			bool
	page				pages.PageName

	viewport			viewport.Model
	viewportWidth   	int
	viewportHeight  	int
	gotoTop				bool
	gotoBottom			bool

	containerWidth  	int
	containerHeight 	int
	contentWidth    	int
	contentHeight   	int

	size				size
	footerKeymaps		[]footerKeymap

	state				state
	game				game.Game
	gameInput			textinput.Model

	settingsPath		string
	settings			*game.Settings
	settingsCopy		game.Settings
	schema				game.SettingsSchema

	fps					int
	animManager			anim.AnimationManager
}

func NewModel(
	debug bool,
	settings game.Settings,
	schema game.SettingsSchema,
	settingsPath string,
) tea.Model {
	title_logo_anim := anim.NewTitleScreenLogoAnim(styles.GetRainbowColors())
	extra_life_anim := anim.NewRainbowScrollAnim(anim.ExtraLife, 30, false, styles.GetRainbowColors())
	validation_msg_dmg_anim := anim.NewDamageShakeAnim(anim.ValidationMessage, 10)
	strike_dmg_anim := anim.NewDamageShakeAnim(anim.StrikeCounter, 8)
	win_anim := anim.NewRainbowScrollAnim(anim.GameOverWin, 0, true, styles.GetRainbowColors())

	mgr := anim.NewAnimationManager(settings.Prefs.AnimationsEnabled)
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

		footerKeymaps: []footerKeymap {
			{key: "ctrl+p", value: "preferences"},
			{key: "q", value: "quit"},
		},

		settingsPath: settingsPath,
		schema: schema,
		settings: &settings,
		settingsCopy: settings,

		page: pages.TitlePage,

		state: state {
			pressPlay: pressPlayState {
				visible: true,
			},

			settings: settingsState {
				selected: 0,
				lastSel: make(map[settingsMenuCategory]int),
			},

			pokemonMenu: pokemonMenuState {
				genList: []int{},
				genState: initSelectedPokemonGens(&settings),
				selected: 1,
			},
		},

		fps: 30,
		animManager: mgr,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.TitleScreenInit(),
		m.tickCmd(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	// TODO: animations/timer only enabled on title screen/game
	// should probably not render other screens at 30fps
	case TickMsg:
		if m.settings.Prefs.AnimationsEnabled {
			m.animManager.Update(msg.Time)
		}
		cmds = append(cmds, m.tickCmd())

	case tea.WindowSizeMsg:
		m.viewportWidth = msg.Width
		m.viewportHeight = msg.Height

		switch {
		case m.viewportWidth < 20 || m.viewportHeight < 10:
			m.size = undersized
			m.containerWidth = m.viewportWidth
			m.containerHeight = m.viewportHeight

		case m.viewportWidth < 50:
			m.size = small
			m.containerWidth = m.viewportWidth
			m.containerHeight = m.viewportHeight

		case m.viewportWidth < 80:
			m.size = medium
			m.containerWidth = 50
			m.containerHeight = min(msg.Height, 30)

		default:
			m.size = large
			m.containerWidth = 80
			m.containerHeight = min(msg.Height, 30)
		}

		m.contentWidth = m.containerWidth - 4
		m = m.updateViewport()

	// Currently these need to stay in main model so input is enabled again on game over screen
	// Better way to do this?
	case EnableInputMsg:
		m.state.game.inputRestricted = false

	case TogglePlayerDamagedMsg:
		m.state.game.playerDamaged = false

	case tea.KeyMsg:
		m.debug_map["keyPress"] = msg.String()
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd

	switch m.page {
	case pages.TitlePage:
		m, cmd = m.TitleScreenUpdate(msg)
	case pages.AboutPage:
		m, cmd = m.AboutUpdate(msg)
	case pages.SettingsPage:
		m, cmd = m.SettingsUpdate(msg)
	case pages.PokemonGenMenuPage:
		m, cmd = m.PokemonGenSelectorUpdate(msg)
	case pages.GamePage:
		m, cmd = m.GameUpdate(msg)
	case pages.GameOverPage:
		m, cmd = m.GameOverUpdate(msg)
	case pages.GameReviewPage:
		m, cmd = m.GameReviewUpdate(msg)
	}

	var header_cmd tea.Cmd
	m, header_cmd = m.HeaderUpdate(msg)
	cmds = append(cmds, header_cmd)

	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if m.switched {
		m = m.updateViewport()
		m.switched = false
	}
	m.viewport.SetContent(m.getContent())
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	if m.gotoTop {
		m.viewport.GotoTop()
		m.gotoTop = false
	} else if m.gotoBottom {
		m.viewport.GotoBottom()
		m.gotoBottom = false
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	start := time.Now()

	header := m.HeaderView()
	footer := m.FooterView()

	height := m.containerHeight - lipgloss.Height(header) - lipgloss.Height(footer)
	content_style := lipgloss.NewStyle().
		Height(height).
		Padding(0, 1).
		AlignVertical(lipgloss.Center) // center all content on screen

	if m.page == pages.AboutPage {
		content_style = content_style.
			Width(m.containerWidth).
			AlignVertical(lipgloss.Top).
			PaddingTop(1)
	}

	has_scroll := false
	if m.page == pages.SettingsPage {
		has_scroll = m.viewport.VisibleLineCount() < m.viewport.TotalLineCount()
	}

	var view string
	if has_scroll {
		view = lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.viewport.View(),
			lipgloss.NewStyle().Foreground(theme.Body).Width(1).Render(), // space between content and scrollbar
			m.getScrollbar(),
		)
	} else {
		view = m.getContent()
	}

	debug_view := m.DebugView()
	child := lipgloss.JoinVertical(
		lipgloss.Center,
		debug_view,
		header,
		content_style.Render(view),
		footer,
	)

	v := lipgloss.Place(
		m.viewportWidth,
		m.viewportHeight,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.NewStyle().
			MaxWidth(m.viewportWidth).
			MaxHeight(m.viewportHeight).
			Render(child),
		)

	if m.debug {
		renderTimeMicros := float64(time.Since(start).Microseconds())
		m.debug_map["viewSize"] = strconv.Itoa(len(v) - len(debug_view))
		m.debug_map["renderTime"] = fmt.Sprintf("renderTime: %.1fms", renderTimeMicros / 1000)
	}

	return v
}

func (m model) SwitchPage(page pages.PageName) model {
	m.page = page
	m.switched = true
	return m
}

func (m model) getContent() string {
	page := "unknown"

	switch m.page {
	case pages.TitlePage:
		page = m.TitleScreenView()
	case pages.AboutPage:
		page = m.AboutView()
	case pages.SettingsPage:
		page = m.SettingsView()
	case pages.PokemonGenMenuPage:
		page = m.PokemonGenSelectorView()
	case pages.StatsPage:
		page = m.StatsView()
	case pages.GamePage:
		page = m.GameView()
	case pages.GameOverPage:
		page = m.GameOverView()
	case pages.GameReviewPage:
		page = m.GameReviewView()
	}

	return page
}

func (m model) updateViewport() model {
    header_height := lipgloss.Height(m.HeaderView())
    footer_height := lipgloss.Height(m.FooterView())
    vertical_margin_height := header_height + footer_height

    width := m.containerWidth - 4
	m.contentHeight = m.containerHeight - vertical_margin_height
    m.contentWidth = m.containerWidth - 4

    if !m.ready {
        m.viewport = viewport.New(width, m.contentHeight)
        m.viewport.YPosition = header_height
        m.viewport.HighPerformanceRendering = false
        m.ready = true

        // m.viewport.YPosition = headerHeight + 1
        m.viewport.YPosition = header_height
    } else {
        m.viewport.Width = width
		m.viewport.Height = m.contentHeight
		m.viewport.GotoTop()
    }

    return m
}

func (m model) getScrollbar() string {
	viewport_height := m.viewport.Height
	content_height := lipgloss.Height(m.getContent())
	if viewport_height >= content_height {
		return ""
	}

	scrollbar_height := (viewport_height * viewport_height) / content_height
	max_scroll := content_height - viewport_height
	scrollbar_pos := 1.0 - (float64(m.viewport.YOffset) / float64(max_scroll))
	if scrollbar_pos <= 0 {
		scrollbar_pos = 1
	} else if scrollbar_pos >= 1 {
		scrollbar_pos = 0
	}

	bar := lipgloss.NewStyle().
		Height(scrollbar_height).
		Width(1).
		Background(theme.Accent).
		Render()

	style := lipgloss.NewStyle().Width(1).Height(viewport_height)
	return style.Render(
		lipgloss.PlaceVertical(
			viewport_height,
			lipgloss.Position(scrollbar_pos),
			bar,
		),
	)
}

//go:embed game_settings_schema.json
var settingsSchemaJson []byte

func LoadSchema() game.SettingsSchema {
	var schema game.SettingsSchema
	if err := json.Unmarshal(settingsSchemaJson, &schema); err != nil {
		slog.Error("Error parsing game_settings_schema.json", "error", err)
		os.Exit(1)
	}
	return schema
}

func LoadSettings(schema game.SettingsSchema) (game.Settings, string) {
	settingsDir, err := os.UserConfigDir()
	if err != nil {
		slog.Error("Config dir not found, using tmp dir to save settings instead", "error", err)
		settingsDir = os.TempDir()
	}
	settingsDir = filepath.Join(settingsDir, "fuzzwords")
	os.MkdirAll(settingsDir, os.ModePerm)

	var settings game.Settings
	path := filepath.Join(settingsDir, "settings.json")
	contents, err := os.ReadFile(path)
    if err != nil {
		settings = game.GetDefaultSettings()
	}

    if err := json.Unmarshal(contents, &settings); err != nil {
		slog.Error("Error parsing settings.json - restoring default settings", "error", err)
		settings = game.GetDefaultSettings()
	}

	settings.ValidateSettings(schema)
	writeSettings(settings, path)

	return settings, path
}

func writeSettings(settings game.Settings, path string) {
	data, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		slog.Error("Error marshaling settings", "error", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		slog.Error("Error writing settings.json", "error", err)
	}
}
