package styles

import (
	"fzwds/pkg/enums"
	"fzwds/pkg/tui/theme"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func CreateBox(content string, selected bool, width int) string {
	padded := lipgloss.PlaceHorizontal(width, lipgloss.Left, content)
	base := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Width(width)

	var style lipgloss.Style
	if selected {
		style = base.BorderForeground(theme.Accent).BorderStyle(lipgloss.RoundedBorder())
	} else {
		style = base.BorderForeground(theme.Border).BorderStyle(lipgloss.RoundedBorder())
	}

	return style.PaddingLeft(1).Render(padded)
}

func CreatePokemonMenuItem(content string, is_selected, apply_bottom_border bool, width int) string {
	padded := lipgloss.PlaceHorizontal(width, lipgloss.Left, content)
	base := lipgloss.NewStyle().
		BorderBottom(apply_bottom_border).
		BorderForeground(theme.Border).
		BorderStyle(lipgloss.NormalBorder()).
		Width(width)

	return base.PaddingLeft(1).Render(padded)
}

func CreateSettingsMenuItem(content string, is_selected, apply_bottom_border bool, width int) string {
	padded := lipgloss.PlaceHorizontal(width, lipgloss.Left, content)
	base := lipgloss.NewStyle().
		BorderBottom(apply_bottom_border).
		BorderForeground(theme.Border).
		BorderStyle(lipgloss.NormalBorder()).
		Width(width)

	return base.PaddingLeft(1).Render(padded)
}

func TextInputBlockBorderStyle(accent_color lipgloss.TerminalColor, width int) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderForeground(theme.InputBg).
		BorderStyle(lipgloss.InnerHalfBlockBorder()).
		BorderLeftForeground(accent_color).
		Width(width)
}

func TextInputRoundedBorderStyle(border_color lipgloss.TerminalColor, width int) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderForeground(border_color).
		BorderStyle(lipgloss.RoundedBorder()).
		Width(width)
}

var (
	TextBlack      = lipgloss.NewStyle().Foreground(theme.Black)
	TextWhite      = lipgloss.NewStyle().Foreground(theme.White)
	TextBackground = lipgloss.NewStyle().Foreground(theme.Background)
	TextBorder     = lipgloss.NewStyle().Foreground(theme.Border)
	TextBody       = lipgloss.NewStyle().Foreground(theme.Body)
	TextAccent     = lipgloss.NewStyle().Foreground(theme.Accent)
	TextDim        = lipgloss.NewStyle().Foreground(theme.Dim)
	TextInputBg    = lipgloss.NewStyle().Foreground(theme.InputBg)
	TextHighlight  = lipgloss.NewStyle().Foreground(theme.Highlight)
	TextRed        = lipgloss.NewStyle().Foreground(theme.Red)
	TextOrange     = lipgloss.NewStyle().Foreground(theme.Orange)
	TextYellow     = lipgloss.NewStyle().Foreground(theme.Yellow)
	TextGreen      = lipgloss.NewStyle().Foreground(theme.Green)
	TextBlue       = lipgloss.NewStyle().Foreground(theme.Blue)
	TextIndigo     = lipgloss.NewStyle().Foreground(theme.Indigo)
	TextPurple     = lipgloss.NewStyle().Foreground(theme.Purple)
)

func GetRainbowColors() []lipgloss.Style {
	return []lipgloss.Style{
		TextRed,
		TextOrange,
		TextYellow,
		TextGreen,
		TextBlue,
		TextIndigo,
		TextPurple,
	}
}
// Highlight prompt letters in current answer
func HighlightPromptAnswer(prompt, answer string, prompt_mode enums.PromptMode) string {
	accent := TextAccent.Render
	highlight := TextHighlight.Render

	prompt_upper := strings.ToUpper(prompt)
	answer_upper := strings.ToUpper(answer)
	var sb strings.Builder

	switch prompt_mode {
	case enums.PromptModeFuzzy:
		prompt_idx := 0
		for _, c := range answer_upper {
			curr_char := string(c)

			if prompt_idx < len(prompt_upper) && curr_char == string(prompt_upper[prompt_idx]) {
				sb.WriteString(highlight(curr_char))
				prompt_idx++
			} else {
				sb.WriteString(accent(curr_char))
			}
		}

	case enums.PromptModeClassic:
		if !strings.Contains(answer_upper, prompt_upper) {
			sb.WriteString(accent(answer_upper))
			return sb.String()
		}

		sub_idx := strings.Index(answer_upper, prompt_upper)
		sb.WriteString(accent(answer_upper[0:sub_idx]))
		sb.WriteString(highlight(answer_upper[sub_idx:sub_idx + len(prompt_upper)]))
		sb.WriteString(accent(answer_upper[sub_idx + len(prompt_upper):]))
	}

	return sb.String()
}


