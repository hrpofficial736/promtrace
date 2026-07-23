package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/util"
)

func RenderTraceInfoContainer(t *store.Trace) string {
	// ── Header row: badge + trace ID ────────────────────────────────────────
	badge := RenderBadge("TRACE", ColorAccent)
	traceID := MutedStyle.Render(" " + t.ID)
	header := lipgloss.JoinHorizontal(lipgloss.Center, badge, traceID)

	// ── Metric chips ─────────────────────────────────────────────────────────
	modelChip := PanelStyle.Render(RenderKeyValue("Model", t.Model))
	latencyChip := PanelStyle.Render(RenderKeyValue("Latency", fmt.Sprintf("%sms", strconv.FormatInt(t.LatencyMs, 10))))
	tokensChip := PanelStyle.Render(RenderKeyValue("Tokens", strconv.Itoa(t.Tokens)))
	costChip := PanelStyle.Render(RenderKeyValue("Cost", util.FmtCost(t.Cost)))

	chips := lipgloss.JoinHorizontal(lipgloss.Top, modelChip, latencyChip, tokensChip, costChip)

	// ── Content sections ─────────────────────────────────────────────────────
	var systemPromptContent string
	if t.SystemPrompt == "" {
		systemPromptContent = MutedStyle.Render("none")
	} else {
		systemPromptContent = BodyStyle.Render(t.SystemPrompt)
	}

	systemPromptSection := RenderSection("System Prompt", systemPromptContent)
	userPromptSection := RenderSection("User Prompt", BodyStyle.Render(t.UserPrompt))
	responseSection := RenderSection("Response", BodyStyle.Render(t.Response))

	// ── Compose everything ───────────────────────────────────────────────────
	container := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		chips,
		"",
		systemPromptSection,
		userPromptSection,
		responseSection,
	)

	return lipgloss.NewStyle().Padding(1, 2).Render(container)
}
