package tui

import (
	"fmt"
	"fzwds/pkg/game"
	"fzwds/pkg/tui/styles"
	"fzwds/pkg/tui/theme"
	"fzwds/pkg/utils"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type turnDisplay struct {
	summaryRow			string
	summaryRowHl		string
	detailView			string
}

type gameReviewState struct {
	summaryRowFmtStr		string
	summaryRowWidth			int
	summaryRowPad			int
	// TODO: store *Turn instead of idx?
	selectedTurn			int
	visibleRowStart			int
	viewCache				map[int]*turnDisplay
}

func (m model) GameReviewSwitch() (model, tea.Cmd) {
	summary_row_width := fmt.Sprintf("%s %d. %s %s %s",
		"v", // validated symbol
		m.game.TurnCount(),
		strings.Repeat("_", m.game.Settings().PromptLenMax),
		"-s", // strike count
		"+l", // extra life
	)
	pad := 2
	summary_row_width = utils.RightPad(summary_row_width, pad)
	summary_row_width = utils.LeftPad(summary_row_width, pad)

	m.state.gameReview.summaryRowWidth = len(summary_row_width)
	m.state.gameReview.summaryRowPad = pad

	m.footerKeymaps = []footerKeymap {
		{key: "↑/↓", value: "scroll"},
		{key: "n/p", value: "next/prev strike"},
        {key: "esc", value: "back"},
	}

	m = m.SwitchPage(gameReviewPage)

	return m, nil
}

func (m model) GameReviewUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "g", "home":
			m.updateSummaryListState(0)

		case "G", "end":
			m.updateSummaryListState(m.game.TurnCount() - 1)

		case "j", "down", "tab":
			if m.state.gameReview.selectedTurn < m.game.TurnCount() - 1 {
				m.updateSummaryListState(m.state.gameReview.selectedTurn + 1)
			}

		case "k", "up", "shift+tab":
			if m.state.gameReview.selectedTurn > 0 {
				m.updateSummaryListState(m.state.gameReview.selectedTurn - 1)
			}

		case "n":
			sel := m.game.NextFailedTurnIdx(m.state.gameReview.selectedTurn)
			m.updateSummaryListState(sel)

		case "p":
			sel := m.game.PrevFailedTurnIdx(m.state.gameReview.selectedTurn)
			m.updateSummaryListState(sel)

		case "ctrl+u", "pgup":
			visible_rows := m.contentHeight - 2
			scroll := int(math.Floor(float64(visible_rows) / 2))
			clamped := utils.Clamp(
				m.state.gameReview.selectedTurn - scroll,
				0,
				m.state.gameReview.selectedTurn - scroll)
			m.updateSummaryListState(clamped)

		case "ctrl+d", "pgdown":
			visible_rows := m.contentHeight - 2
			scroll := int(math.Floor(float64(visible_rows) / 2))
			clamped := utils.Clamp(
				m.state.gameReview.selectedTurn + scroll,
				m.state.gameReview.selectedTurn + scroll,
				m.game.TurnCount() - 1)
			m.updateSummaryListState(clamped)

		case "esc":
			return m.GameOverSwitch()
		}
	}

	return m, nil
}

func (m model) GameReviewView() string {
	height := m.contentHeight - 2 // -2 for top/bottom table border rows
	current_turn := m.game.GetTurn(m.state.gameReview.selectedTurn)
	return lipgloss.NewStyle().Height(m.contentHeight).Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.renderTurnSummaryList(height),
			m.renderTurnDetailView(current_turn, height)))
}

func (m model) renderTurnSummaryList(height int) string {
	border := lipgloss.Border(lipgloss.RoundedBorder())
	border_style := lipgloss.NewStyle().Foreground(theme.Border).Render
	width := m.state.gameReview.summaryRowWidth

	list_title := "Turns"
	list_header := border_style(border.TopLeft + border.Top)
	list_header += styles.TextBody.Render(list_title)
	list_header += border_style(strings.Repeat(border.Top, width - len(list_title) - 1))
	list_header += border_style(border.TopRight)

	last_turn_idx := m.game.TurnCount() - 1
	// 1 space always reserved for final turn
	visible_rows := min(last_turn_idx, height - 1)

	start := m.state.gameReview.visibleRowStart
	// Show divider if last turn not within visible range
	show_divider := start + visible_rows < last_turn_idx
	if show_divider {
		visible_rows--
	}
	end := start + visible_rows

	list_items := make([]string, 0, visible_rows)
	for i := start; i < end; i++ {
		list_items = append(list_items, m.renderReviewSummaryRow(m.game.GetTurn(i)))
	}
	if show_divider {
		list_items = append(list_items, styles.TextDim.Render(strings.Repeat("─", width)))
	}
	// Pin last row to bottom
	list_items = append(list_items,
		m.renderReviewSummaryRow(m.game.GetTurn(last_turn_idx)),
	)

	// TODO: cache bigger styles like this so they only need to be created once
	list := lipgloss.NewStyle().
		Height(height).
		Width(width).
		Border(border).
		BorderTop(false).
		BorderForeground(theme.Border).
		AlignVertical(lipgloss.Top).
		Render(lipgloss.JoinVertical(lipgloss.Top, list_items...))

	return lipgloss.JoinVertical(lipgloss.Top, list_header, list)
}

func (m model) renderTurnDetailView(turn *game.Turn, height int) string {
	td, ok := m.state.gameReview.viewCache[turn.TurnNumber()]
	if ok && td.detailView != "" {
		return td.detailView
	}

	rows := [][]string{}

	if turn.Solved() {
		rows = append(rows, []string{
			"Answer",
			m.highlightPromptAnswer(turn.Prompt(), turn.Answer(), m.game.Settings().PromptMode),
		})
	} else {
		rows = append(rows, []string{
			"Possible answer",
			m.highlightPromptAnswer(turn.Prompt(), turn.SourceWord(), m.game.Settings().PromptMode),
		})
	}

	if turn.TotalTurnDuration() < time.Duration(time.Minute) {
		rows = append(rows, []string{
			"Total duration",
			fmt.Sprintf("%.1fs", turn.TotalTurnDuration().Seconds()),
		})
	} else {
		rows = append(rows, []string{
			"Total duration",
			utils.FormatTime(turn.TotalTurnDuration()),
		})
    }

	rows = append(rows, []string{ "Guesses", fmt.Sprintf("%d", turn.Guesses()) })

	if turn.Strikes() == 0 {
		rows = append(rows, []string{ "Strikes", "-" })
	} else {
		red := styles.TextRed.Render
		rows = append(rows, []string{
			"Strikes",
			red(fmt.Sprintf("%d/%d", turn.Strikes(), m.game.Settings().PromptStrikes)),
		})
	}

	if turn.Solved() {
		rows = append(rows,
			[]string{"Solve length", fmt.Sprintf("%d", len(turn.Answer())) },
			[]string{"Unique letters", fmt.Sprintf("%d", turn.UniqueLetterCount()) },
		)
	} else {
		rows = append(rows,
			[]string{"Solve length", "-" },
			[]string{"Unique letters", "-" },
		)
	}

	rows = append(rows, []string{ "Streak", fmt.Sprintf("%d", turn.Streak()) })

	stats_table := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderColumn(false).
		BorderLeft(false).BorderTop(false).BorderBottom(false).BorderRight(false).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			if row % 2 == 0 {
				style = styles.TextAccent
			} else {
				style = styles.TextBody
			}

			if col == 0 {
				style = style.
					Align(lipgloss.Left).
					Width(len("Possible answer"))
			}
			if col == 1 {
				style = style.
					Align(lipgloss.Left).
					MaxWidth(35).
					PaddingLeft(3)
			}
			return style
		}).
		Render()

	var solved_style lipgloss.Style
	if turn.Solved() {
		solved_style = lipgloss.NewStyle().Foreground(theme.Background).Background(theme.Green).Bold(true)
	} else {
		solved_style = lipgloss.NewStyle().Foreground(theme.Background).Background(theme.Red).Bold(true)
	}

	title_line := solved_style.Render(fmt.Sprintf(" #%d ", turn.TurnNumber()))
	title_line += " "
	title_line += styles.TextAccent.Bold(true).Render(strings.ToUpper(turn.Prompt()))

	detail_table := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		"",
		title_line,
		"",
		m.getTurnBadges(turn),
		"",
		stats_table,
	)

	// TODO: prefer single border style being passed around; maybe state prop
	border := lipgloss.Border(lipgloss.RoundedBorder())

	view := lipgloss.NewStyle().
		Height(height).
		Width(m.contentWidth - m.state.gameReview.summaryRowWidth).
		PaddingLeft(3).
		Border(border).
		BorderForeground(theme.Border).
		Render(detail_table)

	if !ok {
		m.state.gameReview.viewCache[turn.TurnNumber()] = &turnDisplay{}
	}
	m.state.gameReview.viewCache[turn.TurnNumber()].detailView = view

	return view
}

func (m model) renderReviewSummaryRow(turn *game.Turn) string {
	is_turn_selected := m.state.gameReview.selectedTurn == turn.TurnNumber() - 1
	td, ok := m.state.gameReview.viewCache[turn.TurnNumber()]
	if ok {
		if !is_turn_selected && td.summaryRow != "" {
			return td.summaryRow
		}
		if is_turn_selected && td.summaryRowHl != "" {
			return td.summaryRowHl
		}
	}

	strikes_width := " -9"
	extra_lives_width := " +1"

	var solved_indicator_text string
	var solved_indicator_style lipgloss.Style
	if turn.Solved() {
		if turn.FinalTurn() {
			solved_indicator_style = styles.TextYellow.Bold(true)
			solved_indicator_text = "W "
		} else {
			solved_indicator_style = styles.TextGreen.Bold(true)
			solved_indicator_text = "✓ "
		}
	} else {
		solved_indicator_style = styles.TextRed.Bold(true)
		if turn.FinalTurn() {
			if m.state.game.gameQuit {
				solved_indicator_text = "Q "
			} else {
				solved_indicator_text = "L "
			}
		} else {
			solved_indicator_text = "✘ "
		}
	}

	final_turn_num_str := fmt.Sprintf("%d.", m.game.TurnCount())
	turn_num_str := fmt.Sprintf("%d.", turn.TurnNumber())
	turn_num_padding := strings.Repeat(" ", len(final_turn_num_str) - len(turn_num_str))

	prompt := strings.ToLower(turn.Prompt())
	prompt_padding := strings.Repeat(" ", m.game.Settings().PromptLenMax - len(prompt))

	var strikes string
	if turn.Strikes() > 0 {
		strikes = fmt.Sprintf(" -%d", turn.Strikes())
	} else {
		strikes = strings.Repeat(" ", len(strikes_width))
	}

	var extra_life string
	if turn.ExtraLifeGained() {
		extra_life = " +1"
	} else {
		extra_life = strings.Repeat(" ", len(extra_lives_width))
	}

	turn_prompt_style := styles.TextBody
	if is_turn_selected {
		turn_prompt_style = styles.TextAccent.Bold(true)
	} else {
		turn_prompt_style = styles.TextBody
	}

	edge_pad_str := strings.Repeat(" ", m.state.gameReview.summaryRowPad)
	var out strings.Builder

	if is_turn_selected {
		sel_bg := theme.InputBg
		// sel_bg := lipgloss.AdaptiveColor{Dark: "#560cf5", Light: "#560cf5"}
		selection_style := lipgloss.NewStyle().Background(sel_bg)
		highlight := selection_style.Bold(true).Foreground(theme.Highlight)

		out.WriteString(selection_style.Render(edge_pad_str))
		out.WriteString(solved_indicator_style.Background(sel_bg).Render(solved_indicator_text))
		out.WriteString(selection_style.Render(turn_num_padding))
		out.WriteString(turn_prompt_style.Background(sel_bg).Render(turn_num_str, prompt))
		out.WriteString(selection_style.Render(prompt_padding))
		out.WriteString(styles.TextRed.Background(sel_bg).Bold(true).Render(strikes))
		out.WriteString(highlight.Render(extra_life))
		out.WriteString(selection_style.Render(edge_pad_str))

		s := out.String()
		if _, ok := m.state.gameReview.viewCache[turn.TurnNumber()]; !ok {
			m.state.gameReview.viewCache[turn.TurnNumber()] = &turnDisplay{}
		}
		m.state.gameReview.viewCache[turn.TurnNumber()].summaryRowHl = s

		return lipgloss.NewStyle().
			Width(m.state.gameReview.summaryRowWidth).
			Render(s)
	}

	out.WriteString(edge_pad_str)
	out.WriteString(solved_indicator_style.Render(solved_indicator_text))
	out.WriteString(turn_num_padding)
	out.WriteString(turn_prompt_style.Render(turn_num_str, prompt))
	out.WriteString(prompt_padding)
	out.WriteString(styles.TextRed.Render(strikes))
	out.WriteString(styles.TextHighlight.Render(extra_life))
	out.WriteString(edge_pad_str)

	s := out.String()
	if _, ok := m.state.gameReview.viewCache[turn.TurnNumber()]; !ok {
		m.state.gameReview.viewCache[turn.TurnNumber()] = &turnDisplay{}
	}
	m.state.gameReview.viewCache[turn.TurnNumber()].summaryRow = s

	return s
}

func (m model) getTurnBadges(turn *game.Turn) string {
	badges := make([]string, 0)
	base_badge_style := lipgloss.NewStyle().Foreground(theme.Background).Bold(true)

	if turn.ExtraLifeGained() {
		badges = append(badges, base_badge_style.Background(theme.Highlight).Render(" extra life "))
	}

	if turn.Solved() && len(turn.Answer()) == len(m.state.game.stats.LongestSolve()) {
		badges = append(badges, base_badge_style.Background(theme.Yellow).Render(" longest answer "))
	}

	if turn.Solved() && turn.UniqueLetterCount() == m.state.game.stats.MostUniqueCount() {
		badges = append(badges, base_badge_style.Background(theme.Purple).Render(" most unique "))
	}

	if m.state.game.stats.LongestStreak() > 0 && turn.Streak() == m.state.game.stats.LongestStreak() {
		badges = append(badges, base_badge_style.Background(theme.Orange).Render(" longest streak "))
	}

	return strings.Join(badges, " ") // TODO: lipgloss.Wrap on v2 to ensure all badges are styled
}

func (m *model) updateSummaryListState(sel int) {
	if sel == m.state.gameReview.selectedTurn {
		return
	}

	m.state.gameReview.selectedTurn = sel

	scrolloff := 2
	// TODO: this is also calculated in view; need to consolidate/store as struct prop
	max_rows := min(m.game.TurnCount(), m.contentHeight - 2) // -2 rows for top/bottom borders
	scrolloff_clamped := utils.Clamp(scrolloff, 0, int(math.Floor(float64(max_rows / 2))))

	// Scroll up
	if sel < m.state.gameReview.visibleRowStart + scrolloff_clamped {
		m.state.gameReview.visibleRowStart = utils.Clamp(sel - scrolloff_clamped, 0, sel - scrolloff_clamped)
	}

	// Scroll down; add 2 to scrolloff to account for pinned last row (separator + last row)
	if sel >= m.state.gameReview.visibleRowStart + max_rows - (scrolloff_clamped + 2) {
		m.state.gameReview.visibleRowStart = sel - max_rows + 1 + (scrolloff_clamped + 2)
	}

	clamped := utils.Clamp(m.state.gameReview.visibleRowStart, 0, m.game.TurnCount() - max_rows)
	if m.state.gameReview.visibleRowStart > clamped {
		m.state.gameReview.visibleRowStart = clamped
	}
}
