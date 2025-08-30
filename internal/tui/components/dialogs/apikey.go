package dialogs

import (
	"github.com/JyotirmoyDas05/openpilot/internal/config"
	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// ApiKeyDialog is a TUI dialog for updating the OpenRouter API key.
type ApiKeyDialog struct {
	input  textinput.Model
	status string
}

func NewApiKeyDialog(currentKey string) *ApiKeyDialog {
	d := &ApiKeyDialog{}
	d.input = textinput.New()
	d.input.Placeholder = "Enter new OpenRouter API key"
	d.input.CharLimit = 128
	d.input.SetWidth(40)
	d.input.SetValue("")
	d.status = ""
	return d
}

func (d *ApiKeyDialog) Init() tea.Cmd {
	return textinput.Blink
}

func (d *ApiKeyDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			key := d.input.Value()
			if len(key) < 10 {
				d.status = "API key too short."
				return d, nil
			}
			// Save to config
			cfg, err := config.Load("", false)
			if err != nil {
				d.status = "Failed to load config."
				return d, nil
			}
			err = cfg.SetProviderAPIKey("openrouter", key)
			if err != nil {
				d.status = "Failed to save API key."
				return d, nil
			}
			d.status = "API key updated!"
			return d, nil
		case "esc":
			d.status = "Cancelled."
			return d, nil
		}
	}
	var cmd tea.Cmd
	d.input, cmd = d.input.Update(msg)
	return d, cmd
}

func (d *ApiKeyDialog) View() string {
	return lipgloss.NewStyle().
		Width(50).
		Render("Update OpenRouter API Key:\n" + d.input.View() + "\n" + d.status)
}
