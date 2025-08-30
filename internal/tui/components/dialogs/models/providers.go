package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/spinner"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/catwalk/pkg/catwalk"
	"github.com/charmbracelet/lipgloss/v2"

	"github.com/JyotirmoyDas05/openpilot/internal/config"
	"github.com/JyotirmoyDas05/openpilot/internal/tui/components/core"
	"github.com/JyotirmoyDas05/openpilot/internal/tui/components/dialogs"
	"github.com/JyotirmoyDas05/openpilot/internal/tui/exp/list"
	"github.com/JyotirmoyDas05/openpilot/internal/tui/styles"
	"github.com/JyotirmoyDas05/openpilot/internal/tui/util"
)

const (
	ProvidersDialogID     dialogs.DialogID = "providers"
	defaultProvidersWidth int              = 50
)

type providerListModel = list.FilterableList[list.CompletionItem[ModelOption]]

type providerDialogCmp struct {
	width   int
	wWidth  int
	wHeight int

	list   providerListModel
	keyMap KeyMap
	help   help.Model

	// API key state
	needsAPIKey      bool
	apiKeyInput      *APIKeyInput
	spinner          spinner.Model
	selectedProvider *ModelOption
}

func NewProviderDialogCmp() dialogs.DialogModel {
	t := styles.CurrentTheme()
	inputStyle := t.S().Base.PaddingLeft(1).PaddingBottom(1)
	listKeyMap := list.DefaultKeyMap()

	providerList := list.NewFilterableList(
		[]list.CompletionItem[ModelOption]{},
		list.WithFilterInputStyle(inputStyle),
		list.WithFilterListOptions(
			list.WithKeyMap(listKeyMap),
			list.WithWrapNavigation(),
			list.WithResizeByList(),
		),
	)

	apiKeyInput := NewAPIKeyInput()
	apiKeyInput.SetShowTitle(false)
	spin := spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(t.S().Base.Foreground(t.Green)))

	helpM := help.New()
	helpM.Styles = t.S().Help

	cmp := &providerDialogCmp{
		width:       defaultProvidersWidth,
		list:        providerList,
		apiKeyInput: apiKeyInput,
		help:        helpM,
		spinner:     spin,
	}
	cmp.keyMap = DefaultKeyMap()
	return cmp
}

func (p *providerDialogCmp) Init() tea.Cmd {
	providers, err := config.Providers()
	if err == nil {
		items := []list.CompletionItem[ModelOption]{}
		cfg := config.Get()
		for _, pr := range providers {
			name := pr.Name
			if name == "" {
				name = string(pr.ID)
			}
			if _, ok := cfg.Providers.Get(string(pr.ID)); ok {
				// configured providers show a checkmark and allow changing the API key
				items = append(items, list.NewCompletionItem(name, ModelOption{Provider: pr, Model: catwalk.Model{ID: ""}}, list.WithCompletionID(string(pr.ID)), list.WithCompletionShortcut("✅ configured"), list.WithCompletionFocusedShortcut("Change API key...")))
			} else {
				items = append(items, list.NewCompletionItem(name, ModelOption{Provider: pr, Model: catwalk.Model{ID: ""}}, list.WithCompletionID(string(pr.ID)), list.WithCompletionFocusedShortcut("Set API key...")))
			}
		}
		// set all items and try to size the list to show them by default
		p.list.SetItems(items)
		// Try to size the list to fit items; clamp to a sensible min/max
		desired := len(items) + 2
		if desired < 6 {
			desired = 6
		}
		// if we already know window height, clamp to half the window
		if p.wHeight > 0 {
			maxH := p.wHeight / 2
			if desired > maxH {
				desired = maxH
			}
		}
		p.list.SetSize(p.width-2, desired)
	}
	return tea.Batch(p.list.Init(), p.apiKeyInput.Init())
}

func (p *providerDialogCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.wWidth = msg.Width
		p.wHeight = msg.Height
		p.apiKeyInput.SetWidth(p.width - 2)
		p.help.Width = p.width - 2
		return p, p.list.SetSize(p.width-2, p.wHeight/4)
	case APIKeyStateChangeMsg:
		u, cmd := p.apiKeyInput.Update(msg)
		p.apiKeyInput = u.(*APIKeyInput)
		return p, cmd
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, p.keyMap.Select):
			if p.needsAPIKey {
				apiKey := p.apiKeyInput.Value()
				if p.selectedProvider == nil {
					return p, util.ReportError(fmt.Errorf("no provider selected"))
				}
				providerID := string(p.selectedProvider.Provider.ID)
				err := config.Get().SetProviderAPIKey(providerID, apiKey)
				if err != nil {
					return p, util.ReportError(fmt.Errorf("failed to save API key: %w", err))
				}
				// close dialog and show info
				return p, tea.Sequence(
					util.CmdHandler(dialogs.CloseDialogMsg{}),
					util.ReportInfo(fmt.Sprintf("API key saved for %s", p.selectedProvider.Provider.Name)),
				)
			}
			// open API key input for selected provider
			sel := p.list.SelectedItem()
			if sel == nil {
				return p, nil
			}
			v := (*sel).Value()
			p.selectedProvider = &v
			p.apiKeyInput.SetProviderName(v.Provider.Name)
			p.needsAPIKey = true
			return p, nil
		case key.Matches(msg, p.keyMap.Close):
			if p.needsAPIKey {
				if p.apiKeyInput != nil && p.apiKeyInput.Value() != "" {
					// if valid, do nothing here; user should confirm
				}
				p.needsAPIKey = false
				p.selectedProvider = nil
				p.apiKeyInput.Reset()
				return p, nil
			}
			return p, util.CmdHandler(dialogs.CloseDialogMsg{})
		default:
			// allow explicit arrow navigation even when the filter input is present
			if key.Matches(msg, p.keyMap.Next) {
				return p, p.list.SelectItemBelow()
			}
			if key.Matches(msg, p.keyMap.Previous) {
				return p, p.list.SelectItemAbove()
			}

			if p.needsAPIKey {
				u, cmd := p.apiKeyInput.Update(msg)
				p.apiKeyInput = u.(*APIKeyInput)
				return p, cmd
			}
			u, cmd := p.list.Update(msg)
			p.list = u.(providerListModel)
			return p, cmd
		}
	case tea.PasteMsg:
		if p.needsAPIKey {
			u, cmd := p.apiKeyInput.Update(msg)
			p.apiKeyInput = u.(*APIKeyInput)
			return p, cmd
		}
		var cmd tea.Cmd
		u, cmd := p.list.Update(msg)
		p.list = u.(providerListModel)
		return p, cmd
	}
	return p, nil
}

func (p *providerDialogCmp) View() string {
	t := styles.CurrentTheme()
	if p.needsAPIKey {
		p.keyMap.isAPIKeyHelp = true
		p.keyMap.isAPIKeyValid = false
		apiKeyView := p.apiKeyInput.View()
		apiKeyView = t.S().Base.Width(p.width - 3).Height(lipgloss.Height(apiKeyView)).PaddingLeft(1).Render(apiKeyView)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			t.S().Base.Padding(0, 1, 1, 1).Render(core.Title(p.apiKeyInput.GetTitle(), p.width-4)),
			apiKeyView,
			"",
			t.S().Base.Width(p.width-2).PaddingLeft(1).AlignHorizontal(lipgloss.Left).Render(p.help.View(p.keyMap)),
		)
		return p.style().Render(content)
	}

	listView := p.list.View()
	radio := ""
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		t.S().Base.Padding(0, 1, 1, 1).Render(core.Title("Set API keys", p.width-lipgloss.Width(radio)-5)+" "+radio),
		listView,
		"",
		t.S().Base.Width(p.width-2).PaddingLeft(1).AlignHorizontal(lipgloss.Left).Render(p.help.View(p.keyMap)),
	)
	return p.style().Render(content)
}

func (p *providerDialogCmp) Cursor() *tea.Cursor {
	if p.needsAPIKey {
		cursor := p.apiKeyInput.Cursor()
		if cursor != nil {
			cursor = p.moveCursor(cursor)
			return cursor
		}
	} else {
		cursor := p.list.Cursor()
		if cursor != nil {
			cursor = p.moveCursor(cursor)
			return cursor
		}
	}
	return nil
}

func (p *providerDialogCmp) style() lipgloss.Style {
	t := styles.CurrentTheme()
	return t.S().Base.Width(p.width).Border(lipgloss.RoundedBorder()).BorderForeground(t.BorderFocus)
}

func (p *providerDialogCmp) moveCursor(cursor *tea.Cursor) *tea.Cursor {
	row, col := p.Position()
	offset := row + 3
	cursor.Y += offset
	cursor.X = cursor.X + col + 2
	return cursor
}

func (p *providerDialogCmp) ID() dialogs.DialogID {
	return ProvidersDialogID
}

func (p *providerDialogCmp) Position() (int, int) {
	row := p.wHeight/4 - 2
	col := p.wWidth / 2
	col -= p.width / 2
	return row, col
}
