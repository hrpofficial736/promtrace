package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/hrpofficial736/promtrace/internal/diff"
	"github.com/hrpofficial736/promtrace/internal/util"
)

func RenderDiff(d *diff.DiffResponseModel) string {
	title := HeadingStyle.Render("DIFF COMPARISON")

	// ── System Prompt section ────────────────────────────────────────────────
	var sysPromptContent string
	if d.SystemPrompt.Identical {
		sysPromptContent = MutedStyle.Render("identical")
	} else {
		sysPromptContent = lipgloss.JoinVertical(
			lipgloss.Left,
			RemovedStyle.Render("− "+d.SystemPrompt.Old),
			AddedStyle.Render("+ "+d.SystemPrompt.New),
		)
	}
	systemPromptSection := RenderSection("System Prompt", sysPromptContent)

	// ── User Prompt section ──────────────────────────────────────────────────
	var userPromptContent string
	if d.UserPrompt.Identical {
		userPromptContent = MutedStyle.Render("identical")
	} else {
		userPromptContent = lipgloss.JoinVertical(
			lipgloss.Left,
			RemovedStyle.Render("− "+d.UserPrompt.Old),
			AddedStyle.Render("+ "+d.UserPrompt.New),
		)
	}
	userPromptSection := RenderSection("User Prompt", userPromptContent)

	// ── Metrics section ──────────────────────────────────────────────────────
	// Response length
	respDelta := d.ResponseLength.Delta
	respDeltaStr := fmt.Sprintf("%+d", respDelta)
	var respDeltaStyled string
	if respDelta > 0 {
		respDeltaStyled = SuccessStyle.Render(respDeltaStr)
	} else if respDelta < 0 {
		respDeltaStyled = ErrorStyle.Render(respDeltaStr)
	} else {
		respDeltaStyled = MutedStyle.Render(respDeltaStr)
	}
	responseLengthLine := BodyStyle.Render(fmt.Sprintf(
		"Response Length  %d chars → %d chars  ",
		d.ResponseLength.A, d.ResponseLength.B,
	)) + respDeltaStyled

	// Cost
	costDelta := d.Cost.Delta
	costDeltaStr := fmt.Sprintf("%+.6f", costDelta)
	var costDeltaStyled string
	if costDelta > 0 {
		costDeltaStyled = SuccessStyle.Render(costDeltaStr)
	} else if costDelta < 0 {
		costDeltaStyled = ErrorStyle.Render(costDeltaStr)
	} else {
		costDeltaStyled = MutedStyle.Render(costDeltaStr)
	}
	costLine := BodyStyle.Render(fmt.Sprintf(
		"Cost             %s → %s  ",
		util.FmtCost(d.Cost.A), util.FmtCost(d.Cost.B),
	)) + costDeltaStyled

	// Latency
	latencyPct := d.Latency.PctChange
	latencyPctStr := fmt.Sprintf("%+.2f%%", latencyPct)
	var latencyPctStyled string
	if latencyPct > 0 {
		latencyPctStyled = SuccessStyle.Render(latencyPctStr)
	} else if latencyPct < 0 {
		latencyPctStyled = ErrorStyle.Render(latencyPctStr)
	} else {
		latencyPctStyled = MutedStyle.Render(latencyPctStr)
	}
	latencyLine := BodyStyle.Render(fmt.Sprintf(
		"Latency          %dms → %dms  ",
		d.Latency.A, d.Latency.B,
	)) + latencyPctStyled

	metricsContent := lipgloss.JoinVertical(
		lipgloss.Left,
		responseLengthLine,
		costLine,
		latencyLine,
	)
	metricsSection := RenderSection("Metrics", metricsContent)

	// ── Compose ──────────────────────────────────────────────────────────────
	container := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		systemPromptSection,
		userPromptSection,
		metricsSection,
	)

	return lipgloss.NewStyle().Padding(1, 2).Render(container)
}
