package pages

// import (
// 	"fmt"
// 	"fzwds/pkg/game"
// 	"fzwds/pkg/tui/styles"
// 	"fzwds/pkg/utils"
// 	"slices"
// 	"strings"
//
// 	"github.com/charmbracelet/lipgloss"
// )
//
// func (m model) GameReviewHudView() string {
// 	turn := m.getTurn(m.state.gameReview.selectedTurn)
// 	return lipgloss.JoinVertical(
// 		lipgloss.Center,
// 		m.renderTurnInfo(turn),
// 		m.renderReviewRemainingLetters(turn),
// 		"",
// 	)
// }
//
// func (m model) renderTurnInfo(turn game.Turn) string {
// 	health_display := m.renderHealthDisplay(turn.Health())
// 	health_change_info := m.renderHealthChangeInfo(turn)
// 	line := styles.TextBorder.Render(strings.Repeat("─", m.containerWidth))
//
// 	return lipgloss.JoinVertical(
// 		lipgloss.Left,
// 		line,
// 		utils.LeftPad(health_display + health_change_info, 8),
// 		line)
// }
//
// func (m model) renderReviewRemainingLetters(turn game.Turn) string {
// 	var out strings.Builder
//
// 	for i, c := range m.game.Settings().Alphabet.Letters() {
// 		if slices.Contains(turn.NewLettersUsed(), c) {
// 			out.WriteString(styles.TextHighlight.Bold(true).Underline(true).Render(string(c)))
// 		} else if turn.LettersUsed()[c] {
// 			out.WriteString(styles.TextDim.Render(string(c)))
// 		} else {
// 			out.WriteString(styles.TextYellow.Bold(true).Render(string(c)))
// 		}
//
// 		if i < len(m.game.Settings().Alphabet.Letters()) - 1 {
// 			out.WriteRune(' ')
// 		}
// 	}
//
// 	return out.String()
// }
//
// func (m model) renderHealthChangeInfo(turn game.Turn) string {
// 	var health_change_info string
//
// 	if turn.Strikes() > 0 {
// 		health_change_info += styles.TextRed.Render(fmt.Sprintf(" -%d", turn.Strikes()))
// 	}
//
// 	if turn.ExtraLifeGained() {
// 		health_change_info += styles.TextHighlight.Render(" +1")
// 	}
//
// 	return health_change_info
// }
