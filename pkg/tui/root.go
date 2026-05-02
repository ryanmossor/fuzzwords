package tui

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"fzwds/pkg/game"
	anim "fzwds/pkg/tui/animations"
	"fzwds/pkg/tui/commands"
	"fzwds/pkg/tui/figurethisout"
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

// TODO: refactor root model to have context prop that is passed to views
// Root: orchestrator, delegates updates/view renders to pages
// Page: branches of root, interface
//	- Update(...)
//	- View(...) string
// 	- state struct
// Components: leaves of pages, move more complicated rendering here
//	- eg review summary rows/detail tables, scrollable menu items (eg settings)?

type model struct {
	ready               bool
	debug 				bool
	debug_map			map[string]string

	switched			bool
	currentPage			pages.Page
	pages				map[pages.PageName]pages.Page

	viewport			viewport.Model
	gotoTop				bool
	gotoBottom			bool

	uiContext			*figurethisout.UIContext

	helpKeys			[]figurethisout.HelpKeymap

	game				game.Game
	gameInput			textinput.Model
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

	uiContext := figurethisout.UIContext {
		DebugMap: make(map[string]string),
		Size: figurethisout.Large,
		FPS: 30,
		AnimManager: mgr,
		// TODO: idk if this needs to be shared
		InputRestricted: false,

		SettingsPath: settingsPath,
		Settings: &settings,
		Schema: schema,
	}

	appPages := make(map[pages.PageName]pages.Page)
	titlePage := pages.NewTitlePage(&uiContext, &settings)
	appPages[pages.Title] = titlePage
	appPages[pages.About] = pages.NewAboutPage(&uiContext)
	appPages[pages.Stats] = pages.NewStatsPage(&uiContext)
	appPages[pages.Settings] = pages.NewGameSettingsPage(&uiContext)
	appPages[pages.Preferences] = pages.NewPreferencesPage(&uiContext)

	return model {
		debug: debug,
		debug_map: make(map[string]string),

		helpKeys: []figurethisout.HelpKeymap {
			{Key: "ctrl+p", Value: "preferences"},
			{Key: "q", Value: "quit"},
		},

		pages: appPages,
		currentPage: titlePage,
		uiContext: &uiContext,

		// state: state {
		// 	pokemonMenu: pokemonMenuState {
		// 		genList: []int{},
		// 		genState: initSelectedPokemonGens(&settings),
		// 		selected: 1,
		// 	},
		// },
	}
}

func (m model) Init() tea.Cmd {
	return m.currentPage.Switch()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	// TODO: animations/timer only enabled on title screen/game
	// should probably not render other screens at 30fps
	case commands.TickMsg:
		if m.uiContext.Settings.Prefs.AnimationsEnabled {
			m.uiContext.AnimManager.Update(msg.Time)
		}
		cmds = append(cmds, commands.TickCmd(m.uiContext.FPS))

	case pages.SwitchPageMsg:
		p := m.pages[msg.PageName]
		m.currentPage = p
		cmds = append(cmds, p.Switch())
		m.switched = true

	case tea.WindowSizeMsg:
		m.uiContext.ViewportWidth = msg.Width
		m.uiContext.ViewportHeight = msg.Height

		switch {
		case msg.Width < 20 || msg.Height < 10:
			m.uiContext.Size = figurethisout.Undersized
			m.uiContext.ContainerWidth = msg.Width
			m.uiContext.ContainerHeight = msg.Height

		case msg.Width < 50:
			m.uiContext.Size = figurethisout.Small
			m.uiContext.ContainerWidth = msg.Width
			m.uiContext.ContainerHeight = msg.Height

		case msg.Width < 80:
			m.uiContext.Size = figurethisout.Medium
			m.uiContext.ContainerWidth = 50
			m.uiContext.ContainerHeight = min(msg.Height, 30)

		default:
			m.uiContext.Size = figurethisout.Large
			m.uiContext.ContainerWidth = 80
			m.uiContext.ContainerHeight = min(msg.Height, 30)
		}

		m.uiContext.ContentWidth = m.uiContext.ContainerWidth - 4
		m = m.updateViewport()

	// Currently these need to stay in main model so input is enabled again on game over screen
	// Better way to do this?
	case commands.EnableInputMsg:
		// m.state.game.inputRestricted = false
		m.uiContext.InputRestricted = false

	// case commands.TogglePlayerDamagedMsg:
	// 	m.state.game.playerDamaged = false

	case tea.KeyMsg:
		m.debug_map["keyPress"] = msg.String()
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	p, cmd := m.currentPage.Update(msg)
	m.currentPage = p

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
	m.viewport.SetContent(m.currentPage.View())
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

	height := m.uiContext.ContainerHeight - lipgloss.Height(header) - lipgloss.Height(footer)
	content_style := lipgloss.NewStyle().
		// MaxWidth(m.uiContext.ContentWidth).
		Height(height).
		Padding(0, 1)
		// AlignVertical(lipgloss.Center) // center all content on screen

	// TODO: this is not root's concern
	// if m.page == pages.AboutPage {
	// 	content_style = content_style.
	// 		Width(m.containerWidth).
	// 		AlignVertical(lipgloss.Top).
	// 		PaddingTop(1)
	// }

	has_scroll := false
	// TODO: this is not root's concern
	// if m.page == pages.SettingsPage {
	// 	has_scroll = m.viewport.VisibleLineCount() < m.viewport.TotalLineCount()
	// }

	var view string
	if has_scroll {
		view = lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.viewport.View(),
			lipgloss.NewStyle().Foreground(theme.Body).Width(1).Render(), // space between content and scrollbar
			m.getScrollbar(),
		)
	} else {
		view = m.currentPage.View()
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
		m.uiContext.ViewportWidth,
		m.uiContext.ViewportHeight,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.NewStyle().
			MaxWidth(m.uiContext.ViewportWidth).
			MaxHeight(m.uiContext.ViewportHeight).
			Render(child),
		)

	if m.debug {
		renderTimeMicros := float64(time.Since(start).Microseconds())
		m.debug_map["viewSize"] = strconv.Itoa(len(v) - len(debug_view))
		m.debug_map["renderTime"] = fmt.Sprintf("renderTime: %.1fms", renderTimeMicros / 1000)
	}

	return v
}

func (m model) updateViewport() model {
    header_height := lipgloss.Height(m.HeaderView())
    footer_height := lipgloss.Height(m.FooterView())
    vertical_margin_height := header_height + footer_height

    width := m.uiContext.ContainerWidth - 4
	m.uiContext.ContentHeight = m.uiContext.ContainerHeight - vertical_margin_height
    m.uiContext.ContentWidth = m.uiContext.ContainerWidth - 4

    if !m.ready {
        m.viewport = viewport.New(width, m.uiContext.ContentHeight)
        m.viewport.YPosition = header_height
        m.viewport.HighPerformanceRendering = false
        m.ready = true

        // m.viewport.YPosition = headerHeight + 1
        m.viewport.YPosition = header_height
    } else {
        m.viewport.Width = width
		m.viewport.Height = m.uiContext.ContentHeight
		m.viewport.GotoTop()
    }

    return m
}

func (m model) getScrollbar() string {
	viewport_height := m.viewport.Height
	content_height := lipgloss.Height(m.currentPage.View())
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
	settings.WriteSettings(path)

	return settings, path
}
