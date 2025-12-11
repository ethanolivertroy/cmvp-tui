package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethanolivertroy/nist-cmvp-cli/internal/api"
	"github.com/ethanolivertroy/nist-cmvp-cli/internal/model"
)

// ViewState represents the current view of the application
type ViewState int

const (
	ViewList ViewState = iota
	ViewDetail
)

// ModulesLoadedMsg is sent when modules are loaded from the API
type ModulesLoadedMsg struct {
	Modules []list.Item
}

// ErrorMsg is sent when an error occurs
type ErrorMsg struct {
	Err error
}

// Model is the main application model
type Model struct {
	list            list.Model
	allModules      []list.Item
	spinner         spinner.Model
	loading         bool
	err             error
	width           int
	height          int
	view            ViewState
	selectedModule  *model.ModuleItem
	apiClient       *api.Client
	showAlgoDetails bool // Toggle between algorithm categories and detailed list
}

// NewModel creates a new application model
func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(PrimaryColor)

	return Model{
		spinner:   s,
		loading:   true,
		view:      ViewList,
		apiClient: api.NewClient(),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.fetchModules(),
	)
}

func (m Model) fetchModules() tea.Cmd {
	return func() tea.Msg {
		modules, err := m.apiClient.FetchAllModules()
		if err != nil {
			return ErrorMsg{Err: err}
		}

		items := make([]list.Item, len(modules))
		for i, mod := range modules {
			items[i] = model.ModuleItem{Module: mod}
		}
		return ModulesLoadedMsg{Modules: items}
	}
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't handle keys while filtering
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "ctrl+c", "q":
			if m.view == ViewDetail {
				m.view = ViewList
				return m, nil
			}
			return m, tea.Quit
		case "enter":
			if m.view == ViewList && !m.loading {
				if item, ok := m.list.SelectedItem().(model.ModuleItem); ok {
					m.selectedModule = &item
					m.view = ViewDetail
					return m, nil
				}
			}
		case "esc", "backspace":
			if m.view == ViewDetail {
				m.view = ViewList
				return m, nil
			}
		case "d":
			if m.view == ViewDetail {
				m.showAlgoDetails = !m.showAlgoDetails
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.loading {
			m.list.SetSize(msg.Width-4, msg.Height-4)
		}
		return m, nil

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case ModulesLoadedMsg:
		m.loading = false
		m.allModules = msg.Modules

		delegate := NewModuleDelegate()
		m.list = list.New(msg.Modules, delegate, m.width-4, m.height-4)
		m.list.Title = "NIST CMVP Modules"
		m.list.SetShowStatusBar(true)
		m.list.SetFilteringEnabled(true)
		m.list.Styles.Title = TitleStyle
		m.list.FilterInput.Prompt = "Filter: "
		m.list.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(PrimaryColor)

		// Use exact substring matching instead of fuzzy matching
		m.list.Filter = func(term string, targets []string) []list.Rank {
			var ranks []list.Rank
			term = strings.ToLower(term)
			for i, target := range targets {
				if strings.Contains(strings.ToLower(target), term) {
					ranks = append(ranks, list.Rank{Index: i})
				}
			}
			return ranks
		}

		return m, nil

	case ErrorMsg:
		m.loading = false
		m.err = msg.Err
		return m, nil
	}

	// Pass messages to list component when in list view
	if m.view == ViewList && !m.loading && m.err == nil {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the model
func (m Model) View() string {
	if m.loading {
		return AppStyle.Render(
			fmt.Sprintf("\n\n   %s Loading CMVP modules...\n\n", m.spinner.View()),
		)
	}

	if m.err != nil {
		return AppStyle.Render(
			lipgloss.NewStyle().Foreground(ErrorColor).Render(
				fmt.Sprintf("\n\n   Error: %v\n\n   Press q to quit.", m.err),
			),
		)
	}

	switch m.view {
	case ViewDetail:
		return m.renderDetailView()
	default:
		return AppStyle.Render(m.list.View())
	}
}

func (m Model) renderDetailView() string {
	if m.selectedModule == nil {
		return ""
	}

	mod := m.selectedModule

	var b strings.Builder

	// Title with status badge and level badge
	b.WriteString(DetailTitleStyle.Render(mod.ModuleName))
	b.WriteString("  ")
	b.WriteString(StatusBadge(mod.Status))
	if mod.OverallLevel > 0 {
		b.WriteString("  ")
		b.WriteString(LevelBadge(mod.OverallLevel))
	}
	b.WriteString("\n\n")

	// Caveat warning (displayed prominently if present)
	if mod.Caveat != "" {
		b.WriteString(DetailLabelStyle.Render("CAVEAT:"))
		b.WriteString("\n")
		b.WriteString(CaveatStyle.Render(mod.Caveat))
		b.WriteString("\n\n")
	}

	// Details grid
	details := []struct {
		label string
		value string
		isURL bool
	}{
		{"Certificate #:", mod.CertificateNumber, false},
		{"Vendor:", mod.VendorName, false},
		{"Module Type:", mod.ModuleType, false},
		{"Standard:", mod.Standard, false},
		{"Embodiment:", mod.Embodiment, false},
		{"Lab:", mod.Lab, false},
	}

	// Add validation date if available
	if !mod.ValidationDate.IsZero() {
		details = append(details, struct {
			label string
			value string
			isURL bool
		}{"Validation Date:", mod.ValidationDate.Format("January 2, 2006"), false})
	}

	// Add sunset date if available
	if mod.SunsetDate != "" {
		details = append(details, struct {
			label string
			value string
			isURL bool
		}{"Sunset Date:", mod.SunsetDate, false})
	}

	// Add URLs
	if mod.CertificateURL != "" {
		details = append(details, struct {
			label string
			value string
			isURL bool
		}{"NIST URL:", mod.CertificateURL, true})
	}

	if mod.SecurityPolicyURL != "" {
		details = append(details, struct {
			label string
			value string
			isURL bool
		}{"Security Policy:", mod.SecurityPolicyURL, true})
	}

	for _, d := range details {
		if d.value == "" {
			continue
		}
		b.WriteString(DetailLabelStyle.Render(d.label))
		if d.isURL {
			b.WriteString(DetailURLStyle.Render(d.value))
		} else {
			b.WriteString(DetailValueStyle.Render(d.value))
		}
		b.WriteString("\n")
	}

	// Description (if available)
	if mod.Module.Description != "" {
		b.WriteString("\n")
		b.WriteString(DetailLabelStyle.Render("Description:"))
		b.WriteString("\n")
		b.WriteString(DescriptionStyle.Render(mod.Module.Description))
		b.WriteString("\n")
	}

	// Algorithms (toggle between categories and detailed view)
	if m.showAlgoDetails {
		b.WriteString("\n")
		b.WriteString(DetailLabelStyle.Render("Algorithms (Detailed):"))
		b.WriteString("\n")
		if len(mod.AlgorithmsDetailed) > 0 {
			for _, algo := range mod.AlgorithmsDetailed {
				b.WriteString("  • ")
				b.WriteString(DetailValueStyle.Render(algo))
				b.WriteString("\n")
			}
		} else {
			b.WriteString(HelpStyle.Render("  (No detailed algorithm data available yet)"))
			b.WriteString("\n")
		}
	} else if len(mod.Algorithms) > 0 {
		b.WriteString("\n")
		b.WriteString(DetailLabelStyle.Render("Algorithms:"))
		b.WriteString("\n")
		for _, algo := range mod.Algorithms {
			b.WriteString(AlgorithmStyle.Render(algo))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("Press ESC or Backspace to return to list • Press d to toggle algorithm details"))

	return AppStyle.Render(b.String())
}
