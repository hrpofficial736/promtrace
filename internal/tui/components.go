package tui

func RenderStatus(success bool, text string) string {
	var out string
	if success {
		out = "✅ " + text
	} else {
		out = "❌ " + text
	}
	return TitleStyle.Render(out)
}
